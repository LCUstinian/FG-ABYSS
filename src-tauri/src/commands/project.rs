use tauri::State;
use crate::Result;
use crate::state::{AppState, check_locked};
use crate::features::project::models::{Project, CreateProjectInput, UpdateProjectInput};

#[tauri::command]
#[specta::specta]
pub async fn list_projects(state: State<'_, AppState>) -> Result<Vec<Project>> {
    state.project_service.list().await
}

#[tauri::command]
#[specta::specta]
pub async fn create_project(
    state: State<'_, AppState>,
    input: CreateProjectInput,
) -> Result<Project> {
    check_locked(&state)?;
    state.project_service.create(input).await
}

#[tauri::command]
#[specta::specta]
pub async fn update_project(
    state: State<'_, AppState>,
    id:    String,
    input: UpdateProjectInput,
) -> Result<Project> {
    check_locked(&state)?;
    state.project_service.update(&id, input).await
}

#[tauri::command]
#[specta::specta]
pub async fn delete_project(
    state: State<'_, AppState>,
    id:    String,
) -> Result<()> {
    check_locked(&state)?;
    state.project_service.delete(&id).await
}
