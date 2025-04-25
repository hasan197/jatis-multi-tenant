package http

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	
	"sample-stack/internal/modules/user/domain"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
}

// NewUserHandler membuat instance baru UserHandler
func NewUserHandler(userUseCase domain.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// GetUsers menangani request untuk mendapatkan semua user
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userUseCase.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

// GetUser menangani request untuk mendapatkan satu user berdasarkan ID
func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}
	
	user, err := h.userUseCase.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// CreateUser menangani request untuk membuat user baru
func (h *UserHandler) CreateUser(c *gin.Context) {
	log.Println("CreateUser: Mulai memproses request...")
	var userInput struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	
	if err := c.ShouldBindJSON(&userInput); err != nil {
		log.Printf("CreateUser: Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	log.Printf("CreateUser: Data dari client - Name: %s, Email: %s", userInput.Name, userInput.Email)
	
	user := domain.User{
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: userInput.Password,
	}
	
	createdUser, err := h.userUseCase.CreateUser(user)
	if err != nil {
		log.Printf("CreateUser: Error dari usecase: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	log.Printf("CreateUser: User berhasil dibuat dengan ID %d", createdUser.ID)
	c.JSON(http.StatusCreated, gin.H{
		"data": createdUser,
	})
}

// UpdateUser menangani request untuk memperbarui data user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	startTime := time.Now()
	log.Println("UpdateUser: Mulai memproses request...")
	
	idParam := c.Param("id")
	log.Printf("UpdateUser: ID dari URL parameter: %s", idParam)
	
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Printf("UpdateUser: Error parsing ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}
	
	// Metode sederhana, langsung binding JSON
	var userInput struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}
	
	if err := c.ShouldBindJSON(&userInput); err != nil {
		log.Printf("UpdateUser: Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	log.Printf("UpdateUser: Data dari client - Name: %s, Email: %s", userInput.Name, userInput.Email)
	
	user := domain.User{
		ID:    uint(id),
		Name:  userInput.Name,
		Email: userInput.Email,
	}
	
	updatedUser, err := h.userUseCase.UpdateUser(user)
	if err != nil {
		log.Printf("UpdateUser: Error dari usecase: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	duration := time.Since(startTime)
	log.Printf("UpdateUser: User berhasil diupdate. ID: %d, Durasi: %v", updatedUser.ID, duration)
	
	c.JSON(http.StatusOK, gin.H{
		"data": updatedUser,
	})
}

// DeleteUser menangani request untuk menghapus user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}
	
	err = h.userUseCase.DeleteUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully deleted",
	})
} 