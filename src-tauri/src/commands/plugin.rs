use tauri::State;
use crate::Result;
use crate::state::{AppState, check_locked};
use crate::features::plugin::models::Plugin;

#[tauri::command]
#[specta::specta]
pub async fn list_plugins(state: State<'_, AppState>) -> Result<Vec<Plugin>> {
    state.plugin_service.list().await
}

#[tauri::command]
#[specta::specta]
pub async fn enable_plugin(state: State<'_, AppState>, id: String) -> Result<Plugin> {
    check_locked(&state)?;
    state.plugin_service.enable(&id).await
}

#[tauri::command]
#[specta::specta]
pub async fn disable_plugin(state: State<'_, AppState>, id: String) -> Result<Plugin> {
    check_locked(&state)?;
    state.plugin_service.disable(&id).await
}
