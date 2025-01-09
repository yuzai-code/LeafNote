package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"leafnote/internal/model"
	"leafnote/internal/testutil"
)

func setupCategoryTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 自动迁移
	err = db.AutoMigrate(&model.Category{}, &model.Note{})
	assert.NoError(t, err)

	return db
}

func TestCategoryService_CreateCategory(t *testing.T) {
	db := setupCategoryTestDB(t)
	s := NewCategoryService(db)
	ctx := context.Background()

	// 创建父目录
	parent := &model.Category{Name: "父目录"}
	err := s.CreateCategory(ctx, parent)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		cat     *model.Category
		wantErr bool
	}{
		{
			name: "创建成功",
			cat: &model.Category{
				Name: "测试目录",
			},
			wantErr: false,
		},
		{
			name: "名称为空",
			cat: &model.Category{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "创建子目录成功",
			cat: &model.Category{
				Name:     "子目录",
				ParentID: &parent.BaseModel.ID,
			},
			wantErr: false,
		},
		{
			name: "父目录不存在",
			cat: &model.Category{
				Name:     "子目录",
				ParentID: testutil.StringPtr("not-exist"),
			},
			wantErr: true,
		},
		{
			name: "路径重复",
			cat: &model.Category{
				Name: "测试目录", // 与第一个测试用例同名
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.CreateCategory(ctx, tt.cat)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tt.cat.BaseModel.ID)
				assert.NotEmpty(t, tt.cat.Path)
				if tt.cat.ParentID != nil {
					assert.Contains(t, tt.cat.Path, parent.Path)
				}
			}
		})
	}
}

func TestCategoryService_GetCategoryByID(t *testing.T) {
	db := setupCategoryTestDB(t)
	s := NewCategoryService(db)
	ctx := context.Background()

	// 创建测试数据
	parent := &model.Category{Name: "父目录"}
	err := s.CreateCategory(ctx, parent)
	assert.NoError(t, err)

	child := &model.Category{
		Name:     "子目录",
		ParentID: &parent.BaseModel.ID,
	}
	err = s.CreateCategory(ctx, child)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "获取成功",
			id:      parent.BaseModel.ID,
			wantErr: false,
		},
		{
			name:    "获取子目录成功",
			id:      child.BaseModel.ID,
			wantErr: false,
		},
		{
			name:    "目录不存在",
			id:      "not-exist",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetCategoryByID(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, got.BaseModel.ID)
				if tt.id == parent.BaseModel.ID {
					assert.Len(t, got.Children, 1)
				}
			}
		})
	}
}

func TestCategoryService_ListCategories(t *testing.T) {
	db := setupCategoryTestDB(t)
	s := NewCategoryService(db)
	ctx := context.Background()

	// 创建测试数据
	parent1 := &model.Category{Name: "父目录1"}
	err := s.CreateCategory(ctx, parent1)
	assert.NoError(t, err)

	parent2 := &model.Category{Name: "父目录2"}
	err = s.CreateCategory(ctx, parent2)
	assert.NoError(t, err)

	child1 := &model.Category{
		Name:     "子目录1",
		ParentID: &parent1.BaseModel.ID,
	}
	err = s.CreateCategory(ctx, child1)
	assert.NoError(t, err)

	child2 := &model.Category{
		Name:     "子目录2",
		ParentID: &parent2.BaseModel.ID,
	}
	err = s.CreateCategory(ctx, child2)
	assert.NoError(t, err)

	t.Run("列表查询", func(t *testing.T) {
		categories, err := s.ListCategories(ctx)
		assert.NoError(t, err)
		assert.Len(t, categories, 2) // 只返回顶级目录
		for _, cat := range categories {
			assert.Len(t, cat.Children, 1) // 每个顶级目录都有一个子目录
		}
	})
}

func TestCategoryService_UpdateCategory(t *testing.T) {
	db := setupCategoryTestDB(t)
	s := NewCategoryService(db)
	ctx := context.Background()

	// 创建测试数据
	parent := &model.Category{Name: "父目录"}
	err := s.CreateCategory(ctx, parent)
	assert.NoError(t, err)

	category := &model.Category{Name: "原始目录"}
	err = s.CreateCategory(ctx, category)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		cat     *model.Category
		wantErr bool
	}{
		{
			name: "更新成功",
			cat: &model.Category{
				BaseModel: model.BaseModel{ID: category.BaseModel.ID},
				Name:      "新目录名",
			},
			wantErr: false,
		},
		{
			name: "更新为子目录",
			cat: &model.Category{
				BaseModel: model.BaseModel{ID: category.BaseModel.ID},
				Name:      "子目录",
				ParentID:  &parent.BaseModel.ID,
			},
			wantErr: false,
		},
		{
			name: "父目录不存在",
			cat: &model.Category{
				BaseModel: model.BaseModel{ID: category.BaseModel.ID},
				Name:      "子目录",
				ParentID:  testutil.StringPtr("not-exist"),
			},
			wantErr: true,
		},
		{
			name: "父目录是自己",
			cat: &model.Category{
				BaseModel: model.BaseModel{ID: category.BaseModel.ID},
				Name:      "子目录",
				ParentID:  &category.BaseModel.ID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.UpdateCategory(ctx, tt.cat)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// 验证更新结果
				updated, err := s.GetCategoryByID(ctx, tt.cat.BaseModel.ID)
				assert.NoError(t, err)
				assert.Equal(t, tt.cat.Name, updated.Name)
				if tt.cat.ParentID != nil {
					assert.Equal(t, *tt.cat.ParentID, *updated.ParentID)
					assert.Contains(t, updated.Path, parent.Path)
				}
			}
		})
	}
}

func TestCategoryService_DeleteCategory(t *testing.T) {
	db := setupCategoryTestDB(t)
	s := NewCategoryService(db)
	ctx := context.Background()

	// 创建测试数据
	parent := &model.Category{Name: "父目录"}
	err := s.CreateCategory(ctx, parent)
	assert.NoError(t, err)

	child := &model.Category{
		Name:     "子目录",
		ParentID: &parent.BaseModel.ID,
	}
	err = s.CreateCategory(ctx, child)
	assert.NoError(t, err)

	// 创建一个带笔记的目录
	categoryWithNote := &model.Category{Name: "带笔记的目录"}
	err = s.CreateCategory(ctx, categoryWithNote)
	assert.NoError(t, err)

	// 创建笔记
	note := &model.Note{
		Title:      "测试笔记",
		CategoryID: &categoryWithNote.BaseModel.ID,
	}
	err = db.Create(note).Error
	assert.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "删除失败-存在子目录",
			id:      parent.BaseModel.ID,
			wantErr: true,
		},
		{
			name:    "删除失败-存在笔记",
			id:      categoryWithNote.BaseModel.ID,
			wantErr: true,
		},
		{
			name:    "删除成功-子目录",
			id:      child.BaseModel.ID,
			wantErr: false,
		},
		{
			name:    "删除成功-父目录",
			id:      parent.BaseModel.ID,
			wantErr: false,
		},
		{
			name:    "目录不存在",
			id:      "not-exist",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.DeleteCategory(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// 验证删除结果
				_, err := s.GetCategoryByID(ctx, tt.id)
				assert.Error(t, err)
			}
		})
	}
}
