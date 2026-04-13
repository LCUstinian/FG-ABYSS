mod commands;
mod core;
mod infra;
mod plugins;

use std::sync::Arc;
use tokio::sync::Mutex;

pub fn run() {
    // 初始化数据库
    let db_path = infra::database::get_database_path()
        .expect("获取数据库路径失败");
    let conn = infra::database::init_db(&db_path)
        .expect("初始化数据库失败");
    infra::database::run_migrations(&conn)
        .expect("运行迁移失败");
    
    let db = Arc::new(Mutex::new(conn));
    
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_fs::init())
        .plugin(tauri_plugin_http::init())
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_notification::init())
        .manage(db)
        .setup(|_app| {
            // 初始化日志
            core::logger::init();
            
            // 设置系统托盘
            plugins::tray::setup_system_tray(_app.handle())?;
            
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            commands::create_project,
            commands::update_project,
            commands::delete_project,
            commands::restore_project,
            commands::get_all_projects,
            commands::get_deleted_projects,
            commands::create_webshell,
            commands::update_webshell,
            commands::delete_webshell,
            commands::get_webshells_by_project,
            commands::test_connection,
            commands::get_system_info,
            // 加密通信 Commands
            commands::test_webshell_connection,
            commands::send_encrypted_request,
            commands::decrypt_response,
            commands::execute_system_command,
            commands::get_filesystem_info,
            commands::list_directory,
            commands::upload_file,
            commands::download_file,
            commands::generate_encryption_key,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
