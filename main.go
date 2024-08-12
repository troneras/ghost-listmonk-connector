package main

import (
	"flag"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/troneras/ghost-listmonk-connector/database"
	"github.com/troneras/ghost-listmonk-connector/handlers"
	"github.com/troneras/ghost-listmonk-connector/routes"
	"github.com/troneras/ghost-listmonk-connector/services"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

func main() {
	// Define flags
	migrateFlag := flag.String("migrate", "", "Run database migrations. Use 'up' or 'down'")
	migrateSteps := flag.Int("steps", 0, "Number of migration steps to run")
	
	flag.Parse()

	// Check if migrations should be run
	if *migrateFlag != "" {
		if err := runMigrations(*migrateFlag, *migrateSteps); err != nil {
			utils.ErrorLogger.Fatalf("Failed to run migrations: %v", err)
		}
		return
	}

	config := utils.GetConfig()

	// Initialize database connection
	if err := database.InitDB(); err != nil {
		utils.ErrorLogger.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Initialize services
	services, err := services.NewServices(config)
	if err != nil {
		utils.ErrorLogger.Fatalf("Failed to initialize services: %v", err)
	}

	// Start SonExecutor
	go func() {
		if err := services.SonExecutor.Start(); err != nil {
			utils.ErrorLogger.Fatalf("Failed to start SonExecutor: %v", err)
		}
	}()
	defer services.SonExecutor.Stop()

	// Initialize handlers
	handlers := handlers.NewHandlers(services)

	// Set up Gin
	r := setupGin(config)

	// Set up routes
	routes.SetupRoutes(r, handlers, services)

	// Set up static file serving and catch-all route
	setupStaticAndCatchAll(r)

	// Start the server
	utils.InfoLogger.Infof("Server starting on port %s", config.Port)
	if err := r.Run(":" + config.Port); err != nil {
		utils.ErrorLogger.Fatalf("Failed to start server: %v", err)
	}
}

func setupGin(config *utils.Config) *gin.Engine {
	r := gin.New()
	r.RedirectTrailingSlash = true

	// Set up the cookie store
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions("mysession", store))

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.FrontendURL}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))

	return r
}

func setupStaticAndCatchAll(r *gin.Engine) {
	if gin.Mode() == gin.ReleaseMode {
		// Production mode
		r.Use(static.Serve("/", static.LocalFile("./ui/out", false)))

		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path

			// Check for the .html file
			htmlFilePath := filepath.Join("./ui/out", path+".html")
			if info, err := os.Stat(htmlFilePath); err == nil && !info.IsDir() {
				c.File(htmlFilePath)
				return
			}

			// Check for the file without .html extension
			filePath := filepath.Join("./ui/out", path)
			if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
				c.File(filePath)
				return
			}

			// If no specific file is found, serve index.html
			c.File("./ui/out/index.html")
		})
	} else {
		// Development mode
		nextJSURL, err := url.Parse("http://localhost:3000")
		if err != nil {
			utils.ErrorLogger.Fatalf("Failed to parse Next.js URL: %v", err)
		}
		proxy := httputil.NewSingleHostReverseProxy(nextJSURL)

		r.NoRoute(func(c *gin.Context) {
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	}
}

func runMigrations(direction string, steps int) error {
	if err := database.InitDB(); err != nil {
		return err
	}
	defer database.CloseDB()

	return database.RunMigrations(direction, steps)
}
