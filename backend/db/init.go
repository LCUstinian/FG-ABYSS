package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite" // 纯 Go 实现，不需要 CGO
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"fg-abyss/backend/models"
)

// InitDB 初始化数据库
func InitDB() (*gorm.DB, error) {
	// 创建数据目录
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// 构建数据库文件路径
	dbPath := filepath.Join(dataDir, "app.db")

	// 配置GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 连接数据库
	db, err := gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 获取底层的 sql.DB 实例
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 开启外键约束
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// 开启 WAL 模式
	if _, err := sqlDB.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	fmt.Println("开始检查数据库结构...")

	// 自动迁移数据库表结构
	// 这会在表不存在时自动创建，或更新现有表结构
	fmt.Println("正在执行数据库表结构迁移...")
	if err := db.AutoMigrate(&models.Project{}, &models.WebShell{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	fmt.Println("数据库表结构迁移完成")

	// 创建默认项目
	var defaultProject models.Project
	result := db.Where("name = ?", "默认项目").First(&defaultProject)
	if result.Error == gorm.ErrRecordNotFound {
		defaultProject = models.Project{
			Name:        "默认项目",
			Description: "默认项目，无法删除",
			Status:      0,
		}
		if err := db.Create(&defaultProject).Error; err != nil {
			return nil, fmt.Errorf("failed to create default project: %w", err)
		}
		fmt.Println("默认项目创建成功")
	}

	// 创建示例 WebShell 数据
	var webshellCount int64
	db.Model(&models.WebShell{}).Count(&webshellCount)
	if webshellCount == 0 {
		exampleWebShells := []models.WebShell{
			{
				ProjectID: defaultProject.ID,
				Url:       "http://example.com/shell.php",
				Payload:   "php",
				Cryption:  "base64",
				Encoding:  "utf-8",
				ProxyType: "none",
				Remark:    "示例 PHP WebShell，用于测试",
				Status:    "active",
			},
			{
				ProjectID: defaultProject.ID,
				Url:       "http://test.com/backdoor.aspx",
				Payload:   "aspx",
				Cryption:  "xor",
				Encoding:  "utf-8",
				ProxyType: "none",
				Remark:    "示例 ASPX WebShell",
				Status:    "active",
			},
		}

		for _, ws := range exampleWebShells {
			if err := db.Create(&ws).Error; err != nil {
				return nil, fmt.Errorf("failed to create example webshell: %w", err)
			}
		}
		fmt.Println("示例 WebShell 数据创建成功")
	}

	return db, nil
}
