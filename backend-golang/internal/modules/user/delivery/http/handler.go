package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
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
func (h *UserHandler) getRequestContext(c echo.Context) []zap.Field {
	return []zap.Field{
		zap.String("request_id", c.Request().Header.Get("X-Request-ID")),
		zap.String("client_ip", c.RealIP()),
		zap.String("user_agent", c.Request().UserAgent()),
		zap.String("method", c.Request().Method),
		zap.String("path", c.Request().URL.Path),
	}
}

// GetUsers menangani request untuk mendapatkan semua user
func (h *UserHandler) GetUsers(c echo.Context) error {
	ctx := h.getRequestContext(c)
	log := h.logger.With(ctx...)
	
	log.Info("memulai request get users")
	
	users, err := h.userUseCase.GetUsers()
	if err != nil {
		log.Error("gagal mendapatkan users",
			zap.Error(err),
			zap.String("error_type", "internal_error"),
		)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	log.Info("berhasil mendapatkan users",
		zap.Int("total_users", len(users)),
		zap.Duration("duration", time.Since(time.Now())),
	)
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": users,
	})
}

// GetUser menangani request untuk mendapatkan satu user berdasarkan ID
func (h *UserHandler) GetUser(c echo.Context) error {
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid user ID",
		})
	}
	
	user, err := h.userUseCase.GetUser(uint(id))
	if err != nil {
		log.Error("gagal mendapatkan user",
			zap.Error(err),
			zap.String("error_type", "not_found"),
			zap.Uint("user_id", uint(id)),
		)
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	log.Info("berhasil mendapatkan user",
		zap.Uint("user_id", user.ID),
		zap.String("email", user.Email),
		zap.Duration("duration", time.Since(time.Now())),
	)
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": user,
	})
}

// CreateUser menangani request untuk membuat user baru
func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := h.getRequestContext(c)
	log := h.logger.With(ctx...)
	
	log.Info("memulai request create user")
	
	var userInput struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}
	
	if err := c.Bind(&userInput); err != nil {
		log.Error("gagal parse request body",
			zap.Error(err),
			zap.String("error_type", "validation_error"),
		)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	// Validasi input
	if err := c.Validate(&userInput); err != nil {
		log.Error("validasi input gagal",
			zap.Error(err),
			zap.String("error_type", "validation_error"),
		)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	log.Info("user berhasil dibuat",
		zap.Uint("user_id", createdUser.ID),
		zap.String("email", createdUser.Email),
		zap.Duration("duration", time.Since(time.Now())),
	)
	
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": createdUser,
	})
}

// UpdateUser menangani request untuk memperbarui data user
func (h *UserHandler) UpdateUser(c echo.Context) error {
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid user ID",
		})
	}
	
	var userInput struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}
	
	if err := c.Bind(&userInput); err != nil {
		log.Error("gagal parse request body",
			zap.Error(err),
			zap.String("error_type", "validation_error"),
		)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	// Validasi input
	if err := c.Validate(&userInput); err != nil {
		log.Error("validasi input gagal",
			zap.Error(err),
			zap.String("error_type", "validation_error"),
		)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	log.Info("user berhasil diupdate",
		zap.Uint("user_id", updatedUser.ID),
		zap.String("email", updatedUser.Email),
		zap.Duration("duration", time.Since(time.Now())),
	)
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": updatedUser,
	})
}

// DeleteUser menangani request untuk menghapus user
func (h *UserHandler) DeleteUser(c echo.Context) error {
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid user ID",
		})
	}
	
	if err := h.userUseCase.DeleteUser(uint(id)); err != nil {
		log.Error("gagal delete user",
			zap.Error(err),
			zap.String("error_type", "delete_error"),
			zap.Uint("user_id", uint(id)),
		)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	log.Info("user berhasil dihapus",
		zap.Uint("user_id", uint(id)),
		zap.Duration("duration", time.Since(time.Now())),
	)
	
	return c.NoContent(http.StatusNoContent)
} 