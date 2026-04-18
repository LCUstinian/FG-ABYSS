use std::path::Path;
use std::sync::Arc;
use crate::{AppError, Result};

#[derive(Clone)]
pub struct Database(Arc<tokio_rusqlite::Connection>);

impl Database {
    pub async fn open(path: &Path) -> Result<Self> {
        let conn = tokio_rusqlite::Connection::open(path).await
            .map_err(|e| AppError::DbConnect(e.to_string()))?;
        conn.call(|c| {
            c.execute_batch("PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;")?;
            Ok(())
        }).await.map_err(|e| AppError::DbConnect(e.to_string()))?;
        Ok(Self(Arc::new(conn)))
    }

    pub async fn call<F, T>(&self, f: F) -> Result<T>
    where
        F: FnOnce(&mut rusqlite::Connection) -> rusqlite::Result<T> + Send + 'static,
        T: Send + 'static,
    {
        self.0
            .call(move |conn| f(conn).map_err(tokio_rusqlite::Error::Rusqlite))
            .await
            .map_err(|e| AppError::DbConnect(e.to_string()))
    }

    pub async fn migrate(&self) -> Result<()> {
        const MIGRATIONS: &[(&str, &str)] = &[
            ("001_init",    include_str!("../../migrations/001_init.sql")),
            ("002_plugins", include_str!("../../migrations/002_plugins.sql")),
            ("003_payload", include_str!("../../migrations/003_payload.sql")),
            ("004_audit",   include_str!("../../migrations/004_audit.sql")),
        ];

        self.call(|conn| {
            conn.execute_batch(
                "CREATE TABLE IF NOT EXISTS _migrations (
                    name TEXT PRIMARY KEY,
                    applied_at INTEGER NOT NULL
                );"
            )?;
            for (name, sql) in MIGRATIONS {
                let exists: bool = conn.query_row(
                    "SELECT COUNT(*) FROM _migrations WHERE name = ?1",
                    [name],
                    |r| r.get::<_, i64>(0),
                ).map(|n| n > 0).unwrap_or(false);

                if !exists {
                    conn.execute_batch(sql)?;
                    conn.execute(
                        "INSERT INTO _migrations (name, applied_at) VALUES (?1, strftime('%s','now'))",
                        [name],
                    )?;
                }
            }
            Ok(())
        }).await
    }

    pub async fn vacuum_if_needed(&self) -> Result<()> {
        let deleted_count: i64 = self.call(|conn| {
            conn.query_row(
                "SELECT COUNT(*) FROM webshells WHERE deleted_at IS NOT NULL",
                [],
                |r| r.get(0),
            )
        }).await?;
        if deleted_count > 100 {
            self.call(|conn| conn.execute_batch("VACUUM;")).await?;
            tracing::info!("VACUUM completed, {} soft-deleted rows reclaimed", deleted_count);
        }
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    async fn test_db() -> Database {
        Database::open(Path::new(":memory:")).await.unwrap()
    }

    #[tokio::test]
    async fn test_open_and_migrate() {
        let db = test_db().await;
        db.migrate().await.expect("migration should succeed");

        let count: i64 = db.call(|conn| {
            conn.query_row(
                "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='webshells'",
                [],
                |r| r.get(0),
            )
        }).await.unwrap();
        assert_eq!(count, 1);
    }

    #[tokio::test]
    async fn test_migrate_is_idempotent() {
        let db = test_db().await;
        db.migrate().await.unwrap();
        db.migrate().await.expect("second migration should not fail");
    }
}
