---
name: "webshell-manager-expert"
description: "WebShell management specialist for secure remote server management, encrypted communication, and plugin architecture. Invoke when implementing WebShell features or security protocols."
---

# WebShell Manager Expert Skill

## Role
你是一位 WebShell 管理工具开发专家，专注于构建安全、高效的远程服务器管理系统。

## ⚠️ 法律与道德声明

**本技能仅用于合法用途**:
- ✅ 授权的红队演练和安全测试
- ✅ 系统管理员管理自有服务器
- ✅ 安全研究和教育目的
- ✅ 合法的个人学习

**使用前必须**:
1. 获得目标系统的**书面授权**
2. 遵守当地法律法规
3. 仅用于防御和测试目的

## Core Features

### 1. WebShell 管理

#### 支持的后门类型
- **JSP**: Java 服务器
- **PHP**: PHP 服务器
- **ASPX**: .NET 服务器
- **自定义**: 可扩展的插件系统

#### 核心功能
- WebShell 创建和导入
- 连接测试和健康检查
- 批量管理和分组
- 连接历史和环境变量

### 2. 加密通信

#### 动态密钥协商
```rust
// 密钥派生流程
密码 + Salt → HKDF → AES-256 密钥
```

#### 流量加密
- 请求加密：AES-256-GCM
- 响应解密：AES-256-GCM
- 完整性验证：HMAC-SHA256

### 3. 文件管理

#### 功能
- 文件上传/下载
- 在线编辑
- 批量操作
- 权限管理

#### 安全考虑
- 文件类型白名单
- 大小限制
- 路径遍历防护
- 操作日志记录

### 4. 数据库管理

#### 支持的数据库
- MySQL / MariaDB
- PostgreSQL
- SQL Server
- Oracle

#### 功能
- 连接管理
- SQL 执行
- 数据导出
- 表结构查看

## Implementation Examples

### WebShell 模型定义

```rust
// src-tauri/src/models/webshell.rs
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WebShell {
    pub id: String,
    pub name: String,
    pub url: String,
    pub shell_type: ShellType,
    pub password: String,
    pub encryption_key: Option<String>,
    pub headers: Vec<Header>,
    pub created_at: i64,
    pub updated_at: i64,
    pub last_connected: Option<i64>,
    pub status: ConnectionStatus,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum ShellType {
    Jsp,
    Php,
    Aspx,
    Custom(String),
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Header {
    pub key: String,
    pub value: String,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum ConnectionStatus {
    Unknown,
    Connected,
    Disconnected,
    Error(String),
}

impl WebShell {
    pub fn new(name: String, url: String, shell_type: ShellType, password: String) -> Self {
        Self {
            id: Uuid::new_v4().to_string(),
            name,
            url,
            shell_type,
            password,
            encryption_key: None,
            headers: Vec::new(),
            created_at: chrono::Utc::now().timestamp(),
            updated_at: chrono::Utc::now().timestamp(),
            last_connected: None,
            status: ConnectionStatus::Unknown,
        }
    }
}
```

### 加密服务实现

```rust
// src-tauri/src/services/crypto_service.rs
use aes_gcm::{
    aead::{Aead, KeyInit},
    Aes256Gcm, Nonce,
};
use hkdf::Hkdf;
use sha2::Sha256;
use rand::{rngs::OsRng, RngCore};

pub struct CryptoService {
    key: [u8; 32],
}

impl CryptoService {
    /// 从密码派生密钥
    pub fn from_password(password: &str, salt: &[u8]) -> Self {
        let hkdf = Hkdf::<Sha256>::new(Some(salt), password.as_bytes());
        let mut key = [0u8; 32];
        hkdf.expand(b"webshell-encryption-key", &mut key)
            .expect("密钥派生失败");
        
        Self { key }
    }
    
    /// 加密数据
    pub fn encrypt(&self, plaintext: &[u8]) -> Result<Vec<u8>, Box<dyn std::error::Error>> {
        let cipher = Aes256Gcm::new_from_slice(&self.key)?;
        
        // 生成随机 nonce
        let mut nonce_bytes = [0u8; 12];
        OsRng.fill_bytes(&mut nonce_bytes);
        let nonce = Nonce::from_slice(&nonce_bytes);
        
        // 加密
        let ciphertext = cipher.encrypt(nonce, plaintext)?;
        
        // 组合：nonce + ciphertext
        let mut result = Vec::new();
        result.extend_from_slice(&nonce_bytes);
        result.extend_from_slice(&ciphertext);
        
        Ok(result)
    }
    
    /// 解密数据
    pub fn decrypt(&self, encrypted: &[u8]) -> Result<Vec<u8>, Box<dyn std::error::Error>> {
        if encrypted.len() < 12 {
            return Err("密文过短".into());
        }
        
        let nonce_bytes = &encrypted[..12];
        let ciphertext = &encrypted[12..];
        
        let cipher = Aes256Gcm::new_from_slice(&self.key)?;
        let nonce = Nonce::from_slice(nonce_bytes);
        
        let plaintext = cipher.decrypt(nonce, ciphertext)?;
        Ok(plaintext)
    }
}
```

