package model

// Category 分类模型
type Category struct {
	BaseModel
	Name     string     `gorm:"size:50;not null" json:"name"`
	ParentID *uint      `json:"parent_id"`
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Notes    []Note     `gorm:"foreignKey:CategoryID" json:"notes,omitempty"`
}
