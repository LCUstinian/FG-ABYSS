package services

import (
	"context"
	"time"

	"fg-abyss/internal/domain/entity"
)

// ConnectionService 连接服务
type ConnectionService struct {
	manager *ConnectionManager
}

// NewConnectionService 创建连接服务
func NewConnectionService() *ConnectionService {
	return &ConnectionService{
		manager: NewConnectionManager(),
	}
}

// ConnectRequest 连接请求
type ConnectRequest struct {
	// WebShellID WebShell ID
	WebShellID string
	// URL 目标 URL
	URL string
	// Password 密码
	Password string
	// Encoder 编码器类型
	Encoder string
	// EncryptionKey 加密密钥
	EncryptionKey string
	// Timeout 超时时间
	Timeout int
	// ProxyType 代理类型
	ProxyType string
	// ProxyAddress 代理地址
	ProxyAddress string
	// SSLVerify 是否验证 SSL
	SSLVerify bool
}

// ConnectResponse 连接响应
type ConnectResponse struct {
	// Success 是否成功
	Success bool
	// Message 消息
	Message string
	// Connection 连接信息
	Connection *entity.Connection
}

// Connect 建立连接
func (s *ConnectionService) Connect(ctx context.Context, req *ConnectRequest) (*ConnectResponse, error) {
	// 创建连接配置
	config := &entity.ConnectionConfig{
		URL:            req.URL,
		Password:       req.Password,
		Encoder:        req.Encoder,
		EncryptionKey:  req.EncryptionKey,
		Timeout:        req.Timeout,
		ProxyType:      req.ProxyType,
		ProxyAddress:   req.ProxyAddress,
		SSLVerify:      req.SSLVerify,
		RetryCount:     3,
		RetryDelay:     1000,
	}

	// 创建连接
	conn, err := s.manager.CreateConnection(req.WebShellID, config)
	if err != nil {
		return &ConnectResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	// 更新状态为连接中
	conn.UpdateStatus(entity.ConnectionStatusConnecting)

	// 测试连接
	testReq := &HTTPRequestConfig{
		URL:          config.URL,
		Method:       "POST",
		Body:         []byte("test"),
		Timeout:      config.Timeout,
		ProxyType:    config.ProxyType,
		ProxyAddress: config.ProxyAddress,
		SSLVerify:    config.SSLVerify,
	}

	resp, err := s.manager.GetHTTPEngine().Send(ctx, testReq)
	if err != nil {
		conn.UpdateStatus(entity.ConnectionStatusError)
		return &ConnectResponse{
			Success:    false,
			Message:    "连接失败：" + err.Error(),
			Connection: conn,
		}, nil
	}

	// 检查响应
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		conn.UpdateStatus(entity.ConnectionStatusConnected)
		return &ConnectResponse{
			Success:    true,
			Message:    "连接成功",
			Connection: conn,
		}, nil
	}

	conn.UpdateStatus(entity.ConnectionStatusError)
	return &ConnectResponse{
		Success:    false,
		Message:    "连接失败：HTTP " + string(rune(resp.StatusCode)),
		Connection: conn,
	}, nil
}

// Disconnect 断开连接
func (s *ConnectionService) Disconnect(webshellID string) error {
	s.manager.RemoveConnection(webshellID)
	return nil
}

// GetConnection 获取连接
func (s *ConnectionService) GetConnection(webshellID string) (*entity.Connection, error) {
	conn, exists := s.manager.GetConnection(webshellID)
	if !exists {
		return nil, ErrConnectionNotFound
	}
	return conn, nil
}

// GetAllConnections 获取所有连接
func (s *ConnectionService) GetAllConnections() []*entity.Connection {
	return s.manager.GetAllConnections()
}

// ExecuteRequest 执行请求
func (s *ConnectionService) ExecuteRequest(ctx context.Context, webshellID string, payload []byte) (*HTTPResponse, error) {
	conn, exists := s.manager.GetConnection(webshellID)
	if !exists {
		return nil, ErrConnectionNotFound
	}

	// 编码 payload
	encodedPayload, err := EncodeBody(payload, conn.Config.Encoder, conn.Config.EncryptionKey)
	if err != nil {
		return nil, err
	}

	// 创建 HTTP 请求
	req := &HTTPRequestConfig{
		URL:          conn.Config.URL,
		Method:       "POST",
		Body:         encodedPayload,
		Timeout:      conn.Config.Timeout,
		ProxyType:    conn.Config.ProxyType,
		ProxyAddress: conn.Config.ProxyAddress,
		SSLVerify:    conn.Config.SSLVerify,
		Headers:      conn.Config.Headers,
		Cookies:      conn.Config.Cookies,
	}

	// 发送请求
	resp, err := s.manager.GetHTTPEngine().Send(ctx, req)
	if err != nil {
		conn.RecordRequest(false, int64(len(encodedPayload)), 0, float64(resp.Duration))
		return resp, err
	}

	// 解码响应
	decodedBody, err := DecodeBody(resp.Body, conn.Config.Encoder, conn.Config.EncryptionKey)
	if err != nil {
		return resp, err
	}

	resp.Body = decodedBody

	// 记录统计
	conn.RecordRequest(true, int64(len(encodedPayload)), int64(len(decodedBody)), float64(resp.Duration))

	return resp, nil
}

// StartHeartbeat 启动心跳检测
func (s *ConnectionService) StartHeartbeat(ctx context.Context, interval time.Duration) {
	s.manager.StartHeartbeat(ctx, interval)
}

// GetConnectionCount 获取连接数量
func (s *ConnectionService) GetConnectionCount() int {
	return s.manager.GetConnectionCount()
}
