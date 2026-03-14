package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WebShell WebShell 模型
type WebShell struct {
	ID        string         `gorm:"type:text;primaryKey" json:"id"`
	ProjectID string         `gorm:"type:text;index" json:"projectId"`
	Url       string         `gorm:"type:text;not null" json:"url"`
	Payload   string         `gorm:"type:text" json:"payload"`
	Cryption  string         `gorm:"type:text" json:"cryption"`
	Encoding  string         `gorm:"type:text" json:"encoding"`
	ProxyType string         `gorm:"type:text" json:"proxyType"`
	Remark    string         `gorm:"type:text" json:"remark"`
	Status    string         `gorm:"type:text" json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (WebShell) TableName() string {
	return "webshells"
}

// BeforeCreate 创建前生成 UUID
func (w *WebShell) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return nil
}
