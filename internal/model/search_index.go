package model

// SearchIndexType 搜索索引类型
type SearchIndexType string

const (
	SearchIndexTypeTitle   SearchIndexType = "title"   // 标题索引
	SearchIndexTypeContent SearchIndexType = "content" // 内容索引
	SearchIndexTypeTag     SearchIndexType = "tag"     // 标签索引
)

// SearchIndex 搜索索引模型
type SearchIndex struct {
	BaseModel
	NoteID  string          `gorm:"type:varchar(36);not null;index" json:"note_id"` // 关联的笔记ID
	Content string          `gorm:"type:text;not null" json:"content"`              // 加密的搜索索引内容
	Type    SearchIndexType `gorm:"type:varchar(10);not null" json:"type"`          // 索引类型
	Note    Note            `gorm:"foreignKey:NoteID" json:"note"`                  // 关联的笔记
}

// TableName 指定表名
func (SearchIndex) TableName() string {
	return "search_index"
}
