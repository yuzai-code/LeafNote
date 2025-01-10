package handler

import (
	"leafnote/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListNotes 获取笔记列表
func (h *Handler) ListNotes(c *gin.Context) {
	categoryID := c.Query("category_id")
	var categoryIDPtr *string
	if categoryID != "" {
		categoryIDPtr = &categoryID
	}

	noteService := service.NewNoteService(h.db, h.logger)
	notes, err := noteService.ListNotes(categoryIDPtr)
	if err != nil {
		h.logger.Error("Failed to get notes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取笔记列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, notes)
}

// CreateNote 创建笔记
func (h *Handler) CreateNote(c *gin.Context) {
	var req struct {
		Title      string   `json:"title" binding:"required"`
		Content    string   `json:"content"`
		YAMLMeta   string   `json:"yaml_meta"`
		FilePath   string   `json:"file_path" binding:"required"`
		CategoryID *string  `json:"category_id"`
		TagIDs     []string `json:"tag_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
		})
		return
	}

	noteService := service.NewNoteService(h.db, h.logger)
	note, err := noteService.CreateNote(service.CreateNoteInput(req))
	if err != nil {
		h.logger.Error("Failed to create note", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建笔记失败",
		})
		return
	}

	c.JSON(http.StatusCreated, note)
}

// GetNote 获取单个笔记
func (h *Handler) GetNote(c *gin.Context) {
	id := c.Param("id")
	noteService := service.NewNoteService(h.db, h.logger)
	note, err := noteService.GetNote(id)
	if err != nil {
		h.logger.Error("Failed to get note", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{
			"error": "笔记不存在",
		})
		return
	}

	c.JSON(http.StatusOK, note)
}

// UpdateNote 更新笔记
func (h *Handler) UpdateNote(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		YAMLMeta   string   `json:"yaml_meta"`
		CategoryID *string  `json:"category_id"`
		TagIDs     []string `json:"tag_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
		})
		return
	}

	noteService := service.NewNoteService(h.db, h.logger)
	err := noteService.UpdateNote(id, service.UpdateNoteInput(req))
	if err != nil {
		if err.Error() == "笔记不存在" {
			h.logger.Error("Note not found", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.logger.Error("Failed to update note", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新笔记失败",
		})
		return
	}

	// 获取更新后的笔记信息
	note, err := noteService.GetNote(id)
	if err != nil {
		h.logger.Error("Failed to get updated note", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取更新后的笔记失败",
		})
		return
	}

	c.JSON(http.StatusOK, note)
}

// DeleteNote 删除笔记
func (h *Handler) DeleteNote(c *gin.Context) {
	id := c.Param("id")
	noteService := service.NewNoteService(h.db, h.logger)
	err := noteService.DeleteNote(id)
	if err != nil {
		if err.Error() == "笔记不存在" {
			h.logger.Error("Note not found", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.logger.Error("Failed to delete note", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除笔记失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}
