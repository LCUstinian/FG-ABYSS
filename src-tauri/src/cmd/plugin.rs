use tauri::command;
use serde::{Deserialize, Serialize};
use std::path::PathBuf;
use crate::core::plugin::{verify_plugin_signature, PluginSignatureError};

/// 插件签名验证请求
#[derive(Debug, Deserialize)]
pub struct VerifyPluginSignatureRequest {
    pub plugin_path: PathBuf,
    pub signature_path: PathBuf,
}

/// 插件签名验证响应
#[derive(Debug, Serialize)]
pub struct VerifyPluginSignatureResponse {
    pub valid: bool,
    pub message: String,
}

/// 验证插件签名
#[command]
pub fn verify_plugin_signature_command(
    request: VerifyPluginSignatureRequest,
) -> Result<VerifyPluginSignatureResponse, String> {
    // 这里使用默认的公钥，实际应用中应该从配置中读取
    // 注意：在生产环境中，应该使用安全的方式存储公钥
    let public_key = b"\x30\x2a\x30\x05\x06\x03\x2b\x65\x70\x03\x21\x00\x12\x34\x56\x78\x90\xab\xcd\xef\x12\x34\x56\x78\x90\xab\xcd\xef\x12\x34\x56\x78\x90\xab\xcd\xef";

    match verify_plugin_signature(&request.plugin_path, &request.signature_path, public_key) {
        Ok(_) => Ok(VerifyPluginSignatureResponse {
            valid: true,
            message: "插件签名验证通过".to_string(),
        }),
        Err(error) => {
            let message = match error {
                PluginSignatureError::FileReadError => "无法读取插件或签名文件",
                PluginSignatureError::InvalidSignature => "插件签名无效",
                PluginSignatureError::InvalidPublicKey => "无效的公钥",
            };
            Ok(VerifyPluginSignatureResponse {
                valid: false,
                message: message.to_string(),
            })
        }
    }
}
