package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project 项目模型
type Project struct {
	ID          string         `gorm:"type:text;uniqueIndex;not null" json:"id"`
	Name        string         `gorm:"type:text;uniqueIndex;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Status      int            `gorm:"type:integer;default:0" json:"status"` // 0=Active, 1=Archived
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前生成 UUID
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}
