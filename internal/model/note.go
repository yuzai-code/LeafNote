package model

// Note 笔记模型
type Note struct {
	BaseModel
	Title      string    `gorm:"not null" json:"title"`                       // 笔记标题
	Content    string    `gorm:"type:text" json:"content"`                    // 加密后的笔记内容
	YAMLMeta   string    `gorm:"column:yaml_meta;type:text" json:"yaml_meta"` // 加密后的YAML元数据
	FilePath   string    `gorm:"not null;uniqueIndex" json:"file_path"`       // 文件路径
	Version    int       `gorm:"not null;default:1" json:"version"`           // 版本号
	Checksum   string    `gorm:"not null" json:"checksum"`                    // 内容校验和
	CategoryID *string   `gorm:"type:varchar(36)" json:"category_id"`         // 所属目录ID
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category"`       // 所属目录
	Tags       []Tag     `gorm:"many2many:note_tags;" json:"tags"`            // 关联的标签
}

// TableName 指定表名
func (Note) TableName() string {
	return "notes"
}
