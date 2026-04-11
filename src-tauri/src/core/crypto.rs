use ring::aead;
use rand::Rng;
use base64::{Engine as _, engine::general_purpose};
use std::vec::Vec;

/// 加密算法类型
pub enum EncryptionType {
    None,
    XOR,
    AES256GCM,
    ChaCha20Poly1305,
}

/// 加密结果
pub struct EncryptionResult {
    pub ciphertext: Vec<u8>,
    pub nonce: Vec<u8>,
    pub tag: Vec<u8>,
}

/// 从密码派生密钥
/// 使用 HKDF 进行更安全的密钥派生
pub fn derive_key(password: &str) -> Vec<u8> {
    use ring::hkdf;
    let salt = b"fg-abyss-salt";
    let info = &[&b"fg-abyss-key"[..]];
    let password_bytes = password.as_bytes();
    
    let salt = hkdf::Salt::new(hkdf::HKDF_SHA256, salt);
    let prk = salt.extract(password_bytes);
    
    let mut key = vec![0u8; 32]; // 32 bytes for AES-256 and ChaCha20
    let okm = prk.expand(info, hkdf::HKDF_SHA256).unwrap();
    okm.fill(&mut key).unwrap();
    
    key
}

/// 生成随机 nonce
pub fn generate_nonce() -> Vec<u8> {
    let mut rng = rand::thread_rng();
    let mut nonce = vec![0u8; 12]; // AES-GCM 推荐使用 12 字节 nonce
    rng.fill(&mut nonce[..]);
    nonce
}

/// AES-256-GCM 加密
/// ring 会自动把 tag 追加到 ciphertext 末尾，我们需要分离出来
pub fn encrypt_aes256gcm(plaintext: &[u8], key: &[u8]) -> Result<EncryptionResult, String> {
    let nonce = generate_nonce();
    
    let unbound_key = aead::UnboundKey::new(&aead::AES_256_GCM, key)
        .map_err(|e| format!("Failed to create encryption key: {}", e))?;
    let sealing_key = aead::LessSafeKey::new(unbound_key);
    let aad = b"";
    
    let mut ciphertext = plaintext.to_vec();
    let mut nonce_array = [0u8; 12];
    nonce_array.copy_from_slice(&nonce[..]);
    let _ = sealing_key.seal_in_place_append_tag(aead::Nonce::assume_unique_for_key(nonce_array), aead::Aad::from(aad), &mut ciphertext);
    
    // Split tag from end (16 bytes)
    let ciphertext_len = ciphertext.len() - 16;
    let tag = ciphertext.split_off(ciphertext_len);
    
    Ok(EncryptionResult {
        ciphertext,
        nonce,
        tag,
    })
}

/// AES-256-GCM 解密
pub fn decrypt_aes256gcm(ciphertext: &[u8], nonce: &[u8], tag: &[u8], key: &[u8]) -> Result<Vec<u8>, String> {
    let unbound_key = aead::UnboundKey::new(&aead::AES_256_GCM, key)
        .map_err(|e| format!("Failed to create decryption key: {}", e))?;
    let opening_key = aead::LessSafeKey::new(unbound_key);
    let aad = b"";
    
    let mut ciphertext_with_tag = [ciphertext, tag].concat();
    let mut nonce_array = [0u8; 12];
    nonce_array.copy_from_slice(&nonce[..]);
    match opening_key.open_in_place(aead::Nonce::assume_unique_for_key(nonce_array), aead::Aad::from(aad), &mut ciphertext_with_tag) {
        Ok(_) => {
            let plaintext_len = ciphertext_with_tag.len() - 16;
            ciphertext_with_tag.truncate(plaintext_len);
            Ok(ciphertext_with_tag)
        }
        Err(_) => Err("解密失败：认证标签验证失败".to_string()),
    }
}

/// XOR 加密/解密
pub fn xor_encrypt(plaintext: &[u8], key: &[u8]) -> Vec<u8> {
    plaintext.iter().zip(key.iter().cycle()).map(|(p, k)| p ^ k).collect()
}

