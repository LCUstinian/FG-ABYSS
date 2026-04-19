# FG-ABYSS Backend Skeleton Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the complete Rust backend skeleton — all infrastructure layers, feature services (CRUD), command handlers, AppState bootstrap, and Tauri V2 wiring. Protocol Phase1/Phase2, payload generators, and obfuscators are left as `todo!()`.

**Architecture:** Feature-based layered architecture: `commands/` (thin, owns `AppHandle`) → `features/` (business logic, no `AppHandle`, testable via mock traits) → `infra/` (db, crypto, http, logging). `AppState` is the DI root, bootstrapped once, registered via `tauri::Builder::manage()`.

**Tech Stack:** Rust, Tauri V2, tokio-rusqlite, aes-gcm, argon2, reqwest, tracing/tracing-appender, dashmap, tauri-specta 2.x, mockall, uuid, zeroize

---

## File Map

### Create
```
src-tauri/
├── Cargo.toml                                    (modify)
├── migrations/
│   ├── 001_init.sql
│   ├── 002_plugins.sql
│   ├── 003_payload.sql
│   └── 004_audit.sql
├── capabilities/
│   ├── default.json                              (modify)
│   └── console.json                              (create)
├── tauri.conf.json                               (modify)
└── src/
    ├── error.rs
    ├── state.rs
    ├── lib.rs                                    (modify)
    ├── main.rs                                   (modify)
    ├── infra/
    │   ├── mod.rs
    │   ├── paths.rs
    │   ├── logger.rs
    │   ├── db.rs
    │   ├── config.rs
    │   ├── crypto.rs
    │   ├── audit.rs
    │   ├── http.rs
    │   └── c2_profile.rs
    ├── features/
    │   ├── mod.rs
    │   ├── webshell/
    │   │   ├── mod.rs, models.rs, repo.rs, service.rs, session.rs, queue.rs
    │   ├── project/
    │   │   ├── mod.rs, models.rs, repo.rs, service.rs
    │   ├── payload/
    │   │   ├── mod.rs, models.rs, repo.rs, service.rs
    │   │   └── generator/mod.rs
    │   ├── console/
    │   │   └── mod.rs, service.rs
    │   ├── plugin/
    │   │   ├── mod.rs, models.rs, repo.rs, service.rs
    │   └── settings/
    │       └── mod.rs, service.rs
    └── commands/
        ├── mod.rs
        ├── webshell.rs, project.rs, payload.rs, console.rs
        ├── plugin.rs, settings.rs, system.rs, batch.rs
```

---

## Task 1: Cargo.toml Dependencies

**Files:**
- Modify: `src-tauri/Cargo.toml`

- [ ] **Step 1: Replace Cargo.toml with full dependency set**

```toml
[package]
name = "fg-abyss"
version = "0.1.0"
edition = "2021"

[lib]
name = "fg_abyss_lib"
crate-type = ["staticlib", "cdylib", "rlib"]

[build-dependencies]
tauri-build = { version = "2", features = [] }

[dependencies]
tauri       = { version = "2", features = [] }
tauri-specta = { version = "2", features = ["derive", "typescript"] }
specta      = { version = "2" }
specta-typescript = { version = "0.0.7" }

serde       = { version = "1", features = ["derive"] }
serde_json  = "1"
thiserror   = "1"
tokio       = { version = "1", features = ["full"] }
tokio-rusqlite = "0.5"
rusqlite    = { version = "0.31", features = ["bundled"] }
uuid        = { version = "1", features = ["v4"] }
argon2      = "0.5"
aes-gcm     = "0.10"
rand        = "0.8"
base64      = "0.22"
url         = "2"
zeroize     = { version = "1", features = ["derive"] }
reqwest     = { version = "0.12", features = ["json", "multipart"] }
dashmap     = "5"
toml        = "0.8"
chrono      = { version = "0.4", features = ["serde"] }
tracing     = "0.1"
tracing-subscriber = { version = "0.3", features = ["env-filter"] }
tracing-appender   = "0.2"

[dev-dependencies]
mockall    = "0.12"
tokio-test = "0.4"

[profile.release]
codegen-units = 1
lto           = true
opt-level     = "s"
panic         = "abort"
strip         = true
```

- [ ] **Step 2: Run cargo check to verify dependency resolution**

```bash
cargo check --manifest-path src-tauri/Cargo.toml
```

Expected: compiles (may have warnings about unused items — ignore for now).

- [ ] **Step 3: Commit**

```bash
git add src-tauri/Cargo.toml src-tauri/Cargo.lock
git commit -m "chore: add all backend dependencies to Cargo.toml"
```

---

## Task 2: error.rs — Global Error Type

**Files:**
- Create: `src-tauri/src/error.rs`
- Modify: `src-tauri/src/lib.rs`

- [ ] **Step 1: Write error.rs**

```rust
// src-tauri/src/error.rs
use std::collections::HashMap;

#[derive(Debug, thiserror::Error)]
pub enum AppError {
    #[error("数据库错误: {0}")]
    Database(#[from] rusqlite::Error),
    #[error("数据库连接错误: {0}")]
    DbConnect(String),
    #[error("HTTP 错误: {0}")]
    Http(#[from] reqwest::Error),
    #[error("加密错误: {0}")]
    Crypto(String),
    #[error("WebShell 连接失败: {0}")]
    Connection(String),
    #[error("WebShell 响应验证失败: {0}")]
    InvalidResponse(String),
    #[error("熔断器开启: {0}")]
    CircuitOpen(String),
    #[error("未找到: {0}")]
    NotFound(String),
    #[error("参数无效: {0}")]
    InvalidInput(String),
    #[error("IO 错误: {0}")]
    Io(#[from] std::io::Error),
    #[error("序列化错误: {0}")]
    Serialize(#[from] serde_json::Error),
    #[error("应用已锁定")]
    Locked,
    #[error("插件错误: {0}")]
    Plugin(String),
    #[error("WebShell 需要重新部署: {0}")]
    NeedsRedeploy(String),
    #[error("内存马会话已失效（服务端重启导致）")]
    MemShellExpired,
}

impl AppError {
    pub fn kind(&self) -> &'static str {
        match self {
            Self::Database(_) | Self::DbConnect(_) => "Database",
            Self::Http(_)           => "Http",
            Self::Crypto(_)         => "Crypto",
            Self::Connection(_)     => "Connection",
            Self::InvalidResponse(_) => "InvalidResponse",
            Self::CircuitOpen(_)    => "CircuitOpen",
            Self::NotFound(_)       => "NotFound",
            Self::InvalidInput(_)   => "InvalidInput",
            Self::Io(_)             => "Io",
            Self::Serialize(_)      => "Serialize",
            Self::Locked            => "Locked",
            Self::Plugin(_)         => "Plugin",
            Self::NeedsRedeploy(_)  => "NeedsRedeploy",
            Self::MemShellExpired   => "MemShellExpired",
        }
    }
}

impl serde::Serialize for AppError {
    fn serialize<S>(&self, s: S) -> std::result::Result<S::Ok, S::Error>
    where S: serde::Serializer {
        use serde::ser::SerializeStruct;
        let mut state = s.serialize_struct("AppError", 2)?;
        state.serialize_field("kind", self.kind())?;
        state.serialize_field("message", &self.to_string())?;
        state.finish()
    }
}

pub type Result<T> = std::result::Result<T, AppError>;
```

- [ ] **Step 2: Add error module to lib.rs**

```rust
// src-tauri/src/lib.rs
pub mod error;
pub use error::{AppError, Result};
```

- [ ] **Step 3: Run cargo check**

```bash
cargo check --manifest-path src-tauri/Cargo.toml
```

Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add src-tauri/src/error.rs src-tauri/src/lib.rs
git commit -m "feat: add global AppError type with serde serialization"
```

---

## Task 3: Database Migrations (SQL Files)

**Files:**
- Create: `src-tauri/migrations/001_init.sql`
- Create: `src-tauri/migrations/002_plugins.sql`
- Create: `src-tauri/migrations/003_payload.sql`
- Create: `src-tauri/migrations/004_audit.sql`

- [ ] **Step 1: Create 001_init.sql**

```sql
-- migrations/001_init.sql
CREATE TABLE projects (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  INTEGER NOT NULL,
    updated_at  INTEGER NOT NULL,
    deleted_at  INTEGER
);

CREATE TABLE webshells (
    id                TEXT PRIMARY KEY,
    name              TEXT NOT NULL,
    url               TEXT NOT NULL,
    password          TEXT NOT NULL,
    payload_type      TEXT NOT NULL,
    project_id        TEXT REFERENCES projects(id),
    status            TEXT NOT NULL DEFAULT 'inactive'
                          CHECK (status IN ('inactive', 'active', 'needs_redeploy')),
    tags              TEXT NOT NULL DEFAULT '[]',
    custom_headers    TEXT NOT NULL DEFAULT '{}',
    cookies           TEXT NOT NULL DEFAULT '{}',
    proxy_override    TEXT,
    http_method       TEXT NOT NULL DEFAULT 'post',
    c2_profile        TEXT NOT NULL DEFAULT 'default',
    crypto_chain      TEXT NOT NULL DEFAULT '[]',
    fingerprint       TEXT,
    notes             TEXT,
    last_connected_at INTEGER,
    created_at        INTEGER NOT NULL,
    updated_at        INTEGER NOT NULL,
    deleted_at        INTEGER
);
CREATE INDEX idx_webshells_active ON webshells(deleted_at, project_id);
CREATE INDEX idx_webshells_status ON webshells(status, deleted_at);
```

- [ ] **Step 2: Create 002_plugins.sql**

```sql
-- migrations/002_plugins.sql
CREATE TABLE plugins (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    version     TEXT NOT NULL,
    enabled     INTEGER NOT NULL DEFAULT 1,
    config      TEXT NOT NULL DEFAULT '{}',
    source      TEXT NOT NULL DEFAULT 'builtin',
    created_at  INTEGER NOT NULL,
    updated_at  INTEGER NOT NULL
);
```

- [ ] **Step 3: Create 003_payload.sql**

```sql
-- migrations/003_payload.sql
CREATE TABLE payloads (
    id           TEXT PRIMARY KEY,
    name         TEXT NOT NULL,
    payload_type TEXT NOT NULL,
    config       TEXT NOT NULL,
    created_at   INTEGER NOT NULL,
    deleted_at   INTEGER
);

CREATE TABLE payload_history (
    id               TEXT PRIMARY KEY,
    payload_id       TEXT REFERENCES payloads(id),
    webshell_id      TEXT,
    code             TEXT NOT NULL,
    template_version TEXT NOT NULL DEFAULT 'v1',
    created_at       INTEGER NOT NULL
);
CREATE INDEX idx_payload_history_webshell ON payload_history(webshell_id, created_at);
```

- [ ] **Step 4: Create 004_audit.sql**

```sql
-- migrations/004_audit.sql
CREATE TABLE audit_events (
    id           TEXT PRIMARY KEY,
    webshell_id  TEXT,
    action       TEXT NOT NULL,
    detail       TEXT NOT NULL,
    created_at   INTEGER NOT NULL
);
CREATE INDEX idx_audit_webshell ON audit_events(webshell_id, created_at);
```

- [ ] **Step 5: Commit**

```bash
git add src-tauri/migrations/
git commit -m "feat: add SQLite migration files (projects, webshells, plugins, payload, audit)"
```

---

## Task 4: infra/paths.rs + infra/logger.rs

**Files:**
- Create: `src-tauri/src/infra/mod.rs`
- Create: `src-tauri/src/infra/paths.rs`
- Create: `src-tauri/src/infra/logger.rs`

- [ ] **Step 1: Create infra/mod.rs**

```rust
// src-tauri/src/infra/mod.rs
pub mod audit;
pub mod c2_profile;
pub mod config;
pub mod crypto;
pub mod db;
pub mod http;
pub mod logger;
pub mod paths;
```

- [ ] **Step 2: Create infra/paths.rs**

```rust
// src-tauri/src/infra/paths.rs
use std::path::PathBuf;
use tauri::AppHandle;
use crate::Result;

pub struct AppPaths {
    pub data_dir:    PathBuf,
    pub db_path:     PathBuf,
    pub config:      PathBuf,
    pub logs_dir:    PathBuf,
    pub plugins_dir: PathBuf,
    pub exports_dir: PathBuf,
}

