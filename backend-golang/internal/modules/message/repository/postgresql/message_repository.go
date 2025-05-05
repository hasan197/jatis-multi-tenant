package postgresql

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"
	"sample-stack-golang/internal/modules/message/domain"
)

// DBConn adalah interface untuk koneksi database
type DBConn interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

// MessageRepository implements domain.MessageRepository
type MessageRepository struct {
	db DBConn
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(pool *pgxpool.Pool) domain.MessageRepository {
	return &MessageRepository{
		db: pool,
	}
}

// Create creates a new message
func (r *MessageRepository) Create(ctx context.Context, message *domain.Message) error {
	// Buat partisi terlebih dahulu
	_, err := r.db.Exec(ctx, "SELECT create_messages_partition($1)", message.TenantID)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO messages (id, tenant_id, payload, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = r.db.Exec(ctx, query,
		message.ID,
		message.TenantID,
		message.Payload,
		message.CreatedAt,
		message.UpdatedAt,
	)

	return err
}

// FindByID gets a message by ID
func (r *MessageRepository) FindByID(ctx context.Context, tenantID, messageID uuid.UUID) (*domain.Message, error) {
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

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}

	if err != nil {
		return nil, err
	}

	return &message, nil
}

// FindByTenant gets messages by tenant ID
func (r *MessageRepository) FindByTenant(ctx context.Context, filter domain.MessageFilter) ([]*domain.Message, string, error) {
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
func (r *MessageRepository) WithTransaction(ctx context.Context, fn func(domain.MessageRepository) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Buat repository baru dengan transaksi
	txRepo := &MessageRepository{db: tx}

	// Jalankan fungsi dengan repository transaksi
	if err := fn(txRepo); err != nil {
		return err
	}

	return tx.Commit(ctx)
} 