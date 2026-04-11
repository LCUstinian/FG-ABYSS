use ring::{signature, digest};
use std::fs::File;
use std::io::Read;
use std::path::Path;

/// 插件签名验证错误
#[derive(Debug, PartialEq, Eq)]
pub enum PluginSignatureError {
    FileReadError,
    InvalidSignature,
    InvalidPublicKey,
}

/// 验证插件签名
pub fn verify_plugin_signature(plugin_path: &Path, signature_path: &Path, public_key: &[u8]) -> Result<(), PluginSignatureError> {
    // 读取插件文件
    let mut plugin_file = File::open(plugin_path).map_err(|_| PluginSignatureError::FileReadError)?;
    let mut plugin_content = Vec::new();
    plugin_file.read_to_end(&mut plugin_content).map_err(|_| PluginSignatureError::FileReadError)?;

    // 读取签名文件
    let mut signature_file = File::open(signature_path).map_err(|_| PluginSignatureError::FileReadError)?;
    let mut signature = Vec::new();
    signature_file.read_to_end(&mut signature).map_err(|_| PluginSignatureError::FileReadError)?;

    // 验证签名
    let public_key_der = signature::UnparsedPublicKey::new(
        &signature::ED25519,
        public_key,
    );

    public_key_der
        .verify(&plugin_content, &signature)
        .map_err(|_| PluginSignatureError::InvalidSignature)?;

    Ok(())
}

/// 生成插件哈希
pub fn generate_plugin_hash(plugin_path: &Path) -> Result<Vec<u8>, PluginSignatureError> {
    let mut file = File::open(plugin_path).map_err(|_| PluginSignatureError::FileReadError)?;
    let mut content = Vec::new();
    file.read_to_end(&mut content).map_err(|_| PluginSignatureError::FileReadError)?;

    let mut hasher = digest::Context::new(&digest::SHA256);
    hasher.update(&content);
    let result = hasher.finish();
    Ok(result.as_ref().to_vec())
}
