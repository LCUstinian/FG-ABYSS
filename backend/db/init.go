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

	// 开启WAL模式
	if _, err := sqlDB.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&models.Project{}, &models.WebShell{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
