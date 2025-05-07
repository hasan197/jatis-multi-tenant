package setup

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/streadway/amqp"
)

type TestContainer struct {
	Pool     *dockertest.Pool
	Resource *dockertest.Resource
}

type TestContainers struct {
	Postgres *TestContainer
	Redis    *TestContainer
	RabbitMQ *TestContainer
}

type Connections struct {
	DB       *pgxpool.Pool
	Redis    *redis.Client
	RabbitMQ *amqp.Connection
}

func SetupTestContainers() (*TestContainers, *Connections, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct pool: %v", err)
	}

	containers := &TestContainers{}
	connections := &Connections{}

	// Setup PostgreSQL
	pgContainer, pgPool, err := setupPostgres(pool)
	if err != nil {
		return nil, nil, err
	}
	containers.Postgres = pgContainer
	connections.DB = pgPool

	// Setup Redis
	redisContainer, redisClient, err := setupRedis(pool)
	if err != nil {
		return nil, nil, err
	}
	containers.Redis = redisContainer
	connections.Redis = redisClient

	// Setup RabbitMQ
	rabbitContainer, rabbitConn, err := setupRabbitMQ(pool)
	if err != nil {
		return nil, nil, err
	}
	containers.RabbitMQ = rabbitContainer
	connections.RabbitMQ = rabbitConn

	return containers, connections, nil
}

func setupPostgres(pool *dockertest.Pool) (*TestContainer, *pgxpool.Pool, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=testdb",
		},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start postgres: %v", err)
	}

	var db *pgxpool.Pool
	if err := pool.Retry(func() error {
		var err error
		db, err = pgxpool.New(context.Background(), fmt.Sprintf("postgresql://postgres:postgres@localhost:%s/testdb?sslmode=disable", resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping(context.Background())
	}); err != nil {
		return nil, nil, fmt.Errorf("could not connect to postgres: %v", err)
	}

	// Initialize schema
	schema := `
	CREATE TABLE IF NOT EXISTS tenants (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		status VARCHAR(50),
		workers INTEGER DEFAULT 3,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS messages (
		id UUID NOT NULL,
		tenant_id UUID NOT NULL,
		payload JSONB NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (tenant_id, id)
	) PARTITION BY LIST (tenant_id);
	
	CREATE INDEX IF NOT EXISTS idx_messages_tenant_id ON messages(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);

	CREATE OR REPLACE FUNCTION create_messages_partition(tenant_uuid UUID)
	RETURNS void AS $$
	DECLARE
		partition_name TEXT;
	BEGIN
		partition_name := 'messages_' || replace(tenant_uuid::text, '-', '_');
		
		EXECUTE format('CREATE TABLE IF NOT EXISTS %I PARTITION OF messages FOR VALUES IN (%L)',
			partition_name,
			tenant_uuid);
	END;
	$$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION drop_messages_partition(tenant_uuid UUID)
	RETURNS void AS $$
	DECLARE
		partition_name TEXT;
	BEGIN
		partition_name := 'messages_' || replace(tenant_uuid::text, '-', '_');
		
		EXECUTE format('DROP TABLE IF EXISTS %I',
			partition_name);
	END;
	$$ LANGUAGE plpgsql;
	`

	if _, err := db.Exec(context.Background(), schema); err != nil {
		return nil, nil, fmt.Errorf("could not initialize schema: %v", err)
	}

	return &TestContainer{Pool: pool, Resource: resource}, db, nil
}

func setupRedis(pool *dockertest.Pool) (*TestContainer, *redis.Client, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "6",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start redis: %v", err)
	}

	var client *redis.Client
	if err := pool.Retry(func() error {
		client = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
		})
		return client.Ping(context.Background()).Err()
	}); err != nil {
		return nil, nil, fmt.Errorf("could not connect to redis: %v", err)
	}

	return &TestContainer{Pool: pool, Resource: resource}, client, nil
}

func setupRabbitMQ(pool *dockertest.Pool) (*TestContainer, *amqp.Connection, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "rabbitmq",
		Tag:        "3-management",
		Env: []string{
			"RABBITMQ_DEFAULT_USER=guest",
			"RABBITMQ_DEFAULT_PASS=guest",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start rabbitmq: %v", err)
	}

	// Wait for RabbitMQ to be ready
	time.Sleep(10 * time.Second)

	var conn *amqp.Connection
	if err := pool.Retry(func() error {
		var err error
		conn, err = amqp.Dial(fmt.Sprintf("amqp://guest:guest@localhost:%s/", resource.GetPort("5672/tcp")))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, nil, fmt.Errorf("could not connect to rabbitmq: %v", err)
	}

	return &TestContainer{Pool: pool, Resource: resource}, conn, nil
}

func (tc *TestContainers) Cleanup() {
	if tc.Postgres != nil {
		if err := tc.Postgres.Pool.Purge(tc.Postgres.Resource); err != nil {
			log.Printf("Could not purge postgres: %v", err)
		}
	}
	if tc.Redis != nil {
		if err := tc.Redis.Pool.Purge(tc.Redis.Resource); err != nil {
			log.Printf("Could not purge redis: %v", err)
		}
	}
	if tc.RabbitMQ != nil {
		if err := tc.RabbitMQ.Pool.Purge(tc.RabbitMQ.Resource); err != nil {
			log.Printf("Could not purge rabbitmq: %v", err)
		}
	}
}

func (c *Connections) Cleanup() {
	if c.DB != nil {
		c.DB.Close()
	}
	if c.Redis != nil {
		if err := c.Redis.Close(); err != nil {
			log.Printf("Could not close redis connection: %v", err)
		}
	}
	if c.RabbitMQ != nil {
		if err := c.RabbitMQ.Close(); err != nil {
			log.Printf("Could not close rabbitmq connection: %v", err)
		}
	}
} 