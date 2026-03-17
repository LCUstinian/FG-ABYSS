package handlers

import (
	"context"
	"time"

	"fg-abyss/internal/app/services"
)

// EncryptionHandler 加密处理器
type EncryptionHandler struct {
	service *services.TrafficEncryption
}

// NewEncryptionHandler 创建加密处理器
func NewEncryptionHandler() *EncryptionHandler {
	return &EncryptionHandler{
		service: nil, // 需要配置后初始化
	}
}

// EncryptionConfig 加密配置请求
type EncryptionConfig struct {
	// Key 加密密钥（16 进制）
	Key string `json:"key"`
	// IV 初始化向量（16 进制）
	IV string `json:"iv"`
	// Signature 签名密钥（16 进制）
	Signature string `json:"signature"`
}

// EncryptionResponse 加密操作响应
type EncryptionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// InitEncryption 初始化加密服务
func (h *EncryptionHandler) InitEncryption(ctx context.Context, config *EncryptionConfig) (*EncryptionResponse, error) {
	encConfig := &services.EncryptionConfig{
		Key:       config.Key,
		IV:        config.IV,
		Signature: config.Signature,
	}

	if err := services.ValidateEncryptionConfig(encConfig); err != nil {
		return &EncryptionResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	service, err := services.NewTrafficEncryption(encConfig)
	if err != nil {
		return &EncryptionResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	h.service = service

	return &EncryptionResponse{
		Success: true,
		Message: "Encryption initialized successfully",
	}, nil
}

// GenerateConfigRequest 生成配置请求
type GenerateConfigRequest struct {
	// 空请求，仅用于触发配置生成
}

// GenerateConfigResponse 生成配置响应
type GenerateConfigResponse struct {
	// Key 加密密钥
	Key string `json:"key"`
	// IV 初始化向量
	IV string `json:"iv"`
	// Signature 签名密钥
	Signature string `json:"signature"`
}

// GenerateConfig 生成加密配置
func (h *EncryptionHandler) GenerateConfig(ctx context.Context, req *GenerateConfigRequest) (*GenerateConfigResponse, error) {
	config, err := services.GenerateEncryptionConfig()
	if err != nil {
		return &GenerateConfigResponse{}, nil
	}

	return &GenerateConfigResponse{
		Key:       config.Key,
		IV:        config.IV,
		Signature: config.Signature,
	}, nil
}

// EncryptRequestRequest 加密请求请求
type EncryptRequestRequest struct {
	// Data 要加密的数据（Base64）
	Data string `json:"data"`
}

// EncryptRequestResponse 加密请求响应
type EncryptRequestResponse struct {
	// EncryptedData 加密后的数据
	EncryptedData string `json:"encrypted_data"`
	// Timestamp 时间戳
	Timestamp int64 `json:"timestamp"`
	// Signature 签名
	Signature string `json:"signature"`
}

// EncryptRequest 加密数据
func (h *EncryptionHandler) EncryptRequest(ctx context.Context, req *EncryptRequestRequest) (*EncryptRequestResponse, error) {
	if h.service == nil {
		return &EncryptRequestResponse{}, nil
	}

	data := []byte(req.Data)
	encryptedReq, err := h.service.EncryptRequest(data)
	if err != nil {
		return &EncryptRequestResponse{}, nil
	}

	return &EncryptRequestResponse{
		EncryptedData: encryptedReq.Data,
		Timestamp:     encryptedReq.Timestamp,
		Signature:     encryptedReq.Signature,
	}, nil
}

// DecryptRequestRequest 解密请求请求
type DecryptRequestRequest struct {
	// EncryptedData 加密数据
	EncryptedData string `json:"encrypted_data"`
	// Timestamp 时间戳
	Timestamp int64 `json:"timestamp"`
	// Signature 签名
	Signature string `json:"signature"`
	// Window 时间窗口（秒）
	Window int64 `json:"window"`
}

// DecryptRequestResponse 解密请求响应
type DecryptRequestResponse struct {
	// Data 解密后的数据
	Data string `json:"data"`
	// Valid 是否有效
	Valid bool `json:"valid"`
	// Message 消息
	Message string `json:"message"`
}

// DecryptRequest 解密数据
func (h *EncryptionHandler) DecryptRequest(ctx context.Context, req *DecryptRequestRequest) (*DecryptRequestResponse, error) {
	if h.service == nil {
		return &DecryptRequestResponse{
			Valid:   false,
			Message: "Encryption service not initialized",
		}, nil
	}

	encryptedReq := &services.EncryptedRequest{
		Data:      req.EncryptedData,
		Timestamp: req.Timestamp,
		Signature: req.Signature,
	}

	window := time.Duration(req.Window) * time.Second
	data, err := h.service.DecryptRequest(encryptedReq, window)
	if err != nil {
		return &DecryptRequestResponse{
			Valid:   false,
			Message: err.Error(),
			Data:    "",
		}, nil
	}

	return &DecryptRequestResponse{
		Data:    string(data),
		Valid:   true,
		Message: "Decryption successful",
	}, nil
}

// ValidateConfigRequest 验证配置请求
type ValidateConfigRequest struct {
	// Key 加密密钥
	Key string `json:"key"`
	// IV 初始化向量
	IV string `json:"iv"`
	// Signature 签名密钥
	Signature string `json:"signature"`
}

// ValidateConfigResponse 验证配置响应
type ValidateConfigResponse struct {
	// Valid 是否有效
	Valid bool `json:"valid"`
	// Errors 错误列表
	Errors []string `json:"errors"`
}

// ValidateConfig 验证加密配置
func (h *EncryptionHandler) ValidateConfig(ctx context.Context, req *ValidateConfigRequest) (*ValidateConfigResponse, error) {
	config := &services.EncryptionConfig{
		Key:       req.Key,
		IV:        req.IV,
		Signature: req.Signature,
	}

	err := services.ValidateEncryptionConfig(config)
	if err != nil {
		return &ValidateConfigResponse{
			Valid:  false,
			Errors: []string{err.Error()},
		}, nil
	}

	return &ValidateConfigResponse{
		Valid:  true,
		Errors: []string{},
	}, nil
}

// TestEncryptionRequest 测试加密请求
type TestEncryptionRequest struct {
	// Key 加密密钥
	Key string `json:"key"`
	// IV 初始化向量
	IV string `json:"iv"`
	// Signature 签名密钥
	Signature string `json:"signature"`
	// TestData 测试数据
	TestData string `json:"test_data"`
}

// TestEncryptionResponse 测试加密响应
type TestEncryptionResponse struct {
	// Success 是否成功
	Success bool `json:"success"`
	// Message 消息
	Message string `json:"message"`
	// EncryptedData 加密后的数据
	EncryptedData string `json:"encrypted_data"`
	// DecryptedData 解密后的数据
	DecryptedData string `json:"decrypted_data"`
}

// TestEncryption 测试加密解密
func (h *EncryptionHandler) TestEncryption(ctx context.Context, req *TestEncryptionRequest) (*TestEncryptionResponse, error) {
	config := &services.EncryptionConfig{
		Key:       req.Key,
		IV:        req.IV,
		Signature: req.Signature,
	}

	if err := services.ValidateEncryptionConfig(config); err != nil {
		return &TestEncryptionResponse{
			Success: false,
			Message: "Invalid configuration: " + err.Error(),
		}, nil
	}

	service, err := services.NewTrafficEncryption(config)
	if err != nil {
		return &TestEncryptionResponse{
			Success: false,
			Message: "Failed to initialize: " + err.Error(),
		}, nil
	}

	// 测试加密
	testData := []byte(req.TestData)
	encryptedReq, err := service.EncryptRequest(testData)
	if err != nil {
		return &TestEncryptionResponse{
			Success: false,
			Message: "Encryption failed: " + err.Error(),
		}, nil
	}

	// 测试解密
	window := 5 * time.Minute
	decrypted, err := service.DecryptRequest(encryptedReq, window)
	if err != nil {
		return &TestEncryptionResponse{
			Success: false,
			Message: "Decryption failed: " + err.Error(),
		}, nil
	}

	// 验证数据
	if string(decrypted) != req.TestData {
		return &TestEncryptionResponse{
			Success:       false,
			Message:       "Data mismatch after decryption",
			EncryptedData: encryptedReq.Data,
			DecryptedData: string(decrypted),
		}, nil
	}

	return &TestEncryptionResponse{
		Success:       true,
		Message:       "Encryption/Decryption test successful",
		EncryptedData: encryptedReq.Data,
		DecryptedData: string(decrypted),
	}, nil
}
