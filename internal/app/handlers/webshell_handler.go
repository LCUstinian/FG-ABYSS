package handlers

import (
	"fg-abyss/internal/app/services"
	"fg-abyss/internal/domain/entity"
)

// WebShellHandler WebShell 处理器
type WebShellHandler struct {
	webshellService *services.WebShellService
}

// NewWebShellHandler 创建 WebShell 处理器实例
func NewWebShellHandler(webshellService *services.WebShellService) *WebShellHandler {
	return &WebShellHandler{webshellService: webshellService}
}

// GetWebShells 获取 WebShell 列表
func (h *WebShellHandler) GetWebShells(projectID string, page int, pageSize int, searchQuery string, sortField string, sortDir string) ([]entity.WebShell, int64, error) {
	return h.webshellService.GetByProjectIDPaginated(projectID, page, pageSize, searchQuery, sortField, sortDir)
}

// GetWebShellByID 根据 ID 获取 WebShell
func (h *WebShellHandler) GetWebShellByID(id string) (*entity.WebShell, error) {
	return h.webshellService.GetByID(id)
}

// CreateWebShell 创建 WebShell
func (h *WebShellHandler) CreateWebShell(
	projectID, url, payload, cryption, encoding, proxyType, remark, status string,
) (*entity.WebShell, error) {
	return h.webshellService.Create(projectID, url, payload, cryption, encoding, proxyType, remark, status)
}

// UpdateWebShell 更新 WebShell
func (h *WebShellHandler) UpdateWebShell(
	id, projectID, url, payload, cryption, encoding, proxyType, remark, status string,
) (*entity.WebShell, error) {
	return h.webshellService.Update(id, projectID, url, payload, cryption, encoding, proxyType, remark, status)
}

// DeleteWebShell 删除 WebShell
func (h *WebShellHandler) DeleteWebShell(id string) error {
	return h.webshellService.Delete(id)
}

// RecoverWebShell 恢复已删除的 WebShell
func (h *WebShellHandler) RecoverWebShell(id string) error {
	return h.webshellService.Recover(id)
}
