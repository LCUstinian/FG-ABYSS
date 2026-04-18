use aes_gcm::{Aes256Gcm, KeyInit, aead::{Aead, AeadCore, OsRng as AeadOsRng}};
use argon2::{Argon2, PasswordHash, PasswordHasher, PasswordVerifier};
use argon2::password_hash::{rand_core::OsRng, SaltString};
use base64::{Engine as _, engine::general_purpose::STANDARD as BASE64};
use rand::RngCore;
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use zeroize::Zeroize;
use crate::{AppError, Result};

// --- Low-level functions ---

/// Encrypt plaintext with AES-256-GCM. Returns base64(nonce[12] || tag[16] || ciphertext).
pub fn encrypt(plain: &[u8], key: &[u8; 32]) -> Result<String> {
    let cipher = Aes256Gcm::new_from_slice(key)
        .map_err(|e| AppError::Crypto(e.to_string()))?;
    let nonce = Aes256Gcm::generate_nonce(&mut AeadOsRng);
    let mut ciphertext = cipher.encrypt(&nonce, plain)
        .map_err(|e| AppError::Crypto(e.to_string()))?;
    let mut combined = nonce.to_vec();
    combined.append(&mut ciphertext);
    Ok(BASE64.encode(&combined))
}

/// Decrypt base64(nonce[12] || tag[16] || ciphertext) with AES-256-GCM.
pub fn decrypt(cipher_b64: &str, key: &[u8; 32]) -> Result<Vec<u8>> {
    let combined = BASE64.decode(cipher_b64)
        .map_err(|e| AppError::Crypto(e.to_string()))?;
    if combined.len() < 28 {
        return Err(AppError::Crypto("ciphertext too short".into()));
    }
    let (nonce_bytes, ciphertext) = combined.split_at(12);
    let cipher = Aes256Gcm::new_from_slice(key)
        .map_err(|e| AppError::Crypto(e.to_string()))?;
    let nonce = aes_gcm::Nonce::from_slice(nonce_bytes);
    cipher.decrypt(nonce, ciphertext)
        .map_err(|e| AppError::Crypto(e.to_string()))
}

/// Derive 32-byte master key from salt using Argon2id (does NOT store key, caller zeroizes).
pub fn derive_key(salt_b64: &str) -> Result<[u8; 32]> {
    let salt_bytes = BASE64.decode(salt_b64)
        .map_err(|e| AppError::Crypto(e.to_string()))?;
    const APP_SECRET: &[u8] = b"fg-abyss-master-key-v1";
    let mut key = [0u8; 32];
    Argon2::default()
        .hash_password_into(APP_SECRET, &salt_bytes, &mut key)
        .map_err(|e| AppError::Crypto(e.to_string()))?;
    Ok(key)
}

/// Generate a random base64-encoded 32-byte salt.
pub fn generate_salt() -> String {
    let mut salt = [0u8; 32];
    rand::thread_rng().fill_bytes(&mut salt);
    BASE64.encode(&salt)
}

/// Generate a random ASCII string of given length (used for response_mark).
pub fn random_ascii(len: usize) -> String {
    use rand::Rng;
    const CHARS: &[u8] = b"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    let mut rng = rand::thread_rng();
    (0..len).map(|_| CHARS[rng.gen_range(0..CHARS.len())] as char).collect()
}

// --- CryptoContext (Arc-wrapped, zeroize on drop) ---

pub struct CryptoContext {
    master_key: [u8; 32],
}

impl Drop for CryptoContext {
    fn drop(&mut self) {
        self.master_key.zeroize();
    }
}

impl CryptoContext {
    pub fn new(key: [u8; 32]) -> Arc<Self> {
        Arc::new(Self { master_key: key })
    }

    pub fn encrypt(&self, plain: &[u8]) -> Result<String> {
        encrypt(plain, &self.master_key)
    }

    pub fn decrypt(&self, cipher: &str) -> Result<Vec<u8>> {
        decrypt(cipher, &self.master_key)
    }

    pub fn decrypt_str(&self, cipher: &str) -> Result<String> {
        String::from_utf8(self.decrypt(cipher)?)
            .map_err(|e| AppError::Crypto(e.to_string()))
    }

    pub fn verify_password(&self, input: &str, hash: &str) -> bool {
        let Ok(parsed_hash) = PasswordHash::new(hash) else { return false; };
        Argon2::default()
            .verify_password(input.as_bytes(), &parsed_hash)
            .is_ok()
    }

    pub fn hash_password(password: &str) -> Result<String> {
        let salt = SaltString::generate(&mut OsRng);
        Argon2::default()
            .hash_password(password.as_bytes(), &salt)
            .map(|h| h.to_string())
            .map_err(|e| AppError::Crypto(e.to_string()))
    }
}

// --- CryptoChain ---

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize, Default)]
#[serde(rename_all = "snake_case", tag = "type")]
pub enum CodecStep {
    #[default]
    Base64,
    Aes256Gcm,
    XorKey { key: String },
    UrlEncode,
    GzipCompress,
    HexEncode,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize, Default)]
pub struct CryptoChain {
    pub steps: Vec<CodecStep>,
}

impl CryptoChain {
    pub fn encode(&self, _data: &[u8], _key: &[u8; 32]) -> Result<Vec<u8>> {
        todo!("CryptoChain::encode — implement per-step encoding in Phase 2")
    }
    pub fn decode(&self, _data: &[u8], _key: &[u8; 32]) -> Result<Vec<u8>> {
        todo!("CryptoChain::decode — implement per-step decoding in Phase 2")
    }
}

// --- Sensitive<T> wrapper ---

pub struct Sensitive<T>(T);

impl<T> Sensitive<T> {
    pub fn new(val: T) -> Self { Self(val) }
    pub fn inner(&self) -> &T { &self.0 }
    pub fn into_inner(self) -> T { self.0 }
}

impl<T> std::fmt::Debug for Sensitive<T> {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        f.write_str("[REDACTED]")
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_encrypt_decrypt_roundtrip() {
        let key = [42u8; 32];
        let plain = b"hello, FG-ABYSS!";
        let cipher = encrypt(plain, &key).unwrap();
        let decrypted = decrypt(&cipher, &key).unwrap();
        assert_eq!(decrypted, plain);
    }

    #[test]
    fn test_different_keys_fail() {
        let key1 = [1u8; 32];
        let key2 = [2u8; 32];
        let cipher = encrypt(b"secret", &key1).unwrap();
        assert!(decrypt(&cipher, &key2).is_err());
    }

    #[test]
    fn test_generate_salt_unique() {
        let s1 = generate_salt();
        let s2 = generate_salt();
        assert_ne!(s1, s2);
        assert!(!s1.is_empty());
    }

    #[test]
    fn test_derive_key_deterministic() {
        let salt = generate_salt();
        let k1 = derive_key(&salt).unwrap();
        let k2 = derive_key(&salt).unwrap();
        assert_eq!(k1, k2);
    }

    #[test]
    fn test_crypto_context_encrypt_decrypt() {
        let key = [7u8; 32];
        let ctx = CryptoContext::new(key);
        let cipher = ctx.encrypt(b"password123").unwrap();
        let plain = ctx.decrypt_str(&cipher).unwrap();
        assert_eq!(plain, "password123");
    }

    #[test]
    fn test_sensitive_debug_redacted() {
        let s = Sensitive::new("my_secret");
        assert_eq!(format!("{:?}", s), "[REDACTED]");
    }
}
