use tauri::{AppHandle, Manager, State};
use crate::Result;
use crate::state::{AppState, check_locked};

#[tauri::command]
#[specta::specta]
pub async fn open_console(
    app:         AppHandle,
    state:       State<'_, AppState>,
    webshell_id: String,
) -> Result<()> {
    check_locked(&state)?;
    let label = format!("console-{}", webshell_id);

    if let Some(win) = app.get_webview_window(&label) {
        win.set_focus().map_err(|e: tauri::Error| crate::AppError::Io(std::io::Error::new(
            std::io::ErrorKind::Other, e.to_string()
        )))?;
        return Ok(());
    }

    let shell = state.webshell_service.get(&webshell_id).await?;

    tauri::WebviewWindowBuilder::new(
        &app, &label,
        tauri::WebviewUrl::App(format!("/console?id={}", webshell_id).into()),
    )
    .title(format!("{} — {}", shell.name, shell.url))
    .inner_size(1200.0, 800.0)
    .min_inner_size(900.0, 600.0)
    .decorations(false)
    .build()
    .map_err(|e| crate::AppError::Io(std::io::Error::new(
        std::io::ErrorKind::Other, e.to_string()
    )))?;

    Ok(())
}

#[tauri::command]
#[specta::specta]
pub async fn exec_command(
    state:       State<'_, AppState>,
    webshell_id: String,
    cmd:         String,
) -> Result<String> {
    check_locked(&state)?;
    state.console_service.exec_command(&webshell_id, &cmd).await
}

#[tauri::command]
#[specta::specta]
pub async fn list_files(
    state:       State<'_, AppState>,
    webshell_id: String,
    path:        String,
) -> Result<serde_json::Value> {
    check_locked(&state)?;
    state.console_service.list_files(&webshell_id, &path).await
}

#[tauri::command]
#[specta::specta]
pub async fn download_file(
    _app:        AppHandle,
    state:       State<'_, AppState>,
    webshell_id: String,
    path:        String,
    transfer_id: String,
) -> Result<()> {
    check_locked(&state)?;
    let _ = (webshell_id, path, transfer_id);
    todo!("download_file stream — implement in Phase 2")
}

#[tauri::command]
#[specta::specta]
pub async fn upload_file(
    state:       State<'_, AppState>,
    webshell_id: String,
    path:        String,
    data:        Vec<u8>,
) -> Result<()> {
    check_locked(&state)?;
    let _ = (webshell_id, path, data);
    todo!("upload_file — implement in Phase 2")
}

#[tauri::command]
#[specta::specta]
pub async fn connect_database(
    state:       State<'_, AppState>,
    webshell_id: String,
    db_type:     String,
    conn_str:    String,
) -> Result<()> {
    check_locked(&state)?;
    let _ = (webshell_id, db_type, conn_str);
    todo!("connect_database — implement in Phase 2")
}

#[tauri::command]
#[specta::specta]
pub async fn execute_query(
    state:       State<'_, AppState>,
    webshell_id: String,
    query:       String,
) -> Result<serde_json::Value> {
    check_locked(&state)?;
    let _ = (webshell_id, query);
    todo!("execute_query — implement in Phase 2")
}
