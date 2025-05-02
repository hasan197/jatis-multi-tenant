package postgresql

import (
	"database/sql"
	"errors"
	"time"

	"sample-stack-golang/internal/modules/user/domain"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository membuat instance baru dari UserRepository dengan database PostgreSQL
func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

// FindAll mengambil semua user dari database
func (r *userRepository) FindAll() ([]domain.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY id DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// FindByID mengambil user berdasarkan ID
func (r *userRepository) FindByID(id uint) (domain.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return user, nil
}

// FindByEmail mengambil user berdasarkan email
func (r *userRepository) FindByEmail(email string) (domain.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return user, nil
}

// Create membuat user baru di database
func (r *userRepository) Create(user domain.User) (domain.User, error) {
	query := `
		INSERT INTO users (name, email, password, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at, updated_at
	`
	
	now := time.Now()
	err := r.db.QueryRow(
		query, 
		user.Name, 
		user.Email, 
		user.Password, 
		now, 
		now,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// Update memperbarui data user
func (r *userRepository) Update(user domain.User) (domain.User, error) {
	query := `
		UPDATE users 
		SET name = $1, email = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, name, email, created_at, updated_at
	`
	
	now := time.Now()
	row := r.db.QueryRow(
		query, 
		user.Name, 
		user.Email, 
		now, 
		user.ID,
	)
	
	var updatedUser domain.User
	err := row.Scan(
		&updatedUser.ID, 
		&updatedUser.Name, 
		&updatedUser.Email, 
		&updatedUser.CreatedAt, 
		&updatedUser.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return updatedUser, nil
}

// Delete menghapus user dari database
func (r *userRepository) Delete(id uint) error {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	
	return nil
} 