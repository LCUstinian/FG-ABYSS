package services

import (
	"errors"

	"gorm.io/gorm"

	"fg-abyss/internal/domain/entity"
	"fg-abyss/internal/domain/repository"
)

// ProjectService 项目服务
type ProjectService struct {
	projectRepo repository.ProjectRepository
}

// NewProjectService 创建项目服务实例
func NewProjectService(repo repository.ProjectRepository) *ProjectService {
	return &ProjectService{projectRepo: repo}
}

// GetAll 获取所有项目
func (s *ProjectService) GetAll() ([]entity.Project, error) {
	return s.projectRepo.FindAll()
}

// GetByID 根据 ID 获取项目
func (s *ProjectService) GetByID(id string) (*entity.Project, error) {
	return s.projectRepo.FindByID(id)
}

// Create 创建项目
func (s *ProjectService) Create(name, description string) (*entity.Project, error) {
	// 检查名称是否已存在
	existing, err := s.projectRepo.FindByName(name)
	if err == nil && existing != nil {
		return nil, errors.New("项目名称已存在")
	}
	// 如果是记录不存在错误，继续执行
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	project := &entity.Project{
		Name:        name,
		Description: description,
		Status:      0, // Active
	}

	if err := s.projectRepo.Save(project); err != nil {
		return nil, err
	}

	return project, nil
}

// Update 更新项目
func (s *ProjectService) Update(id, name, description string, status int) (*entity.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("项目不存在")
	}

	// 检查名称是否与其他项目重复
	if project.Name != name {
		existing, err := s.projectRepo.FindByName(name)
		if err == nil && existing != nil {
			return nil, errors.New("项目名称已存在")
		}
	}

	project.Name = name
	project.Description = description
	project.Status = status

	if err := s.projectRepo.Save(project); err != nil {
		return nil, err
	}

	return project, nil
}

// Delete 删除项目
func (s *ProjectService) Delete(id string) error {
	// 检查是否为默认项目
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return errors.New("项目不存在")
	}

	if project.Name == "默认项目" {
		return errors.New("默认项目不能删除")
	}

	return s.projectRepo.DeleteSoft(id)
}

// Recover 恢复已删除的项目
func (s *ProjectService) Recover(id string) (*entity.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("项目不存在")
	}

	if project.DeletedAt.Valid {
		// 恢复项目
		if err := s.projectRepo.Recover(id); err != nil {
			return nil, err
		}
		return project, nil
	}

	return nil, errors.New("项目未被删除")
}

// GetDeleted 获取已删除的项目
func (s *ProjectService) GetDeleted() ([]entity.Project, error) {
	return s.projectRepo.FindDeleted()
}
