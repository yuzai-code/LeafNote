package service

import (
	"leafnote/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&model.Note{}, &model.Tag{}, &model.Category{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestNoteService_CreateNote(t *testing.T) {
	db := setupTestDB(t)
	logger, _ := zap.NewDevelopment()
	service := NewNoteService(db, logger)

	tests := []struct {
		name    string
		input   CreateNoteInput
		wantErr bool
	}{
		{
			name: "正常创建笔记",
			input: CreateNoteInput{
				Title:    "测试笔记",
				Content:  "测试内容",
				YAMLMeta: "tags: [test]",
				FilePath: "/test/note.md",
			},
			wantErr: false,
		},
		{
			name: "重复的文件路径",
			input: CreateNoteInput{
				Title:    "测试笔记2",
				Content:  "测试内容2",
				YAMLMeta: "tags: [test]",
				FilePath: "/test/note.md", // 重复的路径
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note, err := service.CreateNote(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, note)
			assert.Equal(t, tt.input.Title, note.Title)
			assert.Equal(t, tt.input.Content, note.Content)
			assert.Equal(t, tt.input.FilePath, note.FilePath)
		})
	}
}

func TestNoteService_GetNote(t *testing.T) {
	db := setupTestDB(t)
	logger, _ := zap.NewDevelopment()
	service := NewNoteService(db, logger)

	// 创建测试数据
	note, err := service.CreateNote(CreateNoteInput{
		Title:    "测试笔记",
		Content:  "测试内容",
		YAMLMeta: "tags: [test]",
		FilePath: "/test/note.md",
	})
	assert.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "获取存在的笔记",
			id:      note.ID,
			wantErr: false,
		},
		{
			name:    "获取不存在的笔记",
			id:      "not-exist",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetNote(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, note.ID, got.ID)
			assert.Equal(t, note.Title, got.Title)
		})
	}
}

func TestNoteService_UpdateNote(t *testing.T) {
	db := setupTestDB(t)
	logger, _ := zap.NewDevelopment()
	service := NewNoteService(db, logger)

	// 创建测试数据
	note, err := service.CreateNote(CreateNoteInput{
		Title:    "测试笔记",
		Content:  "测试内容",
		YAMLMeta: "tags: [test]",
		FilePath: "/test/note.md",
	})
	assert.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		input   UpdateNoteInput
		wantErr bool
	}{
		{
			name: "更新存在的笔记",
			id:   note.ID,
			input: UpdateNoteInput{
				Title:   "更新后的标题",
				Content: "更新后的内容",
			},
			wantErr: false,
		},
		{
			name: "更新不存在的笔记",
			id:   "not-exist",
			input: UpdateNoteInput{
				Title: "更新后的标题",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdateNote(tt.id, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// 验证更新结果
			updated, err := service.GetNote(tt.id)
			assert.NoError(t, err)
			if tt.input.Title != "" {
				assert.Equal(t, tt.input.Title, updated.Title)
			}
			if tt.input.Content != "" {
				assert.Equal(t, tt.input.Content, updated.Content)
			}
		})
	}
}

func TestNoteService_DeleteNote(t *testing.T) {
	db := setupTestDB(t)
	logger, _ := zap.NewDevelopment()
	service := NewNoteService(db, logger)

	// 创建测试数据
	note, err := service.CreateNote(CreateNoteInput{
		Title:    "测试笔记",
		Content:  "测试内容",
		YAMLMeta: "tags: [test]",
		FilePath: "/test/note.md",
	})
	assert.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "删除存在的笔记",
			id:      note.ID,
			wantErr: false,
		},
		{
			name:    "删除不存在的笔记",
			id:      "not-exist",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteNote(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// 验证删除结果
			_, err = service.GetNote(tt.id)
			assert.Error(t, err)
		})
	}
}
