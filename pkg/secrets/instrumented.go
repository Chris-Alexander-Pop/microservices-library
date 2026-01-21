package secrets

import (
	"context"

	"github.com/chris-alexander-pop/system-design-library/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// InstrumentedManager wraps a Manager with logging and tracing.
type InstrumentedManager struct {
	next   Manager
	tracer trace.Tracer
}

// NewInstrumentedManager creates a new InstrumentedManager.
func NewInstrumentedManager(next Manager) *InstrumentedManager {
	return &InstrumentedManager{
		next:   next,
		tracer: otel.Tracer("pkg/secrets"),
	}
}

func (m *InstrumentedManager) GetSecret(ctx context.Context, key string) (string, error) {
	ctx, span := m.tracer.Start(ctx, "secrets.GetSecret", trace.WithAttributes(
		attribute.String("secret.key", key),
	))
	defer span.End()

	logger.L().InfoContext(ctx, "getting secret", "key", key)

	val, err := m.next.GetSecret(ctx, key)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.L().ErrorContext(ctx, "failed to get secret", "key", key, "error", err)
	}
	return val, err
}

func (m *InstrumentedManager) SetSecret(ctx context.Context, key string, value string) error {
	ctx, span := m.tracer.Start(ctx, "secrets.SetSecret", trace.WithAttributes(
		attribute.String("secret.key", key),
	))
	defer span.End()

	logger.L().InfoContext(ctx, "setting secret", "key", key)

	err := m.next.SetSecret(ctx, key, value)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.L().ErrorContext(ctx, "failed to set secret", "key", key, "error", err)
	}
	return err
}

func (m *InstrumentedManager) DeleteSecret(ctx context.Context, key string) error {
	ctx, span := m.tracer.Start(ctx, "secrets.DeleteSecret", trace.WithAttributes(
		attribute.String("secret.key", key),
	))
	defer span.End()

	logger.L().InfoContext(ctx, "deleting secret", "key", key)

	err := m.next.DeleteSecret(ctx, key)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.L().ErrorContext(ctx, "failed to delete secret", "key", key, "error", err)
	}
	return err
}

func (m *InstrumentedManager) Close() error {
	return m.next.Close()
}
