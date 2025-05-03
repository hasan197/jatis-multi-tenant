package utils

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContext returns a context with timeout for testing
func TestContext(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(func() {
		cancel()
	})
	return ctx
}

// AssertNoError is a helper function to assert no error occurred
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	require.NoError(t, err)
}

// AssertError is a helper function to assert error occurred
func AssertError(t *testing.T, err error) {
	t.Helper()
	assert.Error(t, err)
}

// AssertEqual is a helper function to assert two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	assert.Equal(t, expected, actual)
}

// AssertNotNil is a helper function to assert value is not nil
func AssertNotNil(t *testing.T, value interface{}) {
	t.Helper()
	assert.NotNil(t, value)
} 