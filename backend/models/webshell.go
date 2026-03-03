package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WebShell WebShell 模型
type WebShell struct {
	ID         string         `gorm:"type:text;primaryKey" json:"id"`
	ProjectID  string         `gorm:"type:text;index" json:"projectId"`
	Name       string         `gorm:"type:text;not null" json:"name"`
	Url        string         `gorm:"type:text;not null" json:"url"`
	Payload    string         `gorm:"type:text" json:"payload"`
	Cryption   string         `gorm:"type:text" json:"cryption"`
	Encoding   string         `gorm:"type:text" json:"encoding"`
	ProxyType  string         `gorm:"type:text" json:"proxyType"`
	Remark     string         `gorm:"type:text" json:"remark"`
	CreateTime string         `gorm:"type:text" json:"createTime"`
	UpdateTime string         `gorm:"type:text" json:"updateTime"`
	Status     string         `gorm:"type:text" json:"status"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (WebShell) TableName() string {
	return "webshells"
}

// BeforeCreate 创建前生成 UUID 和时间戳
func (w *WebShell) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	if w.CreateTime == "" {
		w.CreateTime = now
	}
	if w.UpdateTime == "" {
		w.UpdateTime = now
	}

	return nil
}

// BeforeUpdate 更新前更新时间戳
func (w *WebShell) BeforeUpdate(tx *gorm.DB) error {
	w.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	return nil
}
