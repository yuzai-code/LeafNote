package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 结构体包含所有处理器依赖
type Handler struct {
	logger *zap.Logger
}

// NewHandler 创建一个新的处理器实例
func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// Health 健康检查处理器
func (h *Handler) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
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
			tags.DELETE("/:id", h.DeleteTag)
		}

		// 目录相关路由
		categories := v1.Group("/categories")
		{
			categories.GET("", h.ListCategories)
			categories.POST("", h.CreateCategory)
			categories.DELETE("/:id", h.DeleteCategory)
		}
	}
}

// 笔记相关处理器
func (h *Handler) ListNotes(c *gin.Context)  { c.JSON(200, gin.H{"message": "not implemented"}) }
func (h *Handler) CreateNote(c *gin.Context) { c.JSON(200, gin.H{"message": "not implemented"}) }
func (h *Handler) GetNote(c *gin.Context)    { c.JSON(200, gin.H{"message": "not implemented"}) }
func (h *Handler) UpdateNote(c *gin.Context) { c.JSON(200, gin.H{"message": "not implemented"}) }
func (h *Handler) DeleteNote(c *gin.Context) { c.JSON(200, gin.H{"message": "not implemented"}) }

// 标签相关处理器
func (h *Handler) ListTags(c *gin.Context)  { c.JSON(200, gin.H{"message": "not implemented"}) }
func (h *Handler) CreateTag(c *gin.Context) { c.JSON(200, gin.H{"message": "not implemented"}) }
func (h *Handler) DeleteTag(c *gin.Context) { c.JSON(200, gin.H{"message": "not implemented"}) }

// 目录相关处理器
func (h *Handler) ListCategories(c *gin.Context) { c.JSON(200, gin.H{"message": "not implemented"}) }
func (h *Handler) CreateCategory(c *gin.Context) { c.JSON(200, gin.H{"message": "not implemented"}) }
func (h *Handler) DeleteCategory(c *gin.Context) { c.JSON(200, gin.H{"message": "not implemented"}) }
