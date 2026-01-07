package internal

import (
	"sre-chat-api/internal/handlers"
	"sre-chat-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRouter configures and returns a Gin router with all routes and middleware
func SetupRouter(logger *zap.Logger) *gin.Engine {
	// Set up Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.LoggerMiddleware(logger))
	router.Use(middleware.RecoveryMiddleware(logger))

	// Initialize handlers
	sseHandler := handlers.NewSSEHandler()
	messageHandler := handlers.NewMessageHandler(sseHandler)
	healthHandler := handlers.NewHealthHandler()

	// Serve web interface
	router.StaticFile("/", "./web/index.html")
	router.Static("/static", "./web")

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/healthcheck", healthHandler.HealthCheck)

		// Server-Sent Events for real-time updates
		v1.GET("/messages/stream", sseHandler.StreamMessages)

		// Message routes
		v1.POST("/messages", messageHandler.CreateMessage)
		v1.GET("/messages", messageHandler.GetMessages)
		v1.GET("/messages/:id", messageHandler.GetMessage)
		v1.PUT("/messages/:id", messageHandler.UpdateMessage)
		v1.DELETE("/messages/:id", messageHandler.DeleteMessage)
	}

	return router
}
