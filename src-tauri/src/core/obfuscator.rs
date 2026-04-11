use rand::{Rng, distributions::Alphanumeric};
use std::collections::HashMap;

/// 支持的脚本语言
#[derive(Debug, Clone, Copy)]
pub enum ScriptLanguage {
    PHP,
    JSP,
    ASP,
    ASPX,
}

/// 混淆强度
#[derive(Debug, Clone, Copy)]
pub enum ObfuscationLevel {
    L1,
    L2,
    L3,
    L4,
}

/// 代码混淆器
pub struct Obfuscator;

impl Obfuscator {
    /// 应用混淆
    pub fn obfuscate(code: &str, language: ScriptLanguage, level: ObfuscationLevel) -> String {
        let mut result = code.to_string();
        
        // 预计算一些常用值，避免重复计算
        
        match level {
            ObfuscationLevel::L1 => {
                result = Self::rename_variables(&result, language);
            }
            ObfuscationLevel::L2 => {
                result = Self::rename_variables(&result, language);
                result = Self::insert_random_garbage(&result, language);
            }
            ObfuscationLevel::L3 => {
                result = Self::rename_variables(&result, language);
                result = Self::insert_random_garbage(&result, language);
                result = Self::add_dead_code(&result, language);
            }
            ObfuscationLevel::L4 => {
                result = Self::rename_variables(&result, language);
                result = Self::insert_random_garbage(&result, language);
                result = Self::add_dead_code(&result, language);
                result = Self::flatten_control_flow(&result, language);
            }
        }
        
        result
    }

    /// 变量重命名
    fn rename_variables(code: &str, language: ScriptLanguage) -> String {
        let mut variable_map: HashMap<String, String> = HashMap::new();
        let mut rng = rand::thread_rng();

        match language {
            ScriptLanguage::PHP => {
                let mut result = String::with_capacity(code.len() * 2);
                let chars = code.chars().collect::<Vec<_>>();
                let mut i = 0;
                
                while i < chars.len() {
                    if chars[i] == '$' && i + 1 < chars.len() && chars[i + 1].is_alphabetic() {
                        // 找到变量名
                        let mut j = i + 1;
                        while j < chars.len() && (chars[j].is_alphanumeric() || chars[j] == '_') {
                            j += 1;
                        }
                        let var_name = chars[i+1..j].iter().collect::<String>();
                        
                        // 生成新变量名
                        let new_name = variable_map.entry(var_name.clone())
                            .or_insert_with(|| {
                                let len = rng.gen_range(2..10);
                                std::iter::repeat(())
                                    .map(|_| rng.sample(Alphanumeric))
                                    .map(char::from)
                                    .take(len)
                                    .collect()
                            });
                        
                        // 写入新变量名
                        result.push('$');
                        result.push_str(new_name);
                        i = j;
                    } else {
                        result.push(chars[i]);
                        i += 1;
                    }
                }
                result
            }
            ScriptLanguage::JSP | ScriptLanguage::ASPX => {
                let mut result = code.to_string();
                result = Self::rename_java_variables(&mut result, &mut rng);
                result
            }
            ScriptLanguage::ASP => {
                let mut result = code.to_string();
                result = Self::rename_vb_variables(&mut result, &mut rng);
                result
            }
        }
    }

    /// 重命名Java/JSP变量
    fn rename_java_variables(code: &mut String, rng: &mut impl Rng) -> String {
        let mut result = code.clone();
        let patterns = &["int ", "String ", "boolean ", "var "];
        for &prefix in patterns {
            let len = prefix.len();
            let indices: Vec<usize> = code.rmatch_indices(prefix)
                .map(|(i, _)| i + len)
                .collect();
            for start in indices {
                let mut end = start;
                let chars: Vec<char> = code.chars().collect();
                while end < chars.len() && chars[end].is_alphanumeric() {
                    end += 1;
                }
                if end > start {
                    let var_name = &chars[start..end].iter().collect::<String>();
                    if !var_name.is_empty() {
                        let new_len = rng.gen_range(3..12);
                        let new_name: String = std::iter::repeat(())
                            .map(|_| rng.sample(Alphanumeric))
                            .map(char::from)
                            .take(new_len)
                            .collect();
                        result = result.replace(var_name, &new_name);
                    }
                }
            }
        }
        result
    }

