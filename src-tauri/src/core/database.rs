use rusqlite::{Connection, Result};
use chrono::NaiveDateTime;
use serde::{Deserialize, Serialize};
use serde_json;
use std::path::Path;


/// 审计日志类型
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum AuditActionType {
    PayloadGenerate,
    WebShellConnect,
    WebShellCommand,
    PluginLoad,
    ProjectCreate,
    ProjectDelete,
    WebShellCreate,
    WebShellDelete,
    SettingsChange,
    Other,
}

/// 审计日志结构体
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AuditLog {
    pub id: i64,
    pub action_type: String,
    pub payload_hash: Option<String>,
    pub encrypted_content: String,
    pub nonce: String,
    pub tag: String,
    pub created_at: NaiveDateTime,
}

/// 数据库管理器
pub struct DatabaseManager {
    conn: Connection,
}

impl DatabaseManager {
    /// 创建数据库连接
    pub fn new(db_path: &Path) -> Result<Self> {
        let conn = Connection::open(db_path)?;
        Self::create_tables(&conn)?;
        Ok(Self { conn })
    }

    /// 创建数据库表
    fn create_tables(conn: &Connection) -> Result<()> {
        // 创建审计日志表
        conn.execute(
            r#"
            CREATE TABLE IF NOT EXISTS audit_logs (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                action_type TEXT NOT NULL,
                payload_hash TEXT,
                encrypted_content TEXT NOT NULL,
                nonce TEXT NOT NULL,
                tag TEXT NOT NULL,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP
            )
            "#,
            [],
        )?;

        // 创建项目表
        conn.execute(
            r#"
            CREATE TABLE IF NOT EXISTS projects (
                uuid TEXT PRIMARY KEY,
                name TEXT NOT NULL,
                desc TEXT,
                is_deleted INTEGER DEFAULT 0,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
            )
            "#,
            [],
        )?;

        // 创建WebShell表
        conn.execute(
            r#"
            CREATE TABLE IF NOT EXISTS webshells (
                uuid TEXT PRIMARY KEY,
                proj_uuid TEXT NOT NULL,
                url TEXT NOT NULL,
                encrypted_passwd TEXT NOT NULL,
                passwd_nonce TEXT NOT NULL,
                passwd_tag TEXT NOT NULL,
                mode TEXT NOT NULL,
                status TEXT NOT NULL,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY (proj_uuid) REFERENCES projects(uuid)
            )
            "#,
            [],
        )?;

        Ok(())
    }

    /// 添加审计日志
    pub fn add_audit_log(
        &self,
        action_type: &str,
        payload_hash: Option<&str>,
        encrypted_content: &str,
        nonce: &str,
        tag: &str,
    ) -> Result<()> {
        self.conn.execute(
            r#"
            INSERT INTO audit_logs (action_type, payload_hash, encrypted_content, nonce, tag)
            VALUES (?, ?, ?, ?, ?)
            "#,
            &[action_type, payload_hash.unwrap_or(""), encrypted_content, nonce, tag],
        )?;
        Ok(())
    }

    /// 获取审计日志列表
    pub fn get_audit_logs(&self, limit: i64, offset: i64) -> Result<Vec<AuditLog>> {
        let mut stmt = self.conn.prepare(
            r#"
            SELECT id, action_type, payload_hash, encrypted_content, nonce, tag, created_at
            FROM audit_logs
            ORDER BY created_at DESC
            LIMIT ? OFFSET ?
            "#,
        )?;

        let log_iter = stmt.query_map(&[&limit, &offset], |row| {
            Ok(AuditLog {
                id: row.get(0)?,
                action_type: row.get(1)?,
                payload_hash: row.get(2)?,
                encrypted_content: row.get(3)?,
                nonce: row.get(4)?,
                tag: row.get(5)?,
                created_at: row.get(6)?,
            })
        })?;

        let mut logs = Vec::new();
        for log in log_iter {
            logs.push(log?);
        }

        Ok(logs)
    }

    /// 获取审计日志总数
    pub fn get_audit_logs_count(&self) -> Result<i64> {
        let count: i64 = self.conn.query_row(
            "SELECT COUNT(*) FROM audit_logs",
            [],
            |row| row.get(0),
        )?;
        Ok(count)
    }

