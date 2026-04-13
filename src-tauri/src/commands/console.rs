use crate::core::crypto::{AesCipher, XorCipher};
use crate::core::http::{HttpClient, HttpClientConfig};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ConsoleCommand {
    pub action: String,
    pub params: HashMap<String, String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ConsoleResponse {
    pub success: bool,
    pub data: Option<String>,
    pub error: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EncryptedPayload {
    pub encrypted_data: Vec<u8>,
    pub encryption_type: String,
}

/// 测试 WebShell 连接
#[tauri::command]
pub async fn test_webshell_connection(
    url: String,
    _password: String,
    _payload_type: String,
) -> Result<ConsoleResponse, String> {
    let config = HttpClientConfig::default();
    let client = HttpClient::new(config)
        .map_err(|e| format!("创建 HTTP 客户端失败：{}", e))?;

    // 构建测试请求
    let response = client
        .post(&url, Some(b"test".to_vec()), None)
        .await
        .map_err(|e| format!("HTTP 请求失败：{}", e))?;

    if response.is_success() {
        Ok(ConsoleResponse {
            success: true,
            data: Some("连接成功".to_string()),
            error: None,
        })
    } else {
        Ok(ConsoleResponse {
            success: false,
            data: None,
            error: Some(format!("连接失败，状态码：{}", response.status_code)),
        })
    }
}

/// 发送加密请求
#[tauri::command]
pub async fn send_encrypted_request(
    url: String,
    _password: String,
    _payload_type: String,
    data: Vec<u8>,
    encryption_type: String,
    encryption_key: Option<String>,
) -> Result<ConsoleResponse, String> {
    let config = HttpClientConfig::default();
    let client = HttpClient::new(config)
        .map_err(|e| format!("创建 HTTP 客户端失败：{}", e))?;

    // 加密数据
    let encrypted_data = match encryption_type.as_str() {
        "aes" => {
            let key_bytes = encryption_key
                .unwrap_or_default()
                .as_bytes()
                .try_into()
                .map_err(|_| "无效的 AES 密钥长度")?;
            let cipher = AesCipher::new(key_bytes);
            cipher.encrypt(&data)
                .map_err(|e| format!("AES 加密失败：{}", e))?
        }
        "xor" => {
            let key = encryption_key.unwrap_or_else(|| "default".to_string());
            let cipher = XorCipher::new(key.as_bytes().to_vec());
            cipher.process(&data)
        }
        _ => return Err(format!("不支持的加密类型：{}", encryption_type)),
    };

    // 构建请求
    let mut headers = HashMap::new();
    headers.insert("Content-Type".to_string(), "application/octet-stream".to_string());

    let response = client
        .post(&url, Some(encrypted_data), Some(headers))
        .await
        .map_err(|e| format!("HTTP 请求失败：{}", e))?;

    if response.is_success() {
        use base64::{Engine as _, engine::general_purpose};
        Ok(ConsoleResponse {
            success: true,
            data: Some(general_purpose::STANDARD.encode(&response.body)),
            error: None,
        })
    } else {
        Ok(ConsoleResponse {
            success: false,
            data: None,
            error: Some(format!("请求失败，状态码：{}", response.status_code)),
        })
    }
}

/// 解密响应数据
#[tauri::command]
pub async fn decrypt_response(
    encrypted_data: Vec<u8>,
    encryption_type: String,
    encryption_key: Option<String>,
) -> Result<ConsoleResponse, String> {
    let decrypted_data = match encryption_type.as_str() {
        "aes" => {
            let key_bytes = encryption_key
                .unwrap_or_default()
                .as_bytes()
                .try_into()
                .map_err(|_| "无效的 AES 密钥长度")?;
            let cipher = AesCipher::new(key_bytes);
            cipher.decrypt(&encrypted_data)
                .map_err(|e| format!("AES 解密失败：{}", e))?
        }
        "xor" => {
            let key = encryption_key.unwrap_or_else(|| "default".to_string());
            let cipher = XorCipher::new(key.as_bytes().to_vec());
            cipher.process(&encrypted_data)
        }
        _ => return Err(format!("不支持的加密类型：{}", encryption_type)),
    };

    let text = String::from_utf8(decrypted_data)
        .map_err(|e| format!("UTF-8 解码失败：{}", e))?;

    Ok(ConsoleResponse {
        success: true,
        data: Some(text),
        error: None,
    })
}

/// 执行系统命令
#[tauri::command]
pub async fn execute_system_command(
    _webshell_id: String,
    command: String,
) -> Result<ConsoleResponse, String> {
    // TODO: 从数据库获取 WebShell 信息
    // 这里只是一个示例实现
    Ok(ConsoleResponse {
        success: true,
        data: Some(format!("执行命令：{}", command)),
        error: None,
    })
}

/// 获取文件系统信息
#[tauri::command]
pub async fn get_filesystem_info(
    _webshell_id: String,
    path: Option<String>,
) -> Result<ConsoleResponse, String> {
    let current_path = path.unwrap_or_else(|| "/".to_string());
    
    // TODO: 实现真实的文件系统信息获取
    Ok(ConsoleResponse {
        success: true,
        data: Some(format!("当前路径：{}", current_path)),
        error: None,
    })
}

/// 列出目录内容
#[tauri::command]
pub async fn list_directory(
    _webshell_id: String,
    path: String,
) -> Result<ConsoleResponse, String> {
    // TODO: 实现真实的目录列表功能
    Ok(ConsoleResponse {
        success: true,
        data: Some(format!("目录列表：{}", path)),
        error: None,
    })
}

/// 上传文件
#[tauri::command]
pub async fn upload_file(
    _webshell_id: String,
    local_path: String,
    remote_path: String,
) -> Result<ConsoleResponse, String> {
    // TODO: 实现真实的文件上传功能
    Ok(ConsoleResponse {
        success: true,
        data: Some(format!("上传文件：{} -> {}", local_path, remote_path)),
        error: None,
    })
}

/// 下载文件
#[tauri::command]
pub async fn download_file(
    _webshell_id: String,
    remote_path: String,
    local_path: String,
) -> Result<ConsoleResponse, String> {
    // TODO: 实现真实的文件下载功能
    Ok(ConsoleResponse {
        success: true,
        data: Some(format!("下载文件：{} -> {}", remote_path, local_path)),
        error: None,
    })
}

/// 生成加密密钥
#[tauri::command]
pub fn generate_encryption_key(
    encryption_type: String,
) -> Result<ConsoleResponse, String> {
    match encryption_type.as_str() {
        "aes" => {
            let key = AesCipher::generate_key();
            let key_hex = hex::encode(&key);
            Ok(ConsoleResponse {
                success: true,
                data: Some(key_hex),
                error: None,
            })
        }
        "xor" => {
            use rand::Rng;
            let mut rng = rand::thread_rng();
            let key: Vec<u8> = (0..32).map(|_| rng.gen()).collect();
            let key_hex = hex::encode(&key);
            Ok(ConsoleResponse {
                success: true,
                data: Some(key_hex),
                error: None,
            })
        }
        _ => Err(format!("不支持的加密类型：{}", encryption_type)),
    }
}
