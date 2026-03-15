package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WebShell WebShell 实体
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

// Validate 验证 WebShell 数据
func (w *WebShell) Validate() error {
	// 领域规则验证
	if w.Url == "" {
		return &ValidationError{Field: "url", Message: "URL 不能为空"}
	}
	if w.ProjectID == "" {
		return &ValidationError{Field: "projectId", Message: "项目 ID 不能为空"}
	}
	return nil
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
