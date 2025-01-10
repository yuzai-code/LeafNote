package handler

import (
	"leafnote/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListTags 获取标签列表
func (h *Handler) ListTags(c *gin.Context) {
	tags, err := h.tagService.ListTags(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get tags", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取标签列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, tags)
}

// CreateTag 创建标签
func (h *Handler) CreateTag(c *gin.Context) {
	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		h.logger.Error("Invalid request parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
		})
		return
	}

	if err := h.tagService.CreateTag(c.Request.Context(), &tag); err != nil {
		h.logger.Error("Failed to create tag", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建标签失败",
		})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// GetTag 获取标签详情
func (h *Handler) GetTag(c *gin.Context) {
	id := c.Param("id")
	tag, err := h.tagService.GetTagByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get tag", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{
			"error": "标签不存在",
		})
		return
	}
	c.JSON(http.StatusOK, tag)
}

// UpdateTag 更新标签
func (h *Handler) UpdateTag(c *gin.Context) {
	id := c.Param("id")
	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		h.logger.Error("Invalid request parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
		})
		return
	}

	// 先检查标签是否存在
	_, err := h.tagService.GetTagByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get tag", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{
			"error": "标签不存在",
		})
		return
	}

	tag.BaseModel.ID = id
	if err := h.tagService.UpdateTag(c.Request.Context(), &tag); err != nil {
		h.logger.Error("Failed to update tag", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新标签失败",
		})
		return
	}

	// 获取更新后的标签信息
	updatedTag, err := h.tagService.GetTagByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get updated tag", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取更新后的标签失败",
		})
		return
	}

	c.JSON(http.StatusOK, updatedTag)
}

// DeleteTag 删除标签
func (h *Handler) DeleteTag(c *gin.Context) {
	id := c.Param("id")
	if err := h.tagService.DeleteTag(c.Request.Context(), id); err != nil {
		if err.Error() == "标签不存在" {
			h.logger.Error("Tag not found", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.logger.Error("Failed to delete tag", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除标签失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}
