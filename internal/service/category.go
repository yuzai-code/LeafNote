package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"leafnote/internal/model"
)

// CategoryService 目录服务
type CategoryService struct {
	db *gorm.DB
}

// NewCategoryService 创建目录服务实例
func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{db: db}
}

// CreateCategory 创建目录
func (s *CategoryService) CreateCategory(ctx context.Context, category *model.Category) error {
	if category.Name == "" {
		return errors.New("目录名称不能为空")
	}

	// 如果有父目录ID，检查父目录是否存在
	if category.ParentID != nil {
		var parent model.Category
		if err := s.db.First(&parent, "id = ?", *category.ParentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("父目录不存在")
			}
			return err
		}
		// 设置完整路径
		category.Path = parent.Path + "/" + category.Name
	} else {
		category.Path = "/" + category.Name
	}

	// 检查路径是否已存在
	var count int64
	if err := s.db.Model(&model.Category{}).Where("path = ?", category.Path).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("目录路径已存在")
	}

	// 生成UUID
	category.BaseModel.ID = uuid.New().String()
	return s.db.Create(category).Error
}

// GetCategoryByID 根据ID获取目录
func (s *CategoryService) GetCategoryByID(ctx context.Context, id string) (*model.Category, error) {
	var category model.Category
	err := s.db.Preload("Children").First(&category, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("目录不存在")
		}
		return nil, err
	}
	return &category, nil
}

// ListCategories 获取目录列表
func (s *CategoryService) ListCategories(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	// 只获取顶级目录（没有父目录的目录）
	err := s.db.Preload("Children").Where("parent_id IS NULL").Find(&categories).Error
	return categories, err
}

// UpdateCategory 更新目录
func (s *CategoryService) UpdateCategory(ctx context.Context, category *model.Category) error {
	if category.BaseModel.ID == "" {
		return errors.New("目录ID不能为空")
	}
	if category.Name == "" {
		return errors.New("目录名称不能为空")
	}

	// 获取原始目录信息
	var oldCategory model.Category
	if err := s.db.First(&oldCategory, "id = ?", category.BaseModel.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("目录不存在")
		}
		return err
	}

	// 如果有父目录ID，检查父目录是否存在且不能是自己
	if category.ParentID != nil {
		if *category.ParentID == category.BaseModel.ID {
			return errors.New("父目录不能是自己")
		}
		var parent model.Category
		if err := s.db.First(&parent, "id = ?", *category.ParentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("父目录不存在")
			}
			return err
		}
		// 设置新的完整路径
		category.Path = parent.Path + "/" + category.Name
	} else {
		category.Path = "/" + category.Name
	}

	// 如果路径发生变化，检查新路径是否已存在
	if category.Path != oldCategory.Path {
		var count int64
		if err := s.db.Model(&model.Category{}).Where("path = ? AND id != ?", category.Path, category.BaseModel.ID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("目录路径已存在")
		}

		// 更新所有子目录的路径
		if err := s.updateChildrenPaths(oldCategory.Path, category.Path); err != nil {
			return err
		}
	}

	return s.db.Model(category).Updates(map[string]interface{}{
		"name":      category.Name,
		"parent_id": category.ParentID,
		"path":      category.Path,
	}).Error
}

// updateChildrenPaths 更新所有子目录的路径
func (s *CategoryService) updateChildrenPaths(oldParentPath, newParentPath string) error {
	// 查找所有以oldParentPath开头的目录
	var categories []model.Category
	if err := s.db.Where("path LIKE ?", oldParentPath+"/%").Find(&categories).Error; err != nil {
		return err
	}

	// 更新每个子目录的路径
	for _, category := range categories {
		newPath := newParentPath + category.Path[len(oldParentPath):]
		if err := s.db.Model(&category).Update("path", newPath).Error; err != nil {
			return err
		}
	}

	return nil
}

// DeleteCategory 删除目录
func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 检查目录是否存在
		var category model.Category
		if err := tx.First(&category, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("目录不存在")
			}
			return err
		}

		// 检查是否有子目录
		var count int64
		if err := tx.Model(&model.Category{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("请先删除子目录")
		}

		// 检查目录下是否有笔记
		if err := tx.Model(&model.Note{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("请先删除目录下的笔记")
		}

		// 删除目录
		if err := tx.Delete(&category).Error; err != nil {
			return err
		}

		return nil
	})
}
