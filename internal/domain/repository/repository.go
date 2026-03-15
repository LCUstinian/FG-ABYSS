package repository

import (
	"fg-abyss/internal/domain/entity"
)

// ProjectRepository 项目仓储接口
type ProjectRepository interface {
	// FindByID 根据 ID 查找项目
	FindByID(id string) (*entity.Project, error)
	// FindByName 根据名称查找项目
	FindByName(name string) (*entity.Project, error)
	// FindAll 查找所有项目
	FindAll() ([]entity.Project, error)
	// Save 保存项目（创建或更新）
	Save(project *entity.Project) error
	// Delete 删除项目
	Delete(id string) error
	// DeleteSoft 软删除项目
	DeleteSoft(id string) error
}

// WebShellRepository WebShell 仓储接口
type WebShellRepository interface {
	// FindByID 根据 ID 查找 WebShell
	FindByID(id string) (*entity.WebShell, error)
	// FindByProjectID 根据项目 ID 查找所有 WebShell
	FindByProjectID(projectID string) ([]entity.WebShell, error)
	// FindByProjectIDPaginated 分页查找 WebShell
	FindByProjectIDPaginated(projectID string, page int, pageSize int, searchQuery string, sortField string, sortDir string) ([]entity.WebShell, int64, error)
	// FindDeletedByProjectID 查找已删除的 WebShell
	FindDeletedByProjectID(projectID string) ([]entity.WebShell, error)
	// Save 保存 WebShell
	Save(webshell *entity.WebShell) error
	// Delete 删除 WebShell
	Delete(id string) error
	// DeleteSoft 软删除 WebShell
	DeleteSoft(id string) error
	// Recover 恢复已删除的 WebShell
	Recover(id string) error
}
