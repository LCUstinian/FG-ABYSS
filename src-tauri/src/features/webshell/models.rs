use std::collections::HashMap;
use serde::{Deserialize, Serialize};
use crate::infra::crypto::CryptoChain;

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "lowercase")]
pub enum PayloadType {
    Php,
    Jsp,
    Asp,
    Aspx,
}

impl PayloadType {
    pub fn extension(&self) -> &'static str {
        match self {
            Self::Php  => "php",
            Self::Jsp  => "jsp",
            Self::Asp  => "asp",
            Self::Aspx => "aspx",
        }
    }

    pub fn is_memshell(&self) -> bool {
        false
    }
}

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct Webshell {
    pub id:               String,
    pub name:             String,
    pub url:              String,
    pub password:         String,
    pub payload_type:     PayloadType,
    pub project_id:       Option<String>,
    pub status:           String,
    pub tags:             Vec<String>,
    pub custom_headers:   HashMap<String, String>,
    pub cookies:          HashMap<String, String>,
    pub proxy_override:   Option<String>,
    pub http_method:      String,
    pub c2_profile:       String,
    pub crypto_chain:     CryptoChain,
    pub fingerprint:      Option<WebshellFingerprint>,
    pub notes:            Option<String>,
    pub last_connected_at: Option<i64>,
    pub created_at:       i64,
    pub updated_at:       i64,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct WebshellFingerprint {
    pub os:             Option<String>,
    pub runtime:        Option<String>,
    pub runtime_ver:    Option<String>,
    pub username:       Option<String>,
    pub cwd:            Option<String>,
    pub disabled_fns:   Vec<String>,
    pub writable:       Option<bool>,
    pub last_probed_at: Option<i64>,
}

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct ConnectionResult {
    pub success:     bool,
    pub latency_ms:  Option<u64>,
    pub error:       Option<String>,
    pub fingerprint: Option<WebshellFingerprint>,
}

#[derive(Debug, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct CreateWebshellInput {
    pub name:             String,
    pub url:              String,
    pub password:         String,
    pub payload_type:     PayloadType,
    pub project_id:       Option<String>,
    #[serde(default)]
    pub tags:             Vec<String>,
    #[serde(default)]
    pub custom_headers:   HashMap<String, String>,
    #[serde(default)]
    pub cookies:          HashMap<String, String>,
    pub proxy_override:   Option<String>,
    #[serde(default = "default_post")]
    pub http_method:      String,
    #[serde(default = "default_c2")]
    pub c2_profile:       String,
    #[serde(default)]
    pub crypto_chain:     CryptoChain,
    pub notes:            Option<String>,
}

fn default_post() -> String { "post".to_string() }
fn default_c2()   -> String { "default".to_string() }

#[derive(Debug, Deserialize, Default, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct UpdateWebshellInput {
    pub name:           Option<String>,
    pub url:            Option<String>,
    pub password:       Option<String>,
    pub project_id:     Option<String>,
    pub tags:           Option<Vec<String>>,
    pub custom_headers: Option<HashMap<String, String>>,
    pub cookies:        Option<HashMap<String, String>>,
    pub proxy_override: Option<String>,
    pub http_method:    Option<String>,
    pub c2_profile:     Option<String>,
    pub crypto_chain:   Option<CryptoChain>,
    pub notes:          Option<String>,
    pub status:         Option<String>,
}
