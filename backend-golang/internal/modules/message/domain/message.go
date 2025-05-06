package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Message represents a message entity
type Message struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	TenantID  uuid.UUID       `json:"tenant_id" db:"tenant_id"`
	Payload   json.RawMessage `json:"payload" db:"payload"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

// MessageFilter represents filter for message search
type MessageFilter struct {
	TenantID uuid.UUID `json:"tenant_id"`
	Cursor   string    `json:"cursor"`
	Limit    int       `json:"limit"`
}

// MessageRepository defines the interface for message data operations
type MessageRepository interface {
	Create(ctx context.Context, message *Message) error
	FindByID(ctx context.Context, tenantID, messageID uuid.UUID) (*Message, error)
	FindByTenant(ctx context.Context, filter MessageFilter) ([]*Message, string, error)
	FindAll(ctx context.Context, cursor string, limit int) ([]*Message, string, error)
	Update(ctx context.Context, message *Message) error
	Delete(ctx context.Context, tenantID, messageID uuid.UUID) error
	WithTransaction(ctx context.Context, fn func(MessageRepository) error) error
}

// MessageUseCase defines the interface for message business logic
type MessageUseCase interface {
	Create(ctx context.Context, message *Message) error
	GetByID(ctx context.Context, tenantID, messageID uuid.UUID) (*Message, error)
	GetByTenant(ctx context.Context, filter MessageFilter) ([]*Message, string, error)
	GetMessages(ctx context.Context, cursor string, limit int) ([]*Message, string, error)
	Update(ctx context.Context, message *Message) error
	Delete(ctx context.Context, tenantID, messageID uuid.UUID) error
} 