package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Note 笔记模型
type Note struct {
	BaseModel
	Title        string    `gorm:"not null" json:"title"`                       // 笔记标题
	Content      string    `gorm:"type:text" json:"content"`                    // 加密后的笔记内容
	YAMLMeta     string    `gorm:"column:yaml_meta;type:text" json:"yaml_meta"` // 加密后的YAML元数据
	FilePath     string    `gorm:"not null" json:"file_path"`                   // 文件路径
	OriginalPath string    `json:"original_path,omitempty"`                     // 原始路径（用于恢复）
	InTrash      bool      `gorm:"default:false" json:"in_trash"`               // 是否在回收站
	TrashTime    time.Time `json:"trash_time,omitempty"`                        // 放入回收站时间
	Version      int       `gorm:"not null;default:1" json:"version"`           // 版本号
	Checksum     string    `gorm:"not null" json:"checksum"`                    // 内容校验和
	CategoryID   *string   `gorm:"type:varchar(36)" json:"category_id"`         // 所属目录ID
	Category     *Category `gorm:"foreignKey:CategoryID" json:"category"`       // 所属目录
	Tags         []Tag     `gorm:"many2many:note_tags;" json:"tags"`            // 关联的标签
}

// TableName 指定表名
func (Note) TableName() string {
	return "notes"
}

// BeforeDelete 在软删除前修改文件路径
func (n *Note) BeforeDelete(tx *gorm.DB) error {
	// 在软删除时添加时间戳后缀，避免路径冲突
	n.FilePath = fmt.Sprintf("%s.deleted.%s",
		n.FilePath,
		time.Now().Format("20060102150405"))
	return tx.Save(n).Error
}

// BeforeUnscoped 在硬删除前清理关联
func (n *Note) BeforeUnscoped(tx *gorm.DB) error {
	// 清理标签关联
	if err := tx.Model(n).Association("Tags").Clear(); err != nil {
		return err
	}
	return nil
}
