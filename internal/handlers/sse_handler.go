package handlers

import (
	"sre-chat-api/internal/database"
	"sre-chat-api/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// SSEMessage represents a message sent via SSE
type SSEMessage struct {
	Type    string      `json:"type"`
	Message models.Message `json:"message,omitempty"`
}

// SSEHandler handles Server-Sent Events for real-time message updates
type SSEHandler struct {
	clients map[chan SSEMessage]bool
	notify  chan models.Message
}

// NewSSEHandler creates a new SSE handler
func NewSSEHandler() *SSEHandler {
	handler := &SSEHandler{
		clients: make(map[chan SSEMessage]bool),
		notify:  make(chan models.Message, 100),
	}
	
	// Start broadcasting goroutine
	go handler.broadcast()
	
	return handler
}

// StreamMessages handles SSE connection for real-time message updates
func (h *SSEHandler) StreamMessages(c *gin.Context) {
	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Create a channel for this client
	messageChan := make(chan SSEMessage, 10)
	h.clients[messageChan] = true

	// Send initial connection message
	initialMsg := SSEMessage{
		Type: "connected",
	}
	sendSSE(c, initialMsg)

	// Clean up when client disconnects
	defer func() {
		delete(h.clients, messageChan)
		close(messageChan)
		log.Println("SSE client disconnected")
	}()

	// Keep connection alive and send messages
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-messageChan:
			if err := sendSSE(c, msg); err != nil {
				return
			}
			c.Writer.Flush()
		case <-ticker.C:
			// Send keepalive ping
			ping := SSEMessage{Type: "ping"}
			if err := sendSSE(c, ping); err != nil {
				return
			}
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}

// NotifyNewMessage sends a new message to all connected clients
func (h *SSEHandler) NotifyNewMessage(message models.Message) {
	select {
	case h.notify <- message:
	default:
		// Channel full, skip (prevents blocking)
	}
}

// broadcast sends new messages to all connected clients
func (h *SSEHandler) broadcast() {
	for {
		message := <-h.notify
		
		// Load group relationship
		database.DB.Preload("Group").First(&message, message.ID)
		
		sseMsg := SSEMessage{
			Type:    "message",
			Message: message,
		}

		// Send to all clients
		for clientChan := range h.clients {
			select {
			case clientChan <- sseMsg:
			default:
				// Client channel full, skip
			}
		}
	}
}

// sendSSE sends a message in SSE format
func sendSSE(c *gin.Context, msg SSEMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	
	_, err = fmt.Fprintf(c.Writer, "data: %s\n\n", data)
	return err
}

