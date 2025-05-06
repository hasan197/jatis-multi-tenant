package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sample-stack-golang/internal/modules/tenant/domain"
)

// TenantRepository implements domain.TenantRepository
type TenantRepository struct {
	db *pgxpool.Pool
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository(db *pgxpool.Pool) domain.TenantRepository {
	return &TenantRepository{
		db: db,
	}
}

// Create creates a new tenant
func (r *TenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	// Generate UUID for new tenant
	tenant.ID = uuid.New().String()

	// Set default workers if not specified
	if tenant.Workers <= 0 {
		tenant.Workers = 3 // Default worker count
	}

	// Start transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert tenant
	query := `
		INSERT INTO tenants (id, name, description, status, workers, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = tx.Exec(ctx, query,
		tenant.ID,
		tenant.Name,
		tenant.Description,
		tenant.Status,
		tenant.Workers,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}

	// Create messages partition
	_, err = tx.Exec(ctx, "SELECT create_messages_partition($1)", tenant.ID)
	if err != nil {
		return fmt.Errorf("failed to create messages partition: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByID gets a tenant by ID
func (r *TenantRepository) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	query := `
		SELECT id, name, description, status, workers, created_at, updated_at
		FROM tenants
		WHERE id = $1`

	var tenant domain.Tenant
	err := r.db.QueryRow(ctx, query, id).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Description,
		&tenant.Status,
		&tenant.Workers,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	return &tenant, nil
}

// Update updates a tenant
func (r *TenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		UPDATE tenants
		SET name = $1, description = $2, status = $3, updated_at = $4
		WHERE id = $5`

	result, err := r.db.Exec(ctx, query,
		tenant.Name,
		tenant.Description,
		tenant.Status,
		time.Now(),
		tenant.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update tenant: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// Delete deletes a tenant
func (r *TenantRepository) Delete(ctx context.Context, id string) error {
	// Start transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Drop the partition for this tenant
	_, err = tx.Exec(ctx, "SELECT drop_messages_partition($1)", id)
	if err != nil {
		return fmt.Errorf("failed to drop messages partition: %w", err)
	}

	// Delete the tenant
	query := `DELETE FROM tenants WHERE id = $1`
	result, err := tx.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// List lists all tenants
func (r *TenantRepository) List(ctx context.Context) ([]*domain.Tenant, error) {
	query := `
		SELECT id, name, description, status, workers, created_at, updated_at
		FROM tenants
		ORDER BY id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}
	defer rows.Close()

	var tenants []*domain.Tenant
	for rows.Next() {
		var tenant domain.Tenant
		err := rows.Scan(
			&tenant.ID,
			&tenant.Name,
			&tenant.Description,
			&tenant.Status,
			&tenant.Workers,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tenant: %w", err)
		}
		tenants = append(tenants, &tenant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tenant rows: %w", err)
	}

	return tenants, nil
}

// UpdateConcurrency updates the concurrency configuration for a tenant
func (r *TenantRepository) UpdateConcurrency(ctx context.Context, id string, workers int) error {
	query := `
		UPDATE tenants
		SET workers = $1, updated_at = $2
		WHERE id = $3`

	result, err := r.db.Exec(ctx, query, workers, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update tenant concurrency: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}