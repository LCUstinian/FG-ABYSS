package services

import (
	"errors"

	"fg-abyss/internal/domain/entity"
	"fg-abyss/internal/domain/repository"
)

// WebShellService WebShell 服务
type WebShellService struct {
	webshellRepo repository.WebShellRepository
}

// NewWebShellService 创建 WebShell 服务实例
func NewWebShellService(repo repository.WebShellRepository) *WebShellService {
	return &WebShellService{webshellRepo: repo}
}

// GetAll 获取所有 WebShell（不分页）
func (s *WebShellService) GetAll(projectID string) ([]entity.WebShell, error) {
	return s.webshellRepo.FindByProjectID(projectID)
}

// GetByID 根据 ID 获取 WebShell
func (s *WebShellService) GetByID(id string) (*entity.WebShell, error) {
	return s.webshellRepo.FindByID(id)
}

// GetPaginated 分页获取 WebShell
func (s *WebShellService) GetPaginated(projectID string, page int, pageSize int, searchQuery string, sortField string, sortDir string) ([]entity.WebShell, int64, error) {
	return s.webshellRepo.FindByProjectIDPaginated(projectID, page, pageSize, searchQuery, sortField, sortDir)
}

// Create 创建 WebShell
func (s *WebShellService) Create(webshell *entity.WebShell) (*entity.WebShell, error) {
	if err := s.webshellRepo.Save(webshell); err != nil {
		return nil, err
	}
	return webshell, nil
}

// Update 更新 WebShell
func (s *WebShellService) Update(webshell *entity.WebShell) (*entity.WebShell, error) {
	// 先检查是否存在
	_, err := s.webshellRepo.FindByID(webshell.ID)
	if err != nil {
		return nil, errors.New("WebShell 不存在")
	}

	if err := s.webshellRepo.Save(webshell); err != nil {
		return nil, err
	}

	return webshell, nil
}

// Delete 删除 WebShell
func (s *WebShellService) Delete(id string) error {
	return s.webshellRepo.DeleteSoft(id)
}

// Recover 恢复已删除的 WebShell
func (s *WebShellService) Recover(id string) error {
	return s.webshellRepo.Recover(id)
}

// GetDeleted 获取已删除的 WebShell
func (s *WebShellService) GetDeleted(projectID string) ([]entity.WebShell, error) {
	return s.webshellRepo.FindDeletedByProjectID(projectID)
}
