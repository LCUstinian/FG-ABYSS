/// 载荷生成模块的类型定义
/// 用于前端和后端之间的 IPC 通信

use serde::{Deserialize, Serialize};

/// 载荷生成模式
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum PayloadMode {
    /// 极简一句话模式 - 仅编码/混淆 (向后兼容)
    Simple,
    /// 高级加密壳模式 - 强加密 (向后兼容)
    Advanced,
    /// 基于文件的载荷模式
    FileBased,
    /// 内存Shell模式
    MemoryShell,
    /// 仅Suo5模式
    Suo5Only,
}

/// 脚本类型
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "lowercase")]
pub enum ScriptType {
    Php,
    Jsp,
    Aspx,
    Asp,
}

impl std::fmt::Display for ScriptType {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            ScriptType::Php => write!(f, "php"),
            ScriptType::Jsp => write!(f, "jsp"),
            ScriptType::Aspx => write!(f, "aspx"),
            ScriptType::Asp => write!(f, "asp"),
        }
    }
}

/// 注入类型
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum InjectionType {
    TomcatFilter,
    SpringInterceptor,
    IisHttpModule,
}

/// Suo5 配置结构体
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Suo5Config {
    pub auth: String,
    pub path: String,
    pub timeout: u32,
}

/// 编码器类型 (仅 Simple 模式使用)
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum EncodeType {
    None,
    Base64,
    Xor,
    GzInflate,
    Hex,
    UrlEncode,
    Rot13,
}

/// 加密算法 (仅 Advanced 模式使用)
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum EncryptAlgo {
    Aes128Cbc,
    Aes256Cbc,
    Xor,
}

/// 功能类型
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum FunctionType {
    Basic,
    FileManager,
    ProcessManager,
    Registry,
    Network,
}

/// 混淆级别
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "lowercase")]
pub enum ObfuscationLevel {
    Low,
    Medium,
    High,
}

/// 载荷生成配置
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PayloadConfig {
    /// 生成模式
    pub mode: PayloadMode,
    /// 脚本类型
    pub script_type: ScriptType,
    /// 功能类型
    pub function_type: FunctionType,
    /// 连接密码
    pub password: String,
    /// 编码器类型 (仅 Simple 模式)
    pub encode_type: Option<EncodeType>,
    /// 加密算法 (仅 Advanced 模式)
    pub encrypt_algo: Option<EncryptAlgo>,
    /// 混淆级别
    pub obfuscation_level: ObfuscationLevel,
    /// 输出文件名 (可选)
    pub output_filename: Option<String>,
    /// 模板名称 (可选，用于自定义模板)
    pub template_name: Option<String>,
    /// 注入类型 (可选)
    pub injection_type: Option<InjectionType>,
    /// Suo5 配置 (可选)
    pub suo5_config: Option<Suo5Config>,
}

/// 客户端配置 (仅 Advanced 模式需要)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ClientConfig {
    /// 加密密钥
    pub key: String,
    /// 初始化向量 (IV)
    pub iv: String,
    /// 加密算法
    pub algorithm: String,
    /// 其他配置项
    pub options: serde_json::Value,
}

/// 载荷生成结果
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PayloadResult {
    /// 生成的源代码
    pub code: String,
    /// 客户端配置 (仅 Advanced 模式)
    pub client_config: Option<ClientConfig>,
    /// 文件名
    pub filename: String,
    /// 文件大小 (字节)
    pub size: u64,
    /// 成功标志
    pub success: bool,
    /// 消息
    pub message: Option<String>,
}

/// 载荷模板结构体
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PayloadTemplate {
    /// 模板名称
    pub name: String,
    /// 脚本类型
    pub script_type: ScriptType,
    /// 功能类型
    pub function_type: FunctionType,
    /// 模板代码
    pub code: String,
    /// 描述
    pub description: String,
    /// 创建时间
    pub created_at: String,
    /// 更新时间
    pub updated_at: String,
}

/// 错误类型
#[derive(Debug, thiserror::Error)]
pub enum GeneratorError {
    #[error("生成失败：{0}")]
    GenerationFailed(String),
    
    #[error("无效的密码：{0}")]
    InvalidPassword(String),
    
    #[error("文件操作失败：{0}")]
    IoError(#[from] std::io::Error),
    
    #[error("JSON 序列化失败：{0}")]
    JsonError(#[from] serde_json::Error),
    
    #[error("模板不存在：{0}")]
    TemplateNotFound(String),
    
    #[error("模板已存在：{0}")]
    TemplateExists(String),
}

pub type Result<T> = std::result::Result<T, GeneratorError>;
