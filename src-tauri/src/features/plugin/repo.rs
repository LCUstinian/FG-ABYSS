use async_trait::async_trait;
use crate::Result;
use super::models::Plugin;
use crate::infra::db::Database;

#[async_trait]
#[cfg_attr(test, mockall::automock)]
pub trait PluginRepo: Send + Sync {
    async fn find_all(&self) -> Result<Vec<Plugin>>;
    async fn find_by_id(&self, id: &str) -> Result<Plugin>;
    async fn insert(&self, p: &Plugin) -> Result<()>;
    async fn set_enabled(&self, id: &str, enabled: bool) -> Result<()>;
}

pub struct DbPluginRepo { db: Database }
impl DbPluginRepo {
    pub fn new(db: Database) -> Self { Self { db } }
}

fn row_to_plugin(r: &rusqlite::Row) -> rusqlite::Result<Plugin> {
    let cfg_str: String = r.get(4)?;
    Ok(Plugin {
        id: r.get(0)?, name: r.get(1)?, version: r.get(2)?,
        enabled: r.get::<_, i32>(3)? != 0,
        config:  serde_json::from_str(&cfg_str).unwrap_or_default(),
        source:  r.get(5)?,
        created_at: r.get(6)?, updated_at: r.get(7)?,
    })
}

#[async_trait]
impl PluginRepo for DbPluginRepo {
    async fn find_all(&self) -> Result<Vec<Plugin>> {
        self.db.call(|conn| {
            let mut stmt = conn.prepare(
                "SELECT id, name, version, enabled, config, source, created_at, updated_at
                 FROM plugins ORDER BY created_at ASC"
            )?;
            let rows = stmt.query_map([], row_to_plugin)?.collect::<rusqlite::Result<Vec<_>>>()?;
            Ok(rows)
        }).await
    }

    async fn find_by_id(&self, id: &str) -> Result<Plugin> {
        let id_owned = id.to_string();
        let id_err = id.to_string();
        self.db.call(move |conn| {
            conn.query_row(
                "SELECT id, name, version, enabled, config, source, created_at, updated_at
                 FROM plugins WHERE id = ?1",
                [&id_owned], row_to_plugin,
            )
        }).await.map_err(|e| {
            if let crate::AppError::Database(ref re) = e {
                if matches!(re, rusqlite::Error::QueryReturnedNoRows) {
                    return crate::AppError::NotFound(format!("plugin {id_err}"));
                }
            }
            e
        })
    }

    async fn insert(&self, p: &Plugin) -> Result<()> {
        let p = p.clone();
        self.db.call(move |conn| {
            conn.execute(
                "INSERT INTO plugins (id, name, version, enabled, config, source, created_at, updated_at)
                 VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8)",
                rusqlite::params![
                    p.id, p.name, p.version, p.enabled as i32,
                    p.config.to_string(), p.source, p.created_at, p.updated_at,
                ],
            )?;
            Ok(())
        }).await
    }

    async fn set_enabled(&self, id: &str, enabled: bool) -> Result<()> {
        let id = id.to_string();
        let v  = enabled as i32;
        self.db.call(move |conn| {
            conn.execute(
                "UPDATE plugins SET enabled = ?1, updated_at = strftime('%s','now') WHERE id = ?2",
                rusqlite::params![v, id],
            )?;
            Ok(())
        }).await
    }
}
