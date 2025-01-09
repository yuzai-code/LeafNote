package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"leafnote/internal/model"
	"leafnote/internal/service"
)

// Handler 结构体包含所有处理器依赖
type Handler struct {
	logger     *zap.Logger
	db         *gorm.DB
	tagService *service.TagService
}

// NewHandler 创建一个新的处理器实例
func NewHandler(logger *zap.Logger, db *gorm.DB) *Handler {
	return &Handler{
		logger:     logger,
		db:         db,
		tagService: service.NewTagService(db),
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
			tags.GET("/:id", h.GetTag)
			tags.PUT("/:id", h.UpdateTag)
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

// ListTags 获取标签列表
func (h *Handler) ListTags(c *gin.Context) {
	tags, err := h.tagService.ListTags(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, tags)
}

// CreateTag 创建标签
func (h *Handler) CreateTag(c *gin.Context) {
	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.tagService.CreateTag(c.Request.Context(), &tag); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tag)
}

// GetTag 获取标签详情
func (h *Handler) GetTag(c *gin.Context) {
	id := c.Param("id")
	tag, err := h.tagService.GetTagByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, tag)
}

// UpdateTag 更新标签
func (h *Handler) UpdateTag(c *gin.Context) {
	id := c.Param("id")
	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	tag.ID = id
	if err := h.tagService.UpdateTag(c.Request.Context(), &tag); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tag)
}

// DeleteTag 删除标签
func (h *Handler) DeleteTag(c *gin.Context) {
	id := c.Param("id")
	if err := h.tagService.DeleteTag(c.Request.Context(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "删除成功"})
}

// 目录相关处理器
func (h *Handler) ListCategories(c *gin.Context) { c.JSON(200, gin.H{"message": "not implemented"}) }
func (h *Handler) CreateCategory(c *gin.Context) { c.JSON(200, gin.H{"message": "not implemented"}) }
func (h *Handler) DeleteCategory(c *gin.Context) { c.JSON(200, gin.H{"message": "not implemented"}) }
