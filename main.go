package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/troneras/ghost-listmonk-connector/handlers"
	"github.com/troneras/ghost-listmonk-connector/middleware"
	"github.com/troneras/ghost-listmonk-connector/services"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		utils.ErrorLogger.Fatalf("Failed to load configuration: %v", err)
	}
	r := gin.Default()

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"} // Add your frontend URL here
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))

	// Initialize services
	sonStorage, err := services.NewSonStorage("ghost_listmonk.db")
	if err != nil {
		utils.ErrorLogger.Fatalf("Failed to initialize son storage: %v", err)
	}

	listmonkClient := services.NewListmonkClient(config)
	sonExecutor := services.NewSonExecutor(listmonkClient)
	sonHandler := handlers.NewSonHandler(sonStorage)
	webhookHandler := handlers.NewWebhookHandler(sonStorage, sonExecutor)
	listmonkHandler := handlers.NewListmonkHandler(listmonkClient)

	// Webhook endpoint (no auth required)
	r.POST("/webhook", webhookHandler.HandleWebhook)

	// Web UI routes (auth required)
	authorized := r.Group("/")
	authorized.Use(middleware.APIKeyAuth(config.APIKey))
	{
		authorized.GET("/", handleHome)

		// Son routes
		sons := authorized.Group("/sons")
		{
			sons.POST("", sonHandler.Create)
			sons.GET("", sonHandler.List)
			sons.GET("/:id", sonHandler.Get)
			sons.PUT("/:id", sonHandler.Update)
			sons.DELETE("/:id", sonHandler.Delete)
		}
		authorized.GET("/lists", listmonkHandler.GetLists)
		authorized.GET("/templates", listmonkHandler.GetTemplates)
	}

	utils.InfoLogger.Infof("Server starting on port %s", config.Port)
	if err := r.Run(":" + config.Port); err != nil {
		utils.ErrorLogger.Fatalf("Failed to start server: %v", err)
	}
}

func handleHome(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Welcome to the Ghost-Listmonk Connector"})
}
