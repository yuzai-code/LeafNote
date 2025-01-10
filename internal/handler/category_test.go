package handler

import (
	"bytes"
	"encoding/json"
	"leafnote/internal/model"
	"leafnote/internal/testutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandler_CreateCategory(t *testing.T) {
	db, err := testutil.SetupTestDB()
	assert.NoError(t, err)
	defer func() {
		err := db.Migrator().DropTable(&model.Category{})
		assert.NoError(t, err)
	}()

	tests := []struct {
		name         string
		requestBody  map[string]interface{}
		wantStatus   int
		wantResponse interface{}
	}{
		{
			name: "创建成功",
			requestBody: map[string]interface{}{
				"name": "测试目录",
			},
			wantStatus: http.StatusCreated,
			wantResponse: map[string]interface{}{
				"name": "测试目录",
			},
		},
		{
			name: "名称为空",
			requestBody: map[string]interface{}{
				"name": "",
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

			req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r := gin.Default()
			h := NewHandler(zap.NewExample(), db)
			r.POST("/api/v1/categories", h.CreateCategory)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusCreated {
				var response model.Category
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.BaseModel.ID)
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["name"], response.Name)
			} else {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["error"], response["error"])
			}
		})
	}
}

func TestHandler_GetCategory(t *testing.T) {
	h, r := setupTestHandler(t)

	// 创建测试数据
	category := &model.Category{
		Name: "测试目录",
	}
	h.db.Create(category)

	tests := []struct {
		name         string
		categoryID   string
		wantStatus   int
		wantResponse interface{}
	}{
		{
			name:       "获取存在的目录",
			categoryID: category.ID,
			wantStatus: http.StatusOK,
			wantResponse: map[string]interface{}{
				"name": "测试目录",
			},
		},
		{
			name:       "获取不存在的目录",
			categoryID: "not-exist",
			wantStatus: http.StatusNotFound,
			wantResponse: map[string]interface{}{
				"error": "目录不存在",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/categories/"+tt.categoryID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["name"], response["name"])
				assert.NotEmpty(t, response["id"])
			} else {
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["error"], response["error"])
			}
		})
	}
}

func TestHandler_UpdateCategory(t *testing.T) {
	h, r := setupTestHandler(t)

	// 创建测试数据
	category := &model.Category{
		Name: "原始目录",
	}
	h.db.Create(category)

	tests := []struct {
		name         string
		categoryID   string
		requestBody  map[string]interface{}
		wantStatus   int
		wantResponse interface{}
	}{
		{
			name:       "更新成功",
			categoryID: category.ID,
			requestBody: map[string]interface{}{
				"name": "新目录名",
			},
			wantStatus: http.StatusOK,
			wantResponse: map[string]interface{}{
				"name": "新目录名",
			},
		},
		{
			name:       "目录不存在",
			categoryID: "not-exist",
			requestBody: map[string]interface{}{
				"name": "新目录名",
			},
			wantStatus: http.StatusNotFound,
			wantResponse: map[string]interface{}{
				"error": "目录不存在",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req := httptest.NewRequest("PUT", "/api/v1/categories/"+tt.categoryID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["name"], response["name"])
				assert.NotEmpty(t, response["id"])
			} else {
				assert.Equal(t, tt.wantResponse.(map[string]interface{})["error"], response["error"])
			}
		})
	}
}

func TestHandler_DeleteCategory(t *testing.T) {
	h, r := setupTestHandler(t)

	// 创建测试数据
	category := &model.Category{
		Name: "测试目录",
	}
	h.db.Create(category)

	tests := []struct {
		name         string
		categoryID   string
		wantStatus   int
		wantResponse interface{}
	}{
		{
			name:       "删除成功",
			categoryID: category.ID,
			wantStatus: http.StatusOK,
			wantResponse: map[string]interface{}{
				"message": "删除成功",
			},
		},
		{
			name:       "目录不存在",
			categoryID: "not-exist",
			wantStatus: http.StatusNotFound,
			wantResponse: map[string]interface{}{
				"error": "目录不存在",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/api/v1/categories/"+tt.categoryID, nil)
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

func TestHandler_ListCategories(t *testing.T) {
	h, r := setupTestHandler(t)

	// 创建测试数据
	parent := &model.Category{
		Name: "父目录",
	}
	h.db.Create(parent)

	child := &model.Category{
		Name:     "子目录",
		ParentID: &parent.ID,
	}
	h.db.Create(child)

	t.Run("列表查询", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/categories", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// 只返回顶级目录
		assert.Len(t, response, 1)
		assert.Equal(t, parent.Name, response[0]["name"])

		// 检查子目录
		children := response[0]["children"].([]interface{})
		assert.Len(t, children, 1)
		childMap := children[0].(map[string]interface{})
		assert.Equal(t, child.Name, childMap["name"])
	})
}
