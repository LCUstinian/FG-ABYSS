use tauri::{command, State};
use serde::{Deserialize, Serialize};
use std::path::PathBuf;
use std::sync::{Arc, Mutex};
use base64::{Engine as _, engine::general_purpose};
use crate::core::database::DatabaseManager;
use crate::core::crypto::{encrypt_data, decrypt_data, EncryptionType};

/// 审计日志配置
#[derive(Debug, Deserialize)]
pub struct AuditLogConfig {
    pub db_path: PathBuf,
    pub encryption_key: String,
}

/// 审计日志状态
pub struct AuditLogState {
    pub db_manager: Arc<Mutex<DatabaseManager>>,
    pub encryption_key: String,
}

/// 添加审计日志请求
#[derive(Debug, Deserialize)]
pub struct AddAuditLogRequest {
    pub action_type: String,
    pub payload_hash: Option<String>,
    pub content: String,
}

/// 获取审计日志列表请求
#[derive(Debug, Deserialize)]
pub struct GetAuditLogsRequest {
    pub limit: i64,
    pub offset: i64,
}

/// 审计日志响应
#[derive(Debug, Serialize)]
pub struct AuditLogResponse {
    pub id: i64,
    pub action_type: String,
    pub payload_hash: Option<String>,
    pub content: String,
    pub created_at: String,
}

/// 审计日志列表响应
#[derive(Debug, Serialize)]
pub struct AuditLogsResponse {
    pub logs: Vec<AuditLogResponse>,
    pub total: i64,
}

/// 添加审计日志
#[command]
pub fn add_audit_log(
    state: State<AuditLogState>,
    request: AddAuditLogRequest,
) -> Result<(), String> {
    // 加密日志内容
    let (ciphertext, nonce, tag) = encrypt_data(
        request.content.as_bytes(),
        &state.encryption_key,
        &EncryptionType::ChaCha20Poly1305,
    )?;

    // 将加密结果转换为Base64
    let encrypted_content = general_purpose::STANDARD.encode(&ciphertext);
    let nonce_str = general_purpose::STANDARD.encode(nonce.unwrap_or_default());
    let tag_str = general_purpose::STANDARD.encode(tag.unwrap_or_default());

    // 添加到数据库
    let db_manager = state.db_manager.lock().map_err(|e| e.to_string())?;
    db_manager.add_audit_log(
        &request.action_type,
        request.payload_hash.as_deref(),
        &encrypted_content,
        &nonce_str,
        &tag_str,
    ).map_err(|e| e.to_string())?;

    Ok(())
}

/// 获取审计日志列表
#[command]
pub fn get_audit_logs(
    state: State<AuditLogState>,
    request: GetAuditLogsRequest,
) -> Result<AuditLogsResponse, String> {
    // 从数据库获取日志
    let db_manager = state.db_manager.lock().map_err(|e| e.to_string())?;
    let logs = db_manager.get_audit_logs(request.limit, request.offset).map_err(|e| e.to_string())?;
    let total = db_manager.get_audit_logs_count().map_err(|e| e.to_string())?;

    // 解密日志内容
    let mut decrypted_logs = Vec::new();
    for log in logs {
        // 解码Base64
        let ciphertext = general_purpose::STANDARD.decode(&log.encrypted_content).map_err(|e| e.to_string())?;
        let nonce = general_purpose::STANDARD.decode(&log.nonce).map_err(|e| e.to_string())?;
        let tag = general_purpose::STANDARD.decode(&log.tag).map_err(|e| e.to_string())?;

        // 解密
        let plaintext = decrypt_data(
            &ciphertext,
            Some(&nonce),
            Some(&tag),
            &state.encryption_key,
            &EncryptionType::ChaCha20Poly1305,
        )?;

        let content = String::from_utf8(plaintext).map_err(|e| e.to_string())?;

        decrypted_logs.push(AuditLogResponse {
            id: log.id,
            action_type: log.action_type,
            payload_hash: log.payload_hash,
            content,
            created_at: log.created_at.to_string(),
        });
    }

    Ok(AuditLogsResponse {
        logs: decrypted_logs,
        total,
    })
}

/// 清空审计日志
#[command]
pub fn clear_audit_logs(
    state: State<AuditLogState>,
) -> Result<(), String> {
    let db_manager = state.db_manager.lock().map_err(|e| e.to_string())?;
    db_manager.clear_audit_logs().map_err(|e| e.to_string())?;
    Ok(())
}
