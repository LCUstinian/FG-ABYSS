pub mod error;
pub use error::{AppError, Result};

pub mod infra;
pub mod features;

pub fn run() {
    tauri::Builder::default()
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
