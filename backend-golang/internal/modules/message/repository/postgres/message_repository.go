package postgres

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"sample-stack/internal/modules/message/model"
)

type messageRepository struct {
	db *sql.DB
}

// NewMessageRepository membuat instance baru MessageRepository
func NewMessageRepository(db *sql.DB) model.MessageRepository {
	return &messageRepository{
		db: db,
	}
}

// Create menyimpan pesan baru
func (r *messageRepository) Create(message *model.Message) error {
	query := `
		INSERT INTO messages (id, tenant_id, payload, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	if message.ID == uuid.Nil {
		message.ID = uuid.New()
	}

	now := time.Now()
	message.CreatedAt = now
	message.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		message.ID,
		message.TenantID,
		message.Payload,
		message.CreatedAt,
		message.UpdatedAt,
	).Scan(&message.ID, &message.CreatedAt, &message.UpdatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// Handle specific PostgreSQL errors
			switch pqErr.Code {
			case "23505": // unique_violation
				return fmt.Errorf("message already exists: %v", err)
			default:
				return fmt.Errorf("failed to create message: %v", err)
			}
		}
		return fmt.Errorf("failed to create message: %v", err)
	}

	return nil
}

// FindByID mencari pesan berdasarkan ID dan tenant ID
func (r *messageRepository) FindByID(tenantID, messageID uuid.UUID) (*model.Message, error) {
	query := `
		SELECT id, tenant_id, payload, created_at, updated_at
		FROM messages
		WHERE tenant_id = $1 AND id = $2`

	message := &model.Message{}
	err := r.db.QueryRow(query, tenantID, messageID).Scan(
		&message.ID,
		&message.TenantID,
		&message.Payload,
		&message.CreatedAt,
		&message.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find message: %v", err)
	}

	return message, nil
}

// FindByTenant mencari pesan berdasarkan tenant dengan pagination
func (r *messageRepository) FindByTenant(filter model.MessageFilter) ([]*model.Message, string, error) {
	var query string
	var args []interface{}
	var err error

	if filter.Cursor != "" {
		decodedCursor, err := base64.StdEncoding.DecodeString(filter.Cursor)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor: %v", err)
		}

		query = `
			SELECT id, tenant_id, payload, created_at, updated_at
			FROM messages
			WHERE tenant_id = $1 AND created_at < $2
			ORDER BY created_at DESC
			LIMIT $3`
		args = []interface{}{filter.TenantID, string(decodedCursor), filter.Limit}
	} else {
		query = `
			SELECT id, tenant_id, payload, created_at, updated_at
			FROM messages
			WHERE tenant_id = $1
			ORDER BY created_at DESC
			LIMIT $2`
		args = []interface{}{filter.TenantID, filter.Limit}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("failed to query messages: %v", err)
	}
	defer rows.Close()

	var messages []*model.Message
	var lastCreatedAt time.Time

	for rows.Next() {
		message := &model.Message{}
		err := rows.Scan(
			&message.ID,
			&message.TenantID,
			&message.Payload,
			&message.CreatedAt,
			&message.UpdatedAt,
		)
		if err != nil {
			return nil, "", fmt.Errorf("failed to scan message: %v", err)
		}
		messages = append(messages, message)
		lastCreatedAt = message.CreatedAt
	}

	if err = rows.Err(); err != nil {
		return nil, "", fmt.Errorf("error iterating messages: %v", err)
	}

	var nextCursor string
	if len(messages) == filter.Limit {
		nextCursor = base64.StdEncoding.EncodeToString([]byte(lastCreatedAt.Format(time.RFC3339Nano)))
	}

	return messages, nextCursor, nil
}

// Update memperbarui pesan yang ada
func (r *messageRepository) Update(message *model.Message) error {
	query := `
		UPDATE messages
		SET payload = $1, updated_at = $2
		WHERE tenant_id = $3 AND id = $4
		RETURNING updated_at`

	message.UpdatedAt = time.Now()

	err := r.db.QueryRow(
		query,
		message.Payload,
		message.UpdatedAt,
		message.TenantID,
		message.ID,
	).Scan(&message.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("message not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update message: %v", err)
	}

	return nil
}

// Delete menghapus pesan
func (r *messageRepository) Delete(tenantID, messageID uuid.UUID) error {
	query := `
		DELETE FROM messages
		WHERE tenant_id = $1 AND id = $2`

	result, err := r.db.Exec(query, tenantID, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
} 