package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"leafnote/internal/model"
	"leafnote/internal/testutil"
)

func setupCategoryTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 自动迁移
	err = db.AutoMigrate(&model.Category{}, &model.Note{})
	assert.NoError(t, err)

	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	handler := NewHandler(logger, db)
	router := gin.Default()
	handler.RegisterRoutes(router)

	return router, db
}

func TestHandler_CreateCategory(t *testing.T) {
	router, db := setupCategoryTestRouter(t)

	// 创建父目录
	parent := &model.Category{Name: "父目录"}
	err := db.Create(parent).Error
	assert.NoError(t, err)

	tests := []struct {
		name       string
		category   model.Category
		wantStatus int
	}{
		{
			name: "创建成功",
			category: model.Category{
				Name: "测试目录",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "创建子目录成功",
			category: model.Category{
				Name:     "子目录",
				ParentID: &parent.BaseModel.ID,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "名称为空",
			category: model.Category{
				Name: "",
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "父目录不存在",
			category: model.Category{
				Name:     "子目录",
				ParentID: testutil.StringPtr("not-exist"),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.category)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var response model.Category
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.BaseModel.ID)
				assert.Equal(t, tt.category.Name, response.Name)
				if tt.category.ParentID != nil {
					assert.Equal(t, *tt.category.ParentID, *response.ParentID)
				}
			}
		})
	}
}

func TestHandler_GetCategory(t *testing.T) {
	router, db := setupCategoryTestRouter(t)

	// 创建测试数据
	category := &model.Category{Name: "测试目录"}
	err := db.Create(category).Error
	assert.NoError(t, err)

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{
			name:       "获取成功",
			id:         category.BaseModel.ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "目录不存在",
			id:         "not-exist",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/categories/"+tt.id, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var response model.Category
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.id, response.BaseModel.ID)
				assert.Equal(t, category.Name, response.Name)
			}
		})
	}
}

func TestHandler_ListCategories(t *testing.T) {
	router, db := setupCategoryTestRouter(t)

	// 创建测试数据
	parent := &model.Category{Name: "父目录"}
	err := db.Create(parent).Error
	assert.NoError(t, err)

	child := &model.Category{
		Name:     "子目录",
		ParentID: &parent.BaseModel.ID,
	}
	err = db.Create(child).Error
	assert.NoError(t, err)

	t.Run("列表查询", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []model.Category
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1) // 只返回顶级目录
		assert.Equal(t, parent.Name, response[0].Name)
		assert.Len(t, response[0].Children, 1) // 包含子目录
	})
}

func TestHandler_UpdateCategory(t *testing.T) {
	router, db := setupCategoryTestRouter(t)

	// 创建测试数据
	parent := &model.Category{Name: "父目录"}
	err := db.Create(parent).Error
	assert.NoError(t, err)

	category := &model.Category{Name: "原始目录"}
	err = db.Create(category).Error
	assert.NoError(t, err)

	tests := []struct {
		name       string
		id         string
		category   model.Category
		wantStatus int
	}{
		{
			name: "更新成功",
			id:   category.BaseModel.ID,
			category: model.Category{
				Name: "新目录名",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "更新为子目录",
			id:   category.BaseModel.ID,
			category: model.Category{
				Name:     "子目录",
				ParentID: &parent.BaseModel.ID,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "目录不存在",
			id:   "not-exist",
			category: model.Category{
				Name: "新目录名",
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "父目录不存在",
			id:   category.BaseModel.ID,
			category: model.Category{
				Name:     "子目录",
				ParentID: testutil.StringPtr("not-exist"),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.category)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/api/v1/categories/"+tt.id, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var response model.Category
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.id, response.BaseModel.ID)
				assert.Equal(t, tt.category.Name, response.Name)
				if tt.category.ParentID != nil {
					assert.Equal(t, *tt.category.ParentID, *response.ParentID)
				}
			}
		})
	}
}

func TestHandler_DeleteCategory(t *testing.T) {
	router, db := setupCategoryTestRouter(t)

	// 创建测试数据
	parent := &model.Category{Name: "父目录"}
	err := db.Create(parent).Error
	assert.NoError(t, err)

	child := &model.Category{
		Name:     "子目录",
		ParentID: &parent.BaseModel.ID,
	}
	err = db.Create(child).Error
	assert.NoError(t, err)

	// 创建一个带笔记的目录
	categoryWithNote := &model.Category{Name: "带笔记的目录"}
	err = db.Create(categoryWithNote).Error
	assert.NoError(t, err)

	note := &model.Note{
		Title:      "测试笔记",
		CategoryID: &categoryWithNote.BaseModel.ID,
	}
	err = db.Create(note).Error
	assert.NoError(t, err)

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{
			name:       "删除失败-存在子目录",
			id:         parent.BaseModel.ID,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "删除失败-存在笔记",
			id:         categoryWithNote.BaseModel.ID,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "删除成功-子目录",
			id:         child.BaseModel.ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "目录不存在",
			id:         "not-exist",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "删除成功-父目录",
			id:         parent.BaseModel.ID,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/categories/"+tt.id, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var response map[string]string
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "删除成功", response["message"])
			}
		})
	}
}
