use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct Plugin {
    pub id:         String,
    pub name:       String,
    pub version:    String,
    pub enabled:    bool,
    pub config:     serde_json::Value,
    pub source:     String,
    pub created_at: i64,
    pub updated_at: i64,
}