    /// 重命名VB/ASP变量
    fn rename_vb_variables(code: &mut String, rng: &mut impl Rng) -> String {
        let mut result = code.clone();
        let patterns = &["Dim ", "Dim ", "Const "];
        for &prefix in patterns {
            let len = prefix.len();
            let indices: Vec<usize> = code.rmatch_indices(prefix)
                .map(|(i, _)| i + len)
                .collect();
            for start in indices {
                let mut end = start;
                let chars: Vec<char> = code.chars().collect();
                while end < chars.len() && chars[end].is_alphanumeric() || chars[end] == '_' {
                    end += 1;
                }
                if end > start {
                    let var_name = &chars[start..end].iter().collect::<String>();
                    if !var_name.is_empty() {
                        let new_len = rng.gen_range(3..10);
                        let new_name: String = std::iter::repeat(())
                            .map(|_| rng.sample(Alphanumeric))
                            .map(char::from)
                            .take(new_len)
                            .collect();
                        result = result.replace(var_name, &new_name);
                    }
                }
            }
        }
        result
    }

    /// 插入随机垃圾代码
    fn insert_random_garbage(code: &str, language: ScriptLanguage) -> String {
        let mut rng = rand::thread_rng();
        let lines: Vec<&str> = code.lines().collect();
        let mut result = String::with_capacity(code.len() * 2);

        for line in lines {
            result.push_str(line);
            result.push('\n');

            if rng.gen_bool(0.3) {
                let garbage = Self::generate_random_garbage(language, &mut rng);
                result.push_str(&garbage);
                result.push('\n');
            }
        }

        result
    }

    /// 生成随机垃圾代码
    fn generate_random_garbage(language: ScriptLanguage, rng: &mut impl Rng) -> String {
        let a: i32 = rng.gen_range(1..1000);
        let b: i32 = rng.gen_range(1..1000);
        let c: i32 = a + b * rng.gen_range(2..50);

        match language {
            ScriptLanguage::PHP => format!("${}_{} = {}; ${}_{} = {}; ${}_{} = ${}_{} * {};", 
                "x", a, a, "y", b, b, "z", a, "y", b, c),
            ScriptLanguage::JSP | ScriptLanguage::ASPX => format!("int {}_{} = {}; int {}_{} = {}; int {}_{} = {}_{} * {};", 
                "x", a, a, "y", b, b, "z", a, "y", b, c),
            ScriptLanguage::ASP => format!("Dim {}_{}: {}_{} = {}: Dim {}_{}: {}_{} = {}: Dim {}_{}: {}_{} = {}_{} * {};", 
                "x", a, "x", a, a, "y", b, "y", b, b, "z", a, "z", a, "y", b, c),
        }
    }

    /// 添加不执行的死代码
    fn add_dead_code(code: &str, language: ScriptLanguage) -> String {
        let mut rng = rand::thread_rng();
        let dead_code = match language {
            ScriptLanguage::PHP => {
                if rng.gen_bool(0.5) {
                    "\nif (false) { $dead = 123; echo 'dead'; }\n"
                } else {
                    "\n$func_dead = function() { return false; }; if ($func_dead()) { echo 'never'; }\n"
                }
            },
            ScriptLanguage::JSP | ScriptLanguage::ASPX => {
                "if (false) { int dead = 123; }\n"
            },
            ScriptLanguage::ASP => {
                "If False Then\n    Dim dead: dead = 123\nEnd If\n"
            },
        };
        // 在随机位置插入死代码
        let mut lines: Vec<&str> = code.lines().collect();
        let insert_pos = rng.gen_range(1..lines.len().max(2));
        lines.insert(insert_pos, dead_code);
        lines.join("\n")
    }

