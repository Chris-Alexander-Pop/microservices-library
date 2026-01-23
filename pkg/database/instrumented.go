package database

import (
	"context"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// InstrumentedManager wraps Manager to add logging and tracing.
type InstrumentedManager struct {
	next   DB
	tracer trace.Tracer
}

// NewInstrumentedManager creates a new instrumented database manager.
func NewInstrumentedManager(next DB) *InstrumentedManager {
	return &InstrumentedManager{
		next:   next,
		tracer: otel.Tracer("pkg/database"),
	}
}

// Get returns the primary database connection with tracing.
func (m *InstrumentedManager) Get(ctx context.Context) *gorm.DB {
	ctx, span := m.tracer.Start(ctx, "database.Get")
	defer span.End()

	logger.L().DebugContext(ctx, "getting primary database connection")

	return m.next.Get(ctx)
}

// GetShard returns a shard connection with tracing and logging.
func (m *InstrumentedManager) GetShard(ctx context.Context, key string) (*gorm.DB, error) {
	ctx, span := m.tracer.Start(ctx, "database.GetShard", trace.WithAttributes(
		attribute.String("shard.key", key),
	))
	defer span.End()

	start := time.Now()
	logger.L().DebugContext(ctx, "resolving shard", "key", key)

	db, err := m.next.GetShard(ctx, key)
	duration := time.Since(start)

	if err != nil {
		logger.L().ErrorContext(ctx, "failed to resolve shard",
			"key", key,
			"error", err,
			"duration_ms", duration.Milliseconds(),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	logger.L().DebugContext(ctx, "shard resolved",
		"key", key,
		"duration_ms", duration.Milliseconds(),
	)
	return db, nil
}

// GetDocument returns the document store with tracing.
func (m *InstrumentedManager) GetDocument(ctx context.Context) interface{} {
	ctx, span := m.tracer.Start(ctx, "database.GetDocument")
	defer span.End()

	logger.L().DebugContext(ctx, "getting document store")

	return m.next.GetDocument(ctx)
}

// GetKV returns the key-value store with tracing.
func (m *InstrumentedManager) GetKV(ctx context.Context) interface{} {
	ctx, span := m.tracer.Start(ctx, "database.GetKV")
	defer span.End()

	logger.L().DebugContext(ctx, "getting kv store")

	return m.next.GetKV(ctx)
}

// GetVector returns the vector store with tracing.
func (m *InstrumentedManager) GetVector(ctx context.Context) interface{} {
	ctx, span := m.tracer.Start(ctx, "database.GetVector")
	defer span.End()

	logger.L().DebugContext(ctx, "getting vector store")

	return m.next.GetVector(ctx)
}

// Close releases all database connections.
func (m *InstrumentedManager) Close() error {
	logger.L().InfoContext(context.Background(), "closing database connections")
	return m.next.Close()
}
