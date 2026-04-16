use chrono::{DateTime, Local};
use serde::{Serialize, Deserialize};
use std::fs::{File, OpenOptions};
use std::io::{Write, BufWriter};
use std::path::PathBuf;
use std::sync::{Arc, Mutex};

/// 日志级别
#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
pub enum LogLevel {
    /// 追踪 - 最详细的日志
    Trace = 0,
    /// 调试 - 调试信息
    Debug = 1,
    /// 信息 - 一般信息
    Info = 2,
    /// 警告 - 警告信息
    Warn = 3,
    /// 错误 - 错误信息
    Error = 4,
}

impl LogLevel {
    pub fn as_str(&self) -> &'static str {
        match self {
            LogLevel::Trace => "TRACE",
            LogLevel::Debug => "DEBUG",
            LogLevel::Info => "INFO",
            LogLevel::Warn => "WARN",
            LogLevel::Error => "ERROR",
        }
    }

    pub fn from_str(s: &str) -> Self {
        match s.to_uppercase().as_str() {
            "TRACE" => LogLevel::Trace,
            "DEBUG" => LogLevel::Debug,
            "INFO" => LogLevel::Info,
            "WARN" => LogLevel::Warn,
            "ERROR" => LogLevel::Error,
            _ => LogLevel::Info,
        }
    }
}

/// 日志条目
#[derive(Debug, Clone)]
pub struct LogEntry {
    /// 时间戳
    pub timestamp: DateTime<Local>,
    /// 日志级别
    pub level: LogLevel,
    /// 目标模块
    pub target: String,
    /// 日志消息
    pub message: String,
    /// 额外数据（可选）
    pub data: Option<serde_json::Value>,
}

impl LogEntry {
    pub fn new(level: LogLevel, target: &str, message: &str) -> Self {
        Self {
            timestamp: Local::now(),
            level,
            target: target.to_string(),
            message: message.to_string(),
            data: None,
        }
    }

    pub fn with_data(mut self, data: serde_json::Value) -> Self {
        self.data = Some(data);
        self
    }

    pub fn format(&self, colored: bool) -> String {
        let timestamp = self.timestamp.format("%Y-%m-%d %H:%M:%S%.3f");
        let level_str = self.level.as_str();
        
        if colored {
            let level_color = match self.level {
                LogLevel::Trace => "\x1b[35m",    // 紫色
                LogLevel::Debug => "\x1b[36m",    // 青色
                LogLevel::Info => "\x1b[32m",     // 绿色
                LogLevel::Warn => "\x1b[33m",     // 黄色
                LogLevel::Error => "\x1b[31m",    // 红色
            };
            let reset = "\x1b[0m";
            format!("{} [{}] [{}] {} - {}{}", timestamp, level_color, level_str, self.target, self.message, reset)
        } else {
            format!("{} [{}] [{}] - {}", timestamp, level_str, self.target, self.message)
        }
    }
}

/// 日志配置
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LoggingConfig {
    /// 日志级别
    pub level: LogLevel,
    /// 是否输出到控制台
    pub console_output: bool,
    /// 是否输出到文件
    pub file_output: bool,
    /// 日志文件目录
    pub log_dir: PathBuf,
    /// 单个日志文件最大大小（MB）
    pub max_file_size_mb: u64,
    /// 保留的日志文件数量
    pub max_files: u32,
    /// 是否启用日志轮转
    pub rotation_enabled: bool,
    /// 是否脱敏敏感信息
    pub sanitize_enabled: bool,
}

impl Default for LoggingConfig {
    fn default() -> Self {
        Self {
            level: LogLevel::Debug,
            console_output: true,
            file_output: true,
            log_dir: PathBuf::from("logs"),
            max_file_size_mb: 10,
            max_files: 7,
            rotation_enabled: true,
            sanitize_enabled: true,
        }
    }
}

/// 日志记录器
pub struct Logger {
    config: LoggingConfig,
    writer: Arc<Mutex<Option<BufWriter<File>>>>,
    current_file_size: Arc<Mutex<u64>>,
}

impl Logger {
    pub fn new(config: LoggingConfig) -> Self {
        let mut logger = Self {
            config,
            writer: Arc::new(Mutex::new(None)),
            current_file_size: Arc::new(Mutex::new(0)),
        };

        // 初始化日志文件
        if logger.config.file_output {
            logger.init_file_writer().expect("初始化日志文件失败");
        }

        logger
    }

    /// 初始化日志文件写入器
    fn init_file_writer(&mut self) -> std::io::Result<()> {
        // 创建日志目录
        std::fs::create_dir_all(&self.config.log_dir)?;

        // 生成日志文件名
        let date = Local::now().format("%Y-%m-%d");
        let log_file = self.config.log_dir.join(format!("fg-abyss-{}.log", date));

        // 打开文件
        let file = OpenOptions::new()
            .create(true)
            .append(true)
            .open(&log_file)?;

        let file_size = file.metadata()?.len();
        *self.current_file_size.lock().unwrap() = file_size;

        *self.writer.lock().unwrap() = Some(BufWriter::new(file));

        Ok(())
    }

    /// 检查是否需要轮转日志文件
    fn check_rotation(&mut self) -> std::io::Result<()> {
        if !self.config.rotation_enabled {
            return Ok(());
        }

        let current_size = *self.current_file_size.lock().unwrap();
        let max_size = self.config.max_file_size_mb * 1024 * 1024;

        if current_size >= max_size {
            self.rotate_logs()?;
        }

        Ok(())
    }

