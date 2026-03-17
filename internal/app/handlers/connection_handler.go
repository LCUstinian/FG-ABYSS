package handlers

import (
	"context"

	"fg-abyss/internal/domain/entity"
	"fg-abyss/internal/app/services"
)

// ConnectionHandler 连接管理处理器（Wails 绑定用）
type ConnectionHandler struct {
	service *services.ConnectionService
}

// NewConnectionHandler 创建连接处理器
func NewConnectionHandler() *ConnectionHandler {
	return &ConnectionHandler{
		service: services.NewConnectionService(),
	}
}

// Connect 连接 WebShell
// @summary 连接到 WebShell
// @param ctx 上下文
// @param webshellID WebShell ID
// @param url 目标 URL
// @param password 密码
// @param encoder 编码器类型（none/base64/rot13/xor）
// @param encryptionKey 加密密钥
// @param timeout 超时时间（秒）
// @param proxyType 代理类型（none/http/https/socks5）
// @param proxyAddress 代理地址
// @param sslVerify 是否验证 SSL
// @return 连接结果
func (h *ConnectionHandler) Connect(
	ctx context.Context,
	webshellID string,
	url string,
	password string,
	encoder string,
	encryptionKey string,
	timeout int,
	proxyType string,
	proxyAddress string,
	sslVerify bool,
) (*ConnectResult, error) {
	req := &services.ConnectRequest{
		WebShellID:     webshellID,
		URL:            url,
		Password:       password,
		Encoder:        encoder,
		EncryptionKey:  encryptionKey,
		Timeout:        timeout,
		ProxyType:      proxyType,
		ProxyAddress:   proxyAddress,
		SSLVerify:      sslVerify,
	}
	
	resp, err := h.service.Connect(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return &ConnectResult{
		Success:    resp.Success,
		Message:    resp.Message,
		Connection: resp.Connection,
	}, nil
}

// Disconnect 断开 WebShell 连接
// @summary 断开 WebShell 连接
// @param webshellID WebShell ID
// @return 错误信息
func (h *ConnectionHandler) Disconnect(webshellID string) error {
	return h.service.Disconnect(webshellID)
}

// GetConnection 获取连接状态
// @summary 获取连接状态
// @param webshellID WebShell ID
// @return 连接信息
func (h *ConnectionHandler) GetConnection(webshellID string) (*entity.Connection, error) {
	return h.service.GetConnection(webshellID)
}

// GetAllConnections 获取所有连接
// @summary 获取所有连接
// @return 连接列表
func (h *ConnectionHandler) GetAllConnections() []*entity.Connection {
	return h.service.GetAllConnections()
}

// GetConnectionCount 获取连接数量
// @summary 获取连接数量
// @return 连接数量
func (h *ConnectionHandler) GetConnectionCount() int {
	return h.service.GetConnectionCount()
}

// ConnectResult 连接结果
type ConnectResult struct {
	Success    bool              `json:"success"`
	Message    string            `json:"message"`
	Connection *entity.Connection `json:"connection"`
}
