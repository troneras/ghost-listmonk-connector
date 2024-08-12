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

func NewSonHandler(storage *services.SonStorage) *SonHandler {
	return &SonHandler{storage: storage}
}

func (h *SonHandler) Create(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.ErrorLogger.Println("User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	currentUser := user.(*models.User)

	var son models.Son
	if err := c.ShouldBindJSON(&son); err != nil {
		utils.ErrorLogger.Errorf("Invalid Son data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	son.ID = utils.GenerateUUID()
	son.UserID = currentUser.ID

	// Check user's subscription level and apply limits
	var maxSons int
	switch currentUser.SubscriptionLevel {
	case models.SubscriptionFree:
		maxSons = 5
	case models.SubscriptionPremium:
		maxSons = 20
	case models.SubscriptionBusiness:
		maxSons = 100
	default:
		maxSons = 1
	}

	existingSons, err := h.storage.List(currentUser.ID)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to list Sons: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing Sons"})
		return
	}

	if len(existingSons) >= maxSons {
		c.JSON(http.StatusForbidden, gin.H{"error": "Son limit reached for your subscription level"})
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
	user, exists := c.Get("user")
	if !exists {
		utils.ErrorLogger.Println("User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	currentUser := user.(*models.User)

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

	if son.UserID != currentUser.ID {
		utils.ErrorLogger.Errorf("Unauthorized access to Son: %s", id)
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access to Son"})
		return
	}

	utils.InfoLogger.Infof("Retrieved Son: %s", utils.PrettyPrint(son))
	c.JSON(http.StatusOK, son)
}

func (h *SonHandler) Update(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.ErrorLogger.Println("User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	currentUser := user.(*models.User)

	id := c.Param("id")
	var son models.Son
	if err := c.ShouldBindJSON(&son); err != nil {
		utils.ErrorLogger.Errorf("Invalid Son data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	son.ID = id
	son.UserID = currentUser.ID
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
	user, exists := c.Get("user")
	if !exists {
		utils.ErrorLogger.Println("User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	currentUser := user.(*models.User)

	id := c.Param("id")
	if err := h.storage.Delete(id, currentUser.ID); err != nil {
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
	user, exists := c.Get("user")
	if !exists {
		utils.ErrorLogger.Println("User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	currentUser := user.(*models.User)

	sons, err := h.storage.List(currentUser.ID)

	if err != nil {
		utils.ErrorLogger.Errorf("Failed to list Sons: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.InfoLogger.Infof("Retrieved %d Sons", len(sons))
	c.JSON(http.StatusOK, sons)
}
