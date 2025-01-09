package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"leafnote/internal/model"
)

// TagService 标签服务
type TagService struct {
	db *gorm.DB
}

// NewTagService 创建标签服务实例
func NewTagService(db *gorm.DB) *TagService {
	return &TagService{db: db}
}

// CreateTag 创建标签
func (s *TagService) CreateTag(ctx context.Context, tag *model.Tag) error {
	if tag.Name == "" {
		return errors.New("标签名称不能为空")
	}

	// 如果有父标签ID，检查父标签是否存在
	if tag.ParentID != nil {
		var parent model.Tag
		if err := s.db.First(&parent, "id = ?", *tag.ParentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("父标签不存在")
			}
			return err
		}
	}

	// 生成UUID
	tag.ID = uuid.New().String()
	return s.db.Create(tag).Error
}

// GetTagByID 根据ID获取标签
func (s *TagService) GetTagByID(ctx context.Context, id string) (*model.Tag, error) {
	var tag model.Tag
	err := s.db.Preload("Children").First(&tag, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("标签不存在")
		}
		return nil, err
	}
	return &tag, nil
}

// ListTags 获取标签列表
func (s *TagService) ListTags(ctx context.Context) ([]model.Tag, error) {
	var tags []model.Tag
	// 只获取顶级标签（没有父标签的标签）
	err := s.db.Preload("Children").Where("parent_id IS NULL").Find(&tags).Error
	return tags, err
}

// UpdateTag 更新标签
func (s *TagService) UpdateTag(ctx context.Context, tag *model.Tag) error {
	if tag.ID == "" {
		return errors.New("标签ID不能为空")
	}
	if tag.Name == "" {
		return errors.New("标签名称不能为空")
	}

	// 如果有父标签ID，检查父标签是否存在且不能是自己
	if tag.ParentID != nil {
		if *tag.ParentID == tag.ID {
			return errors.New("父标签不能是自己")
		}
		var parent model.Tag
		if err := s.db.First(&parent, "id = ?", *tag.ParentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("父标签不存在")
			}
			return err
		}
	}

	return s.db.Model(tag).Updates(map[string]interface{}{
		"name":      tag.Name,
		"parent_id": tag.ParentID,
	}).Error
}

// DeleteTag 删除标签
func (s *TagService) DeleteTag(ctx context.Context, id string) error {
	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 检查是否有子标签
		var count int64
		if err := tx.Model(&model.Tag{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("请先删除子标签")
		}

		// 删除标签与笔记的关联关系
		if err := tx.Exec("DELETE FROM note_tags WHERE tag_id = ?", id).Error; err != nil {
			return err
		}

		// 删除标签
		if err := tx.Delete(&model.Tag{}, "id = ?", id).Error; err != nil {
			return err
		}

		return nil
	})
}
