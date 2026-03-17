package handlers

import (
	"context"

	"fg-abyss/internal/app/services"
	"fg-abyss/internal/domain/entity"
)

// PayloadHandler Payload 处理器
type PayloadHandler struct {
	generator *services.PayloadGenerator
}

// NewPayloadHandler 创建 Payload 处理器
func NewPayloadHandler() *PayloadHandler {
	return &PayloadHandler{
		generator: services.NewPayloadGenerator(),
	}
}

// GenerateRequest Payload 生成请求
type GenerateRequest struct {
	// Type Payload 类型（php/asp/aspx/jsp）
	Type string `json:"type"`
	// Function 功能类型（basic/file/command/database/full）
	Function string `json:"function"`
	// Password 连接密码
	Password string `json:"password"`
	// Encoder 编码器类型
	Encoder string `json:"encoder"`
	// EncryptionKey 加密密钥
	EncryptionKey string `json:"encryption_key"`
	// ObfuscationLevel 混淆级别（none/low/medium/high）
	ObfuscationLevel string `json:"obfuscation_level"`
	// OutputFilename 输出文件名
	OutputFilename string `json:"output_filename"`
	// TemplateName 模板名称
	TemplateName string `json:"template_name"`
}

// GenerateResponse Payload 生成响应
type GenerateResponse struct {
	// Success 是否成功
	Success bool `json:"success"`
	// Content Payload 内容
	Content string `json:"content"`
	// Filename 文件名
	Filename string `json:"filename"`
	// Size 文件大小
	Size int `json:"size"`
	// Message 消息
	Message string `json:"message"`
	// Warnings 警告信息
	Warnings []string `json:"warnings,omitempty"`
}

// Generate 生成 Payload
// @summary 生成 WebShell Payload
// @param ctx 上下文
// @param req 生成请求
// @return 生成结果
func (h *PayloadHandler) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	config := &entity.PayloadConfig{
		Type:             entity.PayloadType(req.Type),
		Function:         entity.PayloadFunction(req.Function),
		Password:         req.Password,
		Encoder:          req.Encoder,
		EncryptionKey:    req.EncryptionKey,
		ObfuscationLevel: entity.ObfuscationLevel(req.ObfuscationLevel),
		OutputFilename:   req.OutputFilename,
	}
	
	result, err := h.generator.Generate(config)
	if err != nil {
		return &GenerateResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	
	// 验证 Payload
	valid, warnings := h.generator.ValidatePayload(result.Content, config.Type)
	if !valid {
		// 添加警告但不阻止生成
		result.Warnings = warnings
	}
	
	return &GenerateResponse{
		Success:  result.Success,
		Content:  result.Content,
		Filename: result.Filename,
		Size:     result.Size,
		Message:  result.Message,
		Warnings: result.Warnings,
	}, nil
}

// GetTemplates 获取所有模板
// @summary 获取所有可用的 Payload 模板
// @return 模板列表
func (h *PayloadHandler) GetTemplates(ctx context.Context) ([]*entity.PayloadTemplate, error) {
	return h.generator.GetTemplates(), nil
}

// GetTemplateByName 按名称获取模板
// @summary 按名称获取模板
// @param name 模板名称
// @return 模板信息
func (h *PayloadHandler) GetTemplateByName(ctx context.Context, name string) (*entity.PayloadTemplate, error) {
	return h.generator.GetTemplateByName(name)
}

// Preview 预览 Payload
// @summary 预览生成的 Payload
// @param ctx 上下文
// @param req 生成请求
// @return Payload 内容
func (h *PayloadHandler) Preview(ctx context.Context, req *GenerateRequest) (string, error) {
	config := &entity.PayloadConfig{
		Type:             entity.PayloadType(req.Type),
		Function:         entity.PayloadFunction(req.Function),
		Password:         req.Password,
		Encoder:          req.Encoder,
		EncryptionKey:    req.EncryptionKey,
		ObfuscationLevel: entity.ObfuscationLevel(req.ObfuscationLevel),
	}
	
	return h.generator.Preview(config)
}

