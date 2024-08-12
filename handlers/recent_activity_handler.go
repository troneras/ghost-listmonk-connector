// handlers/recent_activity_handler.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/services"
)

type RecentActivityHandler struct {
	service *services.RecentActivityService
}

func NewRecentActivityHandler(service *services.RecentActivityService) *RecentActivityHandler {
	return &RecentActivityHandler{service: service}
}

func (h *RecentActivityHandler) GetRecentActivity(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUser := user.(*models.User)

	activities, err := h.service.GetRecentActivity(currentUser.ID, 10) // Get last 10 activities
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recent activities"})
		return
	}

	c.JSON(http.StatusOK, activities)
}
