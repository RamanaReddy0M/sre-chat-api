package handlers

import (
	"net/http"
	"sre-chat-api/internal/database"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck handles GET /api/v1/healthcheck
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Check database connection
	if database.DB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"message": "Database connection not established",
		})
		return
	}

	sqlDB, err := database.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"message": "Failed to get database instance",
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"message": "Database ping failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "Service is operational",
	})
}
