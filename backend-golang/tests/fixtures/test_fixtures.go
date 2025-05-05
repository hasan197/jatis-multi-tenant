package fixtures

import (
	"time"
	"sample-stack-golang/internal/config"
)

// TestUser represents a test user fixture
type TestUser struct {
	ID        string
	Username  string
	Email     string
	CreatedAt time.Time
}

// GetTestUsers returns a slice of test users
func GetTestUsers() []TestUser {
	return []TestUser{
		{
			ID:        "user-1",
			Username:  "testuser1",
			Email:     "test1@example.com",
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        "user-2",
			Username:  "testuser2",
			Email:     "test2@example.com",
			CreatedAt: time.Now().Add(-48 * time.Hour),
		},
	}
}

// TestConfig returns a test configuration
func TestConfig() *config.Config {
	return &config.Config{
		DB: config.DBConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "test",
			Password: "test",
			Name:     "testdb",
		},
		Redis: config.RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
		RabbitMQ: config.RabbitMQConfig{
			Host:     "localhost",
			Port:     5672,
			User:     "guest",
			Password: "guest",
		},
		Server: config.ServerConfig{
			Port: 8080,
		},
		App: config.AppConfig{
			Version: "test",
			Env:     "test",
		},
	}
} 