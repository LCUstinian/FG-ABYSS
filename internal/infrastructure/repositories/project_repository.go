package repositories

import (
	"gorm.io/gorm"

	"fg-abyss/internal/domain/entity"
	"fg-abyss/internal/domain/repository"
)

// ProjectRepositoryImpl 项目仓储实现
type ProjectRepositoryImpl struct {
	db *gorm.DB
}

// NewProjectRepository 创建项目仓储实例
func NewProjectRepository(db *gorm.DB) repository.ProjectRepository {
	return &ProjectRepositoryImpl{db: db}
}

// FindByID 根据 ID 查找项目
func (r *ProjectRepositoryImpl) FindByID(id string) (*entity.Project, error) {
	var project entity.Project
	result := r.db.First(&project, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

// FindByName 根据名称查找项目
func (r *ProjectRepositoryImpl) FindByName(name string) (*entity.Project, error) {
	var project entity.Project
	result := r.db.Where("name = ?", name).First(&project)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

// FindAll 查找所有项目
func (r *ProjectRepositoryImpl) FindAll() ([]entity.Project, error) {
	var projects []entity.Project
	result := r.db.Order("created_at DESC").Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}

// Save 保存项目
func (r *ProjectRepositoryImpl) Save(project *entity.Project) error {
	// 验证
	if err := project.Validate(); err != nil {
		return err
	}

	// 检查是否已存在
	if project.ID != "" {
		var existing entity.Project
		result := r.db.First(&existing, "id = ?", project.ID)
		if result.Error == nil {
			// 更新
			return r.db.Save(project).Error
		}
	}

	// 创建
	return r.db.Create(project).Error
}

// Delete 删除项目
func (r *ProjectRepositoryImpl) Delete(id string) error {
	result := r.db.Delete(&entity.Project{}, "id = ?", id)
	return result.Error
}

// DeleteSoft 软删除项目
func (r *ProjectRepositoryImpl) DeleteSoft(id string) error {
	result := r.db.Where("id = ?", id).Delete(&entity.Project{})
	return result.Error
}