/// ChaCha20-Poly1305 加密
pub fn encrypt_chacha20poly1305(plaintext: &[u8], key: &[u8]) -> Result<EncryptionResult, String> {
    let nonce = generate_nonce();
    
    let unbound_key = aead::UnboundKey::new(&aead::CHACHA20_POLY1305, key)
        .map_err(|e| format!("Failed to create encryption key: {}", e))?;
    let sealing_key = aead::LessSafeKey::new(unbound_key);
    let aad = b"";
    
    let mut ciphertext = plaintext.to_vec();
    let mut nonce_array = [0u8; 12];
    nonce_array.copy_from_slice(&nonce[..]);
    let _ = sealing_key.seal_in_place_append_tag(aead::Nonce::assume_unique_for_key(nonce_array), aead::Aad::from(aad), &mut ciphertext);
    
    // Split tag from end (16 bytes)
    let ciphertext_len = ciphertext.len() - 16;
    let tag = ciphertext.split_off(ciphertext_len);
    
    Ok(EncryptionResult {
        ciphertext,
        nonce,
        tag,
    })
}

/// ChaCha20-Poly1305 解密
pub fn decrypt_chacha20poly1305(ciphertext: &[u8], nonce: &[u8], tag: &[u8], key: &[u8]) -> Result<Vec<u8>, String> {
    let unbound_key = aead::UnboundKey::new(&aead::CHACHA20_POLY1305, key)
        .map_err(|e| format!("Failed to create decryption key: {}", e))?;
    let opening_key = aead::LessSafeKey::new(unbound_key);
    let aad = b"";
    
    let mut ciphertext_with_tag = [ciphertext, tag].concat();
    let mut nonce_array = [0u8; 12];
    nonce_array.copy_from_slice(&nonce[..]);
    match opening_key.open_in_place(aead::Nonce::assume_unique_for_key(nonce_array), aead::Aad::from(aad), &mut ciphertext_with_tag) {
        Ok(_) => {
            let plaintext_len = ciphertext_with_tag.len() - 16;
            ciphertext_with_tag.truncate(plaintext_len);
            Ok(ciphertext_with_tag)
        }
        Err(_) => Err("解密失败：认证标签验证失败".to_string()),
    }
}

/// 加密数据
pub fn encrypt_data(data: &[u8], password: &str, encryption_type: &EncryptionType) -> Result<(Vec<u8>, Option<Vec<u8>>, Option<Vec<u8>>), String> {
    match encryption_type {
        EncryptionType::None => Ok((data.to_vec(), None, None)),
        EncryptionType::XOR => {
            let key = derive_key(password);
            let ciphertext = xor_encrypt(data, &key);
            Ok((ciphertext, None, None))
        },
        EncryptionType::AES256GCM => {
            let key = derive_key(password);
            let result = encrypt_aes256gcm(data, &key)?;
            Ok((result.ciphertext, Some(result.nonce), Some(result.tag)))
        },
        EncryptionType::ChaCha20Poly1305 => {
            let key = derive_key(password);
            let result = encrypt_chacha20poly1305(data, &key)?;
            Ok((result.ciphertext, Some(result.nonce), Some(result.tag)))
        },
    }
}

/// 解密数据
pub fn decrypt_data(ciphertext: &[u8], nonce: Option<&[u8]>, tag: Option<&[u8]>, password: &str, encryption_type: &EncryptionType) -> Result<Vec<u8>, String> {
    match encryption_type {
        EncryptionType::None => Ok(ciphertext.to_vec()),
        EncryptionType::XOR => {
            let key = derive_key(password);
            Ok(xor_encrypt(ciphertext, &key))
        },
        EncryptionType::AES256GCM => {
            let key = derive_key(password);
            let nonce = nonce.ok_or("缺少 nonce".to_string())?;
            let tag = tag.ok_or("缺少 tag".to_string())?;
            decrypt_aes256gcm(ciphertext, nonce, tag, &key)
        },
        EncryptionType::ChaCha20Poly1305 => {
            let key = derive_key(password);
            let nonce = nonce.ok_or("缺少 nonce".to_string())?;
            let tag = tag.ok_or("缺少 tag".to_string())?;
            decrypt_chacha20poly1305(ciphertext, nonce, tag, &key)
        },
    }
}

/// Base64 编码
pub fn base64_encode(data: &[u8]) -> String {
    general_purpose::STANDARD.encode(data)
}

/// Base64 解码
pub fn base64_decode(encoded: &str) -> Result<Vec<u8>, String> {
    general_purpose::STANDARD.decode(encoded).map_err(|e| e.to_string())
}
