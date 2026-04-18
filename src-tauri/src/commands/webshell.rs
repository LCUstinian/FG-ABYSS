use tauri::State;
use crate::Result;
use crate::state::{AppState, check_locked};
use crate::features::webshell::models::{
    Webshell, CreateWebshellInput, UpdateWebshellInput, ConnectionResult,
};

#[tauri::command]
#[specta::specta]
pub async fn list_webshells(
    state:      State<'_, AppState>,
    project_id: Option<String>,
) -> Result<Vec<Webshell>> {
    state.webshell_service.list(project_id.as_deref()).await
}

#[tauri::command]
#[specta::specta]
pub async fn get_webshell(
    state: State<'_, AppState>,
    id:    String,
) -> Result<Webshell> {
    state.webshell_service.get(&id).await
}

#[tauri::command]
#[specta::specta]
pub async fn create_webshell(
    state: State<'_, AppState>,
    input: CreateWebshellInput,
) -> Result<Webshell> {
    check_locked(&state)?;
    state.webshell_service.create(input).await
}

#[tauri::command]
#[specta::specta]
pub async fn update_webshell(
    state: State<'_, AppState>,
    id:    String,
    input: UpdateWebshellInput,
) -> Result<Webshell> {
    check_locked(&state)?;
    state.webshell_service.update(&id, input).await
}

#[tauri::command]
#[specta::specta]
pub async fn delete_webshell(
    state: State<'_, AppState>,
    id:    String,
) -> Result<()> {
    check_locked(&state)?;
    state.webshell_service.delete(&id).await
}

#[tauri::command]
#[specta::specta]
pub async fn test_connection(
    state: State<'_, AppState>,
    id:    String,
) -> Result<ConnectionResult> {
    check_locked(&state)?;
    state.webshell_service.test_connection(&id).await
}

#[tauri::command]
#[specta::specta]
pub async fn reset_redeploy_status(
    state: State<'_, AppState>,
    id:    String,
) -> Result<Webshell> {
    check_locked(&state)?;
    state.webshell_service.reset_redeploy_status(&id).await
}
