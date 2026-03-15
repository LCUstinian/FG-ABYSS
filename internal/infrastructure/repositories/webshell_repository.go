package repositories

import (
	"gorm.io/gorm"

	"fg-abyss/internal/domain/entity"
	"fg-abyss/internal/domain/repository"
)

// WebShellRepositoryImpl WebShell 仓储实现
type WebShellRepositoryImpl struct {
	db *gorm.DB
}

// NewWebShellRepository 创建 WebShell 仓储实例
func NewWebShellRepository(db *gorm.DB) repository.WebShellRepository {
	return &WebShellRepositoryImpl{db: db}
}

// FindByID 根据 ID 查找 WebShell
func (r *WebShellRepositoryImpl) FindByID(id string) (*entity.WebShell, error) {
	var webshell entity.WebShell
	result := r.db.First(&webshell, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &webshell, nil
}

// FindByProjectID 根据项目 ID 查找所有 WebShell
func (r *WebShellRepositoryImpl) FindByProjectID(projectID string) ([]entity.WebShell, error) {
	var webshells []entity.WebShell
	result := r.db.Where("project_id = ? AND deleted_at IS NULL", projectID).
		Order("created_at DESC").
		Find(&webshells)
	if result.Error != nil {
		return nil, result.Error
	}
	return webshells, nil
}

// FindByProjectIDPaginated 分页查找 WebShell
func (r *WebShellRepositoryImpl) FindByProjectIDPaginated(projectID string, page int, pageSize int, searchQuery string, sortField string, sortDir string) ([]entity.WebShell, int64, error) {
	var webshells []entity.WebShell
	var total int64

	// 构建查询
	query := r.db.Model(&entity.WebShell{}).Where("project_id = ?", projectID)

	// 搜索条件
	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		query = query.Where("(url LIKE ? OR remark LIKE ? OR payload LIKE ?)",
			searchPattern, searchPattern, searchPattern)
	}

	// 只查询未删除的记录
	query = query.Where("deleted_at IS NULL")

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	if sortField != "" && sortDir != "" {
		// 验证排序字段
		validFields := map[string]bool{
			"url": true, "created_at": true, "updated_at": true,
			"payload": true, "status": true,
		}
		if validFields[sortField] {
			query = query.Order(fmt.Sprintf("%s %s", sortField, sortDir))
		} else {
			query = query.Order("created_at DESC")
		}
	} else {
		query = query.Order("created_at DESC")
	}

	// 分页
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	query = query.Offset(offset).Limit(pageSize)

	// 执行查询
	result := query.Find(&webshells)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return webshells, total, nil
}

// FindDeletedByProjectID 查找已删除的 WebShell
func (r *WebShellRepositoryImpl) FindDeletedByProjectID(projectID string) ([]entity.WebShell, error) {
	var webshells []entity.WebShell
	result := r.db.Unscoped().
		Where("project_id = ? AND deleted_at IS NOT NULL", projectID).
		Find(&webshells)
	if result.Error != nil {
		return nil, result.Error
	}
	return webshells, nil
}

// Save 保存 WebShell
func (r *WebShellRepositoryImpl) Save(webshell *entity.WebShell) error {
	// 验证
	if err := webshell.Validate(); err != nil {
		return err
	}

	// 检查是否已存在
	if webshell.ID != "" {
		var existing entity.WebShell
		result := r.db.First(&existing, "id = ?", webshell.ID)
		if result.Error == nil {
			// 更新
			return r.db.Save(webshell).Error
		}
	}

	// 创建
	return r.db.Create(webshell).Error
}

// Delete 删除 WebShell
func (r *WebShellRepositoryImpl) Delete(id string) error {
	result := r.db.Unscoped().Delete(&entity.WebShell{}, "id = ?", id)
	return result.Error
}

// DeleteSoft 软删除 WebShell
func (r *WebShellRepositoryImpl) DeleteSoft(id string) error {
	result := r.db.Where("id = ?", id).Delete(&entity.WebShell{})
	return result.Error
}

// Recover 恢复已删除的 WebShell
func (r *WebShellRepositoryImpl) Recover(id string) error {
	result := r.db.Unscoped().Where("id = ?", id).Update("deleted_at", nil)
	return result.Error
}
