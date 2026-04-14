//! FG-ABYSS 集成测试
//! 
//! 包含端到端测试和 API 测试

use fg_abyss::core::crypto::{AesCipher, XorCipher};
use fg_abyss::core::http::{HttpClient, HttpClientConfig};
use fg_abyss::error::{AppError, ErrorCode};

#[cfg(test)]
mod crypto_tests {
    use super::*;

    #[test]
    fn test_aes_encryption_decryption() {
        let key = AesCipher::generate_key();
        let cipher = AesCipher::new(key);
        
        let original = b"Hello, World!";
        let encrypted = cipher.encrypt(original).expect("加密失败");
        let decrypted = cipher.decrypt(&encrypted).expect("解密失败");
        
        assert_eq!(original.to_vec(), decrypted);
    }

    #[test]
    fn test_aes_different_keys() {
        let key1 = AesCipher::generate_key();
        let key2 = AesCipher::generate_key();
        
        let cipher1 = AesCipher::new(key1);
        let cipher2 = AesCipher::new(key2);
        
        let data = b"Test data";
        let encrypted = cipher1.encrypt(data).expect("加密失败");
        
        // 不同密钥无法解密
        let result = cipher2.decrypt(&encrypted);
        assert!(result.is_err());
    }

    #[test]
    fn test_xor_encryption_decryption() {
        let key = b"test_key".to_vec();
        let cipher = XorCipher::new(key);
        
        let original = b"Hello, World!";
        let encrypted = cipher.process(original);
        let decrypted = cipher.process(&encrypted);
        
        assert_eq!(original.to_vec(), decrypted);
    }

    #[test]
    fn test_xor_same_key() {
        let key = b"secret".to_vec();
        let cipher = XorCipher::new(key);
        
        let data1 = b"Message 1";
        let data2 = b"Message 2";
        
        let encrypted1 = cipher.process(data1);
        let encrypted2 = cipher.process(data2);
        
        // 相同密钥加密不同数据，结果不同
        assert_ne!(encrypted1, encrypted2);
    }
}

#[cfg(test)]
mod error_tests {
    use super::*;

    #[test]
    fn test_error_code_conversion() {
        assert_eq!(ErrorCode::DatabaseError.code(), 100);
        assert_eq!(ErrorCode::from_code(100), ErrorCode::DatabaseError);
        
        assert_eq!(ErrorCode::NetworkError.code(), 200);
        assert_eq!(ErrorCode::from_code(201), ErrorCode::HttpRequestFailed);
    }

    #[test]
    fn test_error_severity() {
        let error = AppError::new(ErrorCode::DatabaseConnectionFailed, "连接失败");
        assert_eq!(error.severity, fg_abyss::error::ErrorSeverity::High);
        
        let error = AppError::new(ErrorCode::FileNotFound, "文件未找到");
        assert_eq!(error.severity, fg_abyss::error::ErrorSeverity::Low);
    }

    #[test]
    fn test_error_message() {
        let error = AppError::new(ErrorCode::DatabaseQueryFailed, "查询失败");
        assert_eq!(error.code, 102);
        assert!(error.message.contains("查询失败"));
    }
}

#[cfg(test)]
mod http_tests {
    use super::*;

    #[tokio::test]
    async fn test_http_client_creation() {
        let config = HttpClientConfig::default();
        let client = HttpClient::new(config);
        
        assert!(client.is_ok());
    }

    #[tokio::test]
    async fn test_http_client_timeout() {
        let config = HttpClientConfig {
            timeout: 1,
            connect_timeout: 1,
            ..Default::default()
        };
        
        let client = HttpClient::new(config).expect("创建客户端失败");
        
        // 测试超时（访问一个不存在的地址）
        let result = client.get("http://192.0.2.1/test", None).await;
        assert!(result.is_err());
    }

    #[tokio::test]
    async fn test_http_config_validation() {
        let config = HttpClientConfig::default();
        assert_eq!(config.timeout, 30);
        assert_eq!(config.max_retries, 3);
        assert!(!config.verify_ssl);
    }
}

#[cfg(test)]
mod plugin_tests {
    use fg_abyss::plugins::types::{PluginMetadata, PluginType, PluginError};
    use fg_abyss::plugins::loader::PluginLoader;
    use std::path::PathBuf;
    use std::fs;

    #[test]
    fn test_plugin_type_conversion() {
        assert_eq!(PluginType::from_str("terminal"), PluginType::Terminal);
        assert_eq!(PluginType::from_str("file_manager"), PluginType::FileManager);
        assert_eq!(PluginType::from_str("database"), PluginType::Database);
        
        // 测试 to_str
        assert_eq!(PluginType::Terminal.to_str(), "terminal");
    }

    #[test]
    fn test_plugin_loader_creation() {
        let temp_dir = std::env::temp_dir().join("fg_abyss_test_plugin_loader");
        let _ = fs::remove_dir_all(&temp_dir);
        
        let loader = PluginLoader::new(temp_dir.clone());
        assert!(loader.plugins_dir().exists());
        
        let _ = fs::remove_dir_all(&temp_dir);
    }

    #[test]
    fn test_plugin_metadata_validation() {
        let metadata = PluginMetadata {
            id: "test-plugin".to_string(),
            name: "Test Plugin".to_string(),
            description: "A test plugin".to_string(),
            version: "1.0.0".to_string(),
            author: "Tester".to_string(),
            min_version: "0.1.0".to_string(),
            plugin_type: PluginType::Terminal,
            entry_point: "index.js".to_string(),
            dependencies: vec![],
            permissions: vec!["execute".to_string()],
        };
        
        assert_eq!(metadata.id, "test-plugin");
        assert_eq!(metadata.version, "1.0.0");
    }
}

#[cfg(test)]
mod logging_tests {
    use fg_abyss::logging::{LogLevel, LogEntry, LoggingConfig};

    #[test]
    fn test_log_level_order() {
        assert!(LogLevel::Trace < LogLevel::Debug);
        assert!(LogLevel::Debug < LogLevel::Info);
        assert!(LogLevel::Info < LogLevel::Warn);
        assert!(LogLevel::Warn < LogLevel::Error);
    }

    #[test]
    fn test_log_level_conversion() {
        assert_eq!(LogLevel::from_str("trace"), LogLevel::Trace);
        assert_eq!(LogLevel::from_str("DEBUG"), LogLevel::Debug);
        assert_eq!(LogLevel::from_str("info"), LogLevel::Info);
        assert_eq!(LogLevel::from_str("warn"), LogLevel::Warn);
        assert_eq!(LogLevel::from_str("ERROR"), LogLevel::Error);
        assert_eq!(LogLevel::from_str("unknown"), LogLevel::Info);
    }

    #[test]
    fn test_log_entry_creation() {
        let entry = LogEntry::new(LogLevel::Info, "test", "test message");
        assert_eq!(entry.level, LogLevel::Info);
        assert_eq!(entry.target, "test");
        assert_eq!(entry.message, "test message");
        assert!(entry.data.is_none());
    }

    #[test]
    fn test_log_entry_with_data() {
        let entry = LogEntry::new(LogLevel::Debug, "test", "test message")
            .with_data(serde_json::json!({"key": "value"}));
        
        assert!(entry.data.is_some());
        assert_eq!(entry.data.unwrap(), serde_json::json!({"key": "value"}));
    }

    #[test]
    fn test_logging_config_default() {
        let config = LoggingConfig::default();
        assert_eq!(config.level, LogLevel::Debug);
        assert!(config.console_output);
        assert!(config.file_output);
        assert!(config.rotation_enabled);
        assert!(config.sanitize_enabled);
    }
}
