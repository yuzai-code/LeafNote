package model

import "gorm.io/gorm"

// Category 目录模型
type Category struct {
	BaseModel
	Name     string     `gorm:"not null" json:"name"`                // 目录名称
	Path     string     `gorm:"not null;uniqueIndex" json:"path"`    // 完整路径
	ParentID *string    `gorm:"type:varchar(36)" json:"parent_id"`   // 父目录ID
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent"`   // 父目录
	Children []Category `gorm:"foreignKey:ParentID" json:"children"` // 子目录
	Notes    []Note     `gorm:"foreignKey:CategoryID" json:"notes"`  // 目录下的笔记
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

// BeforeCreate 在创建记录前处理路径
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if err := c.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 如果有父目录，拼接父目录的路径
	if c.ParentID != nil {
		var parent Category
		if err := tx.First(&parent, "id = ?", c.ParentID).Error; err != nil {
			return err
		}
		c.Path = parent.Path + "/" + c.Name
	} else {
		c.Path = "/" + c.Name
	}
	return nil
}

// BeforeUpdate 在更新记录前处理路径
func (c *Category) BeforeUpdate(tx *gorm.DB) error {
	// 如果名称或父目录发生变化，需要更新路径
	if tx.Statement.Changed("Name") || tx.Statement.Changed("ParentID") {
		if c.ParentID != nil {
			var parent Category
			if err := tx.First(&parent, "id = ?", c.ParentID).Error; err != nil {
				return err
			}
			c.Path = parent.Path + "/" + c.Name
		} else {
			c.Path = "/" + c.Name
		}
	}
	return nil
}
