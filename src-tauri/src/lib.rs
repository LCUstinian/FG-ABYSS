mod commands;
mod core;
mod infra;
mod plugins;

pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_fs::init())
        .plugin(tauri_plugin_http::init())
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_notification::init())
        .setup(|app| {
            // 初始化数据库
            infra::init()?;
            
            // 初始化日志
            core::logger::init();
            
            // 设置系统托盘
            plugins::tray::setup_system_tray(app.handle())?;
            
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            commands::ping,
            commands::execute_command,
            commands::get_system_info,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
