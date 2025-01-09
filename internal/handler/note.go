package handler

import (
	"leafnote/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListNotes 获取笔记列表
func (h *Handler) ListNotes(c *gin.Context) {
	noteService := service.NewNoteService(h.db, h.logger)
	notes, err := noteService.ListNotes()
	if err != nil {
		h.logger.Error("Failed to get notes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to get notes",
			"status": "error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   notes,
		"status": "success",
	})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid request body",
			"status": "error",
		})
		return
	}

	noteService := service.NewNoteService(h.db, h.logger)
	note, err := noteService.CreateNote(service.CreateNoteInput(req))
	if err != nil {
		h.logger.Error("Failed to create note", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to create note",
			"status": "error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":   note,
		"status": "success",
	})
}

// GetNote 获取单个笔记
func (h *Handler) GetNote(c *gin.Context) {
	id := c.Param("id")
	noteService := service.NewNoteService(h.db, h.logger)
	note, err := noteService.GetNote(id)
	if err != nil {
		h.logger.Error("Failed to get note", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "Note not found",
			"status": "error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   note,
		"status": "success",
	})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid request body",
			"status": "error",
		})
		return
	}

	noteService := service.NewNoteService(h.db, h.logger)
	err := noteService.UpdateNote(id, service.UpdateNoteInput(req))
	if err != nil {
		h.logger.Error("Failed to update note", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to update note",
			"status": "error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// DeleteNote 删除笔记
func (h *Handler) DeleteNote(c *gin.Context) {
	id := c.Param("id")
	noteService := service.NewNoteService(h.db, h.logger)
	err := noteService.DeleteNote(id)
	if err != nil {
		h.logger.Error("Failed to delete note", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to delete note",
			"status": "error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
