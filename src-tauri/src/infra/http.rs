use std::collections::HashMap;
use std::time::Instant;
use dashmap::DashMap;
use crate::{AppError, Result};
use crate::infra::config::ConnectionConfig;

const FAILURE_THRESHOLD: u32 = 3;
const CIRCUIT_RESET_SECS: u64 = 30;

#[derive(Debug, Clone)]
pub struct WebshellHttpConfig {
    pub url:            String,
    pub method:         String,
    pub request_param:  String,
    pub proxy_override: Option<String>,
    pub custom_headers: HashMap<String, String>,
    pub cookies:        HashMap<String, String>,
    pub timeout_secs:   u64,
    pub jitter_ms:      Option<(u64, u64)>,
}

impl WebshellHttpConfig {
    pub fn proxy_hash(&self) -> u64 {
        use std::hash::{Hash, Hasher};
        use std::collections::hash_map::DefaultHasher;
        let mut h = DefaultHasher::new();
        self.proxy_override.hash(&mut h);
        h.finish()
    }
}

enum CircuitState { Closed, Open(Instant), HalfOpen }

pub struct HttpClientPool {
    clients:       DashMap<u64, reqwest::Client>,
    circuit:       DashMap<String, CircuitState>,
    failure_count: DashMap<String, u32>,
    global_config: ConnectionConfig,
}

impl HttpClientPool {
    pub fn new(config: &ConnectionConfig) -> Self {
        Self {
            clients:       DashMap::new(),
            circuit:       DashMap::new(),
            failure_count: DashMap::new(),
            global_config: config.clone(),
        }
    }

    pub fn get_client(&self, config: &WebshellHttpConfig) -> reqwest::Client {
        let key = config.proxy_hash();
        self.clients
            .entry(key)
            .or_insert_with(|| build_client(config, &self.global_config))
            .clone()
    }

    pub async fn send_raw(
        &self,
        webshell_id: &str,
        req: reqwest::RequestBuilder,
    ) -> Result<reqwest::Response> {
        let is_open = self.circuit.get(webshell_id).map(|s| {
            match *s {
                CircuitState::Open(ref opened_at) => opened_at.elapsed().as_secs() < CIRCUIT_RESET_SECS,
                _ => false,
            }
        }).unwrap_or(false);

        if is_open {
            return Err(AppError::CircuitOpen(webshell_id.to_string()));
        }

        match req.send().await {
            Ok(resp) => {
                self.failure_count.remove(webshell_id);
                self.circuit.insert(webshell_id.to_string(), CircuitState::Closed);
                Ok(resp)
            }
            Err(e) => {
                let count = {
                    let entry = self.failure_count
                        .entry(webshell_id.to_string())
                        .and_modify(|c| *c += 1)
                        .or_insert(1);
                    *entry
                };
                if count >= FAILURE_THRESHOLD {
                    self.circuit.insert(
                        webshell_id.to_string(),
                        CircuitState::Open(Instant::now()),
                    );
                }
                Err(AppError::Http(e))
            }
        }
    }
}

fn build_client(config: &WebshellHttpConfig, global: &ConnectionConfig) -> reqwest::Client {
    let mut builder = reqwest::Client::builder()
        .timeout(std::time::Duration::from_secs(config.timeout_secs));

    if let Some(proxy_url) = &config.proxy_override {
        if let Ok(proxy) = reqwest::Proxy::all(proxy_url) {
            builder = builder.proxy(proxy);
        }
    } else if global.proxy_enabled && !global.proxy_host.is_empty() {
        let proxy_url = format!(
            "{}://{}:{}",
            global.proxy_type, global.proxy_host, global.proxy_port
        );
        if let Ok(proxy) = reqwest::Proxy::all(&proxy_url) {
            builder = builder.proxy(proxy);
        }
    }

    builder.build().unwrap_or_default()
}
