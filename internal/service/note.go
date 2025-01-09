package service

import (
	"crypto/md5"
	"encoding/hex"
	"leafnote/internal/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NoteService 处理笔记相关的业务逻辑
type NoteService struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewNoteService 创建笔记服务实例
func NewNoteService(db *gorm.DB, logger *zap.Logger) *NoteService {
	return &NoteService{
		db:     db,
		logger: logger,
	}
}

// CreateNoteInput 创建笔记的输入参数
type CreateNoteInput struct {
	Title      string
	Content    string
	YAMLMeta   string
	FilePath   string
	CategoryID *string
	TagIDs     []string
}

// UpdateNoteInput 更新笔记的输入参数
type UpdateNoteInput struct {
	Title      string
	Content    string
	YAMLMeta   string
	CategoryID *string
	TagIDs     []string
}

// ListNotes 获取笔记列表
func (s *NoteService) ListNotes() ([]model.Note, error) {
	var notes []model.Note
	err := s.db.Preload("Category").Preload("Tags").Find(&notes).Error
	return notes, err
}

// CreateNote 创建笔记
func (s *NoteService) CreateNote(input CreateNoteInput) (*model.Note, error) {
	note := &model.Note{
		Title:      input.Title,
		Content:    input.Content,
		YAMLMeta:   input.YAMLMeta,
		FilePath:   input.FilePath,
		CategoryID: input.CategoryID,
		Version:    1,
		Checksum:   s.calculateChecksum(input.Content),
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(note).Error; err != nil {
			return err
		}

		if len(input.TagIDs) > 0 {
			var tags []model.Tag
			if err := tx.Find(&tags, "id IN ?", input.TagIDs).Error; err != nil {
				return err
			}
			if err := tx.Model(note).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return note, nil
}

// GetNote 获取单个笔记
func (s *NoteService) GetNote(id string) (*model.Note, error) {
	var note model.Note
	err := s.db.Preload("Category").Preload("Tags").First(&note, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// UpdateNote 更新笔记
func (s *NoteService) UpdateNote(id string, input UpdateNoteInput) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var note model.Note
		if err := tx.First(&note, "id = ?", id).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"version": note.Version + 1,
		}
		if input.Title != "" {
			updates["title"] = input.Title
		}
		if input.Content != "" {
			updates["content"] = input.Content
			updates["checksum"] = s.calculateChecksum(input.Content)
		}
		if input.YAMLMeta != "" {
			updates["yaml_meta"] = input.YAMLMeta
		}
		if input.CategoryID != nil {
			updates["category_id"] = input.CategoryID
		}

		if err := tx.Model(&note).Updates(updates).Error; err != nil {
			return err
		}

		if len(input.TagIDs) > 0 {
			var tags []model.Tag
			if err := tx.Find(&tags, "id IN ?", input.TagIDs).Error; err != nil {
				return err
			}
			if err := tx.Model(&note).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteNote 删除笔记
func (s *NoteService) DeleteNote(id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var note model.Note
		if err := tx.First(&note, "id = ?", id).Error; err != nil {
			return err
		}

		if err := tx.Model(&note).Association("Tags").Clear(); err != nil {
			return err
		}
		return tx.Delete(&note).Error
	})
}

// calculateChecksum 计算内容的校验和
func (s *NoteService) calculateChecksum(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}
