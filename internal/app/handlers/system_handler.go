package handlers

import (
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
