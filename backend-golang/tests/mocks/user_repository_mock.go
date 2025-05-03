package mocks

import (
	"context"
	"sync"

	"github.com/stretchr/testify/mock"
)

// UserRepositoryMock is a mock implementation of UserRepository
type UserRepositoryMock struct {
	mock.Mock
	mu sync.RWMutex
}

// GetUser mocks the GetUser method
func (m *UserRepositoryMock) GetUser(ctx context.Context, id string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, id)
	return args.Get(0), args.Error(1)
}

// CreateUser mocks the CreateUser method
func (m *UserRepositoryMock) CreateUser(ctx context.Context, user interface{}) (interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, user)
	return args.Get(0), args.Error(1)
}

// UpdateUser mocks the UpdateUser method
func (m *UserRepositoryMock) UpdateUser(ctx context.Context, id string, user interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id, user)
	return args.Error(0)
}

// DeleteUser mocks the DeleteUser method
func (m *UserRepositoryMock) DeleteUser(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id)
	return args.Error(0)
} 