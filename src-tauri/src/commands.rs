use tauri::command;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct SystemInfo {
    pub os: String,
    pub arch: String,
    pub hostname: String,
}

#[command]
pub fn ping() -> String {
    "pong".to_string()
}

#[command]
pub async fn execute_command(cmd: String, args: Vec<String>) -> Result<String, String> {
    use std::process::Command;
    
    let output = Command::new(cmd)
        .args(&args)
        .output()
        .map_err(|e| format!("Failed to execute command: {}", e))?;
    
    let result = String::from_utf8_lossy(&output.stdout).to_string();
    Ok(result)
}

#[command]
pub fn get_system_info() -> Result<SystemInfo, String> {
    use sysinfo::System;
    
    let info = SystemInfo {
        os: System::os_version().unwrap_or_else(|| "Unknown".to_string()),
        arch: std::env::consts::ARCH.to_string(),
        hostname: System::host_name().unwrap_or_else(|| "Unknown".to_string()),
    };
    
    Ok(info)
}
