package di

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"sample-stack-golang/internal/config"
	"sample-stack-golang/internal/modules/user/domain"
	"sample-stack-golang/internal/modules/user/repository/postgresql"
	"sample-stack-golang/internal/modules/user/usecase"
	tenantDomain "sample-stack-golang/internal/modules/tenant/domain"
	tenantRepo "sample-stack-golang/internal/modules/tenant/repository/postgresql"
	tenantUsecase "sample-stack-golang/internal/modules/tenant/usecase"
)

// ServiceContainer adalah interface untuk mengakses semua service
type ServiceContainer interface {
	GetUserService() domain.UserUseCase
	GetTenantService() tenantDomain.TenantRepository
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

// GetTenantService mengambil TenantService dari container
func (sc *serviceContainer) GetTenantService() tenantDomain.TenantRepository {
	svc, ok := sc.container.Get("tenant_service")
	if !ok {
		return nil
	}
	return svc.(tenantDomain.TenantRepository)
}

// RegisterServices mendaftarkan semua service ke container
func RegisterServices(container *Container, pool *pgxpool.Pool, cfg *config.Config) {
	// Register repositories
	userRepo := postgresql.NewUserRepository(pool)
	container.Register("user_repository", userRepo)

	tenantRepo := tenantRepo.NewTenantRepository(pool)
	container.Register("tenant_repository", tenantRepo)

	// Register services
	userService := usecase.NewUserUseCase(userRepo)
	container.Register("user_service", userService)

	tenantService := tenantUsecase.NewTenantUsecase(tenantRepo)
	container.Register("tenant_service", tenantService)
} 