use std::path::PathBuf;
use std::sync::Arc;
use crate::Result;
use super::models::Plugin;
use super::repo::PluginRepo;

pub struct PluginService {
    repo:        Arc<dyn PluginRepo>,
    plugins_dir: PathBuf,
}

impl PluginService {
    pub fn new(repo: Arc<dyn PluginRepo>, plugins_dir: PathBuf) -> Self {
        Self { repo, plugins_dir }
    }

    pub async fn list(&self) -> Result<Vec<Plugin>> { self.repo.find_all().await }
    pub async fn get(&self, id: &str) -> Result<Plugin> { self.repo.find_by_id(id).await }

    pub async fn enable(&self, id: &str) -> Result<Plugin> {
        self.repo.set_enabled(id, true).await?;
        self.repo.find_by_id(id).await
    }

    pub async fn disable(&self, id: &str) -> Result<Plugin> {
        self.repo.set_enabled(id, false).await?;
        self.repo.find_by_id(id).await
    }
}
