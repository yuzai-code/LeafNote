package model

// Tag 标签模型
type Tag struct {
	BaseModel
	Name     string  `gorm:"not null" json:"name"`                // 标签名称
	ParentID *string `gorm:"type:varchar(36)" json:"parent_id"`   // 父标签ID
	Parent   *Tag    `gorm:"foreignKey:ParentID" json:"parent"`   // 父标签
	Children []Tag   `gorm:"foreignKey:ParentID" json:"children"` // 子标签
	Notes    []Note  `gorm:"many2many:note_tags;" json:"notes"`   // 关联的笔记
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}
