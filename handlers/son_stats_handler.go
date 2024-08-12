// handlers/son_stats_handler.go

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/services"
)

type SonStatsHandler struct {
	logger *services.SonExecutionLogger
}

func NewSonStatsHandler(logger *services.SonExecutionLogger) *SonStatsHandler {
	return &SonStatsHandler{logger: logger}
}

func (h *SonStatsHandler) GetSonStats(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUser := user.(*models.User)

	timeframe := c.DefaultQuery("timeframe", "24h")

	stats, err := h.logger.GetSonStats(c, currentUser.ID, timeframe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch son statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
