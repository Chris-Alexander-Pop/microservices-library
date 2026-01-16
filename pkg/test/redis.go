package test

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/redis"
)

// StartRedis spins up a Redis container for testing
func StartRedis(t *testing.T) (string, func()) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	redisContainer, err := redis.Run(ctx,
		"redis:7-alpine",
		redis.WithSnapshotting(0, 0),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	if err != nil {
		t.Fatalf("failed to start redis container: %v", err)
	}

	connStr, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	cleanup := func() {
		if err := redisContainer.Terminate(context.Background()); err != nil {
			t.Logf("failed to terminate redis container: %v", err)
		}
	}

	return connStr, cleanup
}
