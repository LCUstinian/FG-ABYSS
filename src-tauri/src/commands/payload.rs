use tauri::State;
use crate::Result;
use crate::state::{AppState, check_locked};
use crate::features::payload::models::{Payload, CreatePayloadInput, PayloadConfig, PayloadHistoryEntry};

#[tauri::command]
#[specta::specta]
pub async fn list_payloads(state: State<'_, AppState>) -> Result<Vec<Payload>> {
    state.payload_service.list().await
}

#[tauri::command]
#[specta::specta]
pub async fn create_payload(
    state: State<'_, AppState>,
    input: CreatePayloadInput,
) -> Result<Payload> {
    check_locked(&state)?;
    state.payload_service.create(input).await
}

#[tauri::command]
#[specta::specta]
pub async fn generate_payload(
    state:       State<'_, AppState>,
    payload_id:  Option<String>,
    webshell_id: Option<String>,
    config:      PayloadConfig,
) -> Result<String> {
    check_locked(&state)?;
    state.payload_service.generate_payload(
        payload_id.as_deref(),
        webshell_id.as_deref(),
        config,
    ).await
}

#[tauri::command]
#[specta::specta]
pub async fn list_payload_history(
    state:      State<'_, AppState>,
    payload_id: String,
) -> Result<Vec<PayloadHistoryEntry>> {
    state.payload_service.list_history(&payload_id).await
}
