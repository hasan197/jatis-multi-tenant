package di

import (
	"database/sql"
	"sample-stack-golang/internal/config"
	"sample-stack-golang/internal/modules/user/domain"
	"sample-stack-golang/internal/modules/user/repository/postgresql"
	"sample-stack-golang/internal/modules/user/usecase"
)

// ServiceContainer adalah interface untuk mengakses semua service
type ServiceContainer interface {
	GetUserService() domain.UserUseCase
}

type serviceContainer struct {
	container *Container
}

// NewServiceContainer membuat instance baru dari ServiceContainer
func NewServiceContainer(container *Container) ServiceContainer {
	return &serviceContainer{
		container: container,
	}
}

// GetUserService mengambil UserService dari container
func (sc *serviceContainer) GetUserService() domain.UserUseCase {
	svc, ok := sc.container.Get("user_service")
	if !ok {
		return nil
	}
	return svc.(domain.UserUseCase)
}

// RegisterServices mendaftarkan semua service ke container
func RegisterServices(container *Container, db *sql.DB, cfg *config.Config) {
	// Register repositories
	userRepo := postgresql.NewUserRepository(db)
	container.Register("user_repository", userRepo)

	// Register services
	userService := usecase.NewUserUseCase(userRepo)
	container.Register("user_service", userService)
} 