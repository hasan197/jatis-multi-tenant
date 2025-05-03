package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	
	"sample-stack-golang/internal/modules/user/domain"
	"sample-stack-golang/pkg/logger"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
	logger      *logrus.Entry
}

// NewUserHandler membuat instance baru UserHandler
func NewUserHandler(userUseCase domain.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		logger:      logger.WithContext(map[string]interface{}{
			"component": "user_handler",
			"version":   "1.0",
		}),
	}
}

// getRequestContext mengambil konteks umum dari request
func (h *UserHandler) getRequestContext(c echo.Context) map[string]interface{} {
	return map[string]interface{}{
		"request_id": c.Request().Header.Get("X-Request-ID"),
		"client_ip":  c.RealIP(),
		"user_agent": c.Request().UserAgent(),
		"method":     c.Request().Method,
		"path":       c.Request().URL.Path,
	}
}

// GetUsers menangani request untuk mendapatkan semua user
func (h *UserHandler) GetUsers(c echo.Context) error {
	ctx := h.getRequestContext(c)
	log := h.logger.WithFields(ctx)
	
	log.Info("memulai request get users")
	
	users, err := h.userUseCase.GetUsers()
	if err != nil {
		log.WithError(err).WithField("error_type", "internal_error").Error("gagal mendapatkan users")
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	log.WithFields(map[string]interface{}{
		"total_users": len(users),
		"duration":    time.Since(time.Now()),
	}).Info("berhasil mendapatkan users")
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": users,
	})
}

// GetUser menangani request untuk mendapatkan satu user berdasarkan ID
func (h *UserHandler) GetUser(c echo.Context) error {
	idParam := c.Param("id")
	ctx := h.getRequestContext(c)
	ctx["user_id"] = idParam
	log := h.logger.WithFields(ctx)
	
	log.Info("memulai request get user")
	
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.WithError(err).WithField("error_type", "validation_error").Error("invalid user id")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid user ID",
		})
	}
	
	user, err := h.userUseCase.GetUser(uint(id))
	if err != nil {
		log.WithError(err).WithFields(map[string]interface{}{
			"error_type": "not_found",
			"user_id":    uint(id),
		}).Error("gagal mendapatkan user")
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	log.WithFields(map[string]interface{}{
		"user_id":  user.ID,
		"email":    user.Email,
		"duration": time.Since(time.Now()),
	}).Info("berhasil mendapatkan user")
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": user,
	})
}

// CreateUser menangani request untuk membuat user baru
func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := h.getRequestContext(c)
	log := h.logger.WithFields(ctx)
	
	log.Info("memulai request create user")
	
	var userInput struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}
	
	if err := c.Bind(&userInput); err != nil {
		log.WithError(err).WithField("error_type", "validation_error").Error("gagal parse request body")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	// Validasi input
	if err := c.Validate(&userInput); err != nil {
		log.WithError(err).WithField("error_type", "validation_error").Error("validasi input gagal")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	// Log data sensitif dengan level debug
	log.WithFields(map[string]interface{}{
		"name":  userInput.Name,
		"email": userInput.Email,
	}).Debug("data user dari request")
	
	user := domain.User{
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: userInput.Password,
	}
	
	createdUser, err := h.userUseCase.CreateUser(user)
	if err != nil {
		log.WithError(err).WithFields(map[string]interface{}{
			"error_type": "creation_error",
			"email":      userInput.Email,
		}).Error("gagal membuat user")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	log.WithFields(map[string]interface{}{
		"user_id":  createdUser.ID,
		"email":    createdUser.Email,
		"duration": time.Since(time.Now()),
	}).Info("user berhasil dibuat")
	
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": createdUser,
	})
}

// UpdateUser menangani request untuk memperbarui data user
func (h *UserHandler) UpdateUser(c echo.Context) error {
	idParam := c.Param("id")
	ctx := h.getRequestContext(c)
	ctx["user_id"] = idParam
	log := h.logger.WithFields(ctx)
	
	log.Info("memulai request update user")
	
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.WithError(err).WithField("error_type", "validation_error").Error("invalid user id")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid user ID",
		})
	}
	
	var userInput struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"omitempty,min=6"`
	}
	
	if err := c.Bind(&userInput); err != nil {
		log.WithError(err).WithField("error_type", "validation_error").Error("gagal parse request body")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	// Validasi input
	if err := c.Validate(&userInput); err != nil {
		log.WithError(err).WithField("error_type", "validation_error").Error("validasi input gagal")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	// Log data sensitif dengan level debug
	log.WithFields(map[string]interface{}{
		"name":  userInput.Name,
		"email": userInput.Email,
	}).Debug("data user dari request")
	
	user := domain.User{
		ID:       uint(id),
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: userInput.Password,
	}
	
	updatedUser, err := h.userUseCase.UpdateUser(user)
	if err != nil {
		log.WithError(err).WithFields(map[string]interface{}{
			"error_type": "update_error",
			"user_id":    uint(id),
			"email":      userInput.Email,
		}).Error("gagal update user")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	log.WithFields(map[string]interface{}{
		"user_id":  updatedUser.ID,
		"email":    updatedUser.Email,
		"duration": time.Since(time.Now()),
	}).Info("user berhasil diupdate")
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": updatedUser,
	})
}

// DeleteUser menangani request untuk menghapus user
func (h *UserHandler) DeleteUser(c echo.Context) error {
	idParam := c.Param("id")
	ctx := h.getRequestContext(c)
	ctx["user_id"] = idParam
	log := h.logger.WithFields(ctx)
	
	log.Info("memulai request delete user")
	
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.WithError(err).WithField("error_type", "validation_error").Error("invalid user id")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid user ID",
		})
	}
	
	if err := h.userUseCase.DeleteUser(uint(id)); err != nil {
		log.WithError(err).WithFields(map[string]interface{}{
			"error_type": "delete_error",
			"user_id":    uint(id),
		}).Error("gagal delete user")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	log.WithFields(map[string]interface{}{
		"user_id":  uint(id),
		"duration": time.Since(time.Now()),
	}).Info("user berhasil dihapus")
	
	return c.NoContent(http.StatusNoContent)
} 