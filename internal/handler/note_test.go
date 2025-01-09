package handler

import (
	"bytes"
	"encoding/json"
	"leafnote/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestHandler(t *testing.T) (*Handler, *gin.Engine) {
	// 设置测试数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&model.Note{}, &model.Tag{}, &model.Category{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// 设置日志
	logger, _ := zap.NewDevelopment()

	// 创建处理器
	h := NewHandler(logger, db)

	// 设置路由
	r := gin.New()
	r.Use(gin.Recovery())
	h.RegisterRoutes(r)

	return h, r
}

func TestHandler_CreateNote(t *testing.T) {
	_, r := setupTestHandler(t)

	tests := []struct {
		name         string
		requestBody  map[string]interface{}
		wantStatus   int
		wantResponse map[string]interface{}
	}{
		{
			name: "正常创建笔记",
			requestBody: map[string]interface{}{
				"title":     "测试笔记",
				"content":   "测试内容",
				"file_path": "/test/note.md",
			},
			wantStatus: http.StatusCreated,
			wantResponse: map[string]interface{}{
				"status": "success",
			},
		},
		{
			name: "缺少必要字段",
			requestBody: map[string]interface{}{
				"content": "测试内容",
			},
			wantStatus: http.StatusBadRequest,
			wantResponse: map[string]interface{}{
				"error":  "Invalid request body",
				"status": "error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/notes", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// 检查响应状态
			assert.Equal(t, tt.wantResponse["status"], response["status"])
			if tt.wantStatus != http.StatusCreated {
				assert.Equal(t, tt.wantResponse["error"], response["error"])
			}
		})
	}
}

func TestHandler_GetNote(t *testing.T) {
	h, r := setupTestHandler(t)

	// 创建测试数据
	note := &model.Note{
		Title:    "测试笔记",
		Content:  "测试内容",
		FilePath: "/test/note.md",
	}
	h.db.Create(note)

	tests := []struct {
		name         string
		noteID       string
		wantStatus   int
		wantResponse map[string]interface{}
	}{
		{
			name:       "获取存在的笔记",
			noteID:     note.ID,
			wantStatus: http.StatusOK,
			wantResponse: map[string]interface{}{
				"status": "success",
			},
		},
		{
			name:       "获取不存在的笔记",
			noteID:     "not-exist",
			wantStatus: http.StatusNotFound,
			wantResponse: map[string]interface{}{
				"error":  "Note not found",
				"status": "error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/notes/"+tt.noteID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantResponse["status"], response["status"])
			if tt.wantStatus != http.StatusOK {
				assert.Equal(t, tt.wantResponse["error"], response["error"])
			}
		})
	}
}

func TestHandler_UpdateNote(t *testing.T) {
	h, r := setupTestHandler(t)

	// 创建测试数据
	note := &model.Note{
		Title:    "测试笔记",
		Content:  "测试内容",
		FilePath: "/test/note.md",
	}
	h.db.Create(note)

	tests := []struct {
		name         string
		noteID       string
		requestBody  map[string]interface{}
		wantStatus   int
		wantResponse map[string]interface{}
	}{
		{
			name:   "更新存在的笔记",
			noteID: note.ID,
			requestBody: map[string]interface{}{
				"title":   "更新后的标题",
				"content": "更新后的内容",
			},
			wantStatus: http.StatusOK,
			wantResponse: map[string]interface{}{
				"status": "success",
			},
		},
		{
			name:   "更新不存在的笔记",
			noteID: "not-exist",
			requestBody: map[string]interface{}{
				"title": "更新后的标题",
			},
			wantStatus: http.StatusInternalServerError,
			wantResponse: map[string]interface{}{
				"error":  "Failed to update note",
				"status": "error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/api/v1/notes/"+tt.noteID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantResponse["status"], response["status"])
			if tt.wantStatus != http.StatusOK {
				assert.Equal(t, tt.wantResponse["error"], response["error"])
			}
		})
	}
}

func TestHandler_DeleteNote(t *testing.T) {
	h, r := setupTestHandler(t)

	// 创建测试数据
	note := &model.Note{
		Title:    "测试笔记",
		Content:  "测试内容",
		FilePath: "/test/note.md",
	}
	h.db.Create(note)

	tests := []struct {
		name         string
		noteID       string
		wantStatus   int
		wantResponse map[string]interface{}
	}{
		{
			name:       "删除存在的笔记",
			noteID:     note.ID,
			wantStatus: http.StatusOK,
			wantResponse: map[string]interface{}{
				"status": "success",
			},
		},
		{
			name:       "删除不存在的笔记",
			noteID:     "not-exist",
			wantStatus: http.StatusInternalServerError,
			wantResponse: map[string]interface{}{
				"error":  "Failed to delete note",
				"status": "error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/api/v1/notes/"+tt.noteID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantResponse["status"], response["status"])
			if tt.wantStatus != http.StatusOK {
				assert.Equal(t, tt.wantResponse["error"], response["error"])
			}
		})
	}
}
