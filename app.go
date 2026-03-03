package main

import (
	"fg-abyss/backend/models"

	"gorm.io/gorm"
)

// App 应用结构体
type App struct {
	db *gorm.DB
}

// NewApp 创建应用实例
func NewApp(db *gorm.DB) *App {
	return &App{
		db: db,
	}
}

// GetWebShells 获取 WebShell 列表
func (a *App) GetWebShells(projectName string, page int, pageSize int, searchQuery string, sortField string, sortDir string) ([]models.WebShell, int64, error) {
	var project models.Project
	var webshells []models.WebShell
	var total int64

	// 先根据项目名称查询项目 ID
	if err := a.db.Where("name = ?", projectName).First(&project).Error; err != nil {
		// 如果项目不存在，返回空列表
		return []models.WebShell{}, 0, nil
	}

	// 构建查询
	query := a.db.Where("project_id = ?", project.ID)

	// 搜索
	if searchQuery != "" {
		query = query.Where(
			"name LIKE ? OR url LIKE ? OR payload LIKE ? OR cryption LIKE ? OR encoding LIKE ? OR proxy_type LIKE ? OR remark LIKE ? OR status LIKE ?",
			"%"+searchQuery+"%",
			"%"+searchQuery+"%",
			"%"+searchQuery+"%",
			"%"+searchQuery+"%",
			"%"+searchQuery+"%",
			"%"+searchQuery+"%",
			"%"+searchQuery+"%",
			"%"+searchQuery+"%",
		)
	}

	// 计算总数
	if err := query.Model(&models.WebShell{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	if sortField != "" {
		order := sortField
		if sortDir == "desc" {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	}

	// 分页
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&webshells).Error; err != nil {
		return nil, 0, err
	}

	return webshells, total, nil
}

// CreateWebShell 创建 WebShell
func (a *App) CreateWebShell(shell models.WebShell) error {
	// 先根据项目名称查询项目 ID
	var project models.Project
	if err := a.db.Where("name = ?", shell.ProjectID).First(&project).Error; err != nil {
		return err
	}

	// 设置项目 ID
	shell.ProjectID = project.ID

	// 创建 WebShell
	return a.db.Create(&shell).Error
}

// UpdateWebShell 更新 WebShell
func (a *App) UpdateWebShell(shell models.WebShell) error {
	// 更新 WebShell
	return a.db.Save(&shell).Error
}

// DeleteWebShell 删除 WebShell
func (a *App) DeleteWebShell(id string) error {
	// 删除 WebShell
	return a.db.Delete(&models.WebShell{}, "id = ?", id).Error
}
