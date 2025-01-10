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
		wantResponse interface{}
	}{
		{
			name: "创建成功",
			requestBody: map[string]interface{}{
				"title":     "测试笔记",
				"content":   "测试内容",
				"file_path": "/test/note.md",
			},
			wantStatus: http.StatusCreated,
			wantResponse: map[string]interface{}{
				"title":     "测试笔记",
				"content":   "测试内容",
				"file_path": "/test/note.md",
			},
		},
		{
			name: "缺少必填字段",
			requestBody: map[string]interface{}{
				"content": "测试内容",
			},
			wantStatus: http.StatusBadRequest,
			wantResponse: map[string]interface{}{
				"error": "无效的请求参数",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/notes", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.wantStatus == http.StatusCreated {
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["title"], response["title"])
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["content"], response["content"])
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["file_path"], response["file_path"])
				assert.NotEmpty(t, response["id"])
			} else {
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["error"], response["error"])
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
		wantResponse interface{}
	}{
		{
			name:       "获取存在的笔记",
			noteID:     note.ID,
			wantStatus: http.StatusOK,
			wantResponse: map[string]interface{}{
				"title":     "测试笔记",
				"content":   "测试内容",
				"file_path": "/test/note.md",
			},
		},
		{
			name:       "获取不存在的笔记",
			noteID:     "not-exist",
			wantStatus: http.StatusNotFound,
			wantResponse: map[string]interface{}{
				"error": "笔记不存在",
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

			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["title"], response["title"])
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["content"], response["content"])
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["file_path"], response["file_path"])
				assert.NotEmpty(t, response["id"])
			} else {
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["error"], response["error"])
			}
		})
	}
}

func TestHandler_UpdateNote(t *testing.T) {
	h, r := setupTestHandler(t)

	// 创建测试数据
	note := &model.Note{
		Title:    "原始标题",
		Content:  "原始内容",
		FilePath: "/test/note.md",
	}
	h.db.Create(note)

	tests := []struct {
		name         string
		noteID       string
		requestBody  map[string]interface{}
		wantStatus   int
		wantResponse interface{}
	}{
		{
			name:   "更新成功",
			noteID: note.ID,
			requestBody: map[string]interface{}{
				"title":   "新标题",
				"content": "新内容",
			},
			wantStatus: http.StatusOK,
			wantResponse: map[string]interface{}{
				"title":     "新标题",
				"content":   "新内容",
				"file_path": "/test/note.md",
			},
		},
		{
			name:   "笔记不存在",
			noteID: "not-exist",
			requestBody: map[string]interface{}{
				"title": "新标题",
			},
			wantStatus: http.StatusInternalServerError,
			wantResponse: map[string]interface{}{
				"error": "更新笔记失败",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req := httptest.NewRequest("PUT", "/api/v1/notes/"+tt.noteID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["title"], response["title"])
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["content"], response["content"])
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["file_path"], response["file_path"])
				assert.NotEmpty(t, response["id"])
			} else {
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["error"], response["error"])
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
		wantResponse interface{}
	}{
		{
			name:       "删除成功",
			noteID:     note.ID,
			wantStatus: http.StatusOK,
			wantResponse: map[string]interface{}{
				"message": "删除成功",
			},
		},
		{
			name:       "笔记不存在",
			noteID:     "not-exist",
			wantStatus: http.StatusInternalServerError,
			wantResponse: map[string]interface{}{
				"error": "删除笔记失败",
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

			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["message"], response["message"])
			} else {
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["error"], response["error"])
			}
		})
	}
}
