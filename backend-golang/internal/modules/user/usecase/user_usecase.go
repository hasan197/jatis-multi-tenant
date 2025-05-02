package usecase

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	
	"sample-stack-golang/internal/modules/user/domain"
)

type userUseCase struct {
	userRepo domain.UserRepository
}

// NewUserUseCase membuat instance baru UserUseCase
func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

// GetUsers mendapatkan semua data user
func (uc *userUseCase) GetUsers() ([]domain.User, error) {
	return uc.userRepo.FindAll()
}

// GetUser mendapatkan user berdasarkan ID
func (uc *userUseCase) GetUser(id uint) (domain.User, error) {
	return uc.userRepo.FindByID(id)
}

// CreateUser membuat user baru
func (uc *userUseCase) CreateUser(user domain.User) (domain.User, error) {
	// Validasi data
	if user.Name == "" {
		return domain.User{}, errors.New("name is required")
	}
	
	if user.Email == "" {
		return domain.User{}, errors.New("email is required")
	}
	
	if user.Password == "" {
		return domain.User{}, errors.New("password is required")
	}
	
	// Cek apakah email sudah digunakan
	existingUser, err := uc.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser.ID != 0 {
		// Email sudah terdaftar
		log.Printf("Email %s sudah terdaftar dengan ID user %d", user.Email, existingUser.ID)
		return domain.User{}, errors.New("email already exists")
	} else if err != nil && !strings.Contains(err.Error(), "not found") {
		// Error selain "not found"
		log.Printf("Error saat mencari user dengan email %s: %v", user.Email, err)
		return domain.User{}, err
	}
	
	// Log untuk debugging
	log.Printf("Mencoba membuat user baru dengan email: %s", user.Email)
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error saat hashing password: %v", err)
		return domain.User{}, err
	}
	
	user.Password = string(hashedPassword)
	
	// Simpan user baru
	newUser, err := uc.userRepo.Create(user)
	if err != nil {
		log.Printf("Error saat menyimpan user: %v", err)
		return domain.User{}, err
	}
	
	log.Printf("User berhasil dibuat dengan ID: %d", newUser.ID)
	return newUser, nil
}

// UpdateUser memperbarui data user
func (uc *userUseCase) UpdateUser(user domain.User) (domain.User, error) {
	// Validasi data
	if user.ID == 0 {
		return domain.User{}, errors.New("user ID is required")
	}
	
	if user.Name == "" {
		return domain.User{}, errors.New("name is required")
	}
	
	if user.Email == "" {
		return domain.User{}, errors.New("email is required")
	}
	
	// Cek apakah user ada
	_, err := uc.userRepo.FindByID(user.ID)
	if err != nil {
		return domain.User{}, err
	}
	
	// Cek apakah email sudah digunakan oleh user lain
	userWithEmail, err := uc.userRepo.FindByEmail(user.Email)
	if err == nil && userWithEmail.ID != 0 && userWithEmail.ID != user.ID {
		return domain.User{}, errors.New("email already exists")
	}
	
	// Update user
	return uc.userRepo.Update(user)
}

// DeleteUser menghapus user berdasarkan ID
func (uc *userUseCase) DeleteUser(id uint) error {
	return uc.userRepo.Delete(id)
} 