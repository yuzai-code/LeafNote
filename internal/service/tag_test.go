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

func setupTagTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 自动迁移
	err = db.AutoMigrate(&model.Tag{}, &model.Note{})
	assert.NoError(t, err)

	return db
}

func TestTagService_CreateTag(t *testing.T) {
	db := setupTagTestDB(t)
	s := NewTagService(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		tag     *model.Tag
		wantErr bool
	}{
		{
			name: "创建成功",
			tag: &model.Tag{
				Name: "测试标签",
			},
			wantErr: false,
		},
		{
			name: "名称为空",
			tag: &model.Tag{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "父标签不存在",
			tag: &model.Tag{
				Name:     "子标签",
				ParentID: testutil.StringPtr("not-exist"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.CreateTag(ctx, tt.tag)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tt.tag.BaseModel.ID)
			}
		})
	}
}

func TestTagService_GetTagByID(t *testing.T) {
	db := setupTagTestDB(t)
	s := NewTagService(db)
	ctx := context.Background()

	// 创建测试数据
	tag := &model.Tag{Name: "测试标签"}
	err := s.CreateTag(ctx, tag)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "获取成功",
			id:      tag.BaseModel.ID,
			wantErr: false,
		},
		{
			name:    "标签不存在",
			id:      "not-exist",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetTagByID(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, got.BaseModel.ID)
			}
		})
	}
}

func TestTagService_ListTags(t *testing.T) {
	db := setupTagTestDB(t)
	s := NewTagService(db)
	ctx := context.Background()

	// 创建测试数据
	parent := &model.Tag{Name: "父标签"}
	err := s.CreateTag(ctx, parent)
	assert.NoError(t, err)

	child := &model.Tag{
		Name:     "子标签",
		ParentID: &parent.BaseModel.ID,
	}
	err = s.CreateTag(ctx, child)
	assert.NoError(t, err)

	t.Run("列表查询", func(t *testing.T) {
		tags, err := s.ListTags(ctx)
		assert.NoError(t, err)
		assert.Len(t, tags, 1) // 只返回顶级标签
		assert.Equal(t, parent.BaseModel.ID, tags[0].BaseModel.ID)
		assert.Len(t, tags[0].Children, 1) // 包含子标签
	})
}

func TestTagService_UpdateTag(t *testing.T) {
	db := setupTagTestDB(t)
	s := NewTagService(db)
	ctx := context.Background()

	// 创建测试数据
	tag := &model.Tag{Name: "原始标签"}
	err := s.CreateTag(ctx, tag)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		tag     *model.Tag
		wantErr bool
	}{
		{
			name: "更新成功",
			tag: &model.Tag{
				BaseModel: model.BaseModel{ID: tag.BaseModel.ID},
				Name:      "新标签名",
			},
			wantErr: false,
		},
		{
			name: "ID为空",
			tag: &model.Tag{
				Name: "新标签名",
			},
			wantErr: true,
		},
		{
			name: "名称为空",
			tag: &model.Tag{
				BaseModel: model.BaseModel{ID: tag.BaseModel.ID},
				Name:      "",
			},
			wantErr: true,
		},
		{
			name: "父标签是自己",
			tag: &model.Tag{
				BaseModel: model.BaseModel{ID: tag.BaseModel.ID},
				Name:      "新标签名",
				ParentID:  &tag.BaseModel.ID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.UpdateTag(ctx, tt.tag)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// 验证更新结果
				updated, err := s.GetTagByID(ctx, tt.tag.BaseModel.ID)
				assert.NoError(t, err)
				assert.Equal(t, tt.tag.Name, updated.Name)
			}
		})
	}
}

func TestTagService_DeleteTag(t *testing.T) {
	db := setupTagTestDB(t)
	s := NewTagService(db)
	ctx := context.Background()

	// 创建测试数据
	parent := &model.Tag{Name: "父标签"}
	err := s.CreateTag(ctx, parent)
	assert.NoError(t, err)

	child := &model.Tag{
		Name:     "子标签",
		ParentID: &parent.BaseModel.ID,
	}
	err = s.CreateTag(ctx, child)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "删除失败-存在子标签",
			id:      parent.BaseModel.ID,
			wantErr: true,
		},
		{
			name:    "删除成功-子标签",
			id:      child.BaseModel.ID,
			wantErr: false,
		},
		{
			name:    "删除成功-父标签",
			id:      parent.BaseModel.ID,
			wantErr: false,
		},
		{
			name:    "标签不存在",
			id:      "not-exist",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.DeleteTag(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// 验证删除结果
				_, err := s.GetTagByID(ctx, tt.id)
				assert.Error(t, err)
			}
		})
	}
}

// 辅助函数：创建字符串指针
func stringPtr(s string) *string {
	return &s
}
