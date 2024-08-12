package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/troneras/ghost-listmonk-connector/services"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type ListmonkHandler struct {
	client *services.ListmonkClient
}

func NewListmonkHandler(client *services.ListmonkClient) *ListmonkHandler {
	return &ListmonkHandler{client: client}
}

func (h *ListmonkHandler) GetLists(c *gin.Context) {
	lists, err := h.client.GetLists()
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to get lists: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": lists})
}

func (h *ListmonkHandler) GetTemplates(c *gin.Context) {
	templates, err := h.client.GetTemplates()
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to get templates: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": templates})
}
