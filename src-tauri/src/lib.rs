pub mod commands;
pub mod error;
pub mod features;
pub mod infra;
pub mod state;

pub use error::{AppError, Result};

use state::bootstrap;
use tauri::Manager;
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
        .on_window_event(|window, event| {
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
        })
        .invoke_handler(specta_builder.invoke_handler())
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
