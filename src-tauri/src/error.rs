use thiserror::Error;
use serde::{Serialize, Deserialize};

/// 错误严重性级别
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum ErrorSeverity {
    /// 低 - 不影响功能
    Low,
    /// 中 - 部分功能受限
    Medium,
    /// 高 - 功能不可用
    High,
    /// 严重 - 系统崩溃
    Critical,
}

/// 统一错误码
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ErrorCode {
    // 通用错误 (000-099)
    Unknown = 0,
    InternalError = 1,
    InvalidArgument = 2,
    Timeout = 3,
    Cancelled = 4,
    
    // 数据库错误 (100-199)
    DatabaseError = 100,
    DatabaseConnectionFailed = 101,
    DatabaseQueryFailed = 102,
    DatabaseMigrationFailed = 103,
    DatabaseNotFound = 104,
    
    // 网络错误 (200-299)
    NetworkError = 200,
    HttpRequestFailed = 201,
    InvalidUrl = 202,
    ProxyError = 203,
    SslError = 204,
    
    // 加密错误 (300-399)
    CryptoError = 300,
    EncryptionFailed = 301,
    DecryptionFailed = 302,
    InvalidKey = 303,
    InvalidCiphertext = 304,
    
    // 认证授权错误 (400-499)
    AuthError = 400,
    AuthenticationFailed = 401,
    AuthorizationFailed = 403,
    TokenExpired = 404,
    InvalidToken = 405,
    
    // 文件系统错误 (500-599)
    FileSystemError = 500,
    FileNotFound = 501,
    FileReadFailed = 502,
    FileWriteFailed = 503,
    FileDeleteFailed = 504,
    PermissionDenied = 505,
    
    // 插件错误 (600-699)
    PluginError = 600,
    PluginLoadFailed = 601,
    PluginUnloadFailed = 602,
    PluginNotFound = 603,
    PluginIncompatible = 604,
    PluginExecutionFailed = 605,
    
    // WebShell 错误 (700-799)
    WebShellError = 700,
    WebShellConnectionFailed = 701,
    WebShellNotFound = 702,
    WebShellExecutionFailed = 703,
    WebShellTimeout = 704,
    
    // 项目错误 (800-899)
    ProjectError = 800,
    ProjectNotFound = 801,
    ProjectCreateFailed = 802,
    ProjectUpdateFailed = 803,
    ProjectDeleteFailed = 804,
    
    // 配置错误 (900-999)
    ConfigError = 900,
    ConfigNotFound = 901,
    ConfigParseFailed = 902,
    ConfigSaveFailed = 903,
}