// ValidatePayload 验证 Payload
// @summary 验证 Payload 的合法性
// @param content Payload 内容
// @param payloadType Payload 类型
// @return 是否有效和警告信息
func (h *PayloadHandler) ValidatePayload(ctx context.Context, content string, payloadType string) (*ValidationResult, error) {
	valid, warnings := h.generator.ValidatePayload(content, entity.PayloadType(payloadType))
	
	return &ValidationResult{
		Valid:    valid,
		Warnings: warnings,
	}, nil
}

// ValidationResult 验证结果
type ValidationResult struct {
	// Valid 是否有效
	Valid bool `json:"valid"`
	// Warnings 警告信息
	Warnings []string `json:"warnings,omitempty"`
}

// AddCustomTemplateRequest 添加模板请求
type AddCustomTemplateRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Function string `json:"function"`
	Content string `json:"content"`
}

// AddTemplate 添加自定义模板
// @summary 添加自定义 Payload 模板
// @param ctx 上下文
// @param req 添加请求
// @return 错误信息
func (h *PayloadHandler) AddTemplate(ctx context.Context, req *AddCustomTemplateRequest) error {
	template := &entity.PayloadTemplate{
		Name:     req.Name,
		Type:     entity.PayloadType(req.Type),
		Function: entity.PayloadFunction(req.Function),
		Content:  req.Content,
		IsBuiltin: false,
	}
	return h.generator.AddCustomTemplate(template)
}

// DeleteTemplate 删除自定义模板
// @summary 删除自定义模板
// @param ctx 上下文
// @param name 模板名称
// @return 错误信息
func (h *PayloadHandler) DeleteTemplate(ctx context.Context, name string) error {
	return h.generator.RemoveCustomTemplate(name)
}

// AddCustomTemplate 添加自定义模板（兼容旧接口）
// @summary 添加自定义 Payload 模板
// @param ctx 上下文
// @param template 模板信息
// @return 错误信息
func (h *PayloadHandler) AddCustomTemplate(ctx context.Context, template *entity.PayloadTemplate) error {
	return h.generator.AddCustomTemplate(template)
}

// RemoveCustomTemplate 移除自定义模板（兼容旧接口）
// @summary 移除自定义模板
// @param name 模板名称
// @return 错误信息
func (h *PayloadHandler) RemoveCustomTemplate(ctx context.Context, name string) error {
	return h.generator.RemoveCustomTemplate(name)
}

// GetTemplateStats 获取模板统计信息
// @summary 获取模板统计信息
// @return 统计信息
func (h *PayloadHandler) GetTemplateStats(ctx context.Context) (*TemplateStats, error) {
	templates := h.generator.GetTemplates()
	
	stats := &TemplateStats{
		Total: len(templates),
	}
	
	for _, tmpl := range templates {
		switch tmpl.Type {
		case entity.PayloadTypePHP:
			stats.PHPCount++
		case entity.PayloadTypeASP:
			stats.ASPCount++
		case entity.PayloadTypeASPX:
			stats.ASPXCount++
		case entity.PayloadTypeJSP:
			stats.JSPCount++
		}
		
		if tmpl.IsBuiltin {
			stats.BuiltinCount++
		} else {
			stats.CustomCount++
		}
	}
	
	return stats, nil
}

// TemplateStats 模板统计
type TemplateStats struct {
	// Total 模板总数
	Total int `json:"total"`
	// PHPCount PHP 模板数量
	PHPCount int `json:"php_count"`
	// ASPCount ASP 模板数量
	ASPCount int `json:"asp_count"`
	// ASPXCount ASPX 模板数量
	ASPXCount int `json:"aspx_count"`
	// JSPCount JSP 模板数量
	JSPCount int `json:"jsp_count"`
	// BuiltinCount 内置模板数量
	BuiltinCount int `json:"builtin_count"`
	// CustomCount 自定义模板数量
	CustomCount int `json:"custom_count"`
}
