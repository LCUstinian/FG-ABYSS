use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct Project {
    pub id:          String,
    pub name:        String,
    pub description: Option<String>,
    pub created_at:  i64,
    pub updated_at:  i64,
}

#[derive(Debug, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct CreateProjectInput {
    pub name:        String,
    pub description: Option<String>,
}

#[derive(Debug, Deserialize, Default, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct UpdateProjectInput {
    pub name:        Option<String>,
    pub description: Option<String>,
}
