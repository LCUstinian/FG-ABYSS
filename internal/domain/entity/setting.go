package entity

import (
	"time"
)

// Setting 设置实体 - 用于存储应用程序的所有设置项
type Setting struct {
	ID          uint      `gorm:"primarykey" json:"id"`                         // 自增 ID
	Key         string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`   // 设置键名（如：theme, language）
	Value       string    `gorm:"type:text" json:"value"`                              // 设置值（统一存储为字符串）
	Type        string    `gorm:"type:varchar(20);default:'string'" json:"type"`       // 数据类型：string, int, bool, json, color
	Category    string    `gorm:"type:varchar(50);index" json:"category"`              // 分类：appearance, network, proxy, general
	Name        string    `gorm:"type:varchar(100)" json:"name"`                       // 显示名称
	Description string    `gorm:"type:text" json:"description"`                        // 描述
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`                     // 更新时间
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`                     // 创建时间
}

// TableName 设置表名
func (Setting) TableName() string {
	return "settings"
}
