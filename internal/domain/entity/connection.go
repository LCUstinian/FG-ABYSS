package entity

import (
	"fmt"
	"time"
)

// generateUUID 生成简单的 UUID（用于开发）
func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// ConnectionStatus 连接状态
type ConnectionStatus string

const (
	// ConnectionStatusConnected 已连接
	ConnectionStatusConnected ConnectionStatus = "connected"
	// ConnectionStatusDisconnected 已断开
	ConnectionStatusDisconnected ConnectionStatus = "disconnected"
	// ConnectionStatusConnecting 连接中
	ConnectionStatusConnecting ConnectionStatus = "connecting"
	// ConnectionStatusError 连接错误
	ConnectionStatusError ConnectionStatus = "error"
)

// Connection WebShell 连接实体
type Connection struct {
	// ID 连接唯一标识
	ID string `json:"id"`
	// WebShellID 关联的 WebShell ID
	WebShellID string `json:"webshell_id"`
	// Status 连接状态
	Status ConnectionStatus `json:"status"`
	// LastActiveTime 最后活跃时间
	LastActiveTime time.Time `json:"last_active_time"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
	// Config 连接配置
	Config *ConnectionConfig `json:"config,omitempty"`
	// Stats 连接统计
	Stats *ConnectionStats `json:"stats,omitempty"`
}

// ConnectionConfig 连接配置
type ConnectionConfig struct {
	// URL 目标 URL
	URL string `json:"url"`
	// Password 连接密码
	Password string `json:"password,omitempty"`
	// Encoder 编码器类型（base64, rot13, xor, none）
	Encoder string `json:"encoder"`
	// EncryptionKey 加密密钥（用于 XOR 等）
	EncryptionKey string `json:"encryption_key,omitempty"`
	// Timeout 超时时间（秒）
	Timeout int `json:"timeout"`
	// RetryCount 重试次数
	RetryCount int `json:"retry_count"`
	// RetryDelay 重试延迟（毫秒）
	RetryDelay int `json:"retry_delay"`
	// Headers 自定义请求头
	Headers map[string]string `json:"headers,omitempty"`
	// Cookies Cookie 列表
	Cookies map[string]string `json:"cookies,omitempty"`
	// ProxyType 代理类型（none, http, https, socks5）
	ProxyType string `json:"proxy_type"`
	// ProxyAddress 代理地址
	ProxyAddress string `json:"proxy_address,omitempty"`
	// SSLVerify 是否验证 SSL 证书
	SSLVerify bool `json:"ssl_verify"`
}

// ConnectionStats 连接统计
type ConnectionStats struct {
	// TotalRequests 总请求数
	TotalRequests int64 `json:"total_requests"`
	// SuccessfulRequests 成功请求数
	SuccessfulRequests int64 `json:"successful_requests"`
	// FailedRequests 失败请求数
	FailedRequests int64 `json:"failed_requests"`
	// TotalBytesSent 发送总字节数
	TotalBytesSent int64 `json:"total_bytes_sent"`
	// TotalBytesReceived 接收总字节数
	TotalBytesReceived int64 `json:"total_bytes_received"`
	// AverageResponseTime 平均响应时间（毫秒）
	AverageResponseTime float64 `json:"average_response_time"`
	// LastRequestTime 最后请求时间
	LastRequestTime time.Time `json:"last_request_time"`
}

// NewConnection 创建新的连接
func NewConnection(webshellID string, config *ConnectionConfig) *Connection {
	now := time.Now()
	return &Connection{
		ID:             generateUUID(),
		WebShellID:     webshellID,
		Status:         ConnectionStatusDisconnected,
		LastActiveTime: now,
		CreatedAt:      now,
		UpdatedAt:      now,
		Config:         config,
		Stats: &ConnectionStats{
			TotalRequests:       0,
			SuccessfulRequests:  0,
			FailedRequests:      0,
			TotalBytesSent:      0,
			TotalBytesReceived:  0,
			AverageResponseTime: 0,
		},
	}
}

// UpdateStatus 更新连接状态
func (c *Connection) UpdateStatus(status ConnectionStatus) {
	c.Status = status
	c.LastActiveTime = time.Now()
	c.UpdatedAt = time.Now()
}

// RecordRequest 记录请求统计
func (c *Connection) RecordRequest(success bool, bytesSent, bytesReceived int64, responseTime float64) {
	c.Stats.TotalRequests++
	if success {
		c.Stats.SuccessfulRequests++
	} else {
		c.Stats.FailedRequests++
	}
	c.Stats.TotalBytesSent += bytesSent
	c.Stats.TotalBytesReceived += bytesReceived
	
	// 更新平均响应时间
	totalTime := c.Stats.AverageResponseTime * float64(c.Stats.TotalRequests-1)
	c.Stats.AverageResponseTime = (totalTime + responseTime) / float64(c.Stats.TotalRequests)
	c.Stats.LastRequestTime = time.Now()
}
