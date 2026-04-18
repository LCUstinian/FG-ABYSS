use async_trait::async_trait;
use crate::Result;
use super::models::{Payload, PayloadHistoryEntry};
use crate::infra::db::Database;

#[async_trait]
#[cfg_attr(test, mockall::automock)]
pub trait PayloadRepo: Send + Sync {
    async fn find_all(&self) -> Result<Vec<Payload>>;
    async fn find_by_id(&self, id: &str) -> Result<Payload>;
    async fn insert(&self, p: &Payload) -> Result<()>;
    async fn soft_delete(&self, id: &str, ts: i64) -> Result<()>;
    async fn insert_history(&self, entry: &PayloadHistoryEntry) -> Result<()>;
    async fn find_history(&self, payload_id: &str) -> Result<Vec<PayloadHistoryEntry>>;
}

pub struct DbPayloadRepo { db: Database }
impl DbPayloadRepo {
    pub fn new(db: Database) -> Self { Self { db } }
}

#[async_trait]
impl PayloadRepo for DbPayloadRepo {
    async fn find_all(&self) -> Result<Vec<Payload>> {
        self.db.call(|conn| {
            let mut stmt = conn.prepare(
                "SELECT id, name, payload_type, config, created_at
                 FROM payloads WHERE deleted_at IS NULL ORDER BY created_at DESC"
            )?;
            let rows = stmt.query_map([], |r| {
                let pt_str: String = r.get(2)?;
                let cfg_str: String = r.get(3)?;
                Ok(Payload {
                    id: r.get(0)?, name: r.get(1)?,
                    payload_type: serde_json::from_str(&pt_str)
                        .unwrap_or(crate::features::webshell::models::PayloadType::Php),
                    config: serde_json::from_str(&cfg_str).unwrap_or_default(),
                    created_at: r.get(4)?,
                })
            })?.collect::<rusqlite::Result<Vec<_>>>()?;
            Ok(rows)
        }).await
    }

    async fn find_by_id(&self, id: &str) -> Result<Payload> {
        let id = id.to_string();
        self.db.call(move |conn| {
            conn.query_row(
                "SELECT id, name, payload_type, config, created_at
                 FROM payloads WHERE id = ?1 AND deleted_at IS NULL",
                [&id],
                |r| {
                    let pt_str: String = r.get(2)?;
                    let cfg_str: String = r.get(3)?;
                    Ok(Payload {
                        id: r.get(0)?, name: r.get(1)?,
                        payload_type: serde_json::from_str(&pt_str)
                            .unwrap_or(crate::features::webshell::models::PayloadType::Php),
                        config: serde_json::from_str(&cfg_str).unwrap_or_default(),
                        created_at: r.get(4)?,
                    })
                },
            )
        }).await
    }

    async fn insert(&self, p: &Payload) -> Result<()> {
        let p = p.clone();
        self.db.call(move |conn| {
            conn.execute(
                "INSERT INTO payloads (id, name, payload_type, config, created_at)
                 VALUES (?1, ?2, ?3, ?4, ?5)",
                rusqlite::params![
                    p.id, p.name,
                    serde_json::to_string(&p.payload_type).unwrap(),
                    p.config.to_string(),
                    p.created_at
                ],
            )?;
            Ok(())
        }).await
    }

    async fn soft_delete(&self, id: &str, ts: i64) -> Result<()> {
        let id = id.to_string();
        self.db.call(move |conn| {
            conn.execute(
                "UPDATE payloads SET deleted_at = ?1 WHERE id = ?2 AND deleted_at IS NULL",
                rusqlite::params![ts, id],
            )?;
            Ok(())
        }).await
    }

    async fn insert_history(&self, e: &PayloadHistoryEntry) -> Result<()> {
        let e = e.clone();
        self.db.call(move |conn| {
            conn.execute(
                "INSERT INTO payload_history (id, payload_id, webshell_id, code, template_version, created_at)
                 VALUES (?1, ?2, ?3, ?4, ?5, ?6)",
                rusqlite::params![e.id, e.payload_id, e.webshell_id, e.code, e.template_version, e.created_at],
            )?;
            Ok(())
        }).await
    }

    async fn find_history(&self, payload_id: &str) -> Result<Vec<PayloadHistoryEntry>> {
        let pid = payload_id.to_string();
        self.db.call(move |conn| {
            let mut stmt = conn.prepare(
                "SELECT id, payload_id, webshell_id, code, template_version, created_at
                 FROM payload_history WHERE payload_id = ?1 ORDER BY created_at DESC"
            )?;
            let rows = stmt.query_map([&pid], |r| Ok(PayloadHistoryEntry {
                id: r.get(0)?, payload_id: r.get(1)?, webshell_id: r.get(2)?,
                code: r.get(3)?, template_version: r.get(4)?, created_at: r.get(5)?,
            }))?.collect::<rusqlite::Result<Vec<_>>>()?;
            Ok(rows)
        }).await
    }
}
