use async_trait::async_trait;
use crate::Result;
use super::models::{Webshell, WebshellFingerprint};
use crate::infra::db::Database;

#[async_trait]
#[cfg_attr(test, mockall::automock)]
pub trait WebshellRepo: Send + Sync {
    async fn find_all(&self, project_id: Option<&str>) -> Result<Vec<Webshell>>;
    async fn find_by_id(&self, id: &str) -> Result<Webshell>;
    async fn find_by_url(&self, url: &str) -> Result<Option<Webshell>>;
    async fn insert(&self, w: &Webshell) -> Result<()>;
    async fn update_fields(&self, id: &str, fields: WebshellUpdateFields) -> Result<()>;
    async fn update_status(
        &self, id: &str, status: &str, fp: Option<&WebshellFingerprint>,
    ) -> Result<()>;
    async fn soft_delete(&self, id: &str, ts: i64) -> Result<()>;
}

pub struct WebshellUpdateFields {
    pub name:           Option<String>,
    pub url:            Option<String>,
    pub password:       Option<String>,
    pub project_id:     Option<Option<String>>,
    pub tags:           Option<String>,
    pub custom_headers: Option<String>,
    pub cookies:        Option<String>,
    pub proxy_override: Option<Option<String>>,
    pub http_method:    Option<String>,
    pub c2_profile:     Option<String>,
    pub crypto_chain:   Option<String>,
    pub notes:          Option<Option<String>>,
    pub status:         Option<String>,
}

pub struct DbWebshellRepo {
    db: Database,
}

impl DbWebshellRepo {
    pub fn new(db: Database) -> Self { Self { db } }
}

#[async_trait]
impl WebshellRepo for DbWebshellRepo {
    async fn find_all(&self, project_id: Option<&str>) -> Result<Vec<Webshell>> {
        let pid = project_id.map(|s| s.to_string());
        self.db.call(move |conn| {
            let sql = if pid.is_some() {
                "SELECT id, name, url, password, payload_type, project_id, status, tags, \
                 custom_headers, cookies, proxy_override, http_method, c2_profile, crypto_chain, \
                 fingerprint, notes, last_connected_at, created_at, updated_at \
                 FROM webshells WHERE deleted_at IS NULL AND project_id = ?1 ORDER BY created_at DESC"
            } else {
                "SELECT id, name, url, password, payload_type, project_id, status, tags, \
                 custom_headers, cookies, proxy_override, http_method, c2_profile, crypto_chain, \
                 fingerprint, notes, last_connected_at, created_at, updated_at \
                 FROM webshells WHERE deleted_at IS NULL ORDER BY created_at DESC"
            };
            let mut stmt = conn.prepare(sql)?;
            let rows = stmt.query_map(
                rusqlite::params_from_iter(pid.iter()),
                row_to_webshell,
            )?.collect::<rusqlite::Result<Vec<_>>>()?;
            Ok(rows)
        }).await
    }

    async fn find_by_id(&self, id: &str) -> Result<Webshell> {
        let id_owned = id.to_string();
        let id_err   = id.to_string();
        self.db.call(move |conn| {
            conn.query_row(
                "SELECT id, name, url, password, payload_type, project_id, status, tags, \
                 custom_headers, cookies, proxy_override, http_method, c2_profile, crypto_chain, \
                 fingerprint, notes, last_connected_at, created_at, updated_at \
                 FROM webshells WHERE id = ?1 AND deleted_at IS NULL",
                [&id_owned],
                row_to_webshell,
            )
        }).await.map_err(|e| {
            if let crate::AppError::Database(ref re) = e {
                if matches!(re, rusqlite::Error::QueryReturnedNoRows) {
                    return crate::AppError::NotFound(format!("webshell {id_err}"));
                }
            }
            e
        })
    }

    async fn find_by_url(&self, url: &str) -> Result<Option<Webshell>> {
        let url = url.to_string();
        self.db.call(move |conn| {
            let mut stmt = conn.prepare(
                "SELECT id, name, url, password, payload_type, project_id, status, tags, \
                 custom_headers, cookies, proxy_override, http_method, c2_profile, crypto_chain, \
                 fingerprint, notes, last_connected_at, created_at, updated_at \
                 FROM webshells WHERE url = ?1 AND deleted_at IS NULL LIMIT 1",
            )?;
            let mut rows = stmt.query_map([&url], row_to_webshell)?;
            Ok(rows.next().transpose()?)
        }).await
    }

