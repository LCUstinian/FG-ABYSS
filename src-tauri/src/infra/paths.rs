// src-tauri/src/infra/paths.rs
use std::path::PathBuf;
use tauri::{AppHandle, Manager};
use crate::Result;

pub struct AppPaths {
    pub data_dir:    PathBuf,
    pub db_path:     PathBuf,
    pub config:      PathBuf,
    pub logs_dir:    PathBuf,
    pub plugins_dir: PathBuf,
    pub exports_dir: PathBuf,
}

impl AppPaths {
    pub fn resolve(app: &AppHandle) -> Result<Self> {
        let base = app.path().app_data_dir()
            .map_err(|e: tauri::Error| crate::AppError::Io(std::io::Error::new(
                std::io::ErrorKind::Other, e.to_string()
            )))?;
        Ok(Self {
            db_path:     base.join("data.db"),
            config:      base.join("config.toml"),
            logs_dir:    base.join("logs"),
            plugins_dir: base.join("plugins"),
            exports_dir: base.join("exports"),
            data_dir:    base,
        })
    }
}
