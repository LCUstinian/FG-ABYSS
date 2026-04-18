use tauri::State;
use crate::Result;
use crate::state::{AppState, check_locked};
use crate::features::webshell::models::ConnectionResult;

#[derive(Debug, serde::Serialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct BatchTestResult {
    pub id:     String,
    pub result: ConnectionResult,
}

#[tauri::command]
#[specta::specta]
pub async fn test_connections(
    state: State<'_, AppState>,
    ids:   Vec<String>,
) -> Result<Vec<BatchTestResult>> {
    check_locked(&state)?;
    let _ = ids;
    todo!("test_connections — implement in Phase 2")
}
