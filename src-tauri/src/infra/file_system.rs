use std::path::{Path, PathBuf};

pub fn get_app_dir() -> PathBuf {
    dirs::config_dir()
        .unwrap_or_else(|| PathBuf::from("."))
        .join("fg-abyss")
}

pub fn ensure_dir_exists(path: &Path) -> std::io::Result<()> {
    std::fs::create_dir_all(path)
}
