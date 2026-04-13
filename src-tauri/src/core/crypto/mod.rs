use aes_gcm::{
    aead::{Aead, KeyInit, OsRng},
    Aes256Gcm, Nonce,
};
use rand::RngCore;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum CryptoError {
    #[error("加密失败：{0}")]
    EncryptionError(String),
    #[error("解密失败：{0}")]
    DecryptionError(String),
    #[error("密钥无效")]
    InvalidKey,
    #[error("密文格式错误")]
    InvalidCiphertext,
}

pub type Result<T> = std::result::Result<T, CryptoError>;

/// AES-256-GCM 加密器
pub struct AesCipher {
    key: [u8; 32],
}

impl AesCipher {
    /// 创建新的 AES 加密器
    pub fn new(key: [u8; 32]) -> Self {
        Self { key }
    }

    /// 生成随机密钥
    pub fn generate_key() -> [u8; 32] {
        let mut key = [0u8; 32];
        OsRng.fill_bytes(&mut key);
        key
    }

    /// 加密数据
    pub fn encrypt(&self, data: &[u8]) -> Result<Vec<u8>> {
        let cipher = Aes256Gcm::new_from_slice(&self.key)
            .map_err(|e| CryptoError::EncryptionError(e.to_string()))?;

        // 生成随机 nonce
        let mut nonce_bytes = [0u8; 12];
        OsRng.fill_bytes(&mut nonce_bytes);
        let nonce = Nonce::from_slice(&nonce_bytes);

        // 加密
        let ciphertext = cipher
            .encrypt(nonce, data)
            .map_err(|e| CryptoError::EncryptionError(e.to_string()))?;

        // 组合：nonce + ciphertext
        let mut result = Vec::new();
        result.extend_from_slice(&nonce_bytes);
        result.extend_from_slice(&ciphertext);

        Ok(result)
    }

    /// 解密数据
    pub fn decrypt(&self, encrypted: &[u8]) -> Result<Vec<u8>> {
        if encrypted.len() < 12 {
            return Err(CryptoError::InvalidCiphertext);
        }

        let nonce_bytes = &encrypted[..12];
        let ciphertext = &encrypted[12..];

        let cipher = Aes256Gcm::new_from_slice(&self.key)
            .map_err(|e| CryptoError::DecryptionError(e.to_string()))?;
        let nonce = Nonce::from_slice(nonce_bytes);

        let plaintext = cipher
            .decrypt(nonce, ciphertext)
            .map_err(|e| CryptoError::DecryptionError(e.to_string()))?;

        Ok(plaintext)
    }
}

/// XOR 加密器（轻量级）
pub struct XorCipher {
    key: Vec<u8>,
}

impl XorCipher {
    /// 创建新的 XOR 加密器
    pub fn new(key: Vec<u8>) -> Self {
        Self { key }
    }

    /// 生成随机密钥
    pub fn generate_key(length: usize) -> Vec<u8> {
        let mut key = vec![0u8; length];
        OsRng.fill_bytes(&mut key);
        key
    }

    /// 加密/解密（XOR 是对称的）
    pub fn process(&self, data: &[u8]) -> Vec<u8> {
        data.iter()
            .zip(self.key.iter().cycle())
            .map(|(&byte, &key_byte)| byte ^ key_byte)
            .collect()
    }

    pub fn encrypt(&self, data: &[u8]) -> Vec<u8> {
        self.process(data)
    }

    pub fn decrypt(&self, data: &[u8]) -> Vec<u8> {
        self.process(data)
    }
}

/// Base64 编码
pub mod base64 {
    use base64::{engine::general_purpose, Engine as _};

    pub fn encode(data: &[u8]) -> String {
        general_purpose::STANDARD.encode(data)
    }

    pub fn decode(s: &str) -> std::result::Result<Vec<u8>, base64::DecodeError> {
        general_purpose::STANDARD.decode(s)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_aes_encrypt_decrypt() {
        let key = AesCipher::generate_key();
        let cipher = AesCipher::new(key);
        let data = b"Hello, World!";

        let encrypted = cipher.encrypt(data).unwrap();
        let decrypted = cipher.decrypt(&encrypted).unwrap();

        assert_eq!(data.to_vec(), decrypted);
    }

    #[test]
    fn test_xor_encrypt_decrypt() {
        let key = XorCipher::generate_key(32);
        let cipher = XorCipher::new(key);
        let data = b"Hello, World!";

        let encrypted = cipher.encrypt(data);
        let decrypted = cipher.decrypt(&encrypted);

        assert_eq!(data.to_vec(), decrypted);
    }
}
