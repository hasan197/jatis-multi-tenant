package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"sample-stack-golang/internal/modules/tenant/domain"
)

// TenantRepository implements domain.TenantRepository
type TenantRepository struct {
	pool *pgxpool.Pool
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository(pool *pgxpool.Pool) *TenantRepository {
	return &TenantRepository{
		pool: pool,
	}
}

// Create creates a new tenant
func (r *TenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		INSERT INTO tenants (id, name, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query,
		tenant.ID,
		tenant.Name,
		tenant.Status,
		tenant.CreatedAt,
		tenant.UpdatedAt,
	)
	return err
}

// GetByID gets a tenant by ID
func (r *TenantRepository) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	query := `
		SELECT id, name, status, created_at, updated_at
		FROM tenants
		WHERE id = $1
	`
	tenant := &domain.Tenant{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Status,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

// Update updates a tenant
func (r *TenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		UPDATE tenants
		SET name = $1, status = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.pool.Exec(ctx, query,
		tenant.Name,
		tenant.Status,
		tenant.UpdatedAt,
		tenant.ID,
	)
	return err
}

// Delete deletes a tenant
func (r *TenantRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tenants WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// List lists all tenants
func (r *TenantRepository) List(ctx context.Context) ([]*domain.Tenant, error) {
	query := `
		SELECT id, name, status, created_at, updated_at
		FROM tenants
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []*domain.Tenant
	for rows.Next() {
		tenant := &domain.Tenant{}
		err := rows.Scan(
			&tenant.ID,
			&tenant.Name,
			&tenant.Status,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}
	return tenants, nil
} 