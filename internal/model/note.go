package model

// Note 笔记模型
type Note struct {
	BaseModel
	Title      string   `gorm:"size:200;not null" json:"title"`
	Content    string   `gorm:"type:text" json:"content"`
	CategoryID *uint    `json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags       string   `gorm:"type:text" json:"tags"`
}
