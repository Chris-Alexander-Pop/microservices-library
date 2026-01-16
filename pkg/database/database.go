package database

import (
	"context"

	"gorm.io/gorm"
)

// Driver enum
const (
	DriverPostgres = "postgres"
	DriverMySQL    = "mysql"
	DriverSQLite   = "sqlite"
)

// Config holds configuration for a single database connection
type Config struct {
	Driver   string `env:"DB_DRIVER" env-default:"postgres"`
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-default:"postgres"`
	Name     string `env:"DB_NAME" env-default:"postgres"`
	SSLMode  string `env:"DB_SSLMODE" env-default:"disable"`
}

// ManagerConfig holds configuration for the entire database layer, including shards
type ManagerConfig struct {
	Primary Config
	Shards  map[string]Config // Keyed by shard ID (e.g., "shard-1", "shard-2")
}

// DB defines the interface for interacting with the database layer
type DB interface {
	// Get returns the primary database connection
	Get(ctx context.Context) *gorm.DB

	// GetShard returns a specific shard connection based on a routing key
	// If the key maps to the primary (or no sharding config), it returns primary.
	GetShard(ctx context.Context, key string) (*gorm.DB, error)

	// Close closes all connections
	Close() error
}
