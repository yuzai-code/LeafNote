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
	"leafnote/internal/service"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 自动迁移
	err = db.AutoMigrate(&model.Tag{})
	assert.NoError(t, err)

	// 创建路由
	r := gin.New()
	r.Use(gin.Recovery())

	// 创建处理器
	logger, _ := zap.NewDevelopment()
	h := NewHandler(logger, db)
	h.RegisterRoutes(r)

	return r, db
}

func TestTagHandler_CreateTag(t *testing.T) {
	r, _ := setupTestRouter(t)

	tests := []struct {
		name       string
		tag        model.Tag
		wantStatus int
	}{
		{
			name: "创建成功",
			tag: model.Tag{
				Name: "测试标签",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "名称为空",
			tag: model.Tag{
				Name: "",
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.tag)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/tags", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var response model.Tag
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.BaseModel.ID)
				assert.Equal(t, tt.tag.Name, response.Name)
			}
		})
	}
}

func TestTagHandler_GetTag(t *testing.T) {
	r, db := setupTestRouter(t)

	// 创建测试数据
	tagService := service.NewTagService(db)
	tag := &model.Tag{Name: "测试标签"}
	err := tagService.CreateTag(nil, tag)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		tagID      string
		wantStatus int
	}{
		{
			name:       "获取成功",
			tagID:      tag.BaseModel.ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "标签不存在",
			tagID:      "not-exist",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/tags/"+tt.tagID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var response model.Tag
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.tagID, response.BaseModel.ID)
			}
		})
	}
}

func TestTagHandler_ListTags(t *testing.T) {
	r, db := setupTestRouter(t)

	// 创建测试数据
	tagService := service.NewTagService(db)
	parent := &model.Tag{Name: "父标签"}
	err := tagService.CreateTag(nil, parent)
	assert.NoError(t, err)

	child := &model.Tag{
		Name:     "子标签",
		ParentID: &parent.BaseModel.ID,
	}
	err = tagService.CreateTag(nil, child)
	assert.NoError(t, err)

	t.Run("列表查询", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tags", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []model.Tag
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1) // 只返回顶级标签
		assert.Equal(t, parent.BaseModel.ID, response[0].BaseModel.ID)
		assert.Len(t, response[0].Children, 1) // 包含子标签
	})
}

func TestTagHandler_UpdateTag(t *testing.T) {
	r, db := setupTestRouter(t)

	// 创建测试数据
	tagService := service.NewTagService(db)
	tag := &model.Tag{Name: "原始标签"}
	err := tagService.CreateTag(nil, tag)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		tagID      string
		updateTag  model.Tag
		wantStatus int
	}{
		{
			name:  "更新成功",
			tagID: tag.BaseModel.ID,
			updateTag: model.Tag{
				Name: "新标签名",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "标签不存在",
			tagID: "not-exist",
			updateTag: model.Tag{
				Name: "新标签名",
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:  "名称为空",
			tagID: tag.BaseModel.ID,
			updateTag: model.Tag{
				Name: "",
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.updateTag)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/api/v1/tags/"+tt.tagID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var response model.Tag
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.tagID, response.BaseModel.ID)
				assert.Equal(t, tt.updateTag.Name, response.Name)
			}
		})
	}
}

func TestTagHandler_DeleteTag(t *testing.T) {
	r, db := setupTestRouter(t)

	// 创建测试数据
	tagService := service.NewTagService(db)
	parent := &model.Tag{Name: "父标签"}
	err := tagService.CreateTag(nil, parent)
	assert.NoError(t, err)

	child := &model.Tag{
		Name:     "子标签",
		ParentID: &parent.BaseModel.ID,
	}
	err = tagService.CreateTag(nil, child)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		tagID      string
		wantStatus int
	}{
		{
			name:       "删除失败-存在子标签",
			tagID:      parent.BaseModel.ID,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "删除成功-子标签",
			tagID:      child.BaseModel.ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "删除成功-父标签",
			tagID:      parent.BaseModel.ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "标签不存在",
			tagID:      "not-exist",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/tags/"+tt.tagID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

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
