package postgresql

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"sample-stack-golang/internal/modules/message/domain"
)

// MessageRepository implements domain.MessageRepository
type MessageRepository struct {
	db *pgx.Conn
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *pgx.Conn) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

// Create creates a new message
func (r *MessageRepository) Create(ctx context.Context, message *domain.Message) error {
	query := `
		INSERT INTO messages (id, tenant_id, payload, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		message.ID,
		message.TenantID,
		message.Payload,
		message.CreatedAt,
		message.UpdatedAt,
	)

	return err
}

// GetByID gets a message by ID
func (r *MessageRepository) GetByID(ctx context.Context, tenantID, messageID uuid.UUID) (*domain.Message, error) {
	query := `
		SELECT id, tenant_id, payload, created_at, updated_at
		FROM messages
		WHERE id = $1 AND tenant_id = $2
	`

	var message domain.Message
	err := r.db.QueryRow(ctx, query, messageID, tenantID).Scan(
		&message.ID,
		&message.TenantID,
		&message.Payload,
		&message.CreatedAt,
		&message.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, sql.ErrNoRows
	}

	if err != nil {
		return nil, err
	}

	return &message, nil
}

// GetByTenant gets messages by tenant ID
func (r *MessageRepository) GetByTenant(ctx context.Context, filter domain.MessageFilter) ([]*domain.Message, string, error) {
	query := `
		SELECT id, tenant_id, payload, created_at, updated_at
		FROM messages
		WHERE tenant_id = $1
	`

	args := []interface{}{filter.TenantID}

	if filter.Cursor != "" {
		query += " AND id > $2"
		args = append(args, filter.Cursor)
	}

	query += " ORDER BY id ASC LIMIT $2"
	args = append(args, filter.Limit+1)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var message domain.Message
		err := rows.Scan(
			&message.ID,
			&message.TenantID,
			&message.Payload,
			&message.CreatedAt,
			&message.UpdatedAt,
		)
		if err != nil {
			return nil, "", err
		}
		messages = append(messages, &message)
	}

	if err = rows.Err(); err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(messages) > filter.Limit {
		nextCursor = messages[filter.Limit-1].ID.String()
		messages = messages[:filter.Limit]
	}

	return messages, nextCursor, nil
}

// Update updates a message
func (r *MessageRepository) Update(ctx context.Context, message *domain.Message) error {
	query := `
		UPDATE messages
		SET payload = $1, updated_at = $2
		WHERE id = $3 AND tenant_id = $4
	`

	message.UpdatedAt = time.Now()

	result, err := r.db.Exec(ctx, query,
		message.Payload,
		message.UpdatedAt,
		message.ID,
		message.TenantID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete deletes a message
func (r *MessageRepository) Delete(ctx context.Context, tenantID, messageID uuid.UUID) error {
	query := `
		DELETE FROM messages
		WHERE id = $1 AND tenant_id = $2
	`

	result, err := r.db.Exec(ctx, query, messageID, tenantID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// WithTransaction executes a function within a transaction
func (r *MessageRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit(ctx)
} 