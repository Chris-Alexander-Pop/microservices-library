package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/chris-alexander-pop/system-design-library/pkg/database/sharding"
	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
	"gorm.io/gorm"
)

// ConnectionFactory is a function that creates a gorm.DB from config
type ConnectionFactory func(cfg Config) (*gorm.DB, error)

// Manager implements the DB interface and manages multiple connections
type Manager struct {
	primary  *gorm.DB
	shards   map[string]*gorm.DB
	strategy sharding.Strategy
	mu       sync.RWMutex
}

// NewManager creates a new database manager
// It accepts a factory to allow dependency injection of the actual dialer (e.g., Postgres adapter)
func NewManager(cfg ManagerConfig, factory ConnectionFactory) (*Manager, error) {
	if factory == nil {
		return nil, errors.New(errors.CodeInternal, "connection factory is required", nil)
	}

	primary, err := factory(cfg.Primary)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to primary db")
	}

	m := &Manager{
		primary: primary,
		shards:  make(map[string]*gorm.DB),
	}

	for id, shardCfg := range cfg.Shards {
		shardDB, err := factory(shardCfg)
		if err != nil {
			// Close primary if shard fails? Or just log?
			// Strict startup: fail.
			return nil, errors.Wrap(err, fmt.Sprintf("failed to connect to shard %s", id))
		}
		m.shards[id] = shardDB
	}

	// Initialize strategy if configured
	if cfg.ShardingStrategy == "consistent_hash" {
		var shardIDs []string
		for id := range m.shards {
			shardIDs = append(shardIDs, id)
		}
		m.strategy = sharding.NewConsistentHash(50, shardIDs) // 50 replicas by default
	}

	return m, nil
}

func (m *Manager) Get(ctx context.Context) *gorm.DB {
	return m.primary.WithContext(ctx)
}

func (m *Manager) GetShard(ctx context.Context, key string) (*gorm.DB, error) {
	// Simple routing strategy: In a real app, this might use consistent hashing on `key`
	// against the list of shard IDs.
	// For this abstraction, we assume `key` IS the shard ID for manual routing,
	// or we can implement a specific ShardingStrategy interface later.

	// Let's assume for now that if the key exists in the map, use it.
	// If connection is not found, return error? Or fallback to primary?
	// Let's implement explicit lookup.

	m.mu.RLock()
	defer m.mu.RUnlock()

	if db, ok := m.shards[key]; ok {
		return db.WithContext(ctx), nil
	}

	// 2. Strategy Routing (if enabled)
	if m.strategy != nil {
		shardID := m.strategy.GetShard(key)
		if db, ok := m.shards[shardID]; ok {
			return db.WithContext(ctx), nil
		}
	}

	// If no shard found, error.
	return nil, errors.New(errors.CodeNotFound, fmt.Sprintf("shard not found for key: %s", key), nil)
}

func (m *Manager) Close() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Helper to close sql.DB
	closeDB := func(g *gorm.DB) error {
		sqlDB, err := g.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	if err := closeDB(m.primary); err != nil {
		return err
	}

	for _, db := range m.shards {
		if err := closeDB(db); err != nil {
			return err
		}
	}
	return nil
}
