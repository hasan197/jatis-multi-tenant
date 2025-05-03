package fixtures

import (
	"time"
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

// TestConfig represents test configuration
type TestConfig struct {
	DatabaseURL string
	APIPort     int
	Environment string
}

// GetTestConfig returns test configuration
func GetTestConfig() TestConfig {
	return TestConfig{
		DatabaseURL: "postgres://test:test@localhost:5432/testdb",
		APIPort:     8080,
		Environment: "test",
	}
} 