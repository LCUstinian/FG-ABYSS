use std::collections::{HashSet, VecDeque};
use std::sync::Arc;
use dashmap::DashMap;
use tokio::task::AbortHandle;
use crate::Result;
use crate::infra::audit::AuditLog;
use crate::infra::http::HttpClientPool;
use crate::features::webshell::service::WebshellService;

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
            webshell_service,
            http_pool,
            audit_log,
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
        let from_transfers: Vec<String> = self.file_transfers.iter().map(|e| e.key().clone()).collect();
        let from_terminals: Vec<String> = self.terminal_sessions.iter().map(|e| e.key().clone()).collect();
        let from_db:        Vec<String> = self.remote_db_sessions.iter().map(|e| e.key().clone()).collect();
        let mut seen = HashSet::new();
        from_transfers
            .into_iter()
            .chain(from_terminals)
            .chain(from_db)
            .filter(move |id| seen.insert(id.clone()))
    }

    pub async fn exec_command(&self, _id: &str, _cmd: &str) -> Result<String> {
        todo!("exec_command — implement in Phase 2")
    }

    pub async fn list_files(&self, _id: &str, _path: &str) -> Result<serde_json::Value> {
        todo!("list_files — implement in Phase 2")
    }
}
