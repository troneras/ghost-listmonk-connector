package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/troneras/ghost-listmonk-connector/services"
)

type Handlers struct {
	Auth            *AuthHandler
	Son             *SonHandler
	Webhook         *WebhookHandler
	Listmonk        *ListmonkHandler
	Home            *HomeHandler
	WebhookLog      *WebhookLogHandler
	SonExecutionLog *SonExecutionLogHandler
	RecentActivity *RecentActivityHandler
	SonStats		*SonStatsHandler
}

func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{
		Auth:            NewAuthHandler(services.User, services.MagicLink, services.Email),
		Son:             NewSonHandler(services.SonStorage),
		Webhook:         NewWebhookHandler(services.SonStorage, services.SonExecutor, services.Webhook, services.WebhookLogger),
		Listmonk:        NewListmonkHandler(services.ListmonkClient),
		Home:            NewHomeHandler(),
		WebhookLog:      NewWebhookLogHandler(services.WebhookLogger),
		SonExecutionLog: NewSonExecutionLogHandler(services.SonExecutionLogger),
		RecentActivity: NewRecentActivityHandler(services.RecentActivity),
		SonStats:		NewSonStatsHandler(services.SonExecutionLogger),
	}
}

// HomeHandler handles the home route
type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) HandleHome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Ghost-Listmonk Connector API"})
}

// Implement NewWebhookLogHandler and WebhookLogHandler methods here
