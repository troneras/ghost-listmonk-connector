// handlers/son_execution_log_handler.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/services"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type SonExecutionLogHandler struct {
	logger *services.SonExecutionLogger
}

func NewSonExecutionLogHandler(logger *services.SonExecutionLogger) *SonExecutionLogHandler {
	return &SonExecutionLogHandler{
		logger: logger,
	}
}

func (h *SonExecutionLogHandler) GetSonExecutionLogs(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.ErrorLogger.Println("User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUser := user.(*models.User)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	logs, total, err := h.logger.GetSonExecutionLogs(currentUser.ID, limit, offset)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to fetch son execution logs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch son execution logs"})
		return
	}

	if logs == nil {
		logs = []models.SonExecutionLog{}
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
		"pagination": gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

func (h *SonExecutionLogHandler) GetActionExecutionLogs(c *gin.Context) {
	executionID := c.Param("executionId")

	logs, err := h.logger.GetActionExecutionLogs(executionID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to fetch action execution logs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch action execution logs"})
		return
	}

	if logs == nil {
		logs = []models.ActionExecutionLog{}
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
