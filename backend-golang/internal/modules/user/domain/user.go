package domain

import (
	"time"
)

// User represents user entity
type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Password tidak ditampilkan dalam JSON response
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository mendefinisikan kontrak untuk repository user
type UserRepository interface {
	FindAll() ([]User, error)
	FindByID(id uint) (User, error)
	FindByEmail(email string) (User, error)
	Create(user User) (User, error)
	Update(user User) (User, error)
	Delete(id uint) error
}

// UserUseCase mendefinisikan kontrak untuk use case user
type UserUseCase interface {
	GetUsers() ([]User, error)
	GetUser(id uint) (User, error)
	CreateUser(user User) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(id uint) error
} 