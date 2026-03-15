package handlers

import (
	"fg-abyss/internal/app/services"
	"fg-abyss/internal/domain/entity"
)

// ProjectHandler 项目处理器
type ProjectHandler struct {
	projectService *services.ProjectService
}

// NewProjectHandler 创建项目处理器实例
func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

// GetProjects 获取所有项目
func (h *ProjectHandler) GetProjects() ([]entity.Project, error) {
	return h.projectService.GetAll()
}

// GetProjectByID 根据 ID 获取项目
func (h *ProjectHandler) GetProjectByID(id string) (*entity.Project, error) {
	return h.projectService.GetByID(id)
}

// CreateProject 创建项目
func (h *ProjectHandler) CreateProject(name, description string) (*entity.Project, error) {
	return h.projectService.Create(name, description)
}

// UpdateProject 更新项目
func (h *ProjectHandler) UpdateProject(id, name, description string, status int) (*entity.Project, error) {
	return h.projectService.Update(id, name, description, status)
}

// DeleteProject 删除项目
func (h *ProjectHandler) DeleteProject(id string) error {
	return h.projectService.Delete(id)
}

// RecoverProject 恢复已删除的项目
func (h *ProjectHandler) RecoverProject(id string) (*entity.Project, error) {
	return h.projectService.Recover(id)
}

// GetDeletedProjects 获取已删除的项目
func (h *ProjectHandler) GetDeletedProjects() ([]entity.Project, error) {
	return h.projectService.GetDeleted()
}
