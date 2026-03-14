package db

import (
	"fmt"
	"gorm.io/gorm"
)

// MigrateRemoveWebShellNameField 删除 WebShell 表的 name 字段
// 这个迁移脚本用于从现有数据库中移除 name 字段
func MigrateRemoveWebShellNameField(db *gorm.DB) error {
	fmt.Println("开始迁移：删除 WebShell 表的 name 字段...")

	// 检查 name 字段是否存在
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM pragma_table_info('webshells') WHERE name='name'").Scan(&count).Error
	if err != nil {
		return fmt.Errorf("检查字段失败：%w", err)
	}

	if count == 0 {
		fmt.Println("name 字段不存在，无需迁移")
		return nil
	}

	fmt.Println("检测到 name 字段，开始迁移...")

	// SQLite 删除列的步骤：
	// 1. 启用外键约束
	// 2. 创建新表（不含 name 字段）
	// 3. 复制数据
	// 4. 删除旧表
	// 5. 重命名新表

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 禁用外键约束（迁移期间）
	if err := tx.Exec("PRAGMA foreign_keys = OFF").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("禁用外键失败：%w", err)
	}

	// 2. 创建临时表（不含 name 字段）
	createTempTableSQL := `
	CREATE TABLE webshells_temp (
		id TEXT PRIMARY KEY,
		projectId TEXT,
		url TEXT NOT NULL,
		payload TEXT,
		cryption TEXT,
		encoding TEXT,
		proxyType TEXT,
		remark TEXT,
		createTime TEXT,
		updateTime TEXT,
		status TEXT,
		createdAt DATETIME,
		updatedAt DATETIME,
		deletedAt DATETIME
	);`

	if err := tx.Exec(createTempTableSQL).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建临时表失败：%w", err)
	}

	// 3. 复制数据（不含 name 字段）
	copyDataSQL := `
	INSERT INTO webshells_temp (
		id, projectId, url, payload, cryption, encoding, 
		proxyType, remark, createTime, updateTime, status,
		createdAt, updatedAt, deletedAt
	)
	SELECT 
		id, projectId, url, payload, cryption, encoding, 
		proxyType, remark, createTime, updateTime, status,
		createdAt, updatedAt, deletedAt
	FROM webshells;`

	if err := tx.Exec(copyDataSQL).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("复制数据失败：%w", err)
	}

	// 4. 删除旧表
	if err := tx.Exec("DROP TABLE webshells").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除旧表失败：%w", err)
	}

	// 5. 重命名临时表
	if err := tx.Exec("ALTER TABLE webshells_temp RENAME TO webshells").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("重命名表失败：%w", err)
	}

	// 6. 重新启用外键约束
	if err := tx.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("启用外键失败：%w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败：%w", err)
	}

	fmt.Println("✓ WebShell 表的 name 字段已成功删除")
	return nil
}
