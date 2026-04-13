pub mod commands;
pub mod core;
pub mod infra;
pub mod plugins;

pub use commands::*;
pub use core::logger;
pub use core::config;
pub use infra::database;
pub use infra::file_system;
pub use plugins::tray;

pub fn run_app() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_fs::init())
        .plugin(tauri_plugin_http::init())
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_notification::init())
        .setup(|app| {
            logger::init();
            tray::setup_system_tray(app.handle())?;
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            ping,
            execute_command,
            get_system_info,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
