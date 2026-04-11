use crate::core::crypto::{encrypt_data, decrypt_data, EncryptionType};

use std::path::Path;

/// 测试加密解密功能
pub fn test_encryption() {
    println!("Testing encryption...");
    
    let test_data = b"Hello, FG-ABYSS!";
    let password = "test-password";
    
    // 测试AES-256-GCM加密
    let (ciphertext, nonce, tag) = encrypt_data(test_data, password, &EncryptionType::AES256GCM).unwrap();
    let plaintext = decrypt_data(&ciphertext, nonce.as_deref(), tag.as_deref(), password, &EncryptionType::AES256GCM).unwrap();
    assert_eq!(test_data, plaintext.as_slice(), "AES-256-GCM encryption/decryption failed");
    println!("✓ AES-256-GCM encryption/decryption test passed");
    
    // 测试ChaCha20-Poly1305加密
    let (ciphertext, nonce, tag) = encrypt_data(test_data, password, &EncryptionType::ChaCha20Poly1305).unwrap();
    let plaintext = decrypt_data(&ciphertext, nonce.as_deref(), tag.as_deref(), password, &EncryptionType::ChaCha20Poly1305).unwrap();
    assert_eq!(test_data, plaintext.as_slice(), "ChaCha20-Poly1305 encryption/decryption failed");
    println!("✓ ChaCha20-Poly1305 encryption/decryption test passed");
    
    println!("All encryption tests passed!");
}

/// 测试插件签名验证功能
pub fn test_plugin_signature() {
    println!("Testing plugin signature verification...");
    
    // 注意：这里需要实际的插件文件和签名文件来测试
    // 目前只是一个框架，实际测试需要准备测试文件
    let _plugin_path = Path::new("test-plugin.dll");
    let _signature_path = Path::new("test-plugin.sig");
    let _public_key = b"\x30\x2a\x30\x05\x06\x03\x2b\x65\x70\x03\x21\x00\x12\x34\x56\x78\x90\xab\xcd\xef\x12\x34\x56\x78\x90\xab\xcd\xef\x12\x34\x56\x78\x90\xab\xcd\xef";
    
    // 这里只是测试函数调用，实际测试需要真实的文件
    // let result = verify_plugin_signature(plugin_path, signature_path, public_key);
    // assert!(result.is_ok(), "Plugin signature verification failed: {:?}", result);
    
    println!("✓ Plugin signature verification test framework ready");
    println!("Note: Actual plugin signature verification requires real test files");
}

/// 运行所有安全测试
pub fn run_security_tests() {
    println!("Running security tests...");
    println!("==============================");
    
    test_encryption();
    test_plugin_signature();
    
    println!("==============================");
    println!("All security tests completed!");
}
