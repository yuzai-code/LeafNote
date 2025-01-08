package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base 模型基础结构
type Base struct {
	ID        string         `gorm:"type:varchar(36);primary_key" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// BeforeCreate 在创建记录前生成UUID
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}
