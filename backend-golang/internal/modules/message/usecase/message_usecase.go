package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"sample-stack-golang/internal/modules/message/domain"
)

// MessageUsecase implements message business logic
type MessageUsecase struct {
	messageRepo domain.MessageRepository
}

// NewMessageUsecase creates a new message usecase
func NewMessageUsecase(messageRepo domain.MessageRepository) *MessageUsecase {
	return &MessageUsecase{
		messageRepo: messageRepo,
	}
}

// Create creates a new message
func (u *MessageUsecase) Create(ctx context.Context, message *domain.Message) error {
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()
	return u.messageRepo.Create(ctx, message)
}

// GetByID gets a message by ID
func (u *MessageUsecase) GetByID(ctx context.Context, tenantID, messageID uuid.UUID) (*domain.Message, error) {
	return u.messageRepo.FindByID(ctx, tenantID, messageID)
}

// GetByTenant gets messages by tenant ID
func (u *MessageUsecase) GetByTenant(ctx context.Context, filter domain.MessageFilter) ([]*domain.Message, string, error) {
	return u.messageRepo.FindByTenant(ctx, filter)
}

// Update updates a message
func (u *MessageUsecase) Update(ctx context.Context, message *domain.Message) error {
	message.UpdatedAt = time.Now()
	return u.messageRepo.Update(ctx, message)
}

// Delete deletes a message
func (u *MessageUsecase) Delete(ctx context.Context, tenantID, messageID uuid.UUID) error {
	return u.messageRepo.Delete(ctx, tenantID, messageID)
}

// WithTransaction executes a function within a transaction
func (u *MessageUsecase) WithTransaction(ctx context.Context, fn func(*MessageUsecase) error) error {
	return u.messageRepo.WithTransaction(ctx, func(repo domain.MessageRepository) error {
		usecase := NewMessageUsecase(repo)
		return fn(usecase)
	})
} 