package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"leafnote/internal/model"
	"leafnote/internal/service"
)

// TagHandler 标签处理器
type TagHandler struct {
	tagService *service.TagService
}

// NewTagHandler 创建标签处理器实例
func NewTagHandler(tagService *service.TagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

// RegisterRoutes 注册路由
func (h *TagHandler) RegisterRoutes(r *gin.RouterGroup) {
	tags := r.Group("/tags")
	{
		tags.GET("", h.ListTags)         // 获取标签列表
		tags.POST("", h.CreateTag)       // 创建标签
		tags.GET("/:id", h.GetTag)       // 获取单个标签
		tags.PUT("/:id", h.UpdateTag)    // 更新标签
		tags.DELETE("/:id", h.DeleteTag) // 删除标签
	}
}

// CreateTag 创建标签
// @Summary 创建标签
// @Description 创建一个新的标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param tag body model.Tag true "标签信息"
// @Success 200 {object} model.Tag
// @Router /api/v1/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.tagService.CreateTag(c.Request.Context(), &tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// GetTag 获取单个标签
// @Summary 获取标签详情
// @Description 根据ID获取标签详情
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path string true "标签ID"
// @Success 200 {object} model.Tag
// @Router /api/v1/tags/{id} [get]
func (h *TagHandler) GetTag(c *gin.Context) {
	id := c.Param("id")
	tag, err := h.tagService.GetTagByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// ListTags 获取标签列表
// @Summary 获取标签列表
// @Description 获取所有标签列表
// @Tags 标签管理
// @Accept json
// @Produce json
// @Success 200 {array} model.Tag
// @Router /api/v1/tags [get]
func (h *TagHandler) ListTags(c *gin.Context) {
	tags, err := h.tagService.ListTags(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// UpdateTag 更新标签
// @Summary 更新标签
// @Description 更新标签信息
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path string true "标签ID"
// @Param tag body model.Tag true "标签信息"
// @Success 200 {object} model.Tag
// @Router /api/v1/tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id := c.Param("id")
	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	tag.ID = id
	if err := h.tagService.UpdateTag(c.Request.Context(), &tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// DeleteTag 删除标签
// @Summary 删除标签
// @Description 删除指定标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path string true "标签ID"
// @Success 200 {string} string "删除成功"
// @Router /api/v1/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id := c.Param("id")
	if err := h.tagService.DeleteTag(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