impl ErrorCode {
    pub fn to_str(&self) -> &'static str {
        match self {
            // 通用错误
            ErrorCode::Unknown => "UNKNOWN",
            ErrorCode::InternalError => "INTERNAL_ERROR",
            ErrorCode::InvalidArgument => "INVALID_ARGUMENT",
            ErrorCode::Timeout => "TIMEOUT",
            ErrorCode::Cancelled => "CANCELLED",
            
            // 数据库错误
            ErrorCode::DatabaseError => "DATABASE_ERROR",
            ErrorCode::DatabaseConnectionFailed => "DATABASE_CONNECTION_FAILED",
            ErrorCode::DatabaseQueryFailed => "DATABASE_QUERY_FAILED",
            ErrorCode::DatabaseMigrationFailed => "DATABASE_MIGRATION_FAILED",
            ErrorCode::DatabaseNotFound => "DATABASE_NOT_FOUND",
            
            // 网络错误
            ErrorCode::NetworkError => "NETWORK_ERROR",
            ErrorCode::HttpRequestFailed => "HTTP_REQUEST_FAILED",
            ErrorCode::InvalidUrl => "INVALID_URL",
            ErrorCode::ProxyError => "PROXY_ERROR",
            ErrorCode::SslError => "SSL_ERROR",
            
            // 加密错误
            ErrorCode::CryptoError => "CRYPTO_ERROR",
            ErrorCode::EncryptionFailed => "ENCRYPTION_FAILED",
            ErrorCode::DecryptionFailed => "DECRYPTION_FAILED",
            ErrorCode::InvalidKey => "INVALID_KEY",
            ErrorCode::InvalidCiphertext => "INVALID_CIPHERTEXT",
            
            // 认证授权错误
            ErrorCode::AuthError => "AUTH_ERROR",
            ErrorCode::AuthenticationFailed => "AUTHENTICATION_FAILED",
            ErrorCode::AuthorizationFailed => "AUTHORIZATION_FAILED",
            ErrorCode::TokenExpired => "TOKEN_EXPIRED",
            ErrorCode::InvalidToken => "INVALID_TOKEN",
            
            // 文件系统错误
            ErrorCode::FileSystemError => "FILESYSTEM_ERROR",
            ErrorCode::FileNotFound => "FILE_NOT_FOUND",
            ErrorCode::FileReadFailed => "FILE_READ_FAILED",
            ErrorCode::FileWriteFailed => "FILE_WRITE_FAILED",
            ErrorCode::FileDeleteFailed => "FILE_DELETE_FAILED",
            ErrorCode::PermissionDenied => "PERMISSION_DENIED",
            
            // 插件错误
            ErrorCode::PluginError => "PLUGIN_ERROR",
            ErrorCode::PluginLoadFailed => "PLUGIN_LOAD_FAILED",
            ErrorCode::PluginUnloadFailed => "PLUGIN_UNLOAD_FAILED",
            ErrorCode::PluginNotFound => "PLUGIN_NOT_FOUND",
            ErrorCode::PluginIncompatible => "PLUGIN_INCOMPATIBLE",
            ErrorCode::PluginExecutionFailed => "PLUGIN_EXECUTION_FAILED",
            
            // WebShell 错误
            ErrorCode::WebShellError => "WEBSHELL_ERROR",
            ErrorCode::WebShellConnectionFailed => "WEBSHELL_CONNECTION_FAILED",
            ErrorCode::WebShellNotFound => "WEBSHELL_NOT_FOUND",
            ErrorCode::WebShellExecutionFailed => "WEBSHELL_EXECUTION_FAILED",
            ErrorCode::WebShellTimeout => "WEBSHELL_TIMEOUT",
            
            // 项目错误
            ErrorCode::ProjectError => "PROJECT_ERROR",
            ErrorCode::ProjectNotFound => "PROJECT_NOT_FOUND",
            ErrorCode::ProjectCreateFailed => "PROJECT_CREATE_FAILED",
            ErrorCode::ProjectUpdateFailed => "PROJECT_UPDATE_FAILED",
            ErrorCode::ProjectDeleteFailed => "PROJECT_DELETE_FAILED",
            
            // 配置错误
            ErrorCode::ConfigError => "CONFIG_ERROR",
            ErrorCode::ConfigNotFound => "CONFIG_NOT_FOUND",
            ErrorCode::ConfigParseFailed => "CONFIG_PARSE_FAILED",
            ErrorCode::ConfigSaveFailed => "CONFIG_SAVE_FAILED",
        }
    }
    
    pub fn from_code(code: u32) -> Self {
        match code {
            0 => ErrorCode::Unknown,
            1 => ErrorCode::InternalError,
            2 => ErrorCode::InvalidArgument,
            3 => ErrorCode::Timeout,
            4 => ErrorCode::Cancelled,
            
            100..=104 => match code {
                100 => ErrorCode::DatabaseError,
                101 => ErrorCode::DatabaseConnectionFailed,
                102 => ErrorCode::DatabaseQueryFailed,
                103 => ErrorCode::DatabaseMigrationFailed,
                104 => ErrorCode::DatabaseNotFound,
                _ => ErrorCode::DatabaseError,
            },
            
            200..=204 => match code {
                200 => ErrorCode::NetworkError,
                201 => ErrorCode::HttpRequestFailed,
                202 => ErrorCode::InvalidUrl,
                203 => ErrorCode::ProxyError,
                204 => ErrorCode::SslError,
                _ => ErrorCode::NetworkError,
            },
            
            300..=304 => match code {
                300 => ErrorCode::CryptoError,
                301 => ErrorCode::EncryptionFailed,
                302 => ErrorCode::DecryptionFailed,
                303 => ErrorCode::InvalidKey,
                304 => ErrorCode::InvalidCiphertext,
                _ => ErrorCode::CryptoError,
            },
            
            400..=405 => match code {
                400 => ErrorCode::AuthError,
                401 => ErrorCode::AuthenticationFailed,
                403 => ErrorCode::AuthorizationFailed,
                404 => ErrorCode::TokenExpired,
                405 => ErrorCode::InvalidToken,
                _ => ErrorCode::AuthError,
            },
            
            500..=505 => match code {
                500 => ErrorCode::FileSystemError,
                501 => ErrorCode::FileNotFound,
                502 => ErrorCode::FileReadFailed,
                503 => ErrorCode::FileWriteFailed,
                504 => ErrorCode::FileDeleteFailed,
                505 => ErrorCode::PermissionDenied,
                _ => ErrorCode::FileSystemError,
            },
            
            600..=605 => match code {
                600 => ErrorCode::PluginError,
                601 => ErrorCode::PluginLoadFailed,
                602 => ErrorCode::PluginUnloadFailed,
                603 => ErrorCode::PluginNotFound,
                604 => ErrorCode::PluginIncompatible,
                605 => ErrorCode::PluginExecutionFailed,
                _ => ErrorCode::PluginError,
            },
            
            700..=704 => match code {
                700 => ErrorCode::WebShellError,
                701 => ErrorCode::WebShellConnectionFailed,
                702 => ErrorCode::WebShellNotFound,
                703 => ErrorCode::WebShellExecutionFailed,
                704 => ErrorCode::WebShellTimeout,
                _ => ErrorCode::WebShellError,
            },
            
            800..=804 => match code {
                800 => ErrorCode::ProjectError,
                801 => ErrorCode::ProjectNotFound,
                802 => ErrorCode::ProjectCreateFailed,
                803 => ErrorCode::ProjectUpdateFailed,
                804 => ErrorCode::ProjectDeleteFailed,
                _ => ErrorCode::ProjectError,
            },
            
            900..=903 => match code {
                900 => ErrorCode::ConfigError,
                901 => ErrorCode::ConfigNotFound,
                902 => ErrorCode::ConfigParseFailed,
                903 => ErrorCode::ConfigSaveFailed,
                _ => ErrorCode::ConfigError,
            },
            
            _ => ErrorCode::Unknown,
        }
    }
    
    pub fn code(&self) -> u32 {
        *self as u32
    }
    
    pub fn get_severity(&self) -> ErrorSeverity {
        match self {
            ErrorCode::Unknown | ErrorCode::InternalError => ErrorSeverity::Critical,
            
            ErrorCode::DatabaseConnectionFailed
            | ErrorCode::NetworkError
            | ErrorCode::CryptoError
            | ErrorCode::PluginLoadFailed => ErrorSeverity::High,
            
            ErrorCode::DatabaseQueryFailed
            | ErrorCode::HttpRequestFailed
            | ErrorCode::EncryptionFailed
            | ErrorCode::DecryptionFailed
            | ErrorCode::WebShellConnectionFailed
            | ErrorCode::ProjectCreateFailed => ErrorSeverity::Medium,
            
            _ => ErrorSeverity::Low,
        }
    }
}

