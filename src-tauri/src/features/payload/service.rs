use std::sync::Arc;
use uuid::Uuid;
use chrono::Utc;
use crate::Result;
use crate::infra::crypto::CryptoContext;
use super::models::{Payload, CreatePayloadInput, PayloadHistoryEntry, PayloadConfig};
use super::repo::PayloadRepo;
use super::generator;

pub struct PayloadService {
    repo:   Arc<dyn PayloadRepo>,
    crypto: Arc<CryptoContext>,
}

impl PayloadService {
    pub fn new(repo: Arc<dyn PayloadRepo>, crypto: Arc<CryptoContext>) -> Self {
        Self { repo, crypto }
    }

    pub async fn list(&self) -> Result<Vec<Payload>> { self.repo.find_all().await }
    pub async fn get(&self, id: &str) -> Result<Payload> { self.repo.find_by_id(id).await }

    pub async fn create(&self, input: CreatePayloadInput) -> Result<Payload> {
        let now = Utc::now().timestamp();
        let p = Payload {
            id:           Uuid::new_v4().to_string(),
            name:         input.name,
            payload_type: input.config.payload_type.clone(),
            config:       serde_json::to_value(&input.config)?,
            created_at:   now,
        };
        self.repo.insert(&p).await?;
        Ok(p)
    }

    pub async fn generate_payload(
        &self,
        payload_id: Option<&str>,
        webshell_id: Option<&str>,
        config: PayloadConfig,
    ) -> Result<String> {
        let code = generator::generate(&config)?;
        let entry = PayloadHistoryEntry {
            id:               Uuid::new_v4().to_string(),
            payload_id:       payload_id.map(|s| s.to_string()),
            webshell_id:      webshell_id.map(|s| s.to_string()),
            code:             code.clone(),
            template_version: "v1".into(),
            created_at:       Utc::now().timestamp(),
        };
        self.repo.insert_history(&entry).await?;
        Ok(code)
    }

    pub async fn list_history(&self, payload_id: &str) -> Result<Vec<PayloadHistoryEntry>> {
        self.repo.find_history(payload_id).await
    }
}
