#[derive(Debug, thiserror::Error)]
pub enum AppError {
    #[error("数据库错误: {0}")]
    Database(#[from] rusqlite::Error),
    #[error("数据库连接错误: {0}")]
    DbConnect(String),
    #[error("HTTP 错误: {0}")]
    Http(#[from] reqwest::Error),
    #[error("加密错误: {0}")]
    Crypto(String),
    #[error("WebShell 连接失败: {0}")]
    Connection(String),
    #[error("WebShell 响应验证失败: {0}")]
    InvalidResponse(String),
    #[error("熔断器开启: {0}")]
    CircuitOpen(String),
    #[error("未找到: {0}")]
    NotFound(String),
    #[error("参数无效: {0}")]
    InvalidInput(String),
    #[error("IO 错误: {0}")]
    Io(#[from] std::io::Error),
    #[error("序列化错误: {0}")]
    Serialize(#[from] serde_json::Error),
    #[error("应用已锁定")]
    Locked,
    #[error("插件错误: {0}")]
    Plugin(String),
    #[error("WebShell 需要重新部署: {0}")]
    NeedsRedeploy(String),
    #[error("内存马会话已失效（服务端重启导致）")]
    MemShellExpired,
}

impl AppError {
    pub fn kind(&self) -> &'static str {
        match self {
            Self::Database(_) | Self::DbConnect(_) => "Database",
            Self::Http(_)            => "Http",
            Self::Crypto(_)          => "Crypto",
            Self::Connection(_)      => "Connection",
            Self::InvalidResponse(_) => "InvalidResponse",
            Self::CircuitOpen(_)     => "CircuitOpen",
            Self::NotFound(_)        => "NotFound",
            Self::InvalidInput(_)    => "InvalidInput",
            Self::Io(_)              => "Io",
            Self::Serialize(_)       => "Serialize",
            Self::Locked             => "Locked",
            Self::Plugin(_)          => "Plugin",
            Self::NeedsRedeploy(_)   => "NeedsRedeploy",
            Self::MemShellExpired    => "MemShellExpired",
        }
    }
}

impl serde::Serialize for AppError {
    fn serialize<S>(&self, s: S) -> std::result::Result<S::Ok, S::Error>
    where S: serde::Serializer {
        use serde::ser::SerializeStruct;
        let mut state = s.serialize_struct("AppError", 2)?;
        state.serialize_field("kind", self.kind())?;
        state.serialize_field("message", &self.to_string())?;
        state.end()
    }
}

pub type Result<T> = std::result::Result<T, AppError>;
