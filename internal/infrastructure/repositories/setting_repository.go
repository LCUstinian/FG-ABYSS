package repositories

import (
	"fg-abyss/internal/domain/entity"
	"fg-abyss/internal/domain/repository"
	"gorm.io/gorm"
)

// SettingRepositoryImpl 设置仓储实现
type SettingRepositoryImpl struct {
	db *gorm.DB
}

// NewSettingRepository 创建设置仓储实例
func NewSettingRepository(db *gorm.DB) repository.SettingRepository {
	return &SettingRepositoryImpl{db: db}
}

// FindAll 查找所有设置
func (r *SettingRepositoryImpl) FindAll() ([]entity.Setting, error) {
	var settings []entity.Setting
	result := r.db.Order("category, id").Find(&settings)
	if result.Error != nil {
		return nil, result.Error
	}
	return settings, nil
}

// FindByCategory 根据分类查找设置
func (r *SettingRepositoryImpl) FindByCategory(category string) ([]entity.Setting, error) {
	var settings []entity.Setting
	result := r.db.Where("category = ?", category).Order("id").Find(&settings)
	if result.Error != nil {
		return nil, result.Error
	}
	return settings, nil
}

// FindByKey 根据键名查找设置
func (r *SettingRepositoryImpl) FindByKey(key string) (*entity.Setting, error) {
	var setting entity.Setting
	result := r.db.Where("key = ?", key).First(&setting)
	if result.Error != nil {
		return nil, result.Error
	}
	return &setting, nil
}

// Save 保存设置（创建或更新）
func (r *SettingRepositoryImpl) Save(setting *entity.Setting) error {
	return r.db.Save(setting).Error
}

// SaveAll 批量保存设置
func (r *SettingRepositoryImpl) SaveAll(settings []entity.Setting) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, setting := range settings {
			if err := tx.Save(&setting).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateValue 更新设置值
func (r *SettingRepositoryImpl) UpdateValue(key string, value string) error {
	return r.db.Model(&entity.Setting{}).Where("key = ?", key).Update("value", value).Error
}

// InitializeDefaults 初始化默认设置
func (r *SettingRepositoryImpl) InitializeDefaults(defaults []entity.Setting) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, setting := range defaults {
			// 检查是否已存在
			var existing entity.Setting
			result := tx.Where("key = ?", setting.Key).First(&existing)
			
			if result.Error == gorm.ErrRecordNotFound {
				// 不存在则创建
				if err := tx.Create(&setting).Error; err != nil {
					return err
				}
			} else if result.Error != nil {
				// 其他错误
				return result.Error
			}
			// 如果已存在，不更新（保留用户自定义的值）
		}
		return nil
	})
}
