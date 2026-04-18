use std::sync::Arc;
use dashmap::DashMap;
use uuid::Uuid;
use chrono::Utc;
use crate::{AppError, Result};
use crate::infra::crypto::CryptoContext;
use super::models::{Webshell, CreateWebshellInput, UpdateWebshellInput, ConnectionResult};
use super::repo::{WebshellRepo, WebshellUpdateFields};
use super::session::WebshellSession;
use super::queue::WebshellQueue;

pub struct WebshellService {
    repo:     Arc<dyn WebshellRepo>,
    crypto:   Arc<CryptoContext>,
    pub queue: Arc<WebshellQueue>,
    sessions: Arc<DashMap<String, WebshellSession>>,
}

impl WebshellService {
    pub fn new(
        repo:   Arc<dyn WebshellRepo>,
        crypto: Arc<CryptoContext>,
        queue:  Arc<WebshellQueue>,
    ) -> Self {
        Self { repo, crypto, queue, sessions: Arc::new(DashMap::new()) }
    }

    pub async fn list(&self, project_id: Option<&str>) -> Result<Vec<Webshell>> {
        let mut shells = self.repo.find_all(project_id).await?;
        for s in &mut shells {
            s.password = self.crypto.decrypt_str(&s.password).unwrap_or_default();
        }
        Ok(shells)
    }

    pub async fn get(&self, id: &str) -> Result<Webshell> {
        let mut w = self.repo.find_by_id(id).await?;
        w.password = self.crypto.decrypt_str(&w.password)?;
        Ok(w)
    }

    pub async fn create(&self, input: CreateWebshellInput) -> Result<Webshell> {
        let url = normalize_url(&input.url)?;

        // Uniqueness check
        if self.repo.find_by_url(&url).await?.is_some() {
            return Err(AppError::InvalidInput(format!("WebShell URL already exists: {url}")));
        }

        let enc_password = self.crypto.encrypt(input.password.as_bytes())?;
        let plain_pw     = input.password.clone();
        let now          = Utc::now().timestamp();
        let id           = Uuid::new_v4().to_string();

        let w = Webshell {
            id, name: input.name, url,
            password:       enc_password,
            payload_type:   input.payload_type,
            project_id:     input.project_id,
            status:         "inactive".into(),
            tags:           input.tags,
            custom_headers: input.custom_headers,
            cookies:        input.cookies,
            proxy_override: input.proxy_override,
            http_method:    input.http_method,
            c2_profile:     input.c2_profile,
            crypto_chain:   input.crypto_chain,
            fingerprint:    None,
            notes:          input.notes,
            last_connected_at: None,
            created_at:     now,
            updated_at:     now,
        };

        self.repo.insert(&w).await?;

        let mut result = w;
        result.password = plain_pw;
        Ok(result)
    }

    pub async fn update(&self, id: &str, input: UpdateWebshellInput) -> Result<Webshell> {
        let current = self.get(id).await?;

        // Detect crypto_chain change → needs_redeploy
        let new_status = match &input.crypto_chain {
            Some(chain) if *chain != current.crypto_chain => "needs_redeploy".to_string(),
            _ => input.status.clone().unwrap_or(current.status.clone()),
        };

        let enc_password = match &input.password {
            Some(p) => Some(self.crypto.encrypt(p.as_bytes())?),
            None    => None,
        };

        let url = match &input.url {
            Some(u) => Some(normalize_url(u)?),
            None => None,
        };

        self.repo.update_fields(id, WebshellUpdateFields {
            name:           input.name,
            url,
            password:       enc_password,
            // Single Option<String> in input → Option<Option<String>> in fields:
            // None (not provided) → None (skip update)
            // Some(val) → Some(Some(val)) (set to val)
            project_id:     input.project_id.map(Some),
            tags:           input.tags.as_ref().and_then(|t| serde_json::to_string(t).ok()),
            custom_headers: input.custom_headers.as_ref().and_then(|h| serde_json::to_string(h).ok()),
            cookies:        input.cookies.as_ref().and_then(|c| serde_json::to_string(c).ok()),
            proxy_override: input.proxy_override.map(Some),
            http_method:    input.http_method,
            c2_profile:     input.c2_profile,
            crypto_chain:   input.crypto_chain.as_ref().and_then(|c| serde_json::to_string(c).ok()),
            notes:          input.notes.map(Some),
            status:         Some(new_status),
        }).await?;

        self.get(id).await
    }

    pub async fn delete(&self, id: &str) -> Result<()> {
        let ts = Utc::now().timestamp();
        self.sessions.remove(id);
        self.queue.cleanup(id);
        self.repo.soft_delete(id, ts).await
    }

    pub async fn test_connection(&self, id: &str) -> Result<ConnectionResult> {
        let shell = self.get(id).await?;
        let start = std::time::Instant::now();
        // TODO: Phase 1 Init protocol — implement in Phase 2
        let _ = shell;
        let _ = start;
        todo!("Phase 1 Init protocol — implement in Phase 2")
    }

    pub async fn exec(
        &self,
        id: &str,
        method: &str,
        args: serde_json::Value,
    ) -> Result<serde_json::Value> {
        let shell = self.get(id).await?;
        if shell.status == "needs_redeploy" {
            return Err(AppError::NeedsRedeploy(
                format!("WebShell {} payload changed, redeploy before exec", id)
            ));
        }
        let _ = (method, args);
        todo!("Phase 2 Exec protocol — implement in Phase 2")
    }

    pub async fn reset_redeploy_status(&self, id: &str) -> Result<Webshell> {
        self.repo.update_status(id, "inactive", None).await?;
        self.sessions.remove(id);
        self.get(id).await
    }
}

fn normalize_url(raw: &str) -> Result<String> {
    let trimmed = raw.trim();
    let parsed  = url::Url::parse(trimmed)
        .map_err(|e| AppError::InvalidInput(format!("invalid URL: {e}")))?;
    let scheme   = parsed.scheme();
    let host     = parsed.host_str().unwrap_or("");
    let port_str = match parsed.port() {
        Some(p) => format!(":{p}"),
        None    => String::new(),
    };
    let path = parsed.path().trim_end_matches('/');
    let path = if path.is_empty() { "/" } else { path };
    Ok(format!("{scheme}://{host}{port_str}{path}"))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_normalize_url() {
        // url crate normalizes scheme and host to lowercase per RFC 3986
        assert_eq!(normalize_url("HTTP://Example.COM/shell.php").unwrap(), "http://example.com/shell.php");
        assert_eq!(normalize_url("https://example.com/shell.php/").unwrap(), "https://example.com/shell.php");
        assert!(normalize_url("not-a-url").is_err());
    }
}
