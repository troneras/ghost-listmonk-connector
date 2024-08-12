package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/troneras/ghost-listmonk-connector/services"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type AuthHandler struct {
	userService      *services.UserService
	magicLinkService *services.MagicLinkService
	emailService     *services.EmailService
}

func NewAuthHandler(userService *services.UserService, magicLinkService *services.MagicLinkService, emailService *services.EmailService) *AuthHandler {
	return &AuthHandler{
		userService:      userService,
		magicLinkService: magicLinkService,
		emailService:     emailService,
	}
}

func (h *AuthHandler) RequestMagicLink(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		// If user doesn't exist, create a new one
		utils.InfoLogger.Printf("User with email %s not found, creating new user", req.Email)
		user, err = h.userService.CreateUser(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	token, err := h.magicLinkService.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create magic link"})
		return
	}

	magicLink := utils.GetConfig().FrontendURL + "/auth/verify?token=" + token

	// Send email with magic link
	err = h.emailService.SendMagicLinkEmail(user.Email, magicLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send magic link email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Magic link sent to your email"})
}

func (h *AuthHandler) VerifyMagicLink(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	userID, err := h.magicLinkService.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Generate JWT token
	jwtToken, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}
