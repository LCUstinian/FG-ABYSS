use serde::{Deserialize, Serialize};
use uuid::Uuid;
use crate::{Result, infra::db::Database};

#[derive(Debug, Clone)]
pub struct AuditLog {
    pub db: Database,
}

#[derive(Debug, Serialize, Deserialize)]
pub enum AuditAction {
    FileRead, FileWrite, FileDelete,
    CommandExec,
    DBQuery, DBExport,
    WebshellConnect, WebshellDisconnect,
}

impl AuditLog {
    pub fn new(db: Database) -> Self { Self { db } }

    pub async fn record(&self, webshell_id: &str, action: AuditAction, detail: &str) -> Result<()> {
        let id         = Uuid::new_v4().to_string();
        let action_str = serde_json::to_string(&action)?;
        let detail     = detail.to_string();
        let ts         = chrono::Utc::now().timestamp();
        let wid        = webshell_id.to_string();
        self.db.call(move |conn| {
            conn.execute(
                "INSERT INTO audit_events (id, webshell_id, action, detail, created_at)
                 VALUES (?1, ?2, ?3, ?4, ?5)",
                rusqlite::params![id, wid, action_str, detail, ts],
            )?;
            Ok(())
        }).await
    }
}
