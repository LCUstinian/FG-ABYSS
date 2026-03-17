package handlers

import (
	"context"
	"fg-abyss/internal/app/services"
)

// SystemHandler 系统处理器
type SystemHandler struct {
	appService *services.AppService
}

// NewSystemHandler 创建系统处理器实例
func NewSystemHandler(appService *services.AppService) *SystemHandler {
	return &SystemHandler{appService: appService}
}

// GetSystemStatus 获取系统状态
func (h *SystemHandler) GetSystemStatus() (map[string]interface{}, error) {
	return h.appService.GetSystemStatus()
}

// OpenWebShellWindowRequest 打开 WebShell 窗口请求
type OpenWebShellWindowRequest struct {
	WebShellID string `json:"webshell_id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
}

// OpenWebShellWindow 打开 WebShell 独立窗口
func (h *SystemHandler) OpenWebShellWindow(ctx context.Context, req *OpenWebShellWindowRequest) error {
	// TODO: 实现打开新窗口的逻辑
	// 这里需要调用 Wails 的窗口管理 API
	return nil
}

// CloseWebShellWindowRequest 关闭 WebShell 窗口请求
type CloseWebShellWindowRequest struct {
	WebShellID string `json:"webshell_id"`
}

// CloseWebShellWindow 关闭 WebShell 独立窗口
func (h *SystemHandler) CloseWebShellWindow(ctx context.Context, req *CloseWebShellWindowRequest) error {
	// TODO: 实现关闭窗口的逻辑
	return nil
}

// GetWebShellWindowStatus 获取 WebShell 窗口状态
func (h *SystemHandler) GetWebShellWindowStatus(ctx context.Context, webshellID string) (map[string]interface{}, error) {
	// TODO: 实现获取窗口状态的逻辑
	return map[string]interface{}{
		"isOpen": false,
	}, nil
}
