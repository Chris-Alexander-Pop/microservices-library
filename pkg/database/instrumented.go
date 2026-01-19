package database

import (
	"context"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// InstrumentedManager wraps Manager to add logging and tracing
type InstrumentedManager struct {
	next   DB
	tracer trace.Tracer
}

func NewInstrumentedManager(next DB) *InstrumentedManager {
	return &InstrumentedManager{
		next:   next,
		tracer: otel.Tracer("pkg/database"),
	}
}

func (m *InstrumentedManager) Get(ctx context.Context) *gorm.DB {
	_, span := m.tracer.Start(ctx, "database.Get")
	defer span.End()
	return m.next.Get(ctx)
}

func (m *InstrumentedManager) GetShard(ctx context.Context, key string) (*gorm.DB, error) {
	ctx, span := m.tracer.Start(ctx, "database.GetShard")
	defer span.End()

	start := time.Now()
	// logger.L().DebugContext(ctx, "resolving shard", "key", key)

	db, err := m.next.GetShard(ctx, key)
	duration := time.Since(start)

	if err != nil {
		logger.L().ErrorContext(ctx, "failed to resolve shard", "key", key, "error", err, "duration", duration)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return db, nil
}

func (m *InstrumentedManager) GetDocument(ctx context.Context) interface{} {
	_, span := m.tracer.Start(ctx, "database.GetDocument")
	defer span.End()
	return m.next.GetDocument(ctx)
}

func (m *InstrumentedManager) GetKV(ctx context.Context) interface{} {
	_, span := m.tracer.Start(ctx, "database.GetKV")
	defer span.End()
	return m.next.GetKV(ctx)
}

func (m *InstrumentedManager) GetVector(ctx context.Context) interface{} {
	_, span := m.tracer.Start(ctx, "database.GetVector")
	defer span.End()
	return m.next.GetVector(ctx)
}

func (m *InstrumentedManager) Close() error {
	logger.L().Info("closing database connections")
	return m.next.Close()
}
