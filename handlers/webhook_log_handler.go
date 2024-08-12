package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/services"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type WebhookLogHandler struct {
	logger *services.WebhookLogger
}

func NewWebhookLogHandler(logger *services.WebhookLogger) *WebhookLogHandler {
	return &WebhookLogHandler{
		logger: logger,
	}
}

// GetLogs retrieves a paginated list of webhook logs for the current user
func (h *WebhookLogHandler) GetLogs(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.ErrorLogger.Println("User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUser := user.(*models.User)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	logs, total, err := h.logger.GetWebhookLogs(currentUser.ID, limit, offset)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to fetch webhook logs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch webhook logs"})
		return
	}

	nextOffset := offset + limit
	if nextOffset >= total {
		nextOffset = -1 // Indicate that there are no more pages
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
		"pagination": gin.H{
			"total":       total,
			"limit":       limit,
			"offset":      offset,
			"next_offset": nextOffset,
		},
	})
}

// GetLogDetails retrieves the details of a specific webhook log
func (h *WebhookLogHandler) GetLogDetails(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.ErrorLogger.Println("User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUser := user.(*models.User)

	logID := c.Param("id")

	log, err := h.logger.GetWebhookLogDetails(logID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to fetch webhook log details: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Webhook log not found"})
		return
	}

	// Ensure the log belongs to the current user
	if log.UserID != currentUser.ID {
		utils.ErrorLogger.Printf("User %s attempted to access log %s belonging to another user", currentUser.ID, logID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, log)
}
