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

	// 检查同级目录下是否存在同名目录
	var count int64
	query := s.db.Model(&model.Category{}).Where("name = ?", category.Name)

	if category.ParentID != nil {
		// 如果是子目录，检查同一父目录下是否有同名目录
		query = query.Where("parent_id = ?", *category.ParentID)

		// 检查父目录是否存在
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
		// 如果是顶级目录，检查是否有同名的顶级目录
		query = query.Where("parent_id IS NULL")
		category.Path = "/" + category.Name
	}

	// 执行查询
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("同级目录下已存在同名目录")
	}

	// 检查路径是否已存在（使用 Unscoped 忽略软删除）
	query = s.db.Unscoped().Model(&model.Category{}).Where("path = ?", category.Path)
	if err := query.Count(&count).Error; err != nil {
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
	err := s.db.Preload("Children").Preload("Notes").First(&category, "id = ?", id).Error
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

	// 首先获取所有顶级目录，并预加载 notes
	if err := s.db.Preload("Notes").Where("parent_id IS NULL").Find(&categories).Error; err != nil {
		return nil, err
	}

	// 递归加载所有子目录
	for i := range categories {
		if err := s.loadChildren(ctx, &categories[i]); err != nil {
			return nil, err
		}
	}

	return categories, nil
}

// loadChildren 递归加载目录的子目录
func (s *CategoryService) loadChildren(ctx context.Context, category *model.Category) error {
	// 查询当前目录的直接子目录，并预加载 notes
	var children []model.Category
	if err := s.db.Preload("Notes").Where("parent_id = ?", category.ID).Find(&children).Error; err != nil {
		return err
	}

	// 如果有子目录，递归加载它们的子目录
	for i := range children {
		if err := s.loadChildren(ctx, &children[i]); err != nil {
			return err
		}
	}

	category.Children = children
	return nil
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
		query := s.db.Unscoped().Model(&model.Category{}).
			Where("path = ? AND id != ?", category.Path, category.BaseModel.ID)

		if category.ParentID != nil {
			// 如果是子目录，还要检查父目录ID
			query = query.Where("parent_id = ?", *category.ParentID)
		} else {
			// 如果是顶级目录，确保没有同名的顶级目录
			query = query.Where("parent_id IS NULL")
		}

		if err := query.Count(&count).Error; err != nil {
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

		// 递归删除所有子目录及其笔记
		if err := s.recursiveDelete(tx, id); err != nil {
			return err
		}

		// 删除当前目录
		if err := tx.Unscoped().Delete(&category).Error; err != nil {
			return err
		}

		return nil
	})
}

// recursiveDelete 递归删除目录及其所有子目录和笔记
func (s *CategoryService) recursiveDelete(tx *gorm.DB, categoryID string) error {
	// 获取所有子目录
	var children []model.Category
	if err := tx.Where("parent_id = ?", categoryID).Find(&children).Error; err != nil {
		return err
	}

	// 递归删除每个子目录
	for _, child := range children {
		if err := s.recursiveDelete(tx, child.ID); err != nil {
			return err
		}
		// 删除子目录
		if err := tx.Unscoped().Delete(&child).Error; err != nil {
			return err
		}
	}

	// 删除当前目录下的所有笔记
	if err := tx.Unscoped().Where("category_id = ?", categoryID).Delete(&model.Note{}).Error; err != nil {
		return err
	}

	return nil
}