    /// 轮转日志文件
    fn rotate_logs(&mut self) -> std::io::Result<()> {
        // 删除最旧的日志文件
        let mut log_files: Vec<_> = std::fs::read_dir(&self.config.log_dir)?
            .filter_map(|e| e.ok())
            .filter(|e| e.path().extension().map_or(false, |ext| ext == "log"))
            .collect();

        log_files.sort_by_key(|e| e.path());

        // 删除超出数量的旧文件
        while log_files.len() >= self.config.max_files as usize {
            if let Some(old_file) = log_files.first() {
                std::fs::remove_file(old_file.path())?;
                log_files.remove(0);
            } else {
                break;
            }
        }

        // 重新初始化写入器
        drop(self.writer.lock().unwrap().take());
        self.init_file_writer()?;

        Ok(())
    }

    /// 脱敏敏感信息
    fn sanitize_message(&self, message: &str) -> String {
        if !self.config.sanitize_enabled {
            return message.to_string();
        }

        let mut sanitized = message.to_string();

        // 脱敏密码
        if let Ok(regex) = regex::Regex::new(r"(?i)(password|passwd|pwd|secret|token|api_key|apikey)\s*[=:]\s*[^\\s]+") {
            sanitized = regex.replace_all(&sanitized, "$1=***REDACTED***").to_string();
        }

        // 脱敏邮箱
        if let Ok(regex) = regex::Regex::new(r"[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}") {
            sanitized = regex.replace_all(&sanitized, "***@***.***").to_string();
        }

        // 脱敏 IP 地址
        if let Ok(regex) = regex::Regex::new(r"\\b\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\b") {
            sanitized = regex.replace_all(&sanitized, "***.***.***.***").to_string();
        }

        sanitized
    }

    /// 记录日志
    pub fn log(&self, level: LogLevel, target: &str, message: &str, data: Option<serde_json::Value>) {
        if level < self.config.level {
            return;
        }

        let mut entry = LogEntry::new(level, target, &self.sanitize_message(message));
        if let Some(d) = data {
            entry = entry.with_data(d);
        }

        // 输出到控制台
        if self.config.console_output {
            println!("{}", entry.format(true));
        }

        // 输出到文件
        if self.config.file_output {
            if let Ok(mut writer_guard) = self.writer.lock() {
                if let Some(writer) = writer_guard.as_mut() {
                    let log_line = format!("{}\n", entry.format(false));
                    let _ = writer.write_all(log_line.as_bytes());
                    *self.current_file_size.lock().unwrap() += log_line.len() as u64;
                }
            }
        }
    }

    pub fn trace(&self, target: &str, message: &str) {
        self.log(LogLevel::Trace, target, message, None);
    }

    pub fn debug(&self, target: &str, message: &str) {
        self.log(LogLevel::Debug, target, message, None);
    }

    pub fn info(&self, target: &str, message: &str) {
        self.log(LogLevel::Info, target, message, None);
    }

    pub fn warn(&self, target: &str, message: &str) {
        self.log(LogLevel::Warn, target, message, None);
    }

    pub fn error(&self, target: &str, message: &str) {
        self.log(LogLevel::Error, target, message, None);
    }

    pub fn flush(&self) {
        if let Ok(mut writer_guard) = self.writer.lock() {
            if let Some(writer) = writer_guard.as_mut() {
                let _ = writer.flush();
            }
        }
    }
}

impl Drop for Logger {
    fn drop(&mut self) {
        self.flush();
    }
}

/// 全局日志记录器
static mut GLOBAL_LOGGER: Option<Logger> = None;

/// 初始化全局日志记录器
pub fn init_global_logger(config: LoggingConfig) {
    unsafe {
        GLOBAL_LOGGER = Some(Logger::new(config));
    }
}

/// 获取全局日志记录器
pub fn get_global_logger() -> Option<&'static Logger> {
    unsafe { GLOBAL_LOGGER.as_ref() }
}

/// 便捷日志宏
#[macro_export]
macro_rules! trace {
    ($target:expr, $message:expr) => {
        if let Some(logger) = $crate::logging::get_global_logger() {
            logger.trace($target, $message);
        }
    };
    ($target:expr, $message:expr, $data:expr) => {
        if let Some(logger) = $crate::logging::get_global_logger() {
            logger.log($crate::logging::LogLevel::Trace, $target, $message, Some($data));
        }
    };
}

#[macro_export]
macro_rules! debug {
    ($target:expr, $message:expr) => {
        if let Some(logger) = $crate::logging::get_global_logger() {
            logger.debug($target, $message);
        }
    };
}

#[macro_export]
macro_rules! info {
    ($target:expr, $message:expr) => {
        if let Some(logger) = $crate::logging::get_global_logger() {
            logger.info($target, $message);
        }
    };
}

#[macro_export]
macro_rules! warn {
    ($target:expr, $message:expr) => {
        if let Some(logger) = $crate::logging::get_global_logger() {
            logger.warn($target, $message);
        }
    };
}

#[macro_export]
macro_rules! error {
    ($target:expr, $message:expr) => {
        if let Some(logger) = $crate::logging::get_global_logger() {
            logger.error($target, $message);
        }
    };
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_log_level_order() {
        assert!(LogLevel::Trace < LogLevel::Debug);
        assert!(LogLevel::Debug < LogLevel::Info);
        assert!(LogLevel::Info < LogLevel::Warn);
        assert!(LogLevel::Warn < LogLevel::Error);
    }

    #[test]
    fn test_log_entry_creation() {
        let entry = LogEntry::new(LogLevel::Info, "test", "test message");
        assert_eq!(entry.level, LogLevel::Info);
        assert_eq!(entry.target, "test");
        assert_eq!(entry.message, "test message");
    }
}
