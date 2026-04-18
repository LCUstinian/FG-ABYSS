use tauri::State;
use crate::Result;
use crate::state::AppState;

#[derive(Debug, serde::Serialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct AppInfo {
    pub version:    String,
    pub data_dir:   String,
    pub db_path:    String,
    pub logs_dir:   String,
}

#[tauri::command]
#[specta::specta]
pub async fn get_app_info(state: State<'_, AppState>) -> Result<AppInfo> {
    Ok(AppInfo {
        version:  env!("CARGO_PKG_VERSION").to_string(),
        data_dir: state.paths.data_dir.to_string_lossy().into(),
        db_path:  state.paths.db_path.to_string_lossy().into(),
        logs_dir: state.paths.logs_dir.to_string_lossy().into(),
    })
}

#[derive(Debug, serde::Serialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct AuditEntry {
    pub id:          String,
    pub webshell_id: Option<String>,
    pub action:      String,
    pub detail:      String,
    pub created_at:  i64,
}

#[tauri::command]
#[specta::specta]
pub async fn get_audit_log(
    state: State<'_, AppState>,
    limit: Option<i64>,
) -> Result<Vec<AuditEntry>> {
    let limit = limit.unwrap_or(100);
    state.audit_log.db.call(move |conn| {
        let mut stmt = conn.prepare(
            "SELECT id, webshell_id, action, detail, created_at
             FROM audit_events ORDER BY created_at DESC LIMIT ?1"
        )?;
        let rows = stmt.query_map([limit], |r| Ok(AuditEntry {
            id:          r.get(0)?,
            webshell_id: r.get(1)?,
            action:      r.get(2)?,
            detail:      r.get(3)?,
            created_at:  r.get(4)?,
        }))?.collect::<rusqlite::Result<Vec<_>>>()?;
        Ok(rows)
    }).await
}
