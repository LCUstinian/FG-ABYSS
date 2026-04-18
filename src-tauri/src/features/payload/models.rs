use serde::{Deserialize, Serialize};
use crate::infra::crypto::CryptoChain;
use crate::features::webshell::models::PayloadType;

#[derive(Debug, Clone, Serialize, Deserialize, Default, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct ObfuscateConfig {
    pub level:          u8,
    pub seed:           Option<u64>,
    pub target_version: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct Payload {
    pub id:           String,
    pub name:         String,
    pub payload_type: PayloadType,
    pub config:       serde_json::Value,
    pub created_at:   i64,
}

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct PayloadConfig {
    pub payload_type:   PayloadType,
    pub password:       String,
    pub crypto_chain:   CryptoChain,
    pub c2_profile:     String,
    pub obfuscate:      ObfuscateConfig,
    pub target_version: Option<String>,
}

#[derive(Debug, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct CreatePayloadInput {
    pub name:   String,
    pub config: PayloadConfig,
}

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct PayloadHistoryEntry {
    pub id:               String,
    pub payload_id:       Option<String>,
    pub webshell_id:      Option<String>,
    pub code:             String,
    pub template_version: String,
    pub created_at:       i64,
}
