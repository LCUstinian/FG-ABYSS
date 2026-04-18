use std::path::Path;
use serde::{Deserialize, Serialize};
use crate::{AppError, Result};

#[derive(Debug, Serialize, Deserialize, Default, Clone)]
pub struct Config {
    pub meta:       MetaConfig,
    pub appearance: AppearanceConfig,
    pub window:     WindowConfig,
    pub connection: ConnectionConfig,
    pub security:   SecurityConfig,
    pub logging:    LoggingConfig,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct MetaConfig {
    pub config_version: u32,
}
impl Default for MetaConfig {
    fn default() -> Self { Self { config_version: 1 } }
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct AppearanceConfig {
    pub theme: String,
    pub language: String,
    pub accent_color: String,
}
impl Default for AppearanceConfig {
    fn default() -> Self {
        Self { theme: "dark".into(), language: "zh-CN".into(), accent_color: "#4f9cff".into() }
    }
}

#[derive(Debug, Serialize, Deserialize, Clone, Default)]
pub struct WindowConfig {
    pub x: i32, pub y: i32,
    pub width: u32, pub height: u32,
    pub maximized: bool,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ConnectionConfig {
    pub timeout_secs: u64,
    pub retry_count: u32,
    pub proxy_enabled: bool,
    pub proxy_type: String,
    pub proxy_host: String,
    pub proxy_port: u16,
}
impl Default for ConnectionConfig {
    fn default() -> Self {
        Self { timeout_secs: 30, retry_count: 3, proxy_enabled: false,
            proxy_type: "http".into(), proxy_host: "".into(), proxy_port: 7890 }
    }
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SecurityConfig {
    pub encryption: String,
    pub key_rotation_days: u32,
    pub salt: String,
    pub master_password_enabled: bool,
    pub master_password_hash: String,
    pub idle_lock_minutes: u32,
}
impl Default for SecurityConfig {
    fn default() -> Self {
        Self { encryption: "aes-256-gcm".into(), key_rotation_days: 30,
            salt: String::new(), master_password_enabled: false,
            master_password_hash: String::new(), idle_lock_minutes: 30 }
    }
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct LoggingConfig {
    pub level: String,
    pub max_file_size_mb: u64,
    pub max_files: usize,
}
impl Default for LoggingConfig {
    fn default() -> Self {
        Self { level: "info".into(), max_file_size_mb: 10, max_files: 5 }
    }
}

pub fn load_or_default(path: &Path) -> Result<Config> {
    if path.exists() {
        let text = std::fs::read_to_string(path)?;
        let mut cfg: Config = toml::from_str(&text)
            .map_err(|e| AppError::InvalidInput(e.to_string()))?;
        migrate_config(&mut cfg);
        Ok(cfg)
    } else {
        Ok(Config::default())
    }
}

pub fn save(path: &Path, cfg: &Config) -> Result<()> {
    let text = toml::to_string_pretty(cfg)
        .map_err(|e| AppError::InvalidInput(e.to_string()))?;
    let tmp = path.with_extension("toml.tmp");
    std::fs::write(&tmp, &text)?;
    std::fs::rename(&tmp, path)?;
    Ok(())
}

fn migrate_config(cfg: &mut Config) {
    let _ = cfg;
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::path::Path;

    #[test]
    fn test_default_config_serializes() {
        let cfg = Config::default();
        let text = toml::to_string_pretty(&cfg).unwrap();
        let _parsed: Config = toml::from_str(&text).unwrap();
    }

    #[test]
    fn test_load_or_default_missing_file() {
        let cfg = load_or_default(Path::new("/nonexistent/config.toml")).unwrap();
        assert_eq!(cfg.meta.config_version, 1);
        assert_eq!(cfg.appearance.theme, "dark");
    }
}
