package handlers

import (
	"context"

	"fg-abyss/internal/app/services"
)

// ProxyHandler 代理处理器
type ProxyHandler struct {
	service *services.ProxyService
}

// NewProxyHandler 创建代理处理器
func NewProxyHandler() *ProxyHandler {
	return &ProxyHandler{
		service: services.NewProxyService(),
	}
}

// ProxyConfig 代理配置请求
type ProxyConfig struct {
	// Type 代理类型（http/https/socks5）
	Type string `json:"type"`
	// Host 代理主机
	Host string `json:"host"`
	// Port 代理端口
	Port int `json:"port"`
	// Username 用户名（可选）
	Username string `json:"username,omitempty"`
	// Password 密码（可选）
	Password string `json:"password,omitempty"`
	// Timeout 超时时间（秒）
	Timeout int `json:"timeout,omitempty"`
}

// ProxyResponse 代理操作响应
type ProxyResponse struct {
	// Success 是否成功
	Success bool `json:"success"`
	// Message 消息
	Message string `json:"message"`
}

// SetProxy 设置代理
func (h *ProxyHandler) SetProxy(ctx context.Context, config *ProxyConfig) (*ProxyResponse, error) {
	proxyConfig := &services.ProxyConfig{
		Type:     services.ProxyType(config.Type),
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		Timeout:  config.Timeout,
	}

	if err := h.service.SetProxy(proxyConfig); err != nil {
		return &ProxyResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &ProxyResponse{
		Success: true,
		Message: "Proxy configured successfully",
	}, nil
}

// TestProxyRequest 测试代理请求
type TestProxyRequest struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  int    `json:"timeout"`
}

// TestProxyResponse 测试代理响应
type TestProxyResponse struct {
	// Success 是否成功
	Success bool   `json:"success"`
	// Message 消息
	Message string `json:"message"`
}

// TestProxy 测试代理连接
func (h *ProxyHandler) TestProxy(ctx context.Context, req *TestProxyRequest) (*TestProxyResponse, error) {
	config := &services.ProxyConfig{
		Type:     services.ProxyType(req.Type),
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Timeout:  req.Timeout,
	}

	success, message, err := h.service.TestProxy(config)
	if err != nil {
		return &TestProxyResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &TestProxyResponse{
		Success: success,
		Message: message,
	}, nil
}

// GetProxyStatusResponse 获取代理状态响应
type GetProxyStatusResponse struct {
	// Enabled 是否启用
	Enabled bool `json:"enabled"`
	// Type 代理类型
	Type string `json:"type"`
	// Host 代理主机
	Host string `json:"host"`
	// Port 代理端口
	Port int `json:"port"`
	// HasAuth 是否有认证
	HasAuth bool `json:"has_auth"`
}

// GetProxyStatus 获取代理状态
func (h *ProxyHandler) GetProxyStatus(ctx context.Context) (*GetProxyStatusResponse, error) {
	stats := h.service.GetProxyStats()

	return &GetProxyStatusResponse{
		Enabled: stats["enabled"].(bool),
		Type:    getString(stats["type"]),
		Host:    getString(stats["host"]),
		Port:    getInt(stats["port"]),
		HasAuth: stats["has_auth"].(bool),
	}, nil
}

// DisableProxy 禁用代理
func (h *ProxyHandler) DisableProxy(ctx context.Context) (*ProxyResponse, error) {
	h.service.ClearProxy()

	return &ProxyResponse{
		Success: true,
		Message: "Proxy disabled",
	}, nil
}

// ValidateProxyRequest 验证代理配置请求
type ValidateProxyRequest struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  int    `json:"timeout"`
}

// ValidateProxyResponse 验证代理配置响应
type ValidateProxyResponse struct {
	// Valid 是否有效
	Valid bool `json:"valid"`
	// Errors 错误列表
	Errors []string `json:"errors"`
}

// ValidateProxy 验证代理配置
func (h *ProxyHandler) ValidateProxy(ctx context.Context, req *ValidateProxyRequest) (*ValidateProxyResponse, error) {
	config := &services.ProxyConfig{
		Type:     services.ProxyType(req.Type),
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Timeout:  req.Timeout,
	}

	err := h.service.ValidateProxyConfig(config)
	if err != nil {
		return &ValidateProxyResponse{
			Valid:  false,
			Errors: []string{err.Error()},
		}, nil
	}

	return &ValidateProxyResponse{
		Valid:  true,
		Errors: []string{},
	}, nil
}

// 辅助函数
func getString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func getInt(v interface{}) int {
	if v == nil {
		return 0
	}
	if i, ok := v.(int); ok {
		return i
	}
	return 0
}
