package handlers

import (
	"net/http"
	"sre-chat-api/internal/database"
	"sre-chat-api/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MessageHandler handles message-related HTTP requests
type MessageHandler struct {
	sseHandler *SSEHandler
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(sseHandler *SSEHandler) *MessageHandler {
	return &MessageHandler{
		sseHandler: sseHandler,
	}
}

// CreateMessage handles POST /api/v1/messages
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var req models.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If no group_id is provided, use the default "SRE Bootcamp" group
	if req.GroupID == 0 {
		var defaultGroup models.Group
		if err := database.DB.Where("name = ?", "SRE Bootcamp").First(&defaultGroup).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Default group not found"})
			return
		}
		req.GroupID = defaultGroup.ID
	}

	// Verify group exists
	var group models.Group
	if err := database.DB.First(&group, req.GroupID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	message := models.Message{
		GroupID: req.GroupID,
		User:    req.User,
		Content: req.Content,
	}

	if err := database.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	// Load group relationship
	database.DB.Preload("Group").First(&message, message.ID)

	// Notify SSE clients about the new message
	if h.sseHandler != nil {
		h.sseHandler.NotifyNewMessage(message)
	}

	c.JSON(http.StatusCreated, message)
}

// GetMessages handles GET /api/v1/messages
func (h *MessageHandler) GetMessages(c *gin.Context) {
	var messages []models.Message

	query := database.DB.Preload("Group")

	// Optional filter by group_id
	if groupID := c.Query("group_id"); groupID != "" {
		id, err := strconv.ParseUint(groupID, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group_id"})
			return
		}
		query = query.Where("group_id = ?", uint(id))
	}

	if err := query.Order("created_at ASC").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// GetMessage handles GET /api/v1/messages/:id
func (h *MessageHandler) GetMessage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var message models.Message
	if err := database.DB.Preload("Group").First(&message, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, message)
}

// UpdateMessage handles PUT /api/v1/messages/:id
func (h *MessageHandler) UpdateMessage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var req models.UpdateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var message models.Message
	if err := database.DB.First(&message, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	message.Content = req.Content
	if err := database.DB.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
		return
	}

	database.DB.Preload("Group").First(&message, message.ID)
	c.JSON(http.StatusOK, message)
}

// DeleteMessage handles DELETE /api/v1/messages/:id
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var message models.Message
	if err := database.DB.First(&message, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := database.DB.Delete(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}
