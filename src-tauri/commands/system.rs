use serde::{Deserialize, Serialize};
use tauri::command;

#[derive(Debug, Serialize, Deserialize)]
pub struct PingResponse {
    pub message: String,
    pub timestamp: u64,
}

#[command]
pub async fn ping() -> Result<PingResponse, String> {
    Ok(PingResponse {
        message: "pong".to_string(),
        timestamp: std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .unwrap()
            .as_secs(),
    })
}

#[command]
pub async fn execute_command(cmd: String, args: Vec<String>) -> Result<String, String> {
    use tokio::process::Command;
    
    let output = Command::new(&cmd)
        .args(&args)
        .output()
        .await
        .map_err(|e| format!("Failed to execute command: {}", e))?;

    let stdout = String::from_utf8_lossy(&output.stdout);
    let stderr = String::from_utf8_lossy(&output.stderr);

    Ok(format!("{}{}", stdout, stderr))
}

#[command]
pub async fn get_system_info() -> Result<SystemInfo, String> {
    use sysinfo::{System, SystemExt};
    
    let mut sys = System::new_all();
    sys.refresh_all();

    Ok(SystemInfo {
        os_name: std::env::consts::OS.to_string(),
        os_version: System::os_version().unwrap_or_default(),
        arch: std::env::consts::ARCH.to_string(),
        cpu_count: sys.cpus().len(),
        memory_total: sys.total_memory(),
        memory_available: sys.available_memory(),
    })
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SystemInfo {
    pub os_name: String,
    pub os_version: String,
    pub arch: String,
    pub cpu_count: usize,
    pub memory_total: u64,
    pub memory_available: u64,
}
