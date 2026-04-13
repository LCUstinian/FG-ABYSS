use serde::{Deserialize, Serialize};

/// 项目实体
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Project {
    pub id: String,
    pub name: String,
    pub description: Option<String>,
    pub created_at: i64,
    pub updated_at: i64,
    pub deleted_at: Option<i64>,
}

impl Project {
    pub fn new(name: String, description: Option<String>) -> Self {
        let now = chrono::Utc::now().timestamp();
        Self {
            id: uuid::Uuid::new_v4().to_string(),
            name,
            description,
            created_at: now,
            updated_at: now,
            deleted_at: None,
        }
    }
}

/// WebShell 实体
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WebShell {
    pub id: String,
    pub project_id: Option<String>,
    pub name: String,
    pub url: String,
    pub password: String,
    pub payload_type: String,
    pub encryption: String,
    pub status: String,
    pub last_connected_at: Option<i64>,
    pub created_at: i64,
    pub updated_at: i64,
    pub deleted_at: Option<i64>,
}

impl WebShell {
    pub fn new(
        name: String,
        url: String,
        password: String,
        payload_type: String,
    ) -> Self {
        let now = chrono::Utc::now().timestamp();
        Self {
            id: uuid::Uuid::new_v4().to_string(),
            project_id: None,
            name,
            url,
            password,
            payload_type,
            encryption: "aes-256-gcm".to_string(),
            status: "unknown".to_string(),
            last_connected_at: None,
            created_at: now,
            updated_at: now,
            deleted_at: None,
        }
    }
}

/// 载荷配置实体
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PayloadConfig {
    pub id: String,
    pub name: String,
    pub url: String,
    pub password: String,
    pub payload_type: String,
    pub encryption: String,
    pub obfuscation_level: i32,
    pub tags: Option<String>,
    pub group: Option<String>,
    pub created_at: i64,
    pub updated_at: i64,
    pub deleted_at: Option<i64>,
}

/// 载荷模板实体
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PayloadTemplate {
    pub id: String,
    pub name: String,
    pub payload_type: String,
    pub code: String,
    pub is_built_in: bool,
    pub created_at: i64,
}

/// 载荷历史实体
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PayloadHistory {
    pub id: String,
    pub config_id: String,
    pub generated_code: String,
    pub generated_at: i64,
}

/// 插件实体
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Plugin {
    pub id: String,
    pub name: String,
    pub version: String,
    pub enabled: bool,
    pub config: Option<String>,
    pub created_at: i64,
    pub updated_at: i64,
}