    /// 控制流平坦化
    fn flatten_control_flow(code: &str, language: ScriptLanguage) -> String {
        match language {
            ScriptLanguage::PHP => Self::flatten_php_control_flow(code),
            ScriptLanguage::JSP | ScriptLanguage::ASPX => Self::flatten_java_control_flow(code),
            ScriptLanguage::ASP => Self::flatten_vb_control_flow(code),
        }
    }

    /// PHP 控制流平坦化
    fn flatten_php_control_flow(code: &str) -> String {
        let mut rng = rand::thread_rng();
        let state_var = format!("${}", Self::generate_random_string(&mut rng, 8));
        let mut result = format!("{state_var} = 0;\n");
        
        let lines: Vec<&str> = code.lines().collect();
        let mut cases = Vec::new();
        
        for (i, line) in lines.iter().enumerate() {
            if line.trim().is_empty() || line.trim().starts_with("//") || line.trim().starts_with("#") {
                result.push_str(&format!("{line}\n"));
                continue;
            }
            
            let case_label = i + 1;
            cases.push(format!("        case {case_label}:\n            {line}\n            {state_var} = {next_case};\n            break;", 
                next_case = if i == lines.len() - 1 { 0 } else { case_label + 1 }
            ));
        }
        
        result.push_str(&format!("while (true) {{
    switch ({state_var}) {{
{}
        default:
            break 2;
    }}
}}", cases.join("\n")));
        
        result
    }

    /// Java/JSP 控制流平坦化
    fn flatten_java_control_flow(code: &str) -> String {
        let mut rng = rand::thread_rng();
        let state_var = Self::generate_random_string(&mut rng, 8);
        let mut result = format!("int {state_var} = 0;\n");
        
        let lines: Vec<&str> = code.lines().collect();
        let mut cases = Vec::new();
        
        for (i, line) in lines.iter().enumerate() {
            if line.trim().is_empty() || line.trim().starts_with("//") {
                result.push_str(&format!("{line}\n"));
                continue;
            }
            
            let case_label = i + 1;
            cases.push(format!("        case {case_label}:\n            {line}\n            {state_var} = {next_case};\n            break;", 
                next_case = if i == lines.len() - 1 { 0 } else { case_label + 1 }
            ));
        }
        
        result.push_str(&format!("while (true) {{
    switch ({state_var}) {{
{}
        default:
            break;
    }}
    break;
}}", cases.join("\n")));
        
        result
    }

    /// VB/ASP 控制流平坦化
    fn flatten_vb_control_flow(code: &str) -> String {
        let mut rng = rand::thread_rng();
        let state_var = Self::generate_random_string(&mut rng, 8);
        let mut result = format!("Dim {state_var}: {state_var} = 0\n");
        
        let lines: Vec<&str> = code.lines().collect();
        let mut cases = Vec::new();
        
        for (i, line) in lines.iter().enumerate() {
            if line.trim().is_empty() || line.trim().starts_with("'") {
                result.push_str(&format!("{line}\n"));
                continue;
            }
            
            let case_label = i + 1;
            cases.push(format!("        Case {case_label}\n            {line}\n            {state_var} = {next_case}", 
                next_case = if i == lines.len() - 1 { 0 } else { case_label + 1 }
            ));
        }
        
        result.push_str(&format!("Do While True\n    Select Case {state_var}\n{}\n    End Select\n    Exit Do\nLoop", cases.join("\n")));
        
        result
    }

    /// 生成随机字符串
    fn generate_random_string(rng: &mut impl Rng, length: usize) -> String {
        std::iter::repeat(())
            .map(|_| rng.sample(Alphanumeric))
            .map(char::from)
            .take(length)
            .collect()
    }

    /// 字符串加密混淆
    pub fn encrypt_strings(code: &str, language: ScriptLanguage) -> String {
        let mut rng = rand::thread_rng();
        let result = code.to_string();
        
        match language {
            ScriptLanguage::PHP => {
                // 简单的字符串加密示例
                let key = rng.gen_range(1..255);
                let encrypted = result.chars().map(|c| (c as u8 ^ key) as char).collect::<String>();
                format!("<?php $k={}; eval(str_rot13('{}')); ?>", key, encrypted)
            },
            _ => result,
        }
    }
}
