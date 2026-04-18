use std::sync::Arc;
use dashmap::DashMap;
use tokio::sync::{mpsc, oneshot};
use tokio::task::JoinHandle;
use crate::{AppError, Result};
use crate::infra::http::{HttpClientPool, WebshellHttpConfig};

pub struct QueuedRequest {
    pub config:     WebshellHttpConfig,
    pub body:       Vec<u8>,
    pub respond_to: oneshot::Sender<Result<reqwest::Response>>,
}

pub struct WebshellQueue {
    workers:   DashMap<String, (mpsc::Sender<QueuedRequest>, JoinHandle<()>)>,
    http_pool: Arc<HttpClientPool>,
}

impl WebshellQueue {
    pub fn new(http_pool: Arc<HttpClientPool>) -> Self {
        Self { workers: DashMap::new(), http_pool }
    }

    pub async fn send(
        &self,
        webshell_id: &str,
        config: WebshellHttpConfig,
        body: Vec<u8>,
    ) -> Result<reqwest::Response> {
        let (resp_tx, resp_rx) = oneshot::channel();

        let tx = if let Some(entry) = self.workers.get(webshell_id) {
            entry.value().0.clone()
        } else {
            let (tx, rx) = mpsc::channel::<QueuedRequest>(32);
            let pool   = self.http_pool.clone();
            let id     = webshell_id.to_string();
            let handle = tokio::spawn(queue_worker(rx, pool, id));
            self.workers.insert(webshell_id.to_string(), (tx.clone(), handle));
            tx
        };

        tx.send(QueuedRequest { config, body, respond_to: resp_tx }).await
            .map_err(|_| AppError::Connection("queue dropped".into()))?;
        resp_rx.await
            .map_err(|_| AppError::Connection("response channel closed".into()))?
    }

    pub fn cleanup(&self, webshell_id: &str) {
        if let Some((_, (_, handle))) = self.workers.remove(webshell_id) {
            handle.abort();
        }
    }

    pub fn shutdown_all(&self) {
        let keys: Vec<String> = self.workers.iter().map(|e| e.key().clone()).collect();
        for key in keys {
            if let Some((_, (_, handle))) = self.workers.remove(&key) {
                handle.abort();
            }
        }
    }
}

async fn queue_worker(
    mut rx: mpsc::Receiver<QueuedRequest>,
    pool:   Arc<HttpClientPool>,
    id:     String,
) {
    while let Some(req) = rx.recv().await {
        if let Some((min, max)) = req.config.jitter_ms {
            if max > min {
                let delay = rand::random::<u64>() % (max - min) + min;
                tokio::time::sleep(tokio::time::Duration::from_millis(delay)).await;
            }
        }

        let client = pool.get_client(&req.config);
        let url    = req.config.url.clone();
        let method = req.config.method.to_uppercase();

        let mut builder = match method.as_str() {
            "GET" => client.get(&url),
            _     => client.post(&url),
        };
        for (k, v) in &req.config.custom_headers {
            builder = builder.header(k, v);
        }
        if !req.config.cookies.is_empty() {
            let cookie_str = req.config.cookies.iter()
                .map(|(k, v)| format!("{k}={v}"))
                .collect::<Vec<_>>()
                .join("; ");
            builder = builder.header("Cookie", cookie_str);
        }
        use base64::Engine as _;
        let b64_body = base64::engine::general_purpose::STANDARD.encode(&req.body);
        builder = builder
            .form(&[(&req.config.request_param, b64_body)])
            .timeout(std::time::Duration::from_secs(req.config.timeout_secs));

        let result = pool.send_raw(&id, builder).await;
        let _ = req.respond_to.send(result);
    }
}
