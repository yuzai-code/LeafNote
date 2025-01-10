package handler

import (
	"leafnote/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListCategories 获取目录列表
func (h *Handler) ListCategories(c *gin.Context) {
	categories, err := h.categoryService.ListCategories(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get categories", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取目录列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// CreateCategory 创建目录
func (h *Handler) CreateCategory(c *gin.Context) {
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		h.logger.Error("Invalid request parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
		})
		return
	}

	if err := h.categoryService.CreateCategory(c.Request.Context(), &category); err != nil {
		h.logger.Error("Failed to create category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建目录失败",
		})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategory 获取目录详情
func (h *Handler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	category, err := h.categoryService.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get category", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{
			"error": "目录不存在",
		})
		return
	}
	c.JSON(http.StatusOK, category)
}

// UpdateCategory 更新目录
func (h *Handler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		h.logger.Error("Invalid request parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
		})
		return
	}

	// 先检查目录是否存在
	_, err := h.categoryService.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get category", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{
			"error": "目录不存在",
		})
		return
	}

	category.BaseModel.ID = id
	if err := h.categoryService.UpdateCategory(c.Request.Context(), &category); err != nil {
		h.logger.Error("Failed to update category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新目录失败",
		})
		return
	}

	// 获取更新后的目录信息
	updatedCategory, err := h.categoryService.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get updated category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取更新后的目录失败",
		})
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
}

// DeleteCategory 删除目录
func (h *Handler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := h.categoryService.DeleteCategory(c.Request.Context(), id); err != nil {
		if err.Error() == "目录不存在" {
			h.logger.Error("Category not found", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.logger.Error("Failed to delete category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除目录失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}