impl AppPaths {
    pub fn resolve(app: &AppHandle) -> Result<Self> {
        let base = app.path().app_data_dir()
            .map_err(|e| crate::AppError::Io(std::io::Error::new(
                std::io::ErrorKind::Other, e.to_string()
            )))?;
        Ok(Self {
            db_path:     base.join("data.db"),
            config:      base.join("config.toml"),
            logs_dir:    base.join("logs"),
            plugins_dir: base.join("plugins"),
            exports_dir: base.join("exports"),
            data_dir:    base,
        })
    }
}
```

- [ ] **Step 3: Create infra/logger.rs**

```rust
// src-tauri/src/infra/logger.rs
use std::path::Path;
use tracing_appender::non_blocking::WorkerGuard;
use tracing_subscriber::prelude::*;
use crate::Result;

pub fn init(logs_dir: &Path) -> Result<WorkerGuard> {
    let file_appender = tracing_appender::rolling::daily(logs_dir, "fg-abyss.log");
    let (non_blocking, guard) = tracing_appender::non_blocking(file_appender);

    tracing_subscriber::registry()
        .with(tracing_subscriber::fmt::layer().with_writer(non_blocking))
        .with(tracing_subscriber::fmt::layer().with_writer(std::io::stderr))
        .with(tracing_subscriber::EnvFilter::new("info"))
        .init();

    Ok(guard)
}
```

- [ ] **Step 4: Add infra module to lib.rs**

```rust
// src-tauri/src/lib.rs  (append)
pub mod infra;
```

- [ ] **Step 5: Run cargo check**

```bash
cargo check --manifest-path src-tauri/Cargo.toml
```

Expected: no errors.

- [ ] **Step 6: Commit**

```bash
git add src-tauri/src/infra/
git commit -m "feat: add infra module with AppPaths and logger initialization"
```

---

## Task 5: infra/db.rs — Async SQLite Wrapper

**Files:**
- Create: `src-tauri/src/infra/db.rs`

- [ ] **Step 1: Write the failing test**

```rust
// at the bottom of src-tauri/src/infra/db.rs (write the whole file):
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
        F: FnOnce(&rusqlite::Connection) -> rusqlite::Result<T> + Send + 'static,
        T: Send + 'static,
    {
        self.0.call(f).await.map_err(|e| AppError::DbConnect(e.to_string()))
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

        // Verify tables exist
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
```

- [ ] **Step 2: Run test to verify it passes**

```bash
cargo test --manifest-path src-tauri/Cargo.toml infra::db::tests
```

Expected: PASS (both tests pass).

- [ ] **Step 3: Commit**

```bash
git add src-tauri/src/infra/db.rs
git commit -m "feat: add async SQLite Database wrapper with migration support"
```

---

## Task 6: infra/config.rs

**Files:**
- Create: `src-tauri/src/infra/config.rs`

- [ ] **Step 1: Write infra/config.rs**

```rust
// src-tauri/src/infra/config.rs
use std::path::Path;
use serde::{Deserialize, Serialize};
use crate::{AppError, Result};

#[derive(Debug, Serialize, Deserialize, Default, Clone)]
pub struct Config {
    pub meta:       MetaConfig,
    pub appearance: AppearanceConfig,
    pub window:     WindowConfig,
    pub connection: ConnectionConfig,
    pub security:   SecurityConfig,
    pub logging:    LoggingConfig,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct MetaConfig {
    pub config_version: u32,
}
impl Default for MetaConfig {
    fn default() -> Self { Self { config_version: 1 } }
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct AppearanceConfig {
    pub theme: String,
    pub language: String,
    pub accent_color: String,
}
impl Default for AppearanceConfig {
    fn default() -> Self {
        Self { theme: "dark".into(), language: "zh-CN".into(), accent_color: "#4f9cff".into() }
    }
}

#[derive(Debug, Serialize, Deserialize, Clone, Default)]
pub struct WindowConfig {
    pub x: i32, pub y: i32,
    pub width: u32, pub height: u32,
    pub maximized: bool,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ConnectionConfig {
    pub timeout_secs: u64,
    pub retry_count: u32,
    pub proxy_enabled: bool,
    pub proxy_type: String,
    pub proxy_host: String,
    pub proxy_port: u16,
}
impl Default for ConnectionConfig {
    fn default() -> Self {
        Self { timeout_secs: 30, retry_count: 3, proxy_enabled: false,
            proxy_type: "http".into(), proxy_host: "".into(), proxy_port: 7890 }
    }
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SecurityConfig {
    pub encryption: String,
    pub key_rotation_days: u32,
    pub salt: String,
    pub master_password_enabled: bool,
    pub master_password_hash: String,
    pub idle_lock_minutes: u32,
}
impl Default for SecurityConfig {
    fn default() -> Self {
        Self { encryption: "aes-256-gcm".into(), key_rotation_days: 30,
            salt: String::new(), master_password_enabled: false,
            master_password_hash: String::new(), idle_lock_minutes: 30 }
    }
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct LoggingConfig {
    pub level: String,
    pub max_file_size_mb: u64,
    pub max_files: usize,
}
impl Default for LoggingConfig {
    fn default() -> Self {
        Self { level: "info".into(), max_file_size_mb: 10, max_files: 5 }
    }
}

pub fn load_or_default(path: &Path) -> Result<Config> {
    if path.exists() {
        let text = std::fs::read_to_string(path)?;
        let mut cfg: Config = toml::from_str(&text)
            .map_err(|e| AppError::InvalidInput(e.to_string()))?;
        migrate_config(&mut cfg);
        Ok(cfg)
    } else {
        Ok(Config::default())
    }
}

pub fn save(path: &Path, cfg: &Config) -> Result<()> {
    let text = toml::to_string_pretty(cfg)
        .map_err(|e| AppError::InvalidInput(e.to_string()))?;
    std::fs::write(path, text)?;
    Ok(())
}

fn migrate_config(cfg: &mut Config) {
    // Fill missing fields for older config versions.
    // Currently a no-op: version 1 is the only version.
    let _ = cfg;
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::path::Path;

    #[test]
    fn test_default_config_serializes() {
        let cfg = Config::default();
        let text = toml::to_string_pretty(&cfg).unwrap();
        let _parsed: Config = toml::from_str(&text).unwrap();
    }

    #[test]
    fn test_load_or_default_missing_file() {
        let cfg = load_or_default(Path::new("/nonexistent/config.toml")).unwrap();
        assert_eq!(cfg.meta.config_version, 1);
        assert_eq!(cfg.appearance.theme, "dark");
    }
}
```

- [ ] **Step 2: Run tests**

```bash
cargo test --manifest-path src-tauri/Cargo.toml infra::config::tests
```

Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add src-tauri/src/infra/config.rs
git commit -m "feat: add Config struct with TOML load/save and migration support"
```

---

## Task 7: infra/crypto.rs — AES-256-GCM + Argon2id

**Files:**
- Create: `src-tauri/src/infra/crypto.rs`

- [ ] **Step 1: Write the failing tests first**

```rust
// src-tauri/src/infra/crypto.rs

use aes_gcm::{Aes256Gcm, KeyInit, aead::{Aead, AeadCore, OsRng as AeadOsRng}};
use argon2::{Argon2, PasswordHash, PasswordHasher, PasswordVerifier};
use argon2::password_hash::{rand_core::OsRng, SaltString};
use base64::{Engine as _, engine::general_purpose::STANDARD as BASE64};
use rand::RngCore;
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use zeroize::Zeroize;
use crate::{AppError, Result};

// --- Low-level functions ---

/// Encrypt plaintext with AES-256-GCM. Returns base64(nonce[12] || tag[16] || ciphertext).
pub fn encrypt(plain: &[u8], key: &[u8; 32]) -> Result<String> {
    let cipher = Aes256Gcm::new_from_slice(key)
        .map_err(|e| AppError::Crypto(e.to_string()))?;
    let nonce = Aes256Gcm::generate_nonce(&mut AeadOsRng);
    let mut ciphertext = cipher.encrypt(&nonce, plain)
        .map_err(|e| AppError::Crypto(e.to_string()))?;
    let mut combined = nonce.to_vec();
    combined.append(&mut ciphertext);
    Ok(BASE64.encode(&combined))
}

/// Decrypt base64(nonce[12] || tag[16] || ciphertext) with AES-256-GCM.
pub fn decrypt(cipher_b64: &str, key: &[u8; 32]) -> Result<Vec<u8>> {
    let combined = BASE64.decode(cipher_b64)
        .map_err(|e| AppError::Crypto(e.to_string()))?;
    if combined.len() < 12 {
        return Err(AppError::Crypto("ciphertext too short".into()));
    }
    let (nonce_bytes, ciphertext) = combined.split_at(12);
    let cipher = Aes256Gcm::new_from_slice(key)
        .map_err(|e| AppError::Crypto(e.to_string()))?;
    let nonce = aes_gcm::Nonce::from_slice(nonce_bytes);
    cipher.decrypt(nonce, ciphertext)
        .map_err(|e| AppError::Crypto(e.to_string()))
}

/// Derive 32-byte master key from salt using Argon2id (does NOT store key, caller zeroizes).
pub fn derive_key(salt_b64: &str) -> Result<[u8; 32]> {
    let salt_bytes = BASE64.decode(salt_b64)
        .map_err(|e| AppError::Crypto(e.to_string()))?;
    // Use a constant application secret mixed with the per-installation salt.
    // This prevents offline attacks if only the salt is stolen.
    const APP_SECRET: &[u8] = b"fg-abyss-master-key-v1";
    let mut key = [0u8; 32];
    Argon2::default()
        .hash_password_into(APP_SECRET, &salt_bytes, &mut key)
        .map_err(|e| AppError::Crypto(e.to_string()))?;
    Ok(key)
}

/// Generate a random base64-encoded 32-byte salt.
pub fn generate_salt() -> String {
    let mut salt = [0u8; 32];
    rand::thread_rng().fill_bytes(&mut salt);
    BASE64.encode(&salt)
}

/// Generate a random ASCII string of given length (used for response_mark).
pub fn random_ascii(len: usize) -> String {
    use rand::Rng;
    const CHARS: &[u8] = b"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    let mut rng = rand::thread_rng();
    (0..len).map(|_| CHARS[rng.gen_range(0..CHARS.len())] as char).collect()
}

// --- CryptoContext (Arc-wrapped, zeroize on drop) ---

pub struct CryptoContext {
    master_key: [u8; 32],
}

impl Drop for CryptoContext {
    fn drop(&mut self) {
        self.master_key.zeroize();
    }
}

impl CryptoContext {
    pub fn new(key: [u8; 32]) -> Arc<Self> {
        Arc::new(Self { master_key: key })
    }

    pub fn encrypt(&self, plain: &[u8]) -> Result<String> {
        encrypt(plain, &self.master_key)
    }

    pub fn decrypt(&self, cipher: &str) -> Result<Vec<u8>> {
        decrypt(cipher, &self.master_key)
    }

    pub fn decrypt_str(&self, cipher: &str) -> Result<String> {
        String::from_utf8(self.decrypt(cipher)?)
            .map_err(|e| AppError::Crypto(e.to_string()))
    }

    pub fn verify_password(&self, input: &str, hash: &str) -> bool {
        let Ok(parsed_hash) = PasswordHash::new(hash) else { return false; };
        Argon2::default()
            .verify_password(input.as_bytes(), &parsed_hash)
            .is_ok()
    }

    pub fn hash_password(password: &str) -> Result<String> {
        let salt = SaltString::generate(&mut OsRng);
        Argon2::default()
            .hash_password(password.as_bytes(), &salt)
            .map(|h| h.to_string())
            .map_err(|e| AppError::Crypto(e.to_string()))
    }
}

// --- CryptoChain ---

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "snake_case", tag = "type")]
pub enum CodecStep {
    Aes256Gcm,
    XorKey(String),
    Base64,
    UrlEncode,
    GzipCompress,
    HexEncode,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize, Default)]
pub struct CryptoChain {
    pub steps: Vec<CodecStep>,
}

impl CryptoChain {
    pub fn encode(&self, _data: &[u8], _key: &[u8; 32]) -> Result<Vec<u8>> {
        todo!("CryptoChain::encode — implement per-step encoding in Phase 2")
    }
    pub fn decode(&self, _data: &[u8], _key: &[u8; 32]) -> Result<Vec<u8>> {
        todo!("CryptoChain::decode — implement per-step decoding in Phase 2")
    }
}

// --- Sensitive<T> wrapper ---

pub struct Sensitive<T>(T);

impl<T> Sensitive<T> {
    pub fn new(val: T) -> Self { Self(val) }
    pub fn inner(&self) -> &T { &self.0 }
    pub fn into_inner(self) -> T { self.0 }
}

impl<T> std::fmt::Debug for Sensitive<T> {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        f.write_str("[REDACTED]")
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_encrypt_decrypt_roundtrip() {
        let key = [42u8; 32];
        let plain = b"hello, FG-ABYSS!";
        let cipher = encrypt(plain, &key).unwrap();
        let decrypted = decrypt(&cipher, &key).unwrap();
        assert_eq!(decrypted, plain);
    }

    #[test]
    fn test_different_keys_fail() {
        let key1 = [1u8; 32];
        let key2 = [2u8; 32];
        let cipher = encrypt(b"secret", &key1).unwrap();
        assert!(decrypt(&cipher, &key2).is_err());
    }

    #[test]
    fn test_generate_salt_unique() {
        let s1 = generate_salt();
        let s2 = generate_salt();
        assert_ne!(s1, s2);
        assert!(!s1.is_empty());
    }

    #[test]
    fn test_derive_key_deterministic() {
        let salt = generate_salt();
        let k1 = derive_key(&salt).unwrap();
        let k2 = derive_key(&salt).unwrap();
        assert_eq!(k1, k2);
    }

    #[test]
    fn test_crypto_context_encrypt_decrypt() {
        let key = [7u8; 32];
        let ctx = CryptoContext::new(key);
        let cipher = ctx.encrypt(b"password123").unwrap();
        let plain = ctx.decrypt_str(&cipher).unwrap();
        assert_eq!(plain, "password123");
    }

    #[test]
    fn test_sensitive_debug_redacted() {
        let s = Sensitive::new("my_secret");
        assert_eq!(format!("{:?}", s), "[REDACTED]");
    }
}
```

- [ ] **Step 2: Run tests**

```bash
cargo test --manifest-path src-tauri/Cargo.toml infra::crypto::tests
```

Expected: PASS (all 6 tests).

- [ ] **Step 3: Commit**

```bash
git add src-tauri/src/infra/crypto.rs
git commit -m "feat: add AES-256-GCM crypto, Argon2id key derivation, CryptoContext"
```

---

## Task 8: infra/audit.rs + infra/http.rs + infra/c2_profile.rs

**Files:**
- Create: `src-tauri/src/infra/audit.rs`
- Create: `src-tauri/src/infra/http.rs`
- Create: `src-tauri/src/infra/c2_profile.rs`

- [ ] **Step 1: Write infra/c2_profile.rs**

```rust
// src-tauri/src/infra/c2_profile.rs
use std::collections::HashMap;

#[derive(Debug, Clone)]
pub struct C2Profile {
    pub name:            String,
    pub request_param:   String,
    pub extra_headers:   HashMap<String, String>,
    pub user_agent:      String,
    pub request_wrapper: Option<String>,
    pub response_prefix: String,
    pub jitter_ms:       (u64, u64),
    pub padding_range:   Option<(usize, usize)>,
}

impl Default for C2Profile {
    fn default() -> Self {
        Self {
            name:            "default".into(),
            request_param:   "pass".into(),
            extra_headers:   HashMap::new(),
            user_agent:      "Mozilla/5.0".into(),
            request_wrapper: None,
            response_prefix: String::new(),
            jitter_ms:       (0, 0),
            padding_range:   None,
        }
    }
}

impl C2Profile {
    pub fn by_name(name: &str) -> Self {
        match name {
            "cdn-callback" => Self {
                name:          "cdn-callback".into(),
                request_param: "data".into(),
                extra_headers: [("Content-Type".into(), "application/octet-stream".into())].into(),
                ..Default::default()
            },
            "api-json" => Self {
                name:            "api-json".into(),
                request_param:   "data".into(),
                request_wrapper: Some(r#"{"code":0,"data":"{DATA}"}"#.into()),
                ..Default::default()
            },
            "form-submit" => Self {
                name:          "form-submit".into(),
                request_param: "_token".into(),
                ..Default::default()
            },
            _ => Self::default(),
        }
    }
}
```

- [ ] **Step 2: Write infra/audit.rs**

```rust
// src-tauri/src/infra/audit.rs
use serde::{Deserialize, Serialize};
use uuid::Uuid;
use crate::{Result, infra::db::Database};

#[derive(Debug, Clone)]
pub struct AuditLog {
    db: Database,
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
```

- [ ] **Step 3: Write infra/http.rs**

```rust
// src-tauri/src/infra/http.rs
use std::collections::HashMap;
use std::sync::Arc;
use std::time::Instant;
use dashmap::DashMap;
use crate::{AppError, Result};
use crate::infra::config::ConnectionConfig;

const FAILURE_THRESHOLD: u32 = 3;
const CIRCUIT_RESET_SECS: u64 = 30;

#[derive(Debug, Clone)]
pub struct WebshellHttpConfig {
    pub url:            String,
    pub method:         String,
    pub request_param:  String,
    pub proxy_override: Option<String>,
    pub custom_headers: HashMap<String, String>,
    pub cookies:        HashMap<String, String>,
    pub timeout_secs:   u64,
    pub jitter_ms:      Option<(u64, u64)>,
}

impl WebshellHttpConfig {
    pub fn proxy_hash(&self) -> u64 {
        use std::hash::{Hash, Hasher};
        use std::collections::hash_map::DefaultHasher;
        let mut h = DefaultHasher::new();
        self.proxy_override.hash(&mut h);
        h.finish()
    }
}

enum CircuitState { Closed, Open(Instant), HalfOpen }

pub struct HttpClientPool {
    clients:       DashMap<u64, reqwest::Client>,
    circuit:       DashMap<String, CircuitState>,
    failure_count: DashMap<String, u32>,
    global_config: ConnectionConfig,
}

impl HttpClientPool {
    pub fn new(config: &ConnectionConfig) -> Self {
        Self {
            clients:       DashMap::new(),
            circuit:       DashMap::new(),
            failure_count: DashMap::new(),
            global_config: config.clone(),
        }
    }

    pub fn get_client(&self, config: &WebshellHttpConfig) -> reqwest::Client {
        let key = config.proxy_hash();
        self.clients
            .entry(key)
            .or_insert_with(|| build_client(config, &self.global_config))
            .clone()
    }

    pub async fn send_raw(
        &self,
        webshell_id: &str,
        req: reqwest::RequestBuilder,
    ) -> Result<reqwest::Response> {
        // Check circuit breaker
        let is_open = self.circuit.get(webshell_id).map(|s| {
            matches!(*s, CircuitState::Open(t) if t.elapsed().as_secs() < CIRCUIT_RESET_SECS)
        }).unwrap_or(false);
        if is_open {
            return Err(AppError::CircuitOpen(webshell_id.to_string()));
        }

        match req.send().await {
            Ok(resp) => {
                self.failure_count.remove(webshell_id);
                self.circuit.insert(webshell_id.to_string(), CircuitState::Closed);
                Ok(resp)
            }
            Err(e) => {
                let count = {
                    let mut entry = self.failure_count
                        .entry(webshell_id.to_string())
                        .and_modify(|c| *c += 1)
                        .or_insert(1);
                    *entry
                };
                if count >= FAILURE_THRESHOLD {
                    self.circuit.insert(
                        webshell_id.to_string(),
                        CircuitState::Open(Instant::now()),
                    );
                }
                Err(AppError::Http(e))
            }
        }
    }
}

fn build_client(config: &WebshellHttpConfig, global: &ConnectionConfig) -> reqwest::Client {
    let mut builder = reqwest::Client::builder()
        .timeout(std::time::Duration::from_secs(config.timeout_secs));

    if let Some(proxy_url) = &config.proxy_override {
        if let Ok(proxy) = reqwest::Proxy::all(proxy_url) {
            builder = builder.proxy(proxy);
        }
    } else if global.proxy_enabled && !global.proxy_host.is_empty() {
        let proxy_url = format!("{}://{}:{}", global.proxy_type, global.proxy_host, global.proxy_port);
        if let Ok(proxy) = reqwest::Proxy::all(&proxy_url) {
            builder = builder.proxy(proxy);
        }
    }

    builder.build().unwrap_or_default()
}
```

- [ ] **Step 4: Run cargo check**

```bash
cargo check --manifest-path src-tauri/Cargo.toml
```

Expected: no errors.

- [ ] **Step 5: Commit**

```bash
git add src-tauri/src/infra/audit.rs src-tauri/src/infra/http.rs src-tauri/src/infra/c2_profile.rs
git commit -m "feat: add AuditLog, HttpClientPool with circuit breaker, C2Profile"
```

---

## Task 9: features/webshell/ — Models, Repo, Session, Queue

**Files:**
- Create: `src-tauri/src/features/mod.rs`
- Create: `src-tauri/src/features/webshell/mod.rs`
- Create: `src-tauri/src/features/webshell/models.rs`
- Create: `src-tauri/src/features/webshell/repo.rs`
- Create: `src-tauri/src/features/webshell/session.rs`
- Create: `src-tauri/src/features/webshell/queue.rs`

- [ ] **Step 1: Create features/mod.rs**

```rust
// src-tauri/src/features/mod.rs
pub mod console;
pub mod payload;
pub mod plugin;
pub mod project;
pub mod settings;
pub mod webshell;
```

- [ ] **Step 2: Create features/webshell/mod.rs**

```rust
// src-tauri/src/features/webshell/mod.rs
pub mod models;
pub mod queue;
pub mod repo;
pub mod service;
pub mod session;
```

- [ ] **Step 3: Create features/webshell/models.rs**

```rust
// src-tauri/src/features/webshell/models.rs
use std::collections::HashMap;
use serde::{Deserialize, Serialize};
use crate::infra::crypto::CryptoChain;

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "lowercase", tag = "type")]
pub enum PayloadType {
    Php, Jsp, Asp, Aspx,
}

impl PayloadType {
    pub fn extension(&self) -> &'static str {
        match self { Self::Php => "php", Self::Jsp => "jsp", Self::Asp => "asp", Self::Aspx => "aspx" }
    }
    pub fn is_memshell(&self) -> bool { false }
}

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct Webshell {
    pub id: String,
    pub name: String,
    pub url: String,
    pub password: String,
    pub payload_type: PayloadType,
    pub project_id: Option<String>,
    pub status: String,
    pub tags: Vec<String>,
    pub custom_headers: HashMap<String, String>,
    pub cookies: HashMap<String, String>,
    pub proxy_override: Option<String>,
    pub http_method: String,
    pub c2_profile: String,
    pub crypto_chain: CryptoChain,
    pub fingerprint: Option<WebshellFingerprint>,
    pub notes: Option<String>,
    pub last_connected_at: Option<i64>,
    pub created_at: i64,
    pub updated_at: i64,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct WebshellFingerprint {
    pub os:           Option<String>,
    pub runtime:      Option<String>,
    pub runtime_ver:  Option<String>,
    pub username:     Option<String>,
    pub cwd:          Option<String>,
    pub disabled_fns: Vec<String>,
    pub writable:     Option<bool>,
    pub last_probed_at: Option<i64>,
}

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct ConnectionResult {
    pub success:     bool,
    pub latency_ms:  Option<u64>,
    pub error:       Option<String>,
    pub fingerprint: Option<WebshellFingerprint>,
}

#[derive(Debug, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct CreateWebshellInput {
    pub name: String,
    pub url: String,
    pub password: String,
    pub payload_type: PayloadType,
    pub project_id: Option<String>,
    #[serde(default)] pub tags: Vec<String>,
    #[serde(default)] pub custom_headers: HashMap<String, String>,
    #[serde(default)] pub cookies: HashMap<String, String>,
    pub proxy_override: Option<String>,
    #[serde(default = "default_post")] pub http_method: String,
    #[serde(default = "default_c2")]   pub c2_profile: String,
    #[serde(default)] pub crypto_chain: CryptoChain,
    pub notes: Option<String>,
}
fn default_post() -> String { "post".to_string() }
fn default_c2()   -> String { "default".to_string() }

#[derive(Debug, Deserialize, Default, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct UpdateWebshellInput {
    pub name: Option<String>,
    pub url: Option<String>,
    pub password: Option<String>,
    pub project_id: Option<String>,
    pub tags: Option<Vec<String>>,
    pub custom_headers: Option<HashMap<String, String>>,
    pub cookies: Option<HashMap<String, String>>,
    pub proxy_override: Option<String>,
    pub http_method: Option<String>,
    pub c2_profile: Option<String>,
    pub crypto_chain: Option<CryptoChain>,
    pub notes: Option<String>,
    pub status: Option<String>,
}
```

- [ ] **Step 4: Create features/webshell/repo.rs**

```rust
// src-tauri/src/features/webshell/repo.rs
use std::sync::Arc;
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
    pub password:       Option<String>,  // already encrypted
    pub project_id:     Option<Option<String>>,
    pub tags:           Option<String>,  // JSON
    pub custom_headers: Option<String>,  // JSON
    pub cookies:        Option<String>,  // JSON
    pub proxy_override: Option<Option<String>>,
    pub http_method:    Option<String>,
    pub c2_profile:     Option<String>,
    pub crypto_chain:   Option<String>,  // JSON
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
                "SELECT * FROM webshells WHERE deleted_at IS NULL AND project_id = ?1 ORDER BY created_at DESC"
            } else {
                "SELECT * FROM webshells WHERE deleted_at IS NULL ORDER BY created_at DESC"
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
        let id = id.to_string();
        self.db.call(move |conn| {
            conn.query_row(
                "SELECT * FROM webshells WHERE id = ?1 AND deleted_at IS NULL",
                [&id], row_to_webshell,
            )
        }).await.map_err(|e| {
            if let crate::AppError::Database(ref re) = e {
                if matches!(re, rusqlite::Error::QueryReturnedNoRows) {
                    return crate::AppError::NotFound(format!("webshell {}", id));
                }
            }
            e
        })
    }

    async fn find_by_url(&self, url: &str) -> Result<Option<Webshell>> {
        let url = url.to_string();
        self.db.call(move |conn| {
            let mut stmt = conn.prepare(
                "SELECT * FROM webshells WHERE url = ?1 AND deleted_at IS NULL LIMIT 1"
            )?;
            let mut rows = stmt.query_map([&url], row_to_webshell)?;
            Ok(rows.next().transpose()?)
        }).await
    }

    async fn insert(&self, w: &Webshell) -> Result<()> {
        let w = w.clone();
        self.db.call(move |conn| {
            conn.execute(
                "INSERT INTO webshells (id, name, url, password, payload_type, project_id, status,
                 tags, custom_headers, cookies, proxy_override, http_method, c2_profile,
                 crypto_chain, fingerprint, notes, last_connected_at, created_at, updated_at)
                 VALUES (?1,?2,?3,?4,?5,?6,?7,?8,?9,?10,?11,?12,?13,?14,?15,?16,?17,?18,?19)",
                rusqlite::params![
                    w.id, w.name, w.url, w.password,
                    serde_json::to_string(&w.payload_type).unwrap(),
                    w.project_id, w.status,
                    serde_json::to_string(&w.tags).unwrap(),
                    serde_json::to_string(&w.custom_headers).unwrap(),
                    serde_json::to_string(&w.cookies).unwrap(),
                    w.proxy_override, w.http_method, w.c2_profile,
                    serde_json::to_string(&w.crypto_chain).unwrap(),
                    w.fingerprint.as_ref().and_then(|f| serde_json::to_string(f).ok()),
                    w.notes, w.last_connected_at, w.created_at, w.updated_at,
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
                  name           = COALESCE(?1, name),
                  url            = COALESCE(?2, url),
                  password       = COALESCE(?3, password),
                  project_id     = COALESCE(?4, project_id),
                  tags           = COALESCE(?5, tags),
                  custom_headers = COALESCE(?6, custom_headers),
                  cookies        = COALESCE(?7, cookies),
                  proxy_override = COALESCE(?8, proxy_override),
                  http_method    = COALESCE(?9, http_method),
                  c2_profile     = COALESCE(?10, c2_profile),
                  crypto_chain   = COALESCE(?11, crypto_chain),
                  notes          = COALESCE(?12, notes),
                  status         = COALESCE(?13, status),
                  updated_at     = strftime('%s','now')
                 WHERE id = ?14 AND deleted_at IS NULL",
                rusqlite::params![
                    f.name, f.url, f.password,
                    f.project_id.flatten(),
                    f.tags, f.custom_headers, f.cookies,
                    f.proxy_override.flatten(),
                    f.http_method, f.c2_profile, f.crypto_chain,
                    f.notes.flatten(),
                    f.status, id,
                ],
            )?;
            Ok(())
        }).await
    }

    async fn update_status(&self, id: &str, status: &str, fp: Option<&WebshellFingerprint>) -> Result<()> {
        let id = id.to_string();
        let status = status.to_string();
        let fp_json = fp.and_then(|f| serde_json::to_string(f).ok());
        self.db.call(move |conn| {
            conn.execute(
                "UPDATE webshells SET status = ?1, fingerprint = COALESCE(?2, fingerprint),
                 last_connected_at = CASE WHEN ?1 = 'active' THEN strftime('%s','now') ELSE last_connected_at END,
                 updated_at = strftime('%s','now')
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
    use crate::infra::crypto::CryptoChain;
    use super::models::{PayloadType, WebshellFingerprint};
    use std::collections::HashMap;

    let payload_type_str: String = row.get(4)?;
    let tags_str:         String = row.get(7)?;
    let headers_str:      String = row.get(8)?;
    let cookies_str:      String = row.get(9)?;
    let chain_str:        String = row.get(12)?;
    let fp_str:           Option<String> = row.get(13)?;

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
        c2_profile:       row.get(12)?, // Note: re-read positionally
        crypto_chain:     serde_json::from_str(&chain_str).unwrap_or_default(),
        fingerprint:      fp_str.and_then(|s| serde_json::from_str(&s).ok()),
        notes:            row.get(15)?,
        last_connected_at: row.get(16)?,
        created_at:       row.get(17)?,
        updated_at:       row.get(18)?,
    })
}
```

> **Note:** `row_to_webshell` uses positional indices. Verify against the INSERT column order: `(id=0, name=1, url=2, password=3, payload_type=4, project_id=5, status=6, tags=7, custom_headers=8, cookies=9, proxy_override=10, http_method=11, c2_profile=12, crypto_chain=13, fingerprint=14, notes=15, last_connected_at=16, created_at=17, updated_at=18)`. Adjust if your SELECT `*` order differs.

- [ ] **Step 5: Create features/webshell/session.rs**

```rust
// src-tauri/src/features/webshell/session.rs
use std::time::Instant;
use zeroize::Zeroizing;

pub struct WebshellSession {
    pub session_id:    String,
    pub session_key:   Zeroizing<[u8; 32]>,
    pub init_time:     Instant,
    pub response_mark: (String, String),
    pub timeout_secs:  u64,
}

impl WebshellSession {
    pub fn is_expired(&self) -> bool {
        self.init_time.elapsed().as_secs() > self.timeout_secs
    }
}
```

- [ ] **Step 6: Create features/webshell/queue.rs**

```rust
// src-tauri/src/features/webshell/queue.rs
use std::sync::Arc;
use dashmap::DashMap;
use tokio::sync::{mpsc, oneshot};
use tokio::task::JoinHandle;
use crate::{AppError, Result};
use crate::infra::http::{HttpClientPool, WebshellHttpConfig};

pub struct QueuedRequest {
    pub config:     WebshellHttpConfig,
    pub body:       Vec<u8>,
    pub respond_to: oneshot::Sender<Result<reqwest::Response>>,
}

pub struct WebshellQueue {
    workers:   DashMap<String, (mpsc::Sender<QueuedRequest>, JoinHandle<()>)>,
    http_pool: Arc<HttpClientPool>,
}

impl WebshellQueue {
    pub fn new(http_pool: Arc<HttpClientPool>) -> Self {
        Self { workers: DashMap::new(), http_pool }
    }

    pub async fn send(
        &self,
        webshell_id: &str,
        config: WebshellHttpConfig,
        body: Vec<u8>,
    ) -> Result<reqwest::Response> {
        let (resp_tx, resp_rx) = oneshot::channel();

        let tx = if let Some(entry) = self.workers.get(webshell_id) {
            entry.0.clone()
        } else {
            let (tx, rx) = mpsc::channel::<QueuedRequest>(32);
            let pool   = self.http_pool.clone();
            let id     = webshell_id.to_string();
            let handle = tokio::spawn(queue_worker(rx, pool, id));
            self.workers.insert(webshell_id.to_string(), (tx.clone(), handle));
            tx
        };

        tx.send(QueuedRequest { config, body, respond_to: resp_tx }).await
            .map_err(|_| AppError::Connection("queue dropped".into()))?;
        resp_rx.await
            .map_err(|_| AppError::Connection("response channel closed".into()))?
    }

    pub fn cleanup(&self, webshell_id: &str) {
        if let Some((_, (_, handle))) = self.workers.remove(webshell_id) {
            handle.abort();
        }
    }

    pub fn shutdown_all(&self) {
        for entry in self.workers.iter() {
            entry.value().1.abort();
        }
        self.workers.clear();
    }
}

async fn queue_worker(
    mut rx: mpsc::Receiver<QueuedRequest>,
    pool:   Arc<HttpClientPool>,
    id:     String,
) {
    while let Some(req) = rx.recv().await {
        if let Some((min, max)) = req.config.jitter_ms {
            let delay = rand::random::<u64>() % (max - min + 1) + min;
            tokio::time::sleep(tokio::time::Duration::from_millis(delay)).await;
        }

        let client = pool.get_client(&req.config);
        let url    = req.config.url.clone();
        let method = req.config.method.to_uppercase();

        let mut builder = match method.as_str() {
            "GET" => client.get(&url),
            _     => client.post(&url),
        };
        for (k, v) in &req.config.custom_headers {
            builder = builder.header(k, v);
        }
        if !req.config.cookies.is_empty() {
            let cookie_str = req.config.cookies.iter()
                .map(|(k, v)| format!("{k}={v}")).collect::<Vec<_>>().join("; ");
            builder = builder.header("Cookie", cookie_str);
        }
        use base64::Engine as _;
        let b64_body = base64::engine::general_purpose::STANDARD.encode(&req.body);
        builder = builder
            .form(&[(&req.config.request_param, b64_body)])
            .timeout(std::time::Duration::from_secs(req.config.timeout_secs));

        let result = pool.send_raw(&id, builder).await;
        let _ = req.respond_to.send(result);
    }
}
```

- [ ] **Step 7: Add async-trait to Cargo.toml and run cargo check**

Add to `[dependencies]` in Cargo.toml:
```toml
async-trait = "0.1"
```

```bash
cargo check --manifest-path src-tauri/Cargo.toml
```

Expected: no errors (warnings about unused items are OK).

- [ ] **Step 8: Commit**

```bash
git add src-tauri/src/features/ src-tauri/Cargo.toml src-tauri/Cargo.lock
git commit -m "feat: add webshell models, repo trait, session, and request queue"
```

---

## Task 10: features/webshell/service.rs

**Files:**
- Create: `src-tauri/src/features/webshell/service.rs`

- [ ] **Step 1: Write the failing test for WebshellService::create**

```rust
// src-tauri/src/features/webshell/service.rs
use std::sync::Arc;
use dashmap::DashMap;
use uuid::Uuid;
use chrono::Utc;
use crate::{AppError, Result};
use crate::infra::crypto::CryptoContext;
use super::models::{Webshell, CreateWebshellInput, UpdateWebshellInput, ConnectionResult};
use super::repo::{WebshellRepo, WebshellUpdateFields};
use super::session::WebshellSession;
use super::queue::WebshellQueue;

pub struct WebshellService {
    repo:     Arc<dyn WebshellRepo>,
    crypto:   Arc<CryptoContext>,
    pub queue: Arc<WebshellQueue>,
    sessions: Arc<DashMap<String, WebshellSession>>,
}

impl WebshellService {
    pub fn new(
        repo:   Arc<dyn WebshellRepo>,
        crypto: Arc<CryptoContext>,
        queue:  Arc<WebshellQueue>,
    ) -> Self {
        Self { repo, crypto, queue, sessions: Arc::new(DashMap::new()) }
    }

    pub async fn list(&self, project_id: Option<&str>) -> Result<Vec<Webshell>> {
        let mut shells = self.repo.find_all(project_id).await?;
        for s in &mut shells {
            s.password = self.crypto.decrypt_str(&s.password).unwrap_or_default();
        }
        Ok(shells)
    }

    pub async fn get(&self, id: &str) -> Result<Webshell> {
        let mut w = self.repo.find_by_id(id).await?;
        w.password = self.crypto.decrypt_str(&w.password)?;
        Ok(w)
    }

    pub async fn create(&self, input: CreateWebshellInput) -> Result<Webshell> {
        let url = normalize_url(&input.url)?;

        // Uniqueness check
        if self.repo.find_by_url(&url).await?.is_some() {
            return Err(AppError::InvalidInput(format!("WebShell URL already exists: {url}")));
        }

        let enc_password = self.crypto.encrypt(input.password.as_bytes())?;
        let plain_pw     = input.password.clone();
        let now          = Utc::now().timestamp();
        let id           = Uuid::new_v4().to_string();

        let w = Webshell {
            id, name: input.name, url,
            password:       enc_password,
            payload_type:   input.payload_type,
            project_id:     input.project_id,
            status:         "inactive".into(),
            tags:           input.tags,
            custom_headers: input.custom_headers,
            cookies:        input.cookies,
            proxy_override: input.proxy_override,
            http_method:    input.http_method,
            c2_profile:     input.c2_profile,
            crypto_chain:   input.crypto_chain,
            fingerprint:    None,
            notes:          input.notes,
            last_connected_at: None,
            created_at:     now,
            updated_at:     now,
        };

        self.repo.insert(&w).await?;

        let mut result = w;
        result.password = plain_pw;
        Ok(result)
    }

    pub async fn update(&self, id: &str, input: UpdateWebshellInput) -> Result<Webshell> {
        let current = self.get(id).await?;

        // Detect crypto_chain change → needs_redeploy
        let new_status = match &input.crypto_chain {
            Some(chain) if *chain != current.crypto_chain => "needs_redeploy".to_string(),
            _ => input.status.clone().unwrap_or(current.status.clone()),
        };

        let enc_password = match &input.password {
            Some(p) => Some(self.crypto.encrypt(p.as_bytes())?),
            None    => None,
        };

        let url = match &input.url {
            Some(u) => Some(normalize_url(u)?),
            None => None,
        };

        self.repo.update_fields(id, WebshellUpdateFields {
            name:           input.name,
            url,
            password:       enc_password,
            project_id:     input.project_id.map(Some),
            tags:           input.tags.as_ref().and_then(|t| serde_json::to_string(t).ok()),
            custom_headers: input.custom_headers.as_ref().and_then(|h| serde_json::to_string(h).ok()),
            cookies:        input.cookies.as_ref().and_then(|c| serde_json::to_string(c).ok()),
            proxy_override: input.proxy_override.map(Some),
            http_method:    input.http_method,
            c2_profile:     input.c2_profile,
            crypto_chain:   input.crypto_chain.as_ref().and_then(|c| serde_json::to_string(c).ok()),
            notes:          input.notes.map(Some),
            status:         Some(new_status),
        }).await?;

        self.get(id).await
    }

    pub async fn delete(&self, id: &str) -> Result<()> {
        let ts = Utc::now().timestamp();
        self.sessions.remove(id);
        self.queue.cleanup(id);
        self.repo.soft_delete(id, ts).await
    }

    pub async fn test_connection(&self, id: &str) -> Result<ConnectionResult> {
        let shell = self.get(id).await?;
        let start = std::time::Instant::now();
        // TODO: Phase 1 Init protocol — implement in Phase 2
        // For now return a stub result
        let _ = shell;
        let _ = start;
        todo!("Phase 1 Init protocol — implement in Phase 2")
    }

    pub async fn exec(
        &self,
        id: &str,
        method: &str,
        args: serde_json::Value,
    ) -> Result<serde_json::Value> {
        let shell = self.get(id).await?;
        if shell.status == "needs_redeploy" {
            return Err(AppError::NeedsRedeploy(
                format!("WebShell {} payload changed, redeploy before exec", id)
            ));
        }
        // TODO: Transparent Init-or-reuse + Phase 2 Exec — implement in Phase 2
        let _ = (method, args);
        todo!("Phase 2 Exec protocol — implement in Phase 2")
    }

    pub async fn reset_redeploy_status(&self, id: &str) -> Result<Webshell> {
        self.repo.update_status(id, "inactive", None).await?;
        self.sessions.remove(id);
        self.get(id).await
    }
}

fn normalize_url(raw: &str) -> Result<String> {
    let trimmed = raw.trim();
    let parsed  = url::Url::parse(trimmed)
        .map_err(|e| AppError::InvalidInput(format!("invalid URL: {e}")))?;
    let scheme   = parsed.scheme();
    let host     = parsed.host_str().unwrap_or("");
    let port_str = match parsed.port() {
        Some(p) => format!(":{p}"),
        None    => String::new(),
    };
    let path = parsed.path().trim_end_matches('/');
    let path = if path.is_empty() { "/" } else { path };
    Ok(format!("{scheme}://{host}{port_str}{path}"))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_normalize_url() {
        assert_eq!(normalize_url("HTTP://Example.COM/shell.php").unwrap(), "http://Example.COM/shell.php");
        assert_eq!(normalize_url("https://example.com/shell.php/").unwrap(), "https://example.com/shell.php");
        assert!(normalize_url("not-a-url").is_err());
    }
}
```

- [ ] **Step 2: Run tests**

```bash
cargo test --manifest-path src-tauri/Cargo.toml features::webshell::service::tests
```

Expected: PASS (normalize_url tests).

- [ ] **Step 3: Commit**

```bash
git add src-tauri/src/features/webshell/service.rs
git commit -m "feat: add WebshellService with CRUD operations (protocol as todo!)"
```

---

## Task 11: features/project/ — Full CRUD

**Files:**
- Create: `src-tauri/src/features/project/mod.rs`
- Create: `src-tauri/src/features/project/models.rs`
- Create: `src-tauri/src/features/project/repo.rs`
- Create: `src-tauri/src/features/project/service.rs`

- [ ] **Step 1: Create all project feature files**

**features/project/mod.rs:**
```rust
pub mod models;
pub mod repo;
pub mod service;
```

**features/project/models.rs:**
```rust
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct Project {
    pub id:          String,
    pub name:        String,
    pub description: Option<String>,
    pub created_at:  i64,
    pub updated_at:  i64,
}

#[derive(Debug, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct CreateProjectInput {
    pub name:        String,
    pub description: Option<String>,
}

#[derive(Debug, Deserialize, Default, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct UpdateProjectInput {
    pub name:        Option<String>,
    pub description: Option<String>,
}
```

**features/project/repo.rs:**
```rust
use std::sync::Arc;
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
        let id = id.to_string();
        self.db.call(move |conn| {
            conn.query_row(
                "SELECT id, name, description, created_at, updated_at
                 FROM projects WHERE id = ?1 AND deleted_at IS NULL",
                [&id],
                |r| Ok(Project {
                    id: r.get(0)?, name: r.get(1)?, description: r.get(2)?,
                    created_at: r.get(3)?, updated_at: r.get(4)?,
                }),
            )
        }).await.map_err(|e| {
            if let crate::AppError::Database(ref re) = e {
                if matches!(re, rusqlite::Error::QueryReturnedNoRows) {
                    return crate::AppError::NotFound(format!("project {}", id));
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
```

**features/project/service.rs:**
```rust
use std::sync::Arc;
use uuid::Uuid;
use chrono::Utc;
use crate::Result;
use super::models::{Project, CreateProjectInput, UpdateProjectInput};
use super::repo::ProjectRepo;

pub struct ProjectService {
    repo: Arc<dyn ProjectRepo>,
}

impl ProjectService {
    pub fn new(repo: Arc<dyn ProjectRepo>) -> Self { Self { repo } }

    pub async fn list(&self) -> Result<Vec<Project>> {
        self.repo.find_all().await
    }

    pub async fn get(&self, id: &str) -> Result<Project> {
        self.repo.find_by_id(id).await
    }

    pub async fn create(&self, input: CreateProjectInput) -> Result<Project> {
        let now = Utc::now().timestamp();
        let p = Project {
            id:          Uuid::new_v4().to_string(),
            name:        input.name,
            description: input.description,
            created_at:  now,
            updated_at:  now,
        };
        self.repo.insert(&p).await?;
        Ok(p)
    }

    pub async fn update(&self, id: &str, input: UpdateProjectInput) -> Result<Project> {
        self.repo.update(id, input.name, input.description).await?;
        self.repo.find_by_id(id).await
    }

    pub async fn delete(&self, id: &str) -> Result<()> {
        let ts = Utc::now().timestamp();
        self.repo.soft_delete_with_webshells(id, ts).await
    }
}
```

- [ ] **Step 2: Run cargo check**

```bash
cargo check --manifest-path src-tauri/Cargo.toml
```

Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add src-tauri/src/features/project/
git commit -m "feat: add project CRUD with atomic soft-delete cascade to webshells"
```

---

## Task 12: features/payload/ + features/plugin/

**Files:**
- Create: `src-tauri/src/features/payload/` (mod.rs, models.rs, repo.rs, service.rs, generator/mod.rs)
- Create: `src-tauri/src/features/plugin/` (mod.rs, models.rs, repo.rs, service.rs)

- [ ] **Step 1: Create all payload and plugin feature files**

**features/payload/mod.rs:**
```rust
pub mod generator;
pub mod models;
pub mod repo;
pub mod service;
```

**features/payload/models.rs:**
```rust
use serde::{Deserialize, Serialize};
use crate::infra::crypto::CryptoChain;
use super::super::webshell::models::PayloadType;
use crate::features::payload::generator::ObfuscateConfig;

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct Payload {
    pub id:           String,
    pub name:         String,
    pub payload_type: PayloadType,
    pub config:       serde_json::Value,
    pub created_at:   i64,
}

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct PayloadConfig {
    pub payload_type:   PayloadType,
    pub password:       String,
    pub crypto_chain:   CryptoChain,
    pub c2_profile:     String,
    pub obfuscate:      ObfuscateConfig,
    pub target_version: Option<String>,
}

#[derive(Debug, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct CreatePayloadInput {
    pub name:   String,
    pub config: PayloadConfig,
}

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct PayloadHistoryEntry {
    pub id:               String,
    pub payload_id:       Option<String>,
    pub webshell_id:      Option<String>,
    pub code:             String,
    pub template_version: String,
    pub created_at:       i64,
}
```

**features/payload/generator/mod.rs:**
```rust
use serde::{Deserialize, Serialize};
use crate::Result;
use super::models::PayloadConfig;

#[derive(Debug, Clone, Serialize, Deserialize, Default, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct ObfuscateConfig {
    pub level:          u8,
    pub seed:           Option<u64>,
    pub target_version: Option<String>,
}

pub fn generate(config: &PayloadConfig) -> Result<String> {
    // TODO: implement PHP/JSP/ASP/ASPX generators in Phase 2
    let _ = config;
    todo!("Payload generator — implement per-language in Phase 2")
}
```

**features/payload/repo.rs:**
```rust
use std::sync::Arc;
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
```

**features/payload/service.rs:**
```rust
use std::sync::Arc;
use uuid::Uuid;
use chrono::Utc;
use crate::{Result};
use crate::infra::crypto::CryptoContext;
use super::models::{Payload, CreatePayloadInput, PayloadHistoryEntry, PayloadConfig};
use super::repo::PayloadRepo;
use super::generator;

pub struct PayloadService {
    repo:   Arc<dyn PayloadRepo>,
    crypto: Arc<CryptoContext>,
}

impl PayloadService {
    pub fn new(repo: Arc<dyn PayloadRepo>, crypto: Arc<CryptoContext>) -> Self {
        Self { repo, crypto }
    }

    pub async fn list(&self) -> Result<Vec<Payload>> { self.repo.find_all().await }
    pub async fn get(&self, id: &str) -> Result<Payload> { self.repo.find_by_id(id).await }

    pub async fn create(&self, input: CreatePayloadInput) -> Result<Payload> {
        let now = Utc::now().timestamp();
        let p = Payload {
            id:           Uuid::new_v4().to_string(),
            name:         input.name,
            payload_type: input.config.payload_type.clone(),
            config:       serde_json::to_value(&input.config)?,
            created_at:   now,
        };
        self.repo.insert(&p).await?;
        Ok(p)
    }

    pub async fn generate_payload(
        &self,
        payload_id: Option<&str>,
        webshell_id: Option<&str>,
        config: PayloadConfig,
    ) -> Result<String> {
        let code = generator::generate(&config)?;
        let entry = PayloadHistoryEntry {
            id:               Uuid::new_v4().to_string(),
            payload_id:       payload_id.map(|s| s.to_string()),
            webshell_id:      webshell_id.map(|s| s.to_string()),
            code:             code.clone(),
            template_version: "v1".into(),
            created_at:       Utc::now().timestamp(),
        };
        self.repo.insert_history(&entry).await?;
        Ok(code)
    }

    pub async fn list_history(&self, payload_id: &str) -> Result<Vec<PayloadHistoryEntry>> {
        self.repo.find_history(payload_id).await
    }
}
```

**features/plugin/mod.rs:**
```rust
pub mod models;
pub mod repo;
pub mod service;
```

**features/plugin/models.rs:**
```rust
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct Plugin {
    pub id:         String,
    pub name:       String,
    pub version:    String,
    pub enabled:    bool,
    pub config:     serde_json::Value,
    pub source:     String,
    pub created_at: i64,
    pub updated_at: i64,
}
```

**features/plugin/repo.rs:**
```rust
use std::sync::Arc;
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
            Ok(stmt.query_map([], row_to_plugin)?.collect::<rusqlite::Result<Vec<_>>>()?)
        }).await
    }

    async fn find_by_id(&self, id: &str) -> Result<Plugin> {
        let id = id.to_string();
        self.db.call(move |conn| {
            conn.query_row(
                "SELECT id, name, version, enabled, config, source, created_at, updated_at
                 FROM plugins WHERE id = ?1",
                [&id], row_to_plugin,
            )
        }).await
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
```

**features/plugin/service.rs:**
```rust
use std::path::PathBuf;
use std::sync::Arc;
use crate::Result;
use super::models::Plugin;
use super::repo::PluginRepo;

pub struct PluginService {
    repo:        Arc<dyn PluginRepo>,
    plugins_dir: PathBuf,
}

impl PluginService {
    pub fn new(repo: Arc<dyn PluginRepo>, plugins_dir: PathBuf) -> Self {
        Self { repo, plugins_dir }
    }

    pub async fn list(&self) -> Result<Vec<Plugin>> { self.repo.find_all().await }
    pub async fn get(&self, id: &str) -> Result<Plugin> { self.repo.find_by_id(id).await }

    pub async fn enable(&self, id: &str) -> Result<Plugin> {
        self.repo.set_enabled(id, true).await?;
        self.repo.find_by_id(id).await
    }

    pub async fn disable(&self, id: &str) -> Result<Plugin> {
        self.repo.set_enabled(id, false).await?;
        self.repo.find_by_id(id).await
    }
}
```

- [ ] **Step 2: Run cargo check**

```bash
cargo check --manifest-path src-tauri/Cargo.toml
```

Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add src-tauri/src/features/payload/ src-tauri/src/features/plugin/
git commit -m "feat: add payload and plugin feature services with CRUD operations"
```

---

## Task 13: features/settings/ + features/console/

**Files:**
- Create: `src-tauri/src/features/settings/mod.rs`
- Create: `src-tauri/src/features/settings/service.rs`
- Create: `src-tauri/src/features/console/mod.rs`
- Create: `src-tauri/src/features/console/service.rs`

- [ ] **Step 1: Create settings feature**

**features/settings/mod.rs:**
```rust
pub mod service;
```

**features/settings/service.rs:**
```rust
use std::path::PathBuf;
use std::sync::Arc;
use tokio::sync::RwLock;
use crate::{Result};
use crate::infra::config::{self, AppearanceConfig, Config, ConnectionConfig};
use crate::infra::crypto::CryptoContext;

pub struct SettingsService {
    config_path: PathBuf,
    config:      Arc<RwLock<Config>>,
}

impl SettingsService {
    pub fn new(config_path: PathBuf, config: Config) -> (Self, Arc<RwLock<Config>>) {
        let shared = Arc::new(RwLock::new(config));
        (Self { config_path, config: shared.clone() }, shared)
    }

    pub async fn get(&self) -> Config {
        self.config.read().await.clone()
    }

    pub async fn update_appearance(
        &self, theme: String, language: String, accent: String,
    ) -> Result<()> {
        let mut cfg = self.config.write().await;
        cfg.appearance = AppearanceConfig { theme, language, accent_color: accent };
        config::save(&self.config_path, &cfg)
    }

    pub async fn update_connection(&self, input: ConnectionConfig) -> Result<()> {
        let mut cfg = self.config.write().await;
        cfg.connection = input;
        config::save(&self.config_path, &cfg)
    }

    pub async fn set_master_password(
        &self,
        crypto: &CryptoContext,
        old_password: Option<&str>,
        new_password: &str,
    ) -> Result<()> {
        let mut cfg = self.config.write().await;
        if cfg.security.master_password_enabled {
            match old_password {
                Some(old) if crypto.verify_password(old, &cfg.security.master_password_hash) => {}
                Some(_)  => return Err(crate::AppError::InvalidInput("old password incorrect".into())),
                None     => return Err(crate::AppError::InvalidInput("old password required".into())),
            }
        }
        cfg.security.master_password_hash    = CryptoContext::hash_password(new_password)?;
        cfg.security.master_password_enabled = true;
        config::save(&self.config_path, &cfg)
    }
}
```

- [ ] **Step 2: Create console feature**

**features/console/mod.rs:**
```rust
pub mod service;
```

**features/console/service.rs:**
```rust
use std::sync::Arc;
use std::collections::VecDeque;
use dashmap::DashMap;
use tokio::task::AbortHandle;
use crate::Result;
use crate::infra::audit::AuditLog;
use crate::infra::http::HttpClientPool;
use super::super::webshell::service::WebshellService;

pub struct TerminalState {
    pub cwd:     String,
    pub history: VecDeque<String>,
}

pub struct RemoteDbState {
    pub db_type:    String,
    pub conn_str:   String,
    pub last_query: Option<String>,
}

pub struct ConsoleService {
    webshell_service:   Arc<WebshellService>,
    http_pool:          Arc<HttpClientPool>,
    pub audit_log:      AuditLog,
    file_transfers:     DashMap<String, DashMap<String, AbortHandle>>,
    terminal_sessions:  DashMap<String, TerminalState>,
    remote_db_sessions: DashMap<String, RemoteDbState>,
}

impl ConsoleService {
    pub fn new(
        webshell_service: Arc<WebshellService>,
        http_pool:        Arc<HttpClientPool>,
        audit_log:        AuditLog,
    ) -> Self {
        Self {
            webshell_service, http_pool, audit_log,
            file_transfers:     DashMap::new(),
            terminal_sessions:  DashMap::new(),
            remote_db_sessions: DashMap::new(),
        }
    }

    pub async fn cleanup(&self, webshell_id: &str) {
        if let Some((_, transfers)) = self.file_transfers.remove(webshell_id) {
            for entry in transfers.iter() {
                entry.value().abort();
            }
        }
        self.terminal_sessions.remove(webshell_id);
        self.remote_db_sessions.remove(webshell_id);
        tracing::info!("console cleanup done: {}", webshell_id);
    }

    pub fn active_webshell_ids(&self) -> impl Iterator<Item = String> + '_ {
        let from_transfers = self.file_transfers.iter().map(|e| e.key().clone());
        let from_terminals = self.terminal_sessions.iter().map(|e| e.key().clone());
        let from_db        = self.remote_db_sessions.iter().map(|e| e.key().clone());
        let mut seen = std::collections::HashSet::new();
        from_transfers.chain(from_terminals).chain(from_db)
            .filter(move |id| seen.insert(id.clone()))
    }

    pub async fn exec_command(&self, _id: &str, _cmd: &str) -> Result<String> {
        todo!("exec_command — implement in Phase 2")
    }

    pub async fn list_files(&self, _id: &str, _path: &str) -> Result<serde_json::Value> {
        todo!("list_files — implement in Phase 2")
    }
}
```

- [ ] **Step 3: Run cargo check**

```bash
cargo check --manifest-path src-tauri/Cargo.toml
```

Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add src-tauri/src/features/settings/ src-tauri/src/features/console/
git commit -m "feat: add settings service (config read/write) and console service (session state)"
```

---

## Task 14: state.rs — AppState + Bootstrap

**Files:**
- Create: `src-tauri/src/state.rs`

- [ ] **Step 1: Write state.rs**

```rust
// src-tauri/src/state.rs
use std::sync::{Arc, atomic::{AtomicBool, Ordering}};
use std::path::PathBuf;
use tauri::AppHandle;
use tracing_appender::non_blocking::WorkerGuard;
use crate::{AppError, Result};
use crate::infra::{
    audit::AuditLog,
    config::{self, Config},
    crypto::CryptoContext,
    db::Database,
    http::HttpClientPool,
    logger,
    paths::AppPaths,
};
use crate::features::{
    console::service::ConsoleService,
    payload::service::PayloadService,
    plugin::service::PluginService,
    project::{repo::DbProjectRepo, service::ProjectService},
    settings::service::SettingsService,
    webshell::{
        queue::WebshellQueue,
        repo::DbWebshellRepo,
        service::WebshellService,
    },
    payload::repo::DbPayloadRepo,
    plugin::repo::DbPluginRepo,
};

pub struct BatchService {
    webshell_service: Arc<WebshellService>,
}

impl BatchService {
    pub fn new(ws: Arc<WebshellService>) -> Self { Self { webshell_service: ws } }

    pub async fn test_connections(&self, ids: Vec<String>) -> Vec<(String, crate::features::webshell::models::ConnectionResult)> {
        // TODO: concurrent test_connection for each id in Phase 2
        let _ = ids;
        todo!("BatchService::test_connections — implement in Phase 2")
    }
}

pub struct AppState {
    pub paths:            AppPaths,
    pub webshell_service: Arc<WebshellService>,
    pub project_service:  ProjectService,
    pub payload_service:  PayloadService,
    pub console_service:  ConsoleService,
    pub batch_service:    BatchService,
    pub plugin_service:   PluginService,
    pub settings_service: SettingsService,
    pub audit_log:        AuditLog,
    pub is_locked:        AtomicBool,
    pub _log_guard:       WorkerGuard,
}

impl AppState {
    pub fn new(
        paths:      AppPaths,
        config:     Config,
        master_key: [u8; 32],
        db:         Database,
        log_guard:  WorkerGuard,
    ) -> Result<Self> {
        let crypto    = CryptoContext::new(master_key);
        let http_pool = Arc::new(HttpClientPool::new(&config.connection));
        let audit_log = AuditLog::new(db.clone());
        let queue     = Arc::new(WebshellQueue::new(http_pool.clone()));
        let is_locked = config.security.master_password_enabled;

        let webshell_repo = Arc::new(DbWebshellRepo::new(db.clone()));
        let project_repo  = Arc::new(DbProjectRepo::new(db.clone()));
        let payload_repo  = Arc::new(DbPayloadRepo::new(db.clone()));
        let plugin_repo   = Arc::new(DbPluginRepo::new(db.clone()));

        let webshell_service = Arc::new(WebshellService::new(webshell_repo, crypto.clone(), queue.clone()));
        let console_service  = ConsoleService::new(webshell_service.clone(), http_pool.clone(), audit_log.clone());
        let plugins_dir      = paths.plugins_dir.clone();
        let (settings_service, _shared_config) = SettingsService::new(paths.config.clone(), config);

        Ok(Self {
            paths,
            batch_service:   BatchService::new(webshell_service.clone()),
            webshell_service,
            project_service: ProjectService::new(project_repo),
            payload_service: PayloadService::new(payload_repo, crypto.clone()),
            console_service,
            plugin_service:  PluginService::new(plugin_repo, plugins_dir),
            settings_service,
            audit_log,
            is_locked:       AtomicBool::new(is_locked),
            _log_guard:      log_guard,
        })
    }

    pub async fn shutdown(&self) {
        tracing::info!("graceful shutdown started");
        self.webshell_service.queue.shutdown_all();
        let ids: Vec<String> = self.console_service.active_webshell_ids().collect();
        for id in ids {
            self.console_service.cleanup(&id).await;
        }
        tracing::info!("graceful shutdown complete");
    }
}

pub fn check_locked(state: &AppState) -> Result<()> {
    if state.is_locked.load(Ordering::Relaxed) {
        Err(AppError::Locked)
    } else {
        Ok(())
    }
}

pub async fn bootstrap(app: &AppHandle) -> Result<AppState> {
    let paths = AppPaths::resolve(app)?;
    std::fs::create_dir_all(&paths.logs_dir)?;
    std::fs::create_dir_all(&paths.plugins_dir)?;
    std::fs::create_dir_all(&paths.exports_dir)?;

    let log_guard  = logger::init(&paths.logs_dir)?;

    // Register panic hook after logger is ready
    let logs_dir_c = paths.logs_dir.clone();
    std::panic::set_hook(Box::new(move |info| {
        let msg = info.to_string();
        tracing::error!("PANIC: {}", msg);
        let path = logs_dir_c.join(format!("crash-{}.txt", chrono::Utc::now().timestamp()));
        let _ = std::fs::write(path, &msg);
    }));

    let mut cfg = config::load_or_default(&paths.config)?;

    if cfg.security.salt.is_empty() {
        cfg.security.salt = crate::infra::crypto::generate_salt();
        config::save(&paths.config, &cfg)?;
    }

    let master_key = crate::infra::crypto::derive_key(&cfg.security.salt)?;
    let db         = Database::open(&paths.db_path).await?;
    db.migrate().await?;

    let db_bg = db.clone();
    tokio::spawn(async move {
        if let Err(e) = db_bg.vacuum_if_needed().await {
            tracing::warn!("VACUUM failed: {}", e);
        }
    });

    AppState::new(paths, cfg, master_key, db, log_guard)
}
```

- [ ] **Step 2: Add state module to lib.rs**

```rust
// src-tauri/src/lib.rs (append)
pub mod state;
```

- [ ] **Step 3: Run cargo check**

```bash
cargo check --manifest-path src-tauri/Cargo.toml
```

Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add src-tauri/src/state.rs src-tauri/src/lib.rs
git commit -m "feat: add AppState, bootstrap sequence, graceful shutdown"
```

---

## Task 15: commands/ Layer — All Thin Command Handlers

**Files:**
- Create: `src-tauri/src/commands/mod.rs`
- Create: `src-tauri/src/commands/webshell.rs`
- Create: `src-tauri/src/commands/project.rs`
- Create: `src-tauri/src/commands/payload.rs`
- Create: `src-tauri/src/commands/console.rs`
- Create: `src-tauri/src/commands/plugin.rs`
- Create: `src-tauri/src/commands/settings.rs`
- Create: `src-tauri/src/commands/system.rs`
- Create: `src-tauri/src/commands/batch.rs`

- [ ] **Step 1: Create commands/mod.rs**

```rust
// src-tauri/src/commands/mod.rs
pub mod batch;
pub mod console;
pub mod payload;
pub mod plugin;
pub mod project;
pub mod settings;
pub mod system;
pub mod webshell;
```

- [ ] **Step 2: Create commands/webshell.rs**

```rust
// src-tauri/src/commands/webshell.rs
use tauri::State;
use crate::Result;
use crate::state::{AppState, check_locked};
use crate::features::webshell::models::{
    Webshell, CreateWebshellInput, UpdateWebshellInput, ConnectionResult,
};

#[tauri_specta::command]
pub async fn list_webshells(
    state:      State<'_, AppState>,
    project_id: Option<String>,
) -> Result<Vec<Webshell>> {
    state.webshell_service.list(project_id.as_deref()).await
}

#[tauri_specta::command]
pub async fn get_webshell(
    state: State<'_, AppState>,
    id:    String,
) -> Result<Webshell> {
    state.webshell_service.get(&id).await
}

#[tauri_specta::command]
pub async fn create_webshell(
    state: State<'_, AppState>,
    input: CreateWebshellInput,
) -> Result<Webshell> {
    check_locked(&state)?;
    state.webshell_service.create(input).await
}

#[tauri_specta::command]
pub async fn update_webshell(
    state: State<'_, AppState>,
    id:    String,
    input: UpdateWebshellInput,
) -> Result<Webshell> {
    check_locked(&state)?;
    state.webshell_service.update(&id, input).await
}

#[tauri_specta::command]
pub async fn delete_webshell(
    state: State<'_, AppState>,
    id:    String,
) -> Result<()> {
    check_locked(&state)?;
    state.webshell_service.delete(&id).await
}

#[tauri_specta::command]
pub async fn test_connection(
    state: State<'_, AppState>,
    id:    String,
) -> Result<ConnectionResult> {
    check_locked(&state)?;
    state.webshell_service.test_connection(&id).await
}

#[tauri_specta::command]
pub async fn reset_redeploy_status(
    state: State<'_, AppState>,
    id:    String,
) -> Result<Webshell> {
    check_locked(&state)?;
    state.webshell_service.reset_redeploy_status(&id).await
}
```

- [ ] **Step 3: Create commands/project.rs**

```rust
// src-tauri/src/commands/project.rs
use tauri::State;
use crate::Result;
use crate::state::{AppState, check_locked};
use crate::features::project::models::{Project, CreateProjectInput, UpdateProjectInput};

#[tauri_specta::command]
pub async fn list_projects(state: State<'_, AppState>) -> Result<Vec<Project>> {
    state.project_service.list().await
}

#[tauri_specta::command]
pub async fn create_project(
    state: State<'_, AppState>,
    input: CreateProjectInput,
) -> Result<Project> {
    check_locked(&state)?;
    state.project_service.create(input).await
}

#[tauri_specta::command]
pub async fn update_project(
    state: State<'_, AppState>,
    id:    String,
    input: UpdateProjectInput,
) -> Result<Project> {
    check_locked(&state)?;
    state.project_service.update(&id, input).await
}

#[tauri_specta::command]
pub async fn delete_project(
    state: State<'_, AppState>,
    id:    String,
) -> Result<()> {
    check_locked(&state)?;
    state.project_service.delete(&id).await
}
```

- [ ] **Step 4: Create commands/payload.rs**

```rust
// src-tauri/src/commands/payload.rs
use tauri::State;
use crate::Result;
use crate::state::{AppState, check_locked};
use crate::features::payload::models::{Payload, CreatePayloadInput, PayloadConfig, PayloadHistoryEntry};

#[tauri_specta::command]
pub async fn list_payloads(state: State<'_, AppState>) -> Result<Vec<Payload>> {
    state.payload_service.list().await
}

#[tauri_specta::command]
pub async fn create_payload(
    state: State<'_, AppState>,
    input: CreatePayloadInput,
) -> Result<Payload> {
    check_locked(&state)?;
    state.payload_service.create(input).await
}

#[tauri_specta::command]
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

#[tauri_specta::command]
pub async fn list_payload_history(
    state:      State<'_, AppState>,
    payload_id: String,
) -> Result<Vec<PayloadHistoryEntry>> {
    state.payload_service.list_history(&payload_id).await
}
```

- [ ] **Step 5: Create commands/console.rs**

```rust
// src-tauri/src/commands/console.rs
use tauri::{AppHandle, State};
use crate::Result;
use crate::state::{AppState, check_locked};

#[tauri_specta::command]
pub async fn open_console(
    app:         AppHandle,
    state:       State<'_, AppState>,
    webshell_id: String,
) -> Result<()> {
    check_locked(&state)?;
    let label = format!("console-{}", webshell_id);

    if let Some(win) = app.get_webview_window(&label) {
        win.set_focus().map_err(|e| crate::AppError::Io(std::io::Error::new(
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

#[tauri_specta::command]
pub async fn exec_command(
    state:       State<'_, AppState>,
    webshell_id: String,
    cmd:         String,
) -> Result<String> {
    check_locked(&state)?;
    state.console_service.exec_command(&webshell_id, &cmd).await
}

#[tauri_specta::command]
pub async fn list_files(
    state:       State<'_, AppState>,
    webshell_id: String,
    path:        String,
) -> Result<serde_json::Value> {
    check_locked(&state)?;
    state.console_service.list_files(&webshell_id, &path).await
}

#[tauri_specta::command]
pub async fn download_file(
    _app:        AppHandle,
    state:       State<'_, AppState>,
    webshell_id: String,
    path:        String,
    transfer_id: String,
) -> Result<()> {
    check_locked(&state)?;
    // TODO: stream download → emit file-progress events in Phase 2
    let _ = (webshell_id, path, transfer_id);
    todo!("download_file stream — implement in Phase 2")
}

#[tauri_specta::command]
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

#[tauri_specta::command]
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

#[tauri_specta::command]
pub async fn execute_query(
    state:       State<'_, AppState>,
    webshell_id: String,
    query:       String,
) -> Result<serde_json::Value> {
    check_locked(&state)?;
    let _ = (webshell_id, query);
    todo!("execute_query — implement in Phase 2")
}
```

- [ ] **Step 6: Create commands/plugin.rs, settings.rs, system.rs, batch.rs**

**commands/plugin.rs:**
```rust
use tauri::State;
use crate::Result;
use crate::state::{AppState, check_locked};
use crate::features::plugin::models::Plugin;

#[tauri_specta::command]
pub async fn list_plugins(state: State<'_, AppState>) -> Result<Vec<Plugin>> {
    state.plugin_service.list().await
}

#[tauri_specta::command]
pub async fn enable_plugin(state: State<'_, AppState>, id: String) -> Result<Plugin> {
    check_locked(&state)?;
    state.plugin_service.enable(&id).await
}

#[tauri_specta::command]
pub async fn disable_plugin(state: State<'_, AppState>, id: String) -> Result<Plugin> {
    check_locked(&state)?;
    state.plugin_service.disable(&id).await
}
```

**commands/settings.rs:**
```rust
use tauri::State;
use serde::{Deserialize, Serialize};
use crate::Result;
use crate::state::{AppState, check_locked};
use crate::infra::config::Config;

#[tauri_specta::command]
pub async fn get_settings(state: State<'_, AppState>) -> Result<Config> {
    Ok(state.settings_service.get().await)
}

#[derive(Debug, Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct UpdateSettingsInput {
    pub theme:       Option<String>,
    pub language:    Option<String>,
    pub accent_color: Option<String>,
}

#[tauri_specta::command]
pub async fn update_settings(
    state: State<'_, AppState>,
    input: UpdateSettingsInput,
) -> Result<Config> {
    check_locked(&state)?;
    let current = state.settings_service.get().await;
    let theme   = input.theme.unwrap_or(current.appearance.theme);
    let lang    = input.language.unwrap_or(current.appearance.language);
    let accent  = input.accent_color.unwrap_or(current.appearance.accent_color);
    state.settings_service.update_appearance(theme, lang, accent).await?;
    Ok(state.settings_service.get().await)
}

#[tauri_specta::command]
pub async fn unlock(
    state:    State<'_, AppState>,
    password: String,
) -> Result<bool> {
    let cfg = state.settings_service.get().await;
    if !cfg.security.master_password_enabled {
        return Ok(true);
    }
    // We can't call crypto through settings_service directly; use a local verify
    // The master_key is already derived from salt; we verify the stored argon2 hash
    use argon2::{Argon2, PasswordHash, PasswordVerifier};
    let hash = cfg.security.master_password_hash;
    let ok = PasswordHash::new(&hash)
        .map(|h| Argon2::default().verify_password(password.as_bytes(), &h).is_ok())
        .unwrap_or(false);
    if ok {
        state.is_locked.store(false, std::sync::atomic::Ordering::Relaxed);
    }
    Ok(ok)
}

#[tauri_specta::command]
pub async fn set_master_password(
    state:        State<'_, AppState>,
    old_password: Option<String>,
    new_password: String,
) -> Result<()> {
    check_locked(&state)?;
    // Reuse the crypto context from webshell_service (they share the same Arc<CryptoContext>)
    // We expose set_master_password via settings_service using a standalone verify
    use argon2::{Argon2, PasswordHash, PasswordHasher, PasswordVerifier};
    use argon2::password_hash::{rand_core::OsRng, SaltString};

    let mut cfg = state.settings_service.get().await;
    if cfg.security.master_password_enabled {
        match old_password.as_deref() {
            Some(old) => {
                let hash_str = &cfg.security.master_password_hash;
                let ok = PasswordHash::new(hash_str)
                    .map(|h| Argon2::default().verify_password(old.as_bytes(), &h).is_ok())
                    .unwrap_or(false);
                if !ok {
                    return Err(crate::AppError::InvalidInput("old password incorrect".into()));
                }
            }
            None => return Err(crate::AppError::InvalidInput("old password required".into())),
        }
    }
    let salt = SaltString::generate(&mut OsRng);
    let new_hash = Argon2::default()
        .hash_password(new_password.as_bytes(), &salt)
        .map(|h| h.to_string())
        .map_err(|e| crate::AppError::Crypto(e.to_string()))?;

    state.settings_service.update_appearance(
        cfg.appearance.theme.clone(),
        cfg.appearance.language.clone(),
        cfg.appearance.accent_color.clone(),
    ).await?;
    // Save security changes via internal access
    // NOTE: settings_service does not expose a generic update_security method in skeleton.
    // Phase 2 will add update_security(). For now update the config directly.
    let path = state.paths.config.clone();
    let mut full_cfg = state.settings_service.get().await;
    full_cfg.security.master_password_hash    = new_hash;
    full_cfg.security.master_password_enabled = true;
    crate::infra::config::save(&path, &full_cfg)?;
    state.is_locked.store(true, std::sync::atomic::Ordering::Relaxed);
    Ok(())
}
```

**commands/system.rs:**
```rust
use tauri::State;
use serde::{Deserialize, Serialize};
use crate::Result;
use crate::state::AppState;

#[derive(Debug, Serialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct AppInfo {
    pub version:    String,
    pub data_dir:   String,
    pub db_path:    String,
    pub logs_dir:   String,
}

#[tauri_specta::command]
pub async fn get_app_info(state: State<'_, AppState>) -> Result<AppInfo> {
    Ok(AppInfo {
        version:  env!("CARGO_PKG_VERSION").to_string(),
        data_dir: state.paths.data_dir.to_string_lossy().into(),
        db_path:  state.paths.db_path.to_string_lossy().into(),
        logs_dir: state.paths.logs_dir.to_string_lossy().into(),
    })
}

#[derive(Debug, Serialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct AuditEntry {
    pub id:          String,
    pub webshell_id: Option<String>,
    pub action:      String,
    pub detail:      String,
    pub created_at:  i64,
}

#[tauri_specta::command]
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
```

**commands/batch.rs:**
```rust
use tauri::State;
use crate::Result;
use crate::state::{AppState, check_locked};
use crate::features::webshell::models::ConnectionResult;

#[derive(Debug, serde::Serialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct BatchTestResult {
    pub id:     String,
    pub result: ConnectionResult,
}

#[tauri_specta::command]
pub async fn test_connections(
    state: State<'_, AppState>,
    ids:   Vec<String>,
) -> Result<Vec<BatchTestResult>> {
    check_locked(&state)?;
    // TODO: implement concurrent batch testing in Phase 2
    let _ = ids;
    todo!("test_connections — implement in Phase 2")
}
```

- [ ] **Step 7: Add commands module to lib.rs**

```rust
// src-tauri/src/lib.rs (append)
pub mod commands;
pub mod features;
```

- [ ] **Step 8: Run cargo check**

```bash
cargo check --manifest-path src-tauri/Cargo.toml
```

Expected: no errors (warnings about `todo!()` branches are OK).

- [ ] **Step 9: Commit**

```bash
git add src-tauri/src/commands/
git commit -m "feat: add all command handlers (thin wrappers, protocol commands as todo!)"
```

---

## Task 16: main.rs + lib.rs — Tauri Builder + tauri-specta

**Files:**
- Modify: `src-tauri/src/main.rs`
- Modify: `src-tauri/src/lib.rs`

- [ ] **Step 1: Rewrite lib.rs with run() using specta builder**

```rust
// src-tauri/src/lib.rs
pub mod commands;
pub mod error;
pub mod features;
pub mod infra;
pub mod state;

pub use error::{AppError, Result};

use state::bootstrap;
use tauri_specta::{Builder, collect_commands};

pub fn run() {
    let specta_builder = Builder::<tauri::Wry>::new()
        .commands(collect_commands![
            commands::webshell::list_webshells,
            commands::webshell::get_webshell,
            commands::webshell::create_webshell,
            commands::webshell::update_webshell,
            commands::webshell::delete_webshell,
            commands::webshell::test_connection,
            commands::webshell::reset_redeploy_status,
            commands::project::list_projects,
            commands::project::create_project,
            commands::project::update_project,
            commands::project::delete_project,
            commands::payload::list_payloads,
            commands::payload::create_payload,
            commands::payload::generate_payload,
            commands::payload::list_payload_history,
            commands::console::open_console,
            commands::console::exec_command,
            commands::console::list_files,
            commands::console::download_file,
            commands::console::upload_file,
            commands::console::connect_database,
            commands::console::execute_query,
            commands::plugin::list_plugins,
            commands::plugin::enable_plugin,
            commands::plugin::disable_plugin,
            commands::settings::get_settings,
            commands::settings::update_settings,
            commands::settings::unlock,
            commands::settings::set_master_password,
            commands::system::get_app_info,
            commands::system::get_audit_log,
            commands::batch::test_connections,
        ]);

    #[cfg(debug_assertions)]
    specta_builder
        .export(
            specta_typescript::Typescript::default(),
            "../src/bindings.ts",
        )
        .expect("failed to export TypeScript bindings");

    tauri::Builder::default()
        .setup(|app| {
            let state = tauri::async_runtime::block_on(bootstrap(app.handle()))
                .map_err(|e| {
                    eprintln!("Fatal: bootstrap failed: {e}");
                    Box::new(e) as Box<dyn std::error::Error>
                })?;
            app.manage(state);

            app.on_window_event(|window, event| {
                if let tauri::WindowEvent::Destroyed = event {
                    if let Some(id) = window.label().strip_prefix("console-") {
                        let webshell_id = id.to_string();
                        let handle = window.app_handle().clone();
                        tauri::async_runtime::spawn(async move {
                            let state: tauri::State<state::AppState> = handle.state();
                            state.console_service.cleanup(&webshell_id).await;
                        });
                    }
                }
            });

            // Register close-requested handler for graceful shutdown
            let handle = app.handle().clone();
            if let Some(main_win) = app.get_webview_window("main") {
                main_win.on_window_event(move |event| {
                    if let tauri::WindowEvent::CloseRequested { .. } = event {
                        let state: tauri::State<state::AppState> = handle.state();
                        tauri::async_runtime::block_on(state.shutdown());
                    }
                });
            }

            Ok(())
        })
        .invoke_handler(specta_builder.invoke_handler())
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
```

- [ ] **Step 2: Run cargo build (debug) to verify full compilation**

```bash
cargo build --manifest-path src-tauri/Cargo.toml
```

Expected: compiles successfully. If there are compile errors (not warnings), fix them before proceeding.

Common issues to fix:
- Type mismatches in command signatures — ensure `specta::Type` is derived on all types used in commands.
- `AuditLog::db` field needs `pub` visibility for `system.rs` access — either make it `pub` in `audit.rs` or add an `AuditLog::query_recent` method.

**Fix for audit log access in system.rs** — add a query method to AuditLog instead of accessing `db` directly:

In `infra/audit.rs`, add:
```rust
use crate::infra::db::Database;

impl AuditLog {
    pub fn db(&self) -> &Database { &self.db }
}
```

Then in `commands/system.rs` use `state.audit_log.db()` instead of `state.audit_log.db`.

- [ ] **Step 3: Commit**

```bash
git add src-tauri/src/lib.rs src-tauri/src/main.rs
git commit -m "feat: wire up tauri-specta builder with all commands and bootstrap"
```

---

## Task 17: Tauri Config + Capabilities

**Files:**
- Modify: `src-tauri/tauri.conf.json`
- Modify: `src-tauri/capabilities/default.json`
- Create: `src-tauri/capabilities/console.json`

- [ ] **Step 1: Update tauri.conf.json**

```json
{
  "$schema": "https://schema.tauri.app/config/2",
  "productName": "fg-abyss",
  "version": "0.1.0",
  "identifier": "com.fg-abyss.app",
  "build": {
    "beforeDevCommand": "pnpm dev",
    "devUrl": "http://localhost:1420",
    "beforeBuildCommand": "pnpm build",
    "frontendDist": "../dist"
  },
  "app": {
    "windows": [
      {
        "label": "main",
        "title": "FG-ABYSS",
        "width": 1400,
        "height": 900,
        "minWidth": 1024,
        "minHeight": 768,
        "decorations": false,
        "resizable": true,
        "center": true
      }
    ],
    "security": {
      "csp": "default-src 'self'; style-src 'self' 'unsafe-inline'; frame-src 'self'"
    }
  },
  "bundle": {
    "active": true,
    "targets": "all",
    "icon": []
  }
}
```

- [ ] **Step 2: Update capabilities/default.json**

```json
{
  "$schema": "../gen/schemas/desktop-schema.json",
  "identifier": "default",
  "description": "Main window capabilities",
  "windows": ["main"],
  "permissions": [
    "core:default",
    "core:window:allow-minimize",
    "core:window:allow-toggle-maximize",
    "core:window:allow-close",
    "core:window:allow-create",
    "core:window:allow-set-focus",
    "http:default"
  ]
}
```

- [ ] **Step 3: Create capabilities/console.json**

```json
{
  "$schema": "../gen/schemas/desktop-schema.json",
  "identifier": "console",
  "description": "Console window capabilities",
  "windows": ["console-*"],
  "permissions": [
    "core:default",
    "http:default"
  ]
}
```

- [ ] **Step 4: Run pnpm tauri dev to verify full dev build**

```bash
pnpm tauri dev
```

Expected: application launches, no runtime panics. Console shows bootstrap log lines. All existing frontend views still work (theme, sidebar, settings).

If `tauri dev` errors on missing TypeScript bindings at `src/bindings.ts`, that's expected on first run — tauri-specta generates it. Re-run after a successful build.

- [ ] **Step 5: Commit**

```bash
git add src-tauri/tauri.conf.json src-tauri/capabilities/
git commit -m "feat: update tauri config with main window label, CSP, and console window capability"
```

---

## Self-Review

### Spec Coverage Check

| Spec Section | Covered? | Task |
|---|---|---|
| Feature-based layered architecture | ✅ | All tasks |
| AppState DI root | ✅ | Task 14 |
| commands/ thin, no business logic | ✅ | Task 15 |
| features/ no AppHandle | ✅ | Tasks 9-13 |
| infra/ via traits + mockall | ✅ | Tasks 5-8 |
| AppError + Result<T> | ✅ | Task 2 |
| tokio-rusqlite connection | ✅ | Task 5 |
| WAL mode + foreign keys | ✅ | Task 5 |
| 4 migration files | ✅ | Task 3 |
| vacuum_if_needed at startup | ✅ | Task 5, 14 |
| soft delete (deleted_at IS NULL) | ✅ | Tasks 9-11 |
| UUID in Service layer | ✅ | Tasks 10-12 |
| AES-256-GCM crypto | ✅ | Task 7 |
| Argon2id key derivation | ✅ | Task 7 |
| zeroize on CryptoContext drop | ✅ | Task 7 |
| Sensitive<T> wrapper | ✅ | Task 7 |
| TOML config with migration | ✅ | Task 6 |
| logger.rs (WorkerGuard in AppState) | ✅ | Task 4, 14 |
| Panic hook + crash file | ✅ | Task 14 |
| bootstrap() sequence | ✅ | Task 14 |
| Graceful shutdown (shutdown_all + cleanup) | ✅ | Task 14 |
| WebshellQueue lazy workers + oneshot | ✅ | Task 9 |
| WebshellSession with session_key zeroize | ✅ | Task 9 |
| response_mark extraction | ❌ | Phase 2 (protocol) |
| Phase1 Init / Phase2 Exec | ❌ | Phase 2 (todo!) |
| HttpClientPool + circuit breaker | ✅ | Task 8 |
| C2Profile presets | ✅ | Task 8 |
| AuditLog (fire-and-forget) | ✅ | Task 8 |
| ConsoleService cleanup + active_ids | ✅ | Task 13 |
| SettingsService Arc<RwLock<Config>> | ✅ | Task 13 |
| tauri-specta builder | ✅ | Task 16 |
| open_console (single window per shell) | ✅ | Task 15 |
| Window destroyed → console cleanup | ✅ | Task 16 |
| check_locked guard | ✅ | Tasks 14, 15 |
| tauri.conf.json window label "main" | ✅ | Task 17 |
| capabilities/console.json | ✅ | Task 17 |
| CryptoChain (encode/decode) | ✅ stub | Phase 2 |
| Payload generators | ✅ stub | Phase 2 |
| PHP obfuscator (tree-sitter) | ❌ | Phase 3 |
| PHP disable_functions bypass chain | ❌ | Phase 2 |
| BatchService::test_connections | ✅ stub | Phase 2 |
| Project soft_delete_with_webshells txn | ✅ | Task 11 |
| normalize_url | ✅ | Task 10 |

**Gaps:** Protocol implementation (Phase1/Phase2), payload generators, obfuscators, console streaming (file/terminal/db) are all left as `todo!()` — this is by design for the skeleton phase.

### Placeholder Scan

All `todo!()` calls are explicitly annotated with a phase label (Phase 2 or Phase 3). No ambiguous "TBD" or "implement later" strings. Each `todo!()` has the feature name so grep finds them easily.

### Type Consistency

- `Webshell` returned by `WebshellService::create/get/update` has `specta::Type` ✅
- `CreateWebshellInput` / `UpdateWebshellInput` have `specta::Type` ✅
- `ConnectionResult` has `specta::Type` ✅
- `PayloadConfig` references `ObfuscateConfig` from `generator/mod.rs` — same type used in models.rs ✅
- `AuditEntry` in `system.rs` is a local struct (not from infra/audit.rs) to avoid coupling; this is consistent ✅
- `WebshellUpdateFields` is internal to repo.rs, not exposed to commands layer ✅

---

**Plan complete and saved to `docs/superpowers/plans/2026-04-18-backend-skeleton.md`.**

**Two execution options:**

**1. Subagent-Driven (recommended)** — Dispatch a fresh subagent per task, review between tasks, fast iteration. Use `superpowers:subagent-driven-development`.

**2. Inline Execution** — Execute tasks in this session using `superpowers:executing-plans`, batch execution with checkpoints.

**Which approach?**
