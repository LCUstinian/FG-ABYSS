package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project 项目实体
type Project struct {
	ID          string         `gorm:"type:text;uniqueIndex;not null" json:"id"`
	Name        string         `gorm:"type:text;uniqueIndex;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Status      int            `gorm:"type:integer;default:0" json:"status"` // 0=Active, 1=Archived
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt"` // 软删除标记
}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}

// BeforeCreate 创建前生成 UUID
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// Validate 验证项目数据
func (p *Project) Validate() error {
	// 领域规则验证
	if p.Name == "" {
		return &ValidationError{Field: "name", Message: "项目名称不能为空"}
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
