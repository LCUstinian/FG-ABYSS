use serde::{Deserialize, Serialize};
use std::path::PathBuf;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum PluginError {
    #[error("插件加载失败：{0}")]
    LoadFailed(String),
    #[error("插件验证失败：{0}")]
    ValidationFailed(String),
    #[error("插件执行失败：{0}")]
    ExecutionFailed(String),
    #[error("插件不兼容：{0}")]
    Incompatible(String),
    #[error("插件已存在：{0}")]
    AlreadyExists(String),
    #[error("插件未找到：{0}")]
    NotFound(String),
}

pub type Result<T> = std::result::Result<T, PluginError>;

/// 插件元数据
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PluginMetadata {
    /// 插件唯一标识
    pub id: String,
    /// 插件名称
    pub name: String,
    /// 插件描述
    pub description: String,
    /// 插件版本
    pub version: String,
    /// 作者
    pub author: String,
    /// 最小 FG-ABYSS 版本要求
    pub min_version: String,
    /// 插件类型
    pub plugin_type: PluginType,
    /// 插件入口文件
    pub entry_point: String,
    /// 依赖的插件 ID 列表
    pub dependencies: Vec<String>,
    /// 权限列表
    pub permissions: Vec<String>,
}

/// 插件类型
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum PluginType {
    /// 终端插件
    Terminal,
    /// 文件管理插件
    FileManager,
    /// 数据库插件
    Database,
    /// 网络工具插件
    NetworkTool,
    /// 加密工具插件
    CryptoTool,
    /// 其他类型
    Other(String),
}

impl PluginType {
    pub fn from_str(s: &str) -> Self {
        match s.to_lowercase().as_str() {
            "terminal" => PluginType::Terminal,
            "file_manager" => PluginType::FileManager,
            "database" => PluginType::Database,
            "network_tool" => PluginType::NetworkTool,
            "crypto_tool" => PluginType::CryptoTool,
            other => PluginType::Other(other.to_string()),
        }
    }

    pub fn to_str(&self) -> String {
        match self {
            PluginType::Terminal => "terminal".to_string(),
            PluginType::FileManager => "file_manager".to_string(),
            PluginType::Database => "database".to_string(),
            PluginType::NetworkTool => "network_tool".to_string(),
            PluginType::CryptoTool => "crypto_tool".to_string(),
            PluginType::Other(s) => s.clone(),
        }
    }
}

/// 插件状态
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum PluginStatus {
    /// 未激活
    Inactive,
    /// 已激活
    Active,
    /// 加载中
    Loading,
    /// 错误状态
    Error(String),
}

/// 插件实例
#[derive(Debug, Clone)]
pub struct Plugin {
    /// 插件元数据
    pub metadata: PluginMetadata,
    /// 插件文件路径
    pub path: PathBuf,
    /// 插件状态
    pub status: PluginStatus,
    /// 插件配置
    pub config: Option<serde_json::Value>,
}

impl Plugin {
    pub fn new(path: PathBuf, metadata: PluginMetadata) -> Self {
        Self {
            metadata,
            path,
            status: PluginStatus::Inactive,
            config: None,
        }
    }

    pub fn is_active(&self) -> bool {
        self.status == PluginStatus::Active
    }

    pub fn validate(&self, app_version: &str) -> Result<()> {
        // 验证版本兼容性
        if !self.is_version_compatible(app_version) {
            return Err(PluginError::Incompatible(format!(
                "插件要求最低版本 {}，当前版本 {}",
                self.metadata.min_version, app_version
            )));
        }

        // 验证入口文件存在
        if !self.path.join(&self.metadata.entry_point).exists() {
            return Err(PluginError::ValidationFailed(format!(
                "入口文件不存在：{}",
                self.metadata.entry_point
            )));
        }

        Ok(())
    }

    fn is_version_compatible(&self, app_version: &str) -> bool {
        // 简单的版本比较（实际应该使用 semver 库）
        app_version.to_string() >= self.metadata.min_version
    }
}

/// 插件配置
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PluginConfig {
    /// 是否启用
    pub enabled: bool,
    /// 自动加载
    pub auto_load: bool,
    /// 优先级
    pub priority: u32,
    /// 自定义配置
    pub settings: serde_json::Value,
}

impl Default for PluginConfig {
    fn default() -> Self {
        Self {
            enabled: true,
            auto_load: false,
            priority: 0,
            settings: serde_json::Value::Null,
        }
    }
}
