use async_trait::async_trait;
use crate::Result;
use super::models::Project;
use crate::infra::db::Database;

#[async_trait]
#[cfg_attr(test, mockall::automock)]
pub trait ProjectRepo: Send + Sync {
    async fn find_all(&self) -> Result<Vec<Project>>;
    async fn find_by_id(&self, id: &str) -> Result<Project>;
    async fn insert(&self, p: &Project) -> Result<()>;
    async fn update(&self, id: &str, name: Option<String>, description: Option<String>) -> Result<()>;
    async fn soft_delete_with_webshells(&self, project_id: &str, ts: i64) -> Result<()>;
}

pub struct DbProjectRepo { db: Database }

impl DbProjectRepo {
    pub fn new(db: Database) -> Self { Self { db } }
}

#[async_trait]
impl ProjectRepo for DbProjectRepo {
    async fn find_all(&self) -> Result<Vec<Project>> {
        self.db.call(|conn| {
            let mut stmt = conn.prepare(
                "SELECT id, name, description, created_at, updated_at
                 FROM projects WHERE deleted_at IS NULL ORDER BY created_at DESC"
            )?;
            let rows = stmt.query_map([], |r| Ok(Project {
                id: r.get(0)?, name: r.get(1)?, description: r.get(2)?,
                created_at: r.get(3)?, updated_at: r.get(4)?,
            }))?.collect::<rusqlite::Result<Vec<_>>>()?;
            Ok(rows)
        }).await
    }

    async fn find_by_id(&self, id: &str) -> Result<Project> {
        let id_owned = id.to_string();
        let id_for_error = id_owned.clone();
        self.db.call(move |conn| {
            conn.query_row(
                "SELECT id, name, description, created_at, updated_at
                 FROM projects WHERE id = ?1 AND deleted_at IS NULL",
                [&id_owned],
                |r| Ok(Project {
                    id: r.get(0)?, name: r.get(1)?, description: r.get(2)?,
                    created_at: r.get(3)?, updated_at: r.get(4)?,
                }),
            )
        }).await.map_err(|e| {
            if let crate::AppError::Database(ref re) = e {
                if matches!(re, rusqlite::Error::QueryReturnedNoRows) {
                    return crate::AppError::NotFound(format!("project {}", id_for_error));
                }
            }
            e
        })
    }

    async fn insert(&self, p: &Project) -> Result<()> {
        let p = p.clone();
        self.db.call(move |conn| {
            conn.execute(
                "INSERT INTO projects (id, name, description, created_at, updated_at)
                 VALUES (?1, ?2, ?3, ?4, ?5)",
                rusqlite::params![p.id, p.name, p.description, p.created_at, p.updated_at],
            )?;
            Ok(())
        }).await
    }

    async fn update(&self, id: &str, name: Option<String>, description: Option<String>) -> Result<()> {
        let id = id.to_string();
        self.db.call(move |conn| {
            conn.execute(
                "UPDATE projects SET name = COALESCE(?1, name),
                 description = COALESCE(?2, description),
                 updated_at = strftime('%s','now')
                 WHERE id = ?3 AND deleted_at IS NULL",
                rusqlite::params![name, description, id],
            )?;
            Ok(())
        }).await
    }

    async fn soft_delete_with_webshells(&self, project_id: &str, ts: i64) -> Result<()> {
        let pid = project_id.to_string();
        self.db.call(move |conn| {
            conn.execute_batch("BEGIN;")?;
            conn.execute(
                "UPDATE webshells SET deleted_at=?1 WHERE project_id=?2 AND deleted_at IS NULL",
                rusqlite::params![ts, &pid],
            )?;
            conn.execute(
                "UPDATE projects SET deleted_at=?1 WHERE id=?2 AND deleted_at IS NULL",
                rusqlite::params![ts, &pid],
            )?;
            conn.execute_batch("COMMIT;")?;
            Ok(())
        }).await
    }
}
