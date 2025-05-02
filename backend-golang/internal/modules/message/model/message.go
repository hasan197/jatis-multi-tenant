package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Message merepresentasikan pesan dalam sistem
type Message struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	TenantID  uuid.UUID       `json:"tenant_id" db:"tenant_id"`
	Payload   json.RawMessage `json:"payload" db:"payload"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

// MessageFilter digunakan untuk filter dalam pencarian pesan
type MessageFilter struct {
	TenantID uuid.UUID `json:"tenant_id"`
	Cursor   string    `json:"cursor"`
	Limit    int       `json:"limit"`
}

// MessageRepository interface untuk akses data message
type MessageRepository interface {
	Create(message *Message) error
	FindByID(tenantID, messageID uuid.UUID) (*Message, error)
	FindByTenant(filter MessageFilter) ([]*Message, string, error)
	Update(message *Message) error
	Delete(tenantID, messageID uuid.UUID) error
	WithTransaction(fn func(MessageRepository) error) error
} 