    async fn insert(&self, w: &Webshell) -> Result<()> {
        let w = w.clone();
        self.db.call(move |conn| {
            conn.execute(
                "INSERT INTO webshells \
                 (id, name, url, password, payload_type, project_id, status, tags, \
                  custom_headers, cookies, proxy_override, http_method, c2_profile, \
                  crypto_chain, fingerprint, notes, last_connected_at, created_at, updated_at) \
                 VALUES (?1,?2,?3,?4,?5,?6,?7,?8,?9,?10,?11,?12,?13,?14,?15,?16,?17,?18,?19)",
                rusqlite::params![
                    w.id,
                    w.name,
                    w.url,
                    w.password,
                    serde_json::to_string(&w.payload_type).unwrap(),
                    w.project_id,
                    w.status,
                    serde_json::to_string(&w.tags).unwrap(),
                    serde_json::to_string(&w.custom_headers).unwrap(),
                    serde_json::to_string(&w.cookies).unwrap(),
                    w.proxy_override,
                    w.http_method,
                    w.c2_profile,
                    serde_json::to_string(&w.crypto_chain).unwrap(),
                    w.fingerprint.as_ref().and_then(|f| serde_json::to_string(f).ok()),
                    w.notes,
                    w.last_connected_at,
                    w.created_at,
                    w.updated_at,
                ],
            )?;
            Ok(())
        }).await
    }

    async fn update_fields(&self, id: &str, f: WebshellUpdateFields) -> Result<()> {
        let id = id.to_string();
        self.db.call(move |conn| {
            conn.execute(
                "UPDATE webshells SET
                  name           = COALESCE(?1,  name),
                  url            = COALESCE(?2,  url),
                  password       = COALESCE(?3,  password),
                  project_id     = COALESCE(?4,  project_id),
                  tags           = COALESCE(?5,  tags),
                  custom_headers = COALESCE(?6,  custom_headers),
                  cookies        = COALESCE(?7,  cookies),
                  proxy_override = COALESCE(?8,  proxy_override),
                  http_method    = COALESCE(?9,  http_method),
                  c2_profile     = COALESCE(?10, c2_profile),
                  crypto_chain   = COALESCE(?11, crypto_chain),
                  notes          = COALESCE(?12, notes),
                  status         = COALESCE(?13, status),
                  updated_at     = strftime('%s','now')
                 WHERE id = ?14 AND deleted_at IS NULL",
                rusqlite::params![
                    f.name,
                    f.url,
                    f.password,
                    f.project_id.flatten(),
                    f.tags,
                    f.custom_headers,
                    f.cookies,
                    f.proxy_override.flatten(),
                    f.http_method,
                    f.c2_profile,
                    f.crypto_chain,
                    f.notes.flatten(),
                    f.status,
                    id,
                ],
            )?;
            Ok(())
        }).await
    }

    async fn update_status(
        &self,
        id: &str,
        status: &str,
        fp: Option<&WebshellFingerprint>,
    ) -> Result<()> {
        let id      = id.to_string();
        let status  = status.to_string();
        let fp_json = fp.and_then(|f| serde_json::to_string(f).ok());
        self.db.call(move |conn| {
            conn.execute(
                "UPDATE webshells SET
                  status           = ?1,
                  fingerprint      = COALESCE(?2, fingerprint),
                  last_connected_at = CASE WHEN ?1 = 'active' \
                                          THEN strftime('%s','now') \
                                          ELSE last_connected_at END,
                  updated_at       = strftime('%s','now')
                 WHERE id = ?3 AND deleted_at IS NULL",
                rusqlite::params![status, fp_json, id],
            )?;
            Ok(())
        }).await
    }

    async fn soft_delete(&self, id: &str, ts: i64) -> Result<()> {
        let id = id.to_string();
        self.db.call(move |conn| {
            conn.execute(
                "UPDATE webshells SET deleted_at = ?1 WHERE id = ?2 AND deleted_at IS NULL",
                rusqlite::params![ts, id],
            )?;
            Ok(())
        }).await
    }
}

fn row_to_webshell(row: &rusqlite::Row) -> rusqlite::Result<Webshell> {
    use super::models::PayloadType;

    let payload_type_str: String        = row.get(4)?;
    let tags_str:         String        = row.get(7)?;
    let headers_str:      String        = row.get(8)?;
    let cookies_str:      String        = row.get(9)?;
    let chain_str:        String        = row.get(13)?;
    let fp_str:           Option<String> = row.get(14)?;

    Ok(Webshell {
        id:               row.get(0)?,
        name:             row.get(1)?,
        url:              row.get(2)?,
        password:         row.get(3)?,
        payload_type:     serde_json::from_str(&payload_type_str).unwrap_or(PayloadType::Php),
        project_id:       row.get(5)?,
        status:           row.get(6)?,
        tags:             serde_json::from_str(&tags_str).unwrap_or_default(),
        custom_headers:   serde_json::from_str(&headers_str).unwrap_or_default(),
        cookies:          serde_json::from_str(&cookies_str).unwrap_or_default(),
        proxy_override:   row.get(10)?,
        http_method:      row.get(11)?,
        c2_profile:       row.get(12)?,
        crypto_chain:     serde_json::from_str(&chain_str).unwrap_or_default(),
        fingerprint:      fp_str.and_then(|s| serde_json::from_str(&s).ok()),
        notes:            row.get(15)?,
        last_connected_at: row.get(16)?,
        created_at:       row.get(17)?,
        updated_at:       row.get(18)?,
    })
}