    /// 根据ID获取审计日志
    pub fn get_audit_log_by_id(&self, id: i64) -> Result<Option<AuditLog>> {
        let mut stmt = self.conn.prepare(
            r#"
            SELECT id, action_type, payload_hash, encrypted_content, nonce, tag, created_at
            FROM audit_logs
            WHERE id = ?
            "#,
        )?;

        stmt.query_row(&[&id], |row| {
            Ok(Some(AuditLog {
                id: row.get(0)?,
                action_type: row.get(1)?,
                payload_hash: row.get(2)?,
                encrypted_content: row.get(3)?,
                nonce: row.get(4)?,
                tag: row.get(5)?,
                created_at: row.get(6)?,
            }))
        })
        .or_else(|_| Ok(None))
    }

    /// 清空审计日志
    pub fn clear_audit_logs(&self) -> Result<()> {
        self.conn.execute("DELETE FROM audit_logs", [])?;
        Ok(())
    }

    /// 添加WebShell
    pub fn add_webshell(
        &self,
        uuid: &str,
        proj_uuid: &str,
        url: &str,
        password: &str,
        mode: &str,
        status: &str,
        _encryption_key: &str,
    ) -> Result<()> {
        // 暂时使用简单的base64编码，不进行加密
        let encrypted_passwd = base64::encode(password.as_bytes());
        let nonce_str = base64::encode(&Vec::<u8>::new());
        let tag_str = base64::encode(&Vec::<u8>::new());

        self.conn.execute(
            r#"
            INSERT INTO webshells (uuid, proj_uuid, url, encrypted_passwd, passwd_nonce, passwd_tag, mode, status)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
            "#,
            &[uuid, proj_uuid, url, &encrypted_passwd, &nonce_str, &tag_str, mode, status],
        )?;
        Ok(())
    }

    /// 获取WebShell列表
    pub fn get_webshells(&self, proj_uuid: &str, _encryption_key: &str) -> Result<Vec<serde_json::Value>> {
        let mut stmt = self.conn.prepare(
            r#"
            SELECT uuid, proj_uuid, url, encrypted_passwd, passwd_nonce, passwd_tag, mode, status, created_at, updated_at
            FROM webshells
            WHERE proj_uuid = ? AND status != 'deleted'
            "#,
        )?;

        let webshell_iter = stmt.query_map(&[proj_uuid], |row| {
            let encrypted_passwd: String = row.get(3)?;

            // 暂时使用简单的base64解码，不进行解密
            let password = base64::decode(&encrypted_passwd).unwrap();
            let password_str = String::from_utf8(password).unwrap();

            Ok(serde_json::json! {
                {
                    "uuid": row.get::<_, String>(0)?,
                    "proj_uuid": row.get::<_, String>(1)?,
                    "url": row.get::<_, String>(2)?,
                    "password": password_str,
                    "mode": row.get::<_, String>(6)?,
                    "status": row.get::<_, String>(7)?,
                    "created_at": row.get::<_, String>(8)?,
                    "updated_at": row.get::<_, String>(9)?
                }
            })
        })?;

        let mut webshells = Vec::new();
        for webshell in webshell_iter {
            webshells.push(webshell?);
        }

        Ok(webshells)
    }

    /// 更新WebShell
    pub fn update_webshell(
        &self,
        uuid: &str,
        url: &str,
        password: &str,
        mode: &str,
        status: &str,
        _encryption_key: &str,
    ) -> Result<()> {
        // 暂时使用简单的base64编码，不进行加密
        let encrypted_passwd = base64::encode(password.as_bytes());
        let nonce_str = base64::encode(&Vec::<u8>::new());
        let tag_str = base64::encode(&Vec::<u8>::new());

        self.conn.execute(
            r#"
            UPDATE webshells
            SET url = ?, encrypted_passwd = ?, passwd_nonce = ?, passwd_tag = ?, mode = ?, status = ?, updated_at = CURRENT_TIMESTAMP
            WHERE uuid = ?
            "#,
            &[url, &encrypted_passwd, &nonce_str, &tag_str, mode, status, uuid],
        )?;
        Ok(())
    }

    /// 删除WebShell（软删除）
    pub fn delete_webshell(&self, uuid: &str) -> Result<()> {
        self.conn.execute(
            "UPDATE webshells SET status = 'deleted' WHERE uuid = ?",
            &[uuid],
        )?;
        Ok(())
    }
}

// 使用外部base64库

mod base64 {
    use ::base64::{Engine as _, engine::general_purpose};
    
    pub fn encode(data: &[u8]) -> String {
        general_purpose::STANDARD.encode(data)
    }

    pub fn decode(encoded: &str) -> Result<Vec<u8>, ::base64::DecodeError> {
        general_purpose::STANDARD.decode(encoded)
    }
}
