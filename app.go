package main

import (
	"fg-abyss/backend/models"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"

	"gorm.io/gorm"
)

// App 应用结构体
type App struct {
	db        *gorm.DB
	startTime time.Time
	process   *process.Process
	processID int32
}

// NewApp 创建应用实例
func NewApp(db *gorm.DB) *App {
	// 获取当前进程信息
	pid := os.Getpid()
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		fmt.Printf("Warning: Failed to get process info: %v\n", err)
	}

	return &App{
		db:        db,
		startTime: time.Now(),
		process:   proc,
		processID: int32(pid),
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
			"url LIKE ? OR payload LIKE ? OR cryption LIKE ? OR encoding LIKE ? OR proxy_type LIKE ? OR remark LIKE ? OR status LIKE ?",
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
		// 将驼峰命名字段转换为下划线命名（GORM 默认行为）
		snakeField := ""
		for i, r := range sortField {
			if i > 0 && r >= 'A' && r <= 'Z' {
				snakeField += "_"
				snakeField += string(r + 32) // 转为小写
			} else {
				snakeField += string(r)
			}
		}

		order := snakeField
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

// CreateProject 创建项目
func (a *App) CreateProject(name string, description string) error {
	// 检查项目是否已存在
	var existingProject models.Project
	if err := a.db.Where("name = ?", name).First(&existingProject).Error; err == nil {
		return fmt.Errorf("项目名称已存在")
	}

	// 创建新项目
	project := models.Project{
		Name:        name,
		Description: description,
		Status:      0,
	}

	return a.db.Create(&project).Error
}

// GetProjects 获取项目列表
func (a *App) GetProjects() ([]models.Project, error) {
	var projects []models.Project
	if err := a.db.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

// DeleteProject 删除项目
func (a *App) DeleteProject(projectName string) error {
	// 检查是否为默认项目
	if projectName == "默认项目" {
		return fmt.Errorf("默认项目无法删除")
	}

	// 先根据项目名称查询项目
	var project models.Project
	if err := a.db.Where("name = ?", projectName).First(&project).Error; err != nil {
		return err
	}

	// 删除项目
	return a.db.Delete(&project).Error
}

// SystemStatus 系统状态结构体
type SystemStatus struct {
	MemoryUsage string `json:"memoryUsage"`
	ProcessID   string `json:"processId"`
	CPUUsage    string `json:"cpuUsage"`
	Uptime      string `json:"uptime"`
}

// GetSystemStatus 获取系统状态信息
func (a *App) GetSystemStatus() SystemStatus {
	status := SystemStatus{
		ProcessID: fmt.Sprintf("%d", a.processID),
		Uptime:    formatDuration(time.Since(a.startTime)),
	}

	// 获取内存使用信息
	if memInfo, err := mem.VirtualMemory(); err == nil {
		usedGB := float64(memInfo.Used) / 1024 / 1024 / 1024
		totalGB := float64(memInfo.Total) / 1024 / 1024 / 1024
		status.MemoryUsage = fmt.Sprintf("%.2f GB / %.2f GB", usedGB, totalGB)
	} else {
		status.MemoryUsage = "N/A"
	}

	// 获取进程 CPU 使用率
	if a.process != nil {
		if cpuPercent, err := a.process.CPUPercent(); err == nil {
			status.CPUUsage = fmt.Sprintf("%.2f%%", cpuPercent)
		} else {
			status.CPUUsage = "N/A"
		}
	} else {
		status.CPUUsage = "N/A"
	}

	return status
}

// formatDuration 格式化时间间隔
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d小时 %d分钟 %d秒", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%d分钟 %d秒", minutes, seconds)
	}
	return fmt.Sprintf("%d秒", seconds)
}

// camelToSnake 将驼峰命名转换为下划线命名
func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
			result.WriteRune(r + 32) // 转为小写
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
