package handlers

import (
	"bytes"
	"sre-chat-api/internal/config"
	"sre-chat-api/internal/database"
	"sre-chat-api/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use test database configuration
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			DBName:   "chat_test_db",
			SSLMode:  "disable",
		},
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping test: database connection failed: %v", err)
		return nil
	}

	// Clean up and migrate
	db.Exec("DROP TABLE IF EXISTS messages CASCADE")
	db.Exec("DROP TABLE IF EXISTS groups CASCADE")
	db.AutoMigrate(&models.Group{}, &models.Message{})

	// Create default group
	db.Create(&models.Group{
		Name:        "SRE Bootcamp",
		Description: "Default group for SRE Bootcamp exercises",
	})

	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	database.DB = db

	router := gin.New()
	messageHandler := NewMessageHandler(nil) // nil SSE handler for tests
	healthHandler := NewHealthHandler()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/healthcheck", healthHandler.HealthCheck)
		v1.POST("/messages", messageHandler.CreateMessage)
		v1.GET("/messages", messageHandler.GetMessages)
		v1.GET("/messages/:id", messageHandler.GetMessage)
		v1.PUT("/messages/:id", messageHandler.UpdateMessage)
		v1.DELETE("/messages/:id", messageHandler.DeleteMessage)
	}

	return router
}

func TestCreateMessage(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	router := setupRouter(db)

	// Test creating a message
	reqBody := models.CreateMessageRequest{
		User:    "testuser",
		Content: "Hello, SRE Bootcamp!",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/v1/messages", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Message
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "testuser", response.User)
	assert.Equal(t, "Hello, SRE Bootcamp!", response.Content)
	assert.NotZero(t, response.ID)
}

func TestGetMessages(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	router := setupRouter(db)

	// Create a test message
	var group models.Group
	db.Where("name = ?", "SRE Bootcamp").First(&group)
	message := models.Message{
		GroupID: group.ID,
		User:    "testuser",
		Content: "Test message",
	}
	db.Create(&message)

	// Test getting all messages
	req, _ := http.NewRequest("GET", "/api/v1/messages", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var messages []models.Message
	json.Unmarshal(w.Body.Bytes(), &messages)
	assert.GreaterOrEqual(t, len(messages), 1)
}

func TestGetMessage(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	router := setupRouter(db)

	// Create a test message
	var group models.Group
	db.Where("name = ?", "SRE Bootcamp").First(&group)
	message := models.Message{
		GroupID: group.ID,
		User:    "testuser",
		Content: "Test message",
	}
	db.Create(&message)

	// Test getting a specific message
	req, _ := http.NewRequest("GET", "/api/v1/messages/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Message
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, message.ID, response.ID)
	assert.Equal(t, "testuser", response.User)
}

func TestUpdateMessage(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	router := setupRouter(db)

	// Create a test message
	var group models.Group
	db.Where("name = ?", "SRE Bootcamp").First(&group)
	message := models.Message{
		GroupID: group.ID,
		User:    "testuser",
		Content: "Original message",
	}
	db.Create(&message)

	// Test updating the message
	updateReq := models.UpdateMessageRequest{
		Content: "Updated message",
	}
	jsonBody, _ := json.Marshal(updateReq)

	req, _ := http.NewRequest("PUT", "/api/v1/messages/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Message
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Updated message", response.Content)
}

func TestDeleteMessage(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	router := setupRouter(db)

	// Create a test message
	var group models.Group
	db.Where("name = ?", "SRE Bootcamp").First(&group)
	message := models.Message{
		GroupID: group.ID,
		User:    "testuser",
		Content: "Message to delete",
	}
	db.Create(&message)

	// Test deleting the message
	req, _ := http.NewRequest("DELETE", "/api/v1/messages/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify message is deleted
	var deletedMessage models.Message
	result := db.First(&deletedMessage, 1)
	assert.Error(t, result.Error)
}

func TestHealthCheck(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	router := setupRouter(db)

	req, _ := http.NewRequest("GET", "/api/v1/healthcheck", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "healthy", response["status"])
}
