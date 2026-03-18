// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/

mod types;
mod core;
mod cmd;

use cmd::payload::{
    generate_payload_cmd,
    get_generated_payloads,
    save_file_cmd,
    export_client_config_cmd,
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
        .manage(cmd::payload::AppState {
            generated_payloads: std::sync::Mutex::new(Vec::new()),
        })
        .invoke_handler(tauri::generate_handler![
            greet,
            generate_payload_cmd,
            get_generated_payloads,
            save_file_cmd,
            export_client_config_cmd,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
