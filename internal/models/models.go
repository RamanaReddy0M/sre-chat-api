package models

import (
	"time"

	"gorm.io/gorm"
)

// Group represents a chat group
type Group struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Messages    []Message `json:"messages,omitempty" gorm:"foreignKey:GroupID"`
}

// Message represents a chat message
type Message struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	GroupID   uint           `json:"group_id" gorm:"not null;index"`
	User      string         `json:"user" gorm:"not null"`
	Content   string         `json:"content" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Group     Group          `json:"group,omitempty" gorm:"foreignKey:GroupID"`
}

// CreateMessageRequest represents the request body for creating a message
type CreateMessageRequest struct {
	User    string `json:"user" binding:"required"`
	Content string `json:"content" binding:"required"`
	GroupID uint   `json:"group_id"`
}

// UpdateMessageRequest represents the request body for updating a message
type UpdateMessageRequest struct {
	Content string `json:"content" binding:"required"`
}
