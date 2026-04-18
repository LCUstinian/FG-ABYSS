use std::path::PathBuf;
use std::sync::Arc;
use tokio::sync::RwLock;
use crate::Result;
use crate::infra::config::{self, AppearanceConfig, Config, ConnectionConfig};
use crate::infra::crypto::CryptoContext;

pub struct SettingsService {
    config_path: PathBuf,
    config:      Arc<RwLock<Config>>,
}

impl SettingsService {
    pub fn new(config_path: PathBuf, config: Config) -> (Self, Arc<RwLock<Config>>) {
        let shared = Arc::new(RwLock::new(config));
        (Self { config_path, config: shared.clone() }, shared)
    }

    pub async fn get(&self) -> Config {
        self.config.read().await.clone()
    }

    pub async fn update_appearance(
        &self,
        theme:    String,
        language: String,
        accent:   String,
    ) -> Result<()> {
        let mut cfg = self.config.write().await;
        cfg.appearance = AppearanceConfig { theme, language, accent_color: accent };
        config::save(&self.config_path, &cfg)
    }

    pub async fn update_connection(&self, input: ConnectionConfig) -> Result<()> {
        let mut cfg = self.config.write().await;
        cfg.connection = input;
        config::save(&self.config_path, &cfg)
    }

    pub async fn set_master_password(
        &self,
        crypto:       &CryptoContext,
        old_password: Option<&str>,
        new_password: &str,
    ) -> Result<()> {
        let mut cfg = self.config.write().await;
        if cfg.security.master_password_enabled {
            match old_password {
                Some(old) if crypto.verify_password(old, &cfg.security.master_password_hash) => {}
                Some(_) => return Err(crate::AppError::InvalidInput("old password incorrect".into())),
                None    => return Err(crate::AppError::InvalidInput("old password required".into())),
            }
        }
        cfg.security.master_password_hash    = CryptoContext::hash_password(new_password)?;
        cfg.security.master_password_enabled = true;
        config::save(&self.config_path, &cfg)
    }
}
