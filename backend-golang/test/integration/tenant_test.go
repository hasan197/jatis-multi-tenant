package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/streadway/amqp"

	"github.com/jatis/sample-stack-golang/internal/config"
	"github.com/jatis/sample-stack-golang/internal/modules/tenant/domain"
	"github.com/jatis/sample-stack-golang/internal/modules/tenant/repository/postgresql"
	"github.com/jatis/sample-stack-golang/internal/modules/tenant/usecase"
	"github.com/jatis/sample-stack-golang/internal/modules/tenant/delivery/messaging/rabbitmq"
	"github.com/jatis/sample-stack-golang/test/integration/setup"
)

func TestTenantIntegration(t *testing.T) {
	// Initialize test logger
	setup.InitTestLogger()
	
	// Setup test containers
	containers, connections, err := setup.SetupTestContainers()
	require.NoError(t, err)
	defer containers.Cleanup()
	defer connections.Cleanup()

	// Setup dependencies
	cfg := &config.Config{
		App: config.AppConfig{
			Workers: 3,
		},
	}

	// Create repositories and services
	tenantRepo := postgresql.NewTenantRepository(connections.DB, cfg)
	tenantManager := rabbitmq.NewTenantManager(connections.RabbitMQ, connections.DB)
	tenantUseCase := usecase.NewTenantUseCase(tenantRepo, tenantManager)

	// Test cases
	t.Run("Create and Get Tenant", func(t *testing.T) {
		ctx := context.Background()
		tenant := &domain.Tenant{
			Name:        "Test Tenant",
			Description: "Test Description",
			Status:      "active",
			Workers:     3,
		}

		// Create tenant
		err := tenantUseCase.Create(ctx, tenant)
		require.NoError(t, err)
		assert.NotEmpty(t, tenant.ID)

		// Get tenant
		retrieved, err := tenantUseCase.GetByID(ctx, tenant.ID)
		require.NoError(t, err)
		assert.Equal(t, tenant.Name, retrieved.Name)
		assert.Equal(t, tenant.Description, retrieved.Description)
		assert.Equal(t, tenant.Status, retrieved.Status)
		assert.Equal(t, tenant.Workers, retrieved.Workers)
	})

	t.Run("Update Tenant", func(t *testing.T) {
		ctx := context.Background()
		tenant := &domain.Tenant{
			Name:        "Update Test",
			Description: "Original Description",
			Status:      "active",
			Workers:     3,
		}

		// Create tenant
		err := tenantUseCase.Create(ctx, tenant)
		require.NoError(t, err)

		// Update tenant
		tenant.Description = "Updated Description"
		err = tenantUseCase.Update(ctx, tenant)
		require.NoError(t, err)

		// Verify update
		updated, err := tenantUseCase.GetByID(ctx, tenant.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Description", updated.Description)
	})

	t.Run("Delete Tenant", func(t *testing.T) {
		ctx := context.Background()
		tenant := &domain.Tenant{
			Name:        "Delete Test",
			Description: "To be deleted",
			Status:      "active",
			Workers:     3,
		}

		// Create tenant
		err := tenantUseCase.Create(ctx, tenant)
		require.NoError(t, err)

		// Delete tenant
		err = tenantUseCase.Delete(ctx, tenant.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = tenantUseCase.GetByID(ctx, tenant.ID)
		assert.Error(t, err)
	})

	t.Run("List Tenants", func(t *testing.T) {
		ctx := context.Background()
		// Create multiple tenants
		tenants := []*domain.Tenant{
			{
				Name:        "Tenant 1",
				Description: "Description 1",
				Status:      "active",
				Workers:     3,
			},
			{
				Name:        "Tenant 2",
				Description: "Description 2",
				Status:      "active",
				Workers:     3,
			},
		}

		for _, tenant := range tenants {
			err := tenantUseCase.Create(ctx, tenant)
			require.NoError(t, err)
		}

		// List tenants
		listed, err := tenantUseCase.List(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(listed), 2)
	})

	t.Run("Update Concurrency", func(t *testing.T) {
		ctx := context.Background()
		tenant := &domain.Tenant{
			Name:        "Concurrency Test",
			Description: "Testing concurrency update",
			Status:      "active",
			Workers:     3,
		}

		// Create tenant
		err := tenantUseCase.Create(ctx, tenant)
		require.NoError(t, err)

		// Update concurrency
		config := &domain.ConcurrencyConfig{
			Workers: 5,
		}
		err = tenantUseCase.UpdateConcurrency(ctx, tenant.ID, config)
		require.NoError(t, err)

		// Verify update
		updated, err := tenantUseCase.GetByID(ctx, tenant.ID)
		require.NoError(t, err)
		assert.Equal(t, 5, updated.Workers)
	})

	t.Run("Consumer Management", func(t *testing.T) {
		ctx := context.Background()
		tenant := &domain.Tenant{
			Name:        "Consumer Test",
			Description: "Testing consumer management",
			Status:      "active",
			Workers:     3,
		}

		// Create tenant
		err := tenantUseCase.Create(ctx, tenant)
		require.NoError(t, err)

		// Start consumer
		err = tenantUseCase.StartConsumer(ctx, tenant.ID)
		require.NoError(t, err)

		// Get consumer
		consumer := tenantUseCase.GetConsumer(tenant.ID)
		assert.NotNil(t, consumer)
		assert.True(t, consumer.IsActive)

		// Stop consumer
		err = tenantUseCase.StopConsumer(ctx, tenant.ID)
		require.NoError(t, err)

		// Verify consumer is stopped
		consumer = tenantUseCase.GetConsumer(tenant.ID)
		assert.Nil(t, consumer)
	})

	t.Run("Message Queue Integration", func(t *testing.T) {
		ctx := context.Background()
		tenant := &domain.Tenant{
			Name:        "Queue Test",
			Description: "Testing message queue",
			Status:      "active",
			Workers:     3,
		}

		// Create tenant
		err := tenantUseCase.Create(ctx, tenant)
		require.NoError(t, err)

		// Start consumer
		err = tenantUseCase.StartConsumer(ctx, tenant.ID)
		require.NoError(t, err)

		// Get channel
		ch, err := tenantUseCase.GetChannel()
		require.NoError(t, err)
		defer ch.Close()

		// Publish test message
		testMessage := map[string]interface{}{
			"test": "message",
			"time": time.Now(),
		}
		messageBytes, err := json.Marshal(testMessage)
		require.NoError(t, err)

		err = ch.Publish(
			"",                        // exchange
			"tenant."+tenant.ID,       // routing key
			false,                     // mandatory
			false,                     // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        messageBytes,
			},
		)
		require.NoError(t, err)

		// Give some time for message processing
		time.Sleep(1 * time.Second)

		// Stop consumer
		err = tenantUseCase.StopConsumer(ctx, tenant.ID)
		require.NoError(t, err)
	})
} 