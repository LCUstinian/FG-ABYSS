/// ASP 纯 ASCII 转换器
/// 将非 ASCII 字符转换为 Chr() 拼接，避免 IIS 乱码问题
pub struct AsciiConverter;

impl AsciiConverter {
    /// 转换 ASP 代码为纯 ASCII (使用 Chr() 拼接)
    pub fn convert(code: &str) -> String {
        code.chars()
            .map(|c| {
                if Self::is_safe_ascii(c) {
                    c.to_string()
                } else {
                    format!("Chr({})", c as u32)
                }
            })
            .collect()
    }

    /// 检查字符是否为安全的可打印 ASCII
    fn is_safe_ascii(c: char) -> bool {
        c.is_ascii() && c.is_ascii_graphic() && c != '\''
    }

    /// 批量转换字符串拼接
    /// 将长字符串转换为多个 Chr() 拼接
    pub fn convert_string_literal(s: &str) -> String {
        s.chars()
            .map(|c| {
                if Self::is_safe_ascii(c) {
                    format!("\"{}\"", c)
                } else {
                    format!("Chr({})", c as u32)
                }
            })
            .collect::<Vec<_>>()
            .join(" & ")
    }
}
