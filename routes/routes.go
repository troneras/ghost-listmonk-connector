package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/troneras/ghost-listmonk-connector/handlers"
	"github.com/troneras/ghost-listmonk-connector/middleware"
	"github.com/troneras/ghost-listmonk-connector/services"
)

func SetupRoutes(r *gin.Engine, handlers *handlers.Handlers, services *services.Services) {
	api := r.Group("/api")
	{
		// Public routes
		api.POST("/auth/magic-link", handlers.Auth.RequestMagicLink)
		api.GET("/auth/verify", handlers.Auth.VerifyMagicLink)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthRequired(services.User))
		{
			protected.GET("/", handlers.Home.HandleHome)

			sons := protected.Group("/sons")
			{
				sons.POST("", handlers.Son.Create)
				sons.GET("", handlers.Son.List)
				sons.GET("/:id", handlers.Son.Get)
				sons.PUT("/:id", handlers.Son.Update)
				sons.DELETE("/:id", handlers.Son.Delete)
			}
			protected.GET("/son-execution-logs", handlers.SonExecutionLog.GetSonExecutionLogs)
			protected.GET("/son-executions/:executionId/action-logs", handlers.SonExecutionLog.GetActionExecutionLogs)

			protected.GET("/webhook-info", handlers.Webhook.GetWebhookInfo)
			protected.GET("/lists", handlers.Listmonk.GetLists)
			protected.GET("/templates", handlers.Listmonk.GetTemplates)

			// Webhook log routes
			protected.GET("/webhook-logs", handlers.WebhookLog.GetLogs)
			protected.GET("/webhook-logs/:id", handlers.WebhookLog.GetLogDetails)
			protected.POST("/webhook-logs/:id/replay", handlers.Webhook.ReplayWebhook)

			protected.GET("/recent-activity", handlers.RecentActivity.GetRecentActivity)
			protected.GET("/son-stats", handlers.SonStats.GetSonStats)
		}
	}

	// Webhook route (public, but requires signature verification)
	r.POST("/webhook/:endpoint", handlers.Webhook.HandleWebhook)
}