/// 应用错误
#[derive(Error, Debug, Serialize, Deserialize)]
pub struct AppError {
    /// 错误码
    pub code: u32,
    /// 错误消息
    pub message: String,
    /// 详细错误信息
    #[serde(skip_serializing_if = "Option::is_none")]
    pub details: Option<String>,
    /// 严重性级别
    pub severity: ErrorSeverity,
}

impl AppError {
    pub fn new(code: ErrorCode, message: impl Into<String>) -> Self {
        Self {
            code: code.code(),
            message: message.into(),
            details: None,
            severity: code.get_severity(),
        }
    }
    
    pub fn with_details(code: ErrorCode, message: impl Into<String>, details: impl Into<String>) -> Self {
        Self {
            code: code.code(),
            message: message.into(),
            details: Some(details.into()),
            severity: code.get_severity(),
        }
    }
    
    pub fn from_error<E: std::error::Error>(code: ErrorCode, error: E) -> Self {
        Self {
            code: code.code(),
            message: error.to_string(),
            details: None,
            severity: code.get_severity(),
        }
    }
    
    pub fn get_code_enum(&self) -> ErrorCode {
        ErrorCode::from_code(self.code)
    }
    
    pub fn user_message(&self) -> String {
        match self.get_code_enum() {
            ErrorCode::Unknown => "发生未知错误",
            ErrorCode::DatabaseConnectionFailed => "数据库连接失败",
            ErrorCode::HttpRequestFailed => "网络请求失败",
            ErrorCode::EncryptionFailed => "加密失败",
            ErrorCode::DecryptionFailed => "解密失败",
            ErrorCode::PluginLoadFailed => "插件加载失败",
            ErrorCode::WebShellConnectionFailed => "WebShell 连接失败",
            _ => &self.message,
        }.to_string()
    }
}

impl std::fmt::Display for AppError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "[{}] {}", ErrorCode::to_str(&ErrorCode::from_code(self.code)), self.message)
    }
}

pub type Result<T> = std::result::Result<T, AppError>;

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_error_code_conversion() {
        assert_eq!(ErrorCode::DatabaseError.code(), 100);
        assert_eq!(ErrorCode::from_code(100), ErrorCode::DatabaseError);
    }
    
    #[test]
    fn test_error_creation() {
        let error = AppError::new(ErrorCode::DatabaseConnectionFailed, "连接失败");
        assert_eq!(error.code, 101);
        assert_eq!(error.severity, ErrorSeverity::High);
    }
}
