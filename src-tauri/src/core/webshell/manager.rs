use rusqlite::Result;
use std::sync::Arc;
use tokio::sync::Mutex;

use crate::infra::database::models::WebShell;

/// WebShell 管理器
pub struct WebShellManager {
    db: Arc<Mutex<rusqlite::Connection>>,
}

impl WebShellManager {
    /// 创建新的 WebShell 管理器
    pub fn new(db: Arc<Mutex<rusqlite::Connection>>) -> Self {
        Self { db }
    }

    /// 创建 WebShell
    pub async fn create(
        &self,
        name: String,
        url: String,
        password: String,
        payload_type: String,
        project_id: Option<String>,
    ) -> Result<WebShell> {
        let webshell = WebShell::new(name, url, password, payload_type);
        
        let conn = self.db.lock().await;
        conn.execute(
            "INSERT INTO webshells (id, project_id, name, url, password, payload_type, created_at, updated_at) 
             VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8)",
            rusqlite::params![
                webshell.id,
                project_id.unwrap_or_default(),
                webshell.name,
                webshell.url,
                webshell.password,
                webshell.payload_type,
                webshell.created_at,
                webshell.updated_at,
            ],
        )?;
        
        Ok(webshell)
    }

    /// 更新 WebShell
    pub async fn update(
        &self,
        id: &str,
        name: String,
        url: String,
        password: String,
        payload_type: String,
        project_id: Option<String>,
    ) -> Result<WebShell> {
        let conn = self.db.lock().await;
        let updated_at = chrono::Utc::now().timestamp();
        
        conn.execute(
            "UPDATE webshells SET name = ?1, url = ?2, password = ?3, payload_type = ?4, 
             project_id = ?5, updated_at = ?6 WHERE id = ?7 AND deleted_at IS NULL",
            rusqlite::params![name, url, password, payload_type, project_id.unwrap_or_default(), updated_at, id],
        )?;
        
        Ok(WebShell {
            id: id.to_string(),
            project_id: None,
            name,
            url,
            password,
            payload_type,
            encryption: "aes-256-gcm".to_string(),
            status: "unknown".to_string(),
            last_connected_at: None,
            created_at: 0,
            updated_at,
            deleted_at: None,
        })
    }

    /// 删除 WebShell（软删除）
    pub async fn delete(&self, id: &str) -> Result<()> {
        let conn = self.db.lock().await;
        let deleted_at = chrono::Utc::now().timestamp();
        
        conn.execute(
            "UPDATE webshells SET deleted_at = ?1 WHERE id = ?2",
            rusqlite::params![deleted_at, id],
        )?;
        
        Ok(())
    }

    /// 获取项目的所有 WebShell
    pub async fn get_by_project(&self, project_id: &str) -> Result<Vec<WebShell>> {
        let conn = self.db.lock().await;
        
        let mut stmt = conn.prepare(
            "SELECT id, project_id, name, url, password, payload_type, encryption, status, 
                    last_connected_at, created_at, updated_at 
             FROM webshells 
             WHERE project_id = ?1 AND deleted_at IS NULL"
        )?;
        
        let webshells = stmt.query_map(rusqlite::params![project_id], |row| {
            Ok(WebShell {
                id: row.get(0)?,
                project_id: row.get(1)?,
                name: row.get(2)?,
                url: row.get(3)?,
                password: row.get(4)?,
                payload_type: row.get(5)?,
                encryption: row.get(6)?,
                status: row.get(7)?,
                last_connected_at: row.get(8)?,
                created_at: row.get(9)?,
                updated_at: row.get(10)?,
                deleted_at: None,
            })
        })?;
        
        let mut result = Vec::new();
        for webshell in webshells {
            result.push(webshell?);
        }
        
        Ok(result)
    }

    /// 更新 WebShell 状态
    pub async fn update_status(&self, id: &str, status: String) -> Result<()> {
        let conn = self.db.lock().await;
        let last_connected_at = chrono::Utc::now().timestamp();
        
        conn.execute(
            "UPDATE webshells SET status = ?1, last_connected_at = ?2 WHERE id = ?3",
            rusqlite::params![status, last_connected_at, id],
        )?;
        
        Ok(())
    }
}
