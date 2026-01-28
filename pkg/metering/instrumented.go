package metering

import (
	"context"

	"github.com/chris-alexander-pop/system-design-library/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// InstrumentedMeter wraps a Meter with logging and tracing.
type InstrumentedMeter struct {
	next   Meter
	tracer trace.Tracer
}

// NewInstrumentedMeter creates a new instrumented meter.
func NewInstrumentedMeter(next Meter) *InstrumentedMeter {
	return &InstrumentedMeter{
		next:   next,
		tracer: otel.Tracer("pkg/metering"),
	}
}

func (m *InstrumentedMeter) RecordUsage(ctx context.Context, event UsageEvent) error {
	ctx, span := m.tracer.Start(ctx, "metering.RecordUsage", trace.WithAttributes(
		attribute.String("tenant.id", event.TenantID),
		attribute.String("resource.type", event.ResourceType),
		attribute.Float64("quantity", event.Quantity),
	))
	defer span.End()

	logger.L().InfoContext(ctx, "recording usage", "tenant_id", event.TenantID, "resource", event.ResourceType, "quantity", event.Quantity)

	err := m.next.RecordUsage(ctx, event)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.L().ErrorContext(ctx, "failed to record usage", "error", err)
		return err
	}

	return nil
}

func (m *InstrumentedMeter) GetUsage(ctx context.Context, filter UsageFilter) ([]UsageEvent, error) {
	ctx, span := m.tracer.Start(ctx, "metering.GetUsage")
	defer span.End()

	events, err := m.next.GetUsage(ctx, filter)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.Int("result.count", len(events)))
	return events, nil
}
