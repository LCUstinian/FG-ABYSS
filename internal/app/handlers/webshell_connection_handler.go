package handlers

import (
	"context"

	"fg-abyss/internal/app/services"
	"fg-abyss/internal/domain/entity"
)

// WebShellConnectionHandler WebShell 连接处理器
type WebShellConnectionHandler struct {
	connectionService *services.ConnectionService
	commandService    *services.CommandService
	fileService       *services.FileService
}

// NewWebShellConnectionHandler 创建 WebShell 连接处理器
func NewWebShellConnectionHandler() *WebShellConnectionHandler {
	connService := services.NewConnectionService()
	return &WebShellConnectionHandler{
		connectionService: connService,
		commandService:    services.NewCommandService(connService),
		fileService:       services.NewFileService(connService),
	}
}

// WebShellConnectRequest WebShell 连接请求参数
type WebShellConnectRequest struct {
	WebShellID     string            `json:"webshell_id"`
	URL            string            `json:"url"`
	Password       string            `json:"password"`
	Encoder        string            `json:"encoder"`
	EncryptionKey  string            `json:"encryption_key"`
	Timeout        int               `json:"timeout"`
	ProxyType      string            `json:"proxy_type"`
	ProxyAddress   string            `json:"proxy_address"`
	SSLVerify      bool              `json:"ssl_verify"`
	Headers        map[string]string `json:"headers"`
	Cookies        map[string]string `json:"cookies"`
}

// Connect 连接 WebShell
func (h *WebShellConnectionHandler) Connect(req *WebShellConnectRequest) (*services.ConnectResponse, error) {
	ctx := context.Background()
	
	connReq := &services.ConnectRequest{
		WebShellID:     req.WebShellID,
		URL:            req.URL,
		Password:       req.Password,
		Encoder:        req.Encoder,
		EncryptionKey:  req.EncryptionKey,
		Timeout:        req.Timeout,
		ProxyType:      req.ProxyType,
		ProxyAddress:   req.ProxyAddress,
		SSLVerify:      req.SSLVerify,
	}
	
	return h.connectionService.Connect(ctx, connReq)
}

// Disconnect 断开 WebShell 连接
func (h *WebShellConnectionHandler) Disconnect(webshellID string) error {
	return h.connectionService.Disconnect(webshellID)
}

// GetConnection 获取连接状态
func (h *WebShellConnectionHandler) GetConnection(webshellID string) (*entity.Connection, error) {
	return h.connectionService.GetConnection(webshellID)
}

// GetAllConnections 获取所有连接
func (h *WebShellConnectionHandler) GetAllConnections() []*entity.Connection {
	return h.connectionService.GetAllConnections()
}

// ExecuteCommand 执行命令
func (h *WebShellConnectionHandler) ExecuteCommand(webshellID, command string, timeout int) (*services.CommandResponse, error) {
	ctx := context.Background()
	
	cmdReq := &services.CommandRequest{
		WebShellID: webshellID,
		Command:    command,
		Timeout:    timeout,
	}
	
	return h.commandService.Execute(ctx, cmdReq)
}

// ListFiles 列出目录内容
func (h *WebShellConnectionHandler) ListFiles(webshellID, path string) (*services.FileListResponse, error) {
	ctx := context.Background()
	
	req := &services.FileListRequest{
		WebShellID: webshellID,
		Path:       path,
	}
	
	return h.fileService.ListFiles(ctx, req)
}

// ReadFile 读取文件
func (h *WebShellConnectionHandler) ReadFile(webshellID, path string, limit int) (string, error) {
	ctx := context.Background()
	
	req := &services.FileReadRequest{
		WebShellID: webshellID,
		Path:       path,
		Limit:      limit,
	}
	
	return h.fileService.ReadFile(ctx, req)
}

// WriteFile 写入文件
func (h *WebShellConnectionHandler) WriteFile(webshellID, path, content string, isBase64 bool) error {
	ctx := context.Background()
	
	req := &services.FileWriteRequest{
		WebShellID: webshellID,
		Path:       path,
		Content:    content,
		IsBase64:   isBase64,
	}
	
	return h.fileService.WriteFile(ctx, req)
}

// UploadFile 上传文件
func (h *WebShellConnectionHandler) UploadFile(webshellID, path, filename, content string) error {
	ctx := context.Background()
	
	req := &services.FileUploadRequest{
		WebShellID: webshellID,
		Path:       path,
		Filename:   filename,
		Content:    content,
	}
	
	return h.fileService.UploadFile(ctx, req)
}

// DownloadFile 下载文件
func (h *WebShellConnectionHandler) DownloadFile(webshellID, path string) ([]byte, error) {
	ctx := context.Background()
	
	req := &services.FileDownloadRequest{
		WebShellID: webshellID,
		Path:       path,
	}
	
	return h.fileService.DownloadFile(ctx, req)
}

// DeleteFile 删除文件
func (h *WebShellConnectionHandler) DeleteFile(webshellID, path string) error {
	return h.fileService.DeleteFile(context.Background(), webshellID, path)
}

// CreateDirectory 创建目录
func (h *WebShellConnectionHandler) CreateDirectory(webshellID, path string) error {
	return h.fileService.CreateDirectory(context.Background(), webshellID, path)
}

// RenameFile 重命名文件
func (h *WebShellConnectionHandler) RenameFile(webshellID, oldPath, newPath string) error {
	return h.fileService.RenameFile(context.Background(), webshellID, oldPath, newPath)
}

// GetSystemInfo 获取系统信息
func (h *WebShellConnectionHandler) GetSystemInfo(webshellID string) (map[string]interface{}, error) {
	ctx := context.Background()
	return h.commandService.GetSystemInfo(ctx, webshellID)
}

// GetCurrentUser 获取当前用户
func (h *WebShellConnectionHandler) GetCurrentUser(webshellID string) (string, error) {
	ctx := context.Background()
	return h.commandService.GetCurrentUser(ctx, webshellID)
}

// GetWorkingDirectory 获取工作目录
func (h *WebShellConnectionHandler) GetWorkingDirectory(webshellID string) (string, error) {
	ctx := context.Background()
	return h.commandService.GetWorkingDirectory(ctx, webshellID)
}

// SearchFiles 搜索文件
func (h *WebShellConnectionHandler) SearchFiles(webshellID, path, pattern string) ([]services.FileInfo, error) {
	ctx := context.Background()
	return h.fileService.SearchFiles(ctx, webshellID, path, pattern)
}

// GetConnectionCount 获取连接数量
func (h *WebShellConnectionHandler) GetConnectionCount() int {
	return h.connectionService.GetConnectionCount()
}
