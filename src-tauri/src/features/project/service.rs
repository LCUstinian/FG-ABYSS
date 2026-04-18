use std::sync::Arc;
use uuid::Uuid;
use chrono::Utc;
use crate::Result;
use super::models::{Project, CreateProjectInput, UpdateProjectInput};
use super::repo::ProjectRepo;

pub struct ProjectService {
    repo: Arc<dyn ProjectRepo>,
}

impl ProjectService {
    pub fn new(repo: Arc<dyn ProjectRepo>) -> Self { Self { repo } }

    pub async fn list(&self) -> Result<Vec<Project>> {
        self.repo.find_all().await
    }

    pub async fn get(&self, id: &str) -> Result<Project> {
        self.repo.find_by_id(id).await
    }

    pub async fn create(&self, input: CreateProjectInput) -> Result<Project> {
        let now = Utc::now().timestamp();
        let p = Project {
            id:          Uuid::new_v4().to_string(),
            name:        input.name,
            description: input.description,
            created_at:  now,
            updated_at:  now,
        };
        self.repo.insert(&p).await?;
        Ok(p)
    }

    pub async fn update(&self, id: &str, input: UpdateProjectInput) -> Result<Project> {
        self.repo.update(id, input.name, input.description).await?;
        self.repo.find_by_id(id).await
    }

    pub async fn delete(&self, id: &str) -> Result<()> {
        let ts = Utc::now().timestamp();
        self.repo.soft_delete_with_webshells(id, ts).await
    }
}
