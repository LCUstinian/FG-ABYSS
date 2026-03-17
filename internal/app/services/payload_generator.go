package services

import (
	"fg-abyss/internal/domain/entity"
	"fmt"
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
