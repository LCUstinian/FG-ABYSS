use serde::{Deserialize, Serialize};
use std::path::PathBuf;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AppConfig {
    pub theme: String,
    pub language: String,
    pub window_width: u32,
    pub window_height: u32,
}

impl Default for AppConfig {
    fn default() -> Self {
        Self {
            theme: "dark".to_string(),
            language: "zh".to_string(),
            window_width: 1280,
            window_height: 720,
        }
    }
}

pub fn get_config_path() -> PathBuf {
    dirs::config_dir()
        .unwrap_or_else(|| PathBuf::from("."))
        .join("fg-abyss")
        .join("config.json")
}
