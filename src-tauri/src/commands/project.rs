use tauri::command;
use std::sync::Arc;
use tokio::sync::Mutex;

use crate::core::project::ProjectManager;
use crate::infra::database::models::Project;

/// 创建项目
#[command]
pub async fn create_project(
    name: String,
    description: Option<String>,
    db: tauri::State<'_, Arc<Mutex<rusqlite::Connection>>>,
) -> Result<Project, String> {
    let manager = ProjectManager::new(db.inner().clone());
    manager.create(name, description)
        .await
        .map_err(|e| e.to_string())
}

/// 更新项目
#[command]
pub async fn update_project(
    id: String,
    name: String,
    description: Option<String>,
    db: tauri::State<'_, Arc<Mutex<rusqlite::Connection>>>,
) -> Result<Project, String> {
    let manager = ProjectManager::new(db.inner().clone());
    manager.update(&id, name, description)
        .await
        .map_err(|e| e.to_string())
}

/// 删除项目（软删除）
#[command]
pub async fn delete_project(
    id: String,
    db: tauri::State<'_, Arc<Mutex<rusqlite::Connection>>>,
) -> Result<(), String> {
    let manager = ProjectManager::new(db.inner().clone());
    manager.delete(&id)
        .await
        .map_err(|e| e.to_string())
}

/// 恢复项目
#[command]
pub async fn restore_project(
    id: String,
    db: tauri::State<'_, Arc<Mutex<rusqlite::Connection>>>,
) -> Result<(), String> {
    let manager = ProjectManager::new(db.inner().clone());
    manager.restore(&id)
        .await
        .map_err(|e| e.to_string())
}

/// 获取所有项目
#[command]
pub async fn get_all_projects(
    db: tauri::State<'_, Arc<Mutex<rusqlite::Connection>>>,
) -> Result<Vec<Project>, String> {
    let manager = ProjectManager::new(db.inner().clone());
    manager.get_all()
        .await
        .map_err(|e| e.to_string())
}

/// 获取已删除的项目
#[command]
pub async fn get_deleted_projects(
    db: tauri::State<'_, Arc<Mutex<rusqlite::Connection>>>,
) -> Result<Vec<Project>, String> {
    let manager = ProjectManager::new(db.inner().clone());
    manager.get_deleted()
        .await
        .map_err(|e| e.to_string())
}
