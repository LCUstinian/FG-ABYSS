use rusqlite::Result;
use std::sync::Arc;
use tokio::sync::Mutex;

use crate::infra::database::models::Project;

/// 项目管理器
pub struct ProjectManager {
    db: Arc<Mutex<rusqlite::Connection>>,
}

impl ProjectManager {
    /// 创建新的项目管理器
    pub fn new(db: Arc<Mutex<rusqlite::Connection>>) -> Self {
        Self { db }
    }

    /// 创建项目
    pub async fn create(&self, name: String, description: Option<String>) -> Result<Project> {
        let project = Project::new(name, description);
        let project_id = project.id.clone();
        let project_name = project.name.clone();
        let project_description = project.description.clone().unwrap_or_default();
        let project_created_at = project.created_at;
        let project_updated_at = project.updated_at;
        
        let conn = self.db.lock().await;
        conn.execute(
            "INSERT INTO projects (id, name, description, created_at, updated_at) VALUES (?1, ?2, ?3, ?4, ?5)",
            rusqlite::params![
                project_id,
                project_name,
                project_description,
                project_created_at,
                project_updated_at,
            ],
        )?;
        
        Ok(project)
    }

    /// 更新项目
    pub async fn update(&self, id: &str, name: String, description: Option<String>) -> Result<Project> {
        let conn = self.db.lock().await;
        let updated_at = chrono::Utc::now().timestamp();
        
        conn.execute(
            "UPDATE projects SET name = ?1, description = ?2, updated_at = ?3 WHERE id = ?4 AND deleted_at IS NULL",
            rusqlite::params![name, description.unwrap_or_default(), updated_at, id],
        )?;
        
        Ok(Project {
            id: id.to_string(),
            name,
            description: None,
            created_at: 0,
            updated_at,
            deleted_at: None,
        })
    }

    /// 删除项目（软删除）
    pub async fn delete(&self, id: &str) -> Result<()> {
        let conn = self.db.lock().await;
        let deleted_at = chrono::Utc::now().timestamp();
        
        conn.execute(
            "UPDATE projects SET deleted_at = ?1 WHERE id = ?2",
            rusqlite::params![deleted_at, id],
        )?;
        
        Ok(())
    }

    /// 恢复项目
    pub async fn restore(&self, id: &str) -> Result<()> {
        let conn = self.db.lock().await;
        
        conn.execute(
            "UPDATE projects SET deleted_at = NULL WHERE id = ?1",
            rusqlite::params![id],
        )?;
        
        Ok(())
    }

    /// 获取所有项目（未删除的）
    pub async fn get_all(&self) -> Result<Vec<Project>> {
        let conn = self.db.lock().await;
        
        let mut stmt = conn.prepare(
            "SELECT id, name, description, created_at, updated_at FROM projects WHERE deleted_at IS NULL"
        )?;
        
        let projects = stmt.query_map([], |row| {
            Ok(Project {
                id: row.get(0)?,
                name: row.get(1)?,
                description: row.get(2)?,
                created_at: row.get(3)?,
                updated_at: row.get(4)?,
                deleted_at: None,
            })
        })?;
        
        let mut result = Vec::new();
        for project in projects {
            result.push(project?);
        }
        
        Ok(result)
    }

    /// 获取已删除的项目
    pub async fn get_deleted(&self) -> Result<Vec<Project>> {
        let conn = self.db.lock().await;
        
        let mut stmt = conn.prepare(
            "SELECT id, name, description, created_at, updated_at FROM projects WHERE deleted_at IS NOT NULL"
        )?;
        
        let projects = stmt.query_map([], |row| {
            Ok(Project {
                id: row.get(0)?,
                name: row.get(1)?,
                description: row.get(2)?,
                created_at: row.get(3)?,
                updated_at: row.get(4)?,
                deleted_at: row.get(4)?,
            })
        })?;
        
        let mut result = Vec::new();
        for project in projects {
            result.push(project?);
        }
        
        Ok(result)
    }
}
