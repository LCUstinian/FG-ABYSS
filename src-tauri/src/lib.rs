// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/

mod types;
mod core;
mod cmd;
mod templates;

use tauri::Manager;

use cmd::payload::{
    generate_payload_cmd,
    get_generated_payloads,
    save_file_cmd,
    export_client_config_cmd,
    clear_payload_history,
    get_payload_templates,
    add_payload_template,
    update_payload_template,
    delete_payload_template,
    get_payload_template,
    AppState,
};

use cmd::audit::{
    add_audit_log,
    get_audit_logs,
    clear_audit_logs,
    AuditLogState,
};

use cmd::plugin::{
    verify_plugin_signature_command,
};

#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .plugin(tauri_plugin_fs::init())
        .plugin(tauri_plugin_dialog::init())
        .setup(|app| {
            // 使用当前工作目录作为数据目录，避免沙箱权限限制
            let data_dir = std::env::current_dir().unwrap_or_else(|_| std::path::PathBuf::from("."));
            
            // 确保数据目录存在
            std::fs::create_dir_all(&data_dir).unwrap_or(());
            
            // 创建并管理 AppState
            let app_state = AppState::new(data_dir.clone());
            app.manage(app_state);
            
            // 初始化审计日志服务
            let db_path = data_dir.join("audit.db");
            let encryption_key = "fg-abyss-audit-key".to_string(); // 实际应用中应该使用更安全的密钥管理
            
            let db_manager = crate::core::database::DatabaseManager::new(&db_path).unwrap_or_else(|e| {
                eprintln!("Failed to initialize database: {}", e);
                panic!("Database initialization failed");
            });
            
            let audit_state = AuditLogState {
                db_manager: std::sync::Arc::new(std::sync::Mutex::new(db_manager)),
                encryption_key,
            };
            app.manage(audit_state);
            
            Ok(())
        })
        .invoke_handler(tauri::generate_handler!(
            greet,
            generate_payload_cmd,
            get_generated_payloads,
            save_file_cmd,
            export_client_config_cmd,
            clear_payload_history,
            get_payload_templates,
            add_payload_template,
            update_payload_template,
            delete_payload_template,
            get_payload_template,
            add_audit_log,
            get_audit_logs,
            clear_audit_logs,
            verify_plugin_signature_command,
        ))
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