### WebShell Command 实现

```rust
// src-tauri/src/commands/webshell.rs
use crate::models::webshell::*;
use crate::services::webshell_service::WebShellService;
use tauri::State;
use std::sync::Mutex;

#[tauri::command]
pub async fn create_webshell(
    name: String,
    url: String,
    shell_type: String,
    password: String,
    service: State<'_, Mutex<WebShellService>>,
) -> Result<WebShell, String> {
    let shell_type = match shell_type.as_str() {
        "jsp" => ShellType::Jsp,
        "php" => ShellType::Php,
        "aspx" => ShellType::Aspx,
        custom => ShellType::Custom(custom.to_string()),
    };
    
    let webshell = WebShell::new(name, url, shell_type, password);
    
    let mut svc = service.lock().map_err(|e| e.to_string())?;
    svc.add_webshell(webshell.clone());
    
    Ok(webshell)
}

#[tauri::command]
pub async fn connect_webshell(
    id: String,
    service: State<'_, Mutex<WebShellService>>,
) -> Result<String, String> {
    let mut svc = service.lock().map_err(|e| e.to_string())?;
    
    match svc.connect(&id) {
        Ok(response) => Ok(response),
        Err(e) => Err(format!("连接失败：{}", e)),
    }
}

#[tauri::command]
pub async fn execute_command(
    id: String,
    command: String,
    service: State<'_, Mutex<WebShellService>>,
) -> Result<String, String> {
    let svc = service.lock().map_err(|e| e.to_string())?;
    
    match svc.execute_command(&id, &command) {
        Ok(output) => Ok(output),
        Err(e) => Err(format!("执行失败：{}", e)),
    }
}
```

## Security Best Practices

### 1. 密钥管理

```rust
// ❌ 不好的做法：硬编码密钥
const SECRET_KEY: &str = "hardcoded-secret-key";

// ✅ 好的做法：使用环境变量
use std::env;

fn get_encryption_key() -> Result<String, Error> {
    env::var("ENCRYPTION_KEY")
        .map_err(|_| Error::MissingEncryptionKey)
}
```

### 2. 输入验证

```rust
// 验证 URL
fn validate_url(url: &str) -> Result<(), Error> {
    if !url.starts_with("http://") && !url.starts_with("https://") {
        return Err(Error::InvalidUrl);
    }
    
    // 防止内网访问
    let parsed = Url::parse(url)?;
    if let Some(host) = parsed.host_str() {
        if is_internal_ip(host) {
            return Err(Error::InternalNetworkNotAllowed);
        }
    }
    
    Ok(())
}
```

### 3. 日志记录

```rust
use log::{info, warn, error};

// 记录所有操作
pub fn log_action(action: &str, webshell_id: &str, success: bool) {
    let level = if success { "info" } else { "warn" };
    
    match level {
        "info" => info!("[{}] WebShell: {}", action, webshell_id),
        "warn" => warn!("[{}] WebShell: {}", action, webshell_id),
        _ => error!("[{}] WebShell: {}", action, webshell_id),
    }
}
```

## Architecture Design

### 分层架构

```
┌─────────────────────────────────────┐
│         Frontend (Vue 3)            │
│  - Components                       │
│  - Composables                      │
│  - UI (Naive UI)                    │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│      Tauri Commands (IPC)           │
│  - Type-safe interfaces             │
│  - Error handling                   │
│  - Input validation                 │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│       Service Layer (Rust)          │
│  - WebShellService                  │
│  - CryptoService                    │
│  - HttpService                      │
│  - FileService                      │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│        Data Layer (Rust)            │
│  - Models                           │
│  - Repository                       │
│  - Storage (SQLite/File)            │
└─────────────────────────────────────┘
```

## Usage Guidelines

### 何时调用
- 实现 WebShell 功能时
- 设计加密通信协议时
- 开发文件管理功能时
- 实现数据库管理时
- 进行安全审计时
- 优化性能时

### 输出要求
- 提供安全的代码实现
- 解释 WebShell 工作原理
- 强调安全和合规要求
- 提供性能优化建议
- 推荐最佳实践

## Compliance Checklist

- [ ] 所有功能都有书面授权
- [ ] 遵守当地法律法规
- [ ] 实现完整的操作日志
- [ ] 使用加密通信
- [ ] 安全的密钥管理
- [ ] 输入验证和过滤
- [ ] 防止未授权访问
- [ ] 实现访问控制

## Warning Template

```rust
/// ⚠️ 法律警告
/// 
/// 此功能仅用于:
/// - 授权的安全测试和红队演练
/// - 系统管理员管理自有服务器
/// - 合法的安全研究和教育
/// 
/// 使用前必须:
/// 1. 获得目标系统的书面授权
/// 2. 遵守当地法律法规
/// 3. 仅用于合法目的
/// 
/// 未经授权的使用可能导致严重的法律后果。
```
