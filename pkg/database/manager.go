package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/database/sharding"
	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
	"github.com/gocql/gocql"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// ConnectionFactory is a function that creates a connection from config
// Use interface{} return type to support *gorm.DB, *mongo.Database, etc.
type ConnectionFactory func(cfg Config) (interface{}, error)

// Manager implements the DB interface and manages multiple connections (SQL, NoSQL, etc.)
type Manager struct {
	primary  interface{}
	shards   map[string]interface{}
	strategy sharding.Strategy
	config   ManagerConfig
	mu       sync.RWMutex
}

// NewManager creates a new database manager
func NewManager(cfg ManagerConfig, factory ConnectionFactory) (*Manager, error) {
	if factory == nil {
		return nil, errors.New(errors.CodeInternal, "connection factory is required", nil)
	}

	m := &Manager{
		config: cfg,
		shards: make(map[string]interface{}),
	}

	// Initialize Primary
	if cfg.Primary.Driver != "" {
		conn, err := factory(cfg.Primary)
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to primary db")
		}
		m.primary = conn
	}

	// Initialize Shards
	for id, shardCfg := range cfg.Shards {
		conn, err := factory(shardCfg)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to connect to shard %s", id))
		}
		m.shards[id] = conn
	}

	// Initialize Strategy
	if cfg.ShardingStrategy == "consistent_hash" {
		var shardIDs []string
		for id := range m.shards {
			shardIDs = append(shardIDs, id)
		}
		m.strategy = sharding.NewConsistentHash(50, shardIDs)
	}

	return m, nil
}

// Get returns the primary database connection (SQL)
func (m *Manager) Get(ctx context.Context) *gorm.DB {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if db, ok := m.primary.(*gorm.DB); ok {
		return db.WithContext(ctx)
	}
	return nil
}

// GetShard returns SQL shard
func (m *Manager) GetShard(ctx context.Context, key string) (*gorm.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Internal helper to get generic shard
	getGeneric := func(k string) (interface{}, error) {
		if conn, ok := m.shards[k]; ok {
			return conn, nil
		}
		if m.strategy != nil {
			shardID := m.strategy.GetShard(k)
			if conn, ok := m.shards[shardID]; ok {
				return conn, nil
			}
		}
		return nil, errors.New(errors.CodeNotFound, fmt.Sprintf("shard not found for key: %s", k), nil)
	}

	conn, err := getGeneric(key)
	if err != nil {
		return nil, err
	}

	if db, ok := conn.(*gorm.DB); ok {
		return db.WithContext(ctx), nil
	}
	return nil, errors.New(errors.CodeInternal, "shard is not a SQL connection", nil)
}

// GetDocument returns the primary document store (Mongo, Dynamo)
func (m *Manager) GetDocument(ctx context.Context) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.primary
}

// GetKV returns the primary KV store (Redis)
func (m *Manager) GetKV(ctx context.Context) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.primary
}

// GetVector returns the primary vector store (Pinecone)
func (m *Manager) GetVector(ctx context.Context) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.primary
}

func (m *Manager) Close() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Helper to close specific types
	closeConn := func(c interface{}) error {
		if c == nil {
			return nil
		}

		switch v := c.(type) {
		case *gorm.DB:
			sqlDB, err := v.DB()
			if err != nil {
				return err
			}
			return sqlDB.Close()

		case *mongo.Database:
			// Mongo Disconnect is on Client
			if client := v.Client(); client != nil {
				// Timeout context for disconnect
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				return client.Disconnect(ctx)
			}
			return nil

		case *gocql.Session:
			v.Close()
			return nil

		case *redis.Client:
			return v.Close()

		default:
			// Check if it implements io.Closer
			if closer, ok := c.(interface{ Close() error }); ok {
				return closer.Close()
			}
			return nil
		}
		return nil
	}

	var firstErr error
	// Close primary
	if err := closeConn(m.primary); err != nil && firstErr == nil {
		firstErr = err
	}
	// Close shards
	for _, conn := range m.shards {
		if err := closeConn(conn); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
