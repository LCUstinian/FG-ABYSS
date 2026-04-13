use tauri::command;
use std::sync::Arc;
use tokio::sync::Mutex;

use crate::core::webshell::WebShellManager;
use crate::infra::database::models::WebShell;

/// 创建 WebShell
#[command]
pub async fn create_webshell(
    name: String,
    url: String,
    password: String,
    payload_type: String,
    project_id: Option<String>,
    db: tauri::State<'_, Arc<Mutex<rusqlite::Connection>>>,
) -> Result<WebShell, String> {
    let manager = WebShellManager::new(db.inner().clone());
    manager.create(name, url, password, payload_type, project_id)
        .await
        .map_err(|e| e.to_string())
}

/// 更新 WebShell
#[command]
pub async fn update_webshell(
    id: String,
    name: String,
    url: String,
    password: String,
    payload_type: String,
    project_id: Option<String>,
    db: tauri::State<'_, Arc<Mutex<rusqlite::Connection>>>,
) -> Result<WebShell, String> {
    let manager = WebShellManager::new(db.inner().clone());
    manager.update(id.as_str(), name, url, password, payload_type, project_id)
        .await
        .map_err(|e| e.to_string())
}

/// 删除 WebShell（软删除）
#[command]
pub async fn delete_webshell(
    id: String,
    db: tauri::State<'_, Arc<Mutex<rusqlite::Connection>>>,
) -> Result<(), String> {
    let manager = WebShellManager::new(db.inner().clone());
    manager.delete(&id)
        .await
        .map_err(|e| e.to_string())
}

/// 获取项目的所有 WebShell
#[command]
pub async fn get_webshells_by_project(
    project_id: String,
    db: tauri::State<'_, Arc<Mutex<rusqlite::Connection>>>,
) -> Result<Vec<WebShell>, String> {
    let manager = WebShellManager::new(db.inner().clone());
    manager.get_by_project(&project_id)
        .await
        .map_err(|e| e.to_string())
}

/// 测试 WebShell 连接
#[command]
pub async fn test_connection(
    id: String,
    db: tauri::State<'_, Arc<Mutex<rusqlite::Connection>>>,
) -> Result<bool, String> {
    // TODO: 实现实际的连接测试逻辑
    let manager = WebShellManager::new(db.inner().clone());
    manager.update_status(&id, "online".to_string())
        .await
        .map_err(|e| e.to_string())?;
    Ok(true)
}
