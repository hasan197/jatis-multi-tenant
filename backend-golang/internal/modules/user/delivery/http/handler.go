package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	
	"sample-stack-golang/internal/modules/user/domain"
	"sample-stack-golang/pkg/logger"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
	logger      *zap.Logger
}

// NewUserHandler membuat instance baru UserHandler
func NewUserHandler(userUseCase domain.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		logger:      logger.WithContext(
			zap.String("component", "user_handler"),
			zap.String("version", "1.0"),
		),
	}
}

// getRequestContext mengambil konteks umum dari request
func (h *UserHandler) getRequestContext(c *gin.Context) []zap.Field {
	return []zap.Field{
		zap.String("request_id", c.GetString("request_id")),
		zap.String("client_ip", c.ClientIP()),
		zap.String("user_agent", c.GetHeader("User-Agent")),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	}
}

// GetUsers menangani request untuk mendapatkan semua user
func (h *UserHandler) GetUsers(c *gin.Context) {
	ctx := h.getRequestContext(c)
	log := h.logger.With(ctx...)
	
	log.Info("memulai request get users")
	
	users, err := h.userUseCase.GetUsers()
	if err != nil {
		log.Error("gagal mendapatkan users",
			zap.Error(err),
			zap.String("error_type", "internal_error"),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	log.Info("berhasil mendapatkan users",
		zap.Int("total_users", len(users)),
		zap.Duration("duration", time.Since(c.GetTime("start_time"))),
	)
	
	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

// GetUser menangani request untuk mendapatkan satu user berdasarkan ID
func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	ctx := h.getRequestContext(c)
	log := h.logger.With(append(ctx, zap.String("user_id", idParam))...)
	
	log.Info("memulai request get user")
	
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Error("invalid user id",
			zap.Error(err),
			zap.String("error_type", "validation_error"),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}
	
	user, err := h.userUseCase.GetUser(uint(id))
	if err != nil {
		log.Error("gagal mendapatkan user",
			zap.Error(err),
			zap.String("error_type", "not_found"),
			zap.Uint("user_id", uint(id)),
		)
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	log.Info("berhasil mendapatkan user",
		zap.Uint("user_id", user.ID),
		zap.String("email", user.Email),
		zap.Duration("duration", time.Since(c.GetTime("start_time"))),
	)
	
	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// CreateUser menangani request untuk membuat user baru
func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx := h.getRequestContext(c)
	log := h.logger.With(ctx...)
	
	log.Info("memulai request create user")
	
	var userInput struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	
	if err := c.ShouldBindJSON(&userInput); err != nil {
		log.Error("gagal parse request body",
			zap.Error(err),
			zap.String("error_type", "validation_error"),
			zap.Any("request_body", c.Request.Body),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	// Log data sensitif dengan level debug
	log.Debug("data user dari request",
		zap.String("name", userInput.Name),
		zap.String("email", userInput.Email),
	)
	
	user := domain.User{
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: userInput.Password,
	}
	
	createdUser, err := h.userUseCase.CreateUser(user)
	if err != nil {
		log.Error("gagal membuat user",
			zap.Error(err),
			zap.String("error_type", "creation_error"),
			zap.String("email", userInput.Email),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	log.Info("user berhasil dibuat",
		zap.Uint("user_id", createdUser.ID),
		zap.String("email", createdUser.Email),
		zap.Duration("duration", time.Since(c.GetTime("start_time"))),
	)
	
	c.JSON(http.StatusCreated, gin.H{
		"data": createdUser,
	})
}

// UpdateUser menangani request untuk memperbarui data user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	ctx := h.getRequestContext(c)
	log := h.logger.With(append(ctx, zap.String("user_id", idParam))...)
	
	log.Info("memulai request update user")
	
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Error("invalid user id",
			zap.Error(err),
			zap.String("error_type", "validation_error"),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}
	
	var userInput struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}
	
	if err := c.ShouldBindJSON(&userInput); err != nil {
		log.Error("gagal parse request body",
			zap.Error(err),
			zap.String("error_type", "validation_error"),
			zap.Any("request_body", c.Request.Body),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	// Log data sensitif dengan level debug
	log.Debug("data update dari request",
		zap.String("name", userInput.Name),
		zap.String("email", userInput.Email),
	)
	
	user := domain.User{
		ID:    uint(id),
		Name:  userInput.Name,
		Email: userInput.Email,
	}
	
	updatedUser, err := h.userUseCase.UpdateUser(user)
	if err != nil {
		log.Error("gagal update user",
			zap.Error(err),
			zap.String("error_type", "update_error"),
			zap.Uint("user_id", uint(id)),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	log.Info("user berhasil diupdate",
		zap.Uint("user_id", updatedUser.ID),
		zap.String("email", updatedUser.Email),
		zap.Duration("duration", time.Since(c.GetTime("start_time"))),
	)
	
	c.JSON(http.StatusOK, gin.H{
		"data": updatedUser,
	})
}

// DeleteUser menangani request untuk menghapus user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	ctx := h.getRequestContext(c)
	log := h.logger.With(append(ctx, zap.String("user_id", idParam))...)
	
	log.Info("memulai request delete user")
	
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Error("invalid user id",
			zap.Error(err),
			zap.String("error_type", "validation_error"),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}
	
	err = h.userUseCase.DeleteUser(uint(id))
	if err != nil {
		log.Error("gagal menghapus user",
			zap.Error(err),
			zap.String("error_type", "deletion_error"),
			zap.Uint("user_id", uint(id)),
		)
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	log.Info("user berhasil dihapus",
		zap.Uint("user_id", uint(id)),
		zap.Duration("duration", time.Since(c.GetTime("start_time"))),
	)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully deleted",
	})
} 