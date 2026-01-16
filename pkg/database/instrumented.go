package database

import (
	"context"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/logger"
	"gorm.io/gorm"
)

// InstrumentedManager wraps Manager to add logging for connection acquisition
// Note: GORM has its own logger for SQL queries. This is high-level manager instrumentation.
type InstrumentedManager struct {
	next DB
}

func NewInstrumentedManager(next DB) *InstrumentedManager {
	return &InstrumentedManager{next: next}
}

func (m *InstrumentedManager) Get(ctx context.Context) *gorm.DB {
	// Not logging every Get() as it's too noisy, but we could trace it if needed.
	// For now, we assume GORM plugins handle the Query Tracing.
	return m.next.Get(ctx)
}

func (m *InstrumentedManager) GetShard(ctx context.Context, key string) (*gorm.DB, error) {
	start := time.Now()
	logger.L().DebugContext(ctx, "resolving shard", "key", key)

	db, err := m.next.GetShard(ctx, key)
	duration := time.Since(start)

	if err != nil {
		logger.L().ErrorContext(ctx, "failed to resolve shard", "key", key, "error", err, "duration", duration)
		return nil, err
	}

	// We might want to log which shard was selected if Manager exposed it,
	// but GetShard just returns *DB.

	return db, nil
}

func (m *InstrumentedManager) Close() error {
	logger.L().Info("closing database connections")
	return m.next.Close()
}
