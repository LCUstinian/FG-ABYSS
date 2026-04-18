use std::sync::{Arc, atomic::{AtomicBool, Ordering}};
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
    payload::{repo::DbPayloadRepo, service::PayloadService},
    plugin::{repo::DbPluginRepo, service::PluginService},
    project::{repo::DbProjectRepo, service::ProjectService},
    settings::service::SettingsService,
    webshell::{
        queue::WebshellQueue,
        repo::DbWebshellRepo,
        service::WebshellService,
    },
};

pub struct BatchService {
    webshell_service: Arc<WebshellService>,
}

impl BatchService {
    pub fn new(ws: Arc<WebshellService>) -> Self { Self { webshell_service: ws } }

    pub async fn test_connections(&self, ids: Vec<String>) -> Vec<(String, crate::features::webshell::models::ConnectionResult)> {
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

    let log_guard = logger::init(&paths.logs_dir)?;

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
