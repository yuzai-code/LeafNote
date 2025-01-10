package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"leafnote/internal/model"
	"os"
	"path/filepath"
	"strings"

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
func (s *NoteService) ListNotes(categoryID *string) ([]model.Note, error) {
	var notes []model.Note
	query := s.db.Preload("Category").Preload("Tags")

	// 如果指定了分类ID，则只获取该分类下的笔记
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	err := query.Find(&notes).Error
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
	// 检查文件是否存在
	var count int64
	if err := s.db.Model(&model.Note{}).Where("file_path = ?", input.FilePath).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("文件路径已存在")
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
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("笔记不存在")
			}
			return err
		}

		if err := tx.Model(&note).Association("Tags").Clear(); err != nil {
			return err
		}
		return tx.Delete(&note).Error
	})
}

// RestoreNote 恢复笔记
func (s *NoteService) RestoreNote(id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var note model.Note
		// Unscoped 可以查询到已软删除的记录
		if err := tx.Unscoped().First(&note, "id = ?", id).Error; err != nil {
			return err
		}

		// 移除 .deleted 后缀
		originalPath := strings.TrimSuffix(
			note.FilePath,
			note.FilePath[strings.LastIndex(note.FilePath, ".deleted."):],
		)

		// 检查原始路径是否可用
		var count int64
		if err := tx.Model(&model.Note{}).
			Where("file_path = ? AND deleted_at IS NULL", originalPath).
			Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			// 如果原始路径被占用，生成新路径
			note.FilePath = s.generateUniqueFilePath(originalPath)
		} else {
			note.FilePath = originalPath
		}

		// 恢复笔记（清除删除标记）
		return tx.Unscoped().Model(&note).Updates(map[string]interface{}{
			"deleted_at": nil,
			"file_path":  note.FilePath,
		}).Error
	})
}

// 生成唯一的文件路径
func (s *NoteService) generateUniqueFilePath(originalPath string) string {
	ext := filepath.Ext(originalPath)
	basePath := strings.TrimSuffix(originalPath, ext)
	counter := 1
	newPath := originalPath

	for {
		var count int64
		s.db.Model(&model.Note{}).
			Where("file_path = ? AND deleted_at IS NULL", newPath).
			Count(&count)
		if count == 0 {
			break
		}
		newPath = fmt.Sprintf("%s_%d%s", basePath, counter, ext)
		counter++
	}

	return newPath
}

// calculateChecksum 计算内容的校验和
func (s *NoteService) calculateChecksum(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}

// ExportToMarkdown 导出笔记为 Markdown 文件
func (s *NoteService) ExportToMarkdown(id string, exportPath string) error {
	var note model.Note
	if err := s.db.First(&note, "id = ?", id).Error; err != nil {
		return err
	}

	// 构建完整的文件路径
	fullPath := filepath.Join(exportPath, note.FilePath)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	// 构建 Markdown 内容
	content := fmt.Sprintf("---\n%s\n---\n\n%s", note.YAMLMeta, note.Content)

	// 写入文件
	return os.WriteFile(fullPath, []byte(content), 0644)
}
