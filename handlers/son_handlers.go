package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/services"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type SonHandler struct {
	storage *services.SonStorage
}

type ListmonkHandler struct {
	client *services.ListmonkClient
}

func NewSonHandler(storage *services.SonStorage) *SonHandler {
	return &SonHandler{storage: storage}
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

func (h *SonHandler) Create(c *gin.Context) {
	var son models.Son
	if err := c.ShouldBindJSON(&son); err != nil {
		utils.ErrorLogger.Errorf("Invalid Son data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.storage.Create(&son); err != nil {
		utils.ErrorLogger.Errorf("Failed to create Son: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.InfoLogger.Infof("Created Son: %s", utils.PrettyPrint(son))
	c.JSON(http.StatusCreated, son)
}

func (h *SonHandler) Get(c *gin.Context) {
	id := c.Param("id")
	son, err := h.storage.Get(id)
	if err != nil {
		if err == services.ErrSonNotFound {
			utils.ErrorLogger.Errorf("Son not found: %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Son not found"})
		} else {
			utils.ErrorLogger.Errorf("Failed to get Son: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	utils.InfoLogger.Infof("Retrieved Son: %s", utils.PrettyPrint(son))
	c.JSON(http.StatusOK, son)
}

func (h *SonHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var son models.Son
	if err := c.ShouldBindJSON(&son); err != nil {
		utils.ErrorLogger.Errorf("Invalid Son data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	son.ID = id
	if err := h.storage.Update(son); err != nil {
		if err == services.ErrSonNotFound {
			utils.ErrorLogger.Errorf("Son not found for update: %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Son not found"})
		} else {
			utils.ErrorLogger.Errorf("Failed to update Son: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	utils.InfoLogger.Infof("Updated Son: %s", utils.PrettyPrint(son))
	c.JSON(http.StatusOK, son)
}

func (h *SonHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.storage.Delete(id); err != nil {
		if err == services.ErrSonNotFound {
			utils.ErrorLogger.Errorf("Son not found for deletion: %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Son not found"})
		} else {
			utils.ErrorLogger.Errorf("Failed to delete Son: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	utils.InfoLogger.Infof("Deleted Son: %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "Son deleted successfully"})
}

func (h *SonHandler) List(c *gin.Context) {
	sons, err := h.storage.List()

	if err != nil {
		utils.ErrorLogger.Errorf("Failed to list Sons: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.InfoLogger.Infof("Retrieved %d Sons", len(sons))
	c.JSON(http.StatusOK, sons)
}
