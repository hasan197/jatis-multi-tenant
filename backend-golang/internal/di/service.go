package di

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/streadway/amqp"
	"sample-stack-golang/internal/config"
	"sample-stack-golang/internal/modules/user/domain"
	userRepo "sample-stack-golang/internal/modules/user/repository/postgresql"
	userUsecase "sample-stack-golang/internal/modules/user/usecase"
	tenantDomain "sample-stack-golang/internal/modules/tenant/domain"
	tenantRepo "sample-stack-golang/internal/modules/tenant/repository/postgresql"
	tenantUsecase "sample-stack-golang/internal/modules/tenant/usecase"
	tenantRabbitMQ "sample-stack-golang/internal/modules/tenant/delivery/messaging/rabbitmq"
	messageRepo "sample-stack-golang/internal/modules/message/repository/postgresql"
	messageUsecase "sample-stack-golang/internal/modules/message/usecase"
)

// ServiceContainer adalah interface untuk mengakses service
type ServiceContainer interface {
	GetConfig() *config.Config
	GetDB() *pgxpool.Pool
}

type serviceContainer struct {
	config *config.Config
	pool   *pgxpool.Pool
}

// NewServiceContainer membuat instance baru dari ServiceContainer
func NewServiceContainer(cfg *config.Config, pool *pgxpool.Pool) ServiceContainer {
	return &serviceContainer{
		config: cfg,
		pool:   pool,
	}
}

// GetConfig mengambil konfigurasi dari container
func (sc *serviceContainer) GetConfig() *config.Config {
	return sc.config
}

// GetDB mengambil database connection dari container
func (sc *serviceContainer) GetDB() *pgxpool.Pool {
	return sc.pool
}

// InitializeService menginisialisasi service dan dependencies yang diperlukan
func InitializeService() (ServiceContainer, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Build database URL
	dbURL := cfg.DB.DatabaseURL()

	// Initialize database connection
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return NewServiceContainer(cfg, pool), nil
}

// Service holds all dependencies
type Service struct {
	Config        *config.Config
	Pool          *pgxpool.Pool
	Redis         *redis.Client
	RabbitMQ      *amqp.Connection
	UserUseCase   domain.UserUseCase
	TenantUseCase tenantDomain.TenantUseCase
	MessageUseCase *messageUsecase.MessageUsecase
}

// NewService creates a new service with all dependencies
func NewService(cfg *config.Config) (*Service, error) {
	// Initialize database
	pool, err := initDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	// Initialize Redis
	redis, err := initRedis(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %v", err)
	}

	// Initialize RabbitMQ
	rabbitmq, err := initRabbitMQ(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize RabbitMQ: %v", err)
	}

	// Initialize repositories
	userRepo := userRepo.NewUserRepository(pool)
	tenantRepo := tenantRepo.NewTenantRepository(pool)
	messageRepo := messageRepo.NewMessageRepository(pool)

	// Initialize RabbitMQ tenant manager
	tenantManager := tenantRabbitMQ.NewTenantManager(rabbitmq)

	// Initialize usecases
	userUseCase := userUsecase.NewUserUseCase(userRepo)
	tenantUseCase := tenantUsecase.NewTenantUseCase(tenantRepo, tenantManager)
	messageUseCase := messageUsecase.NewMessageUsecase(messageRepo)

	// Start tenant manager
	if err := tenantManager.Start(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to start tenant manager: %v", err)
	}

	// Start consumers for existing tenants
	tenants, err := tenantRepo.List(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %v", err)
	}

	for _, tenant := range tenants {
		if err := tenantManager.StartConsumer(context.Background(), tenant.ID); err != nil {
			fmt.Printf("Warning: failed to start consumer for tenant %s: %v\n", tenant.ID, err)
		}
	}

	return &Service{
		Config:        cfg,
		Pool:          pool,
		Redis:         redis,
		RabbitMQ:      rabbitmq,
		UserUseCase:   userUseCase,
		TenantUseCase: tenantUseCase,
		MessageUseCase: messageUseCase,
	}, nil
}

// Close closes all connections
func (s *Service) Close() error {
	// Close database
	if s.Pool != nil {
		s.Pool.Close()
	}

	// Close Redis
	if err := s.Redis.Close(); err != nil {
		return fmt.Errorf("failed to close Redis: %v", err)
	}

	// Close RabbitMQ
	if err := s.RabbitMQ.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ: %v", err)
	}

	return nil
}

// initDB initializes database connection
func initDB(cfg *config.Config) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.DB.DatabaseURL())
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}

// initRedis initializes Redis connection
func initRedis(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

// initRabbitMQ initializes RabbitMQ connection
func initRabbitMQ(cfg *config.Config) (*amqp.Connection, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
} 