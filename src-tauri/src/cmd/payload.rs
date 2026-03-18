/// Tauri 命令模块 - 载荷生成相关命令

use crate::core::generator::generate_payload;
use crate::types::payload::*;
use tauri::State;
use std::sync::Mutex;

// 使用标准 Result 类型，避免与 types::payload::Result 冲突
type CommandResult<T> = std::result::Result<T, String>;

/// 应用状态
pub struct AppState {
    pub generated_payloads: Mutex<Vec<PayloadResult>>,
}

/// 生成载荷
#[tauri::command]
pub async fn generate_payload_cmd(
    config: PayloadConfig,
    state: State<'_, AppState>,
) -> CommandResult<PayloadResult> {
    // 调用生成器
    let result = generate_payload(&config)
        .map_err(|e| e.to_string())?;
    
    // 保存到历史记录
    {
        let mut payloads = state.generated_payloads.lock().map_err(|e| e.to_string())?;
        payloads.push(result.clone());
    }
    
    Ok(result)
}

/// 获取已生成的载荷列表
#[tauri::command]
pub async fn get_generated_payloads(
    state: State<'_, AppState>,
) -> CommandResult<Vec<PayloadResult>> {
    let payloads = state.generated_payloads.lock().map_err(|e| e.to_string())?;
    Ok(payloads.clone())
}

/// 保存文件
#[tauri::command]
pub async fn save_file_cmd(
    path: String,
    content: String,
) -> CommandResult<()> {
    use std::fs;
    
    // 使用标准库写入文件
    fs::write(&path, &content)
        .map_err(|e| format!("Failed to write file: {}", e))?;
    
    Ok(())
}

/// 导出客户端配置
#[tauri::command]
pub async fn export_client_config_cmd(
    config: ClientConfig,
    path: String,
) -> CommandResult<()> {
    use std::fs;
    
    let json = serde_json::to_string_pretty(&config)
        .map_err(|e| format!("Failed to serialize config: {}", e))?;
    
    fs::write(&path, &json)
        .map_err(|e| format!("Failed to write config file: {}", e))?;
    
    Ok(())
}
