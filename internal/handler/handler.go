package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"leafnote/internal/service"
)

// Handler 结构体包含所有处理器依赖
type Handler struct {
	logger          *zap.Logger
	db              *gorm.DB
	tagService      *service.TagService
	categoryService *service.CategoryService
}

// NewHandler 创建一个新的处理器实例
func NewHandler(logger *zap.Logger, db *gorm.DB) *Handler {
	return &Handler{
		logger:          logger,
		db:              db,
		tagService:      service.NewTagService(db),
		categoryService: service.NewCategoryService(db),
	}
}

// Health 健康检查处理器
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "服务正常",
	})
}

// RegisterRoutes 注册所有路由
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// API 版本组
	v1 := r.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/health", h.Health)

		// 笔记相关路由
		notes := v1.Group("/notes")
		{
			notes.GET("", h.ListNotes)
			notes.POST("", h.CreateNote)
			notes.GET("/:id", h.GetNote)
			notes.PUT("/:id", h.UpdateNote)
			notes.DELETE("/:id", h.DeleteNote)
		}

		// 标签相关路由
		tags := v1.Group("/tags")
		{
			tags.GET("", h.ListTags)
			tags.POST("", h.CreateTag)
			tags.GET("/:id", h.GetTag)
			tags.PUT("/:id", h.UpdateTag)
			tags.DELETE("/:id", h.DeleteTag)
		}

		// 目录相关路由
		categories := v1.Group("/categories")
		{
			categories.GET("", h.ListCategories)
			categories.POST("", h.CreateCategory)
			categories.GET("/:id", h.GetCategory)
			categories.PUT("/:id", h.UpdateCategory)
			categories.DELETE("/:id", h.DeleteCategory)
		}
	}
}
