package events

import (
	"context"

	"github.com/chris-alexander-pop/system-design-library/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// InstrumentedBus wraps a Bus to add logging and tracing
type InstrumentedBus struct {
	next   Bus
	tracer trace.Tracer
}

func NewInstrumentedBus(next Bus) *InstrumentedBus {
	return &InstrumentedBus{
		next:   next,
		tracer: otel.Tracer("pkg/events"),
	}
}

func (b *InstrumentedBus) Publish(ctx context.Context, topic string, event Event) error {
	ctx, span := b.tracer.Start(ctx, "events.Publish", trace.WithAttributes(
		attribute.String("event.topic", topic),
		attribute.String("event.type", event.Type),
		attribute.String("event.id", event.ID),
	))
	defer span.End()

	logger.L().InfoContext(ctx, "publishing event", "topic", topic, "type", event.Type, "id", event.ID)

	err := b.next.Publish(ctx, topic, event)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.L().ErrorContext(ctx, "failed to publish event", "topic", topic, "error", err)
		return err
	}
	return nil
}

func (b *InstrumentedBus) Subscribe(ctx context.Context, topic string, handler Handler) error {
	ctx, span := b.tracer.Start(ctx, "events.Subscribe", trace.WithAttributes(
		attribute.String("event.topic", topic),
	))
	defer span.End()

	logger.L().InfoContext(ctx, "subscribing to topic", "topic", topic)

	// Wrap the handler to trace processing
	instrumentedHandler := func(ctx context.Context, event Event) error {
		ctx, span := b.tracer.Start(ctx, "events.Handle", trace.WithAttributes(
			attribute.String("event.topic", topic),
			attribute.String("event.type", event.Type),
			attribute.String("event.id", event.ID),
		))
		defer span.End()

		logger.L().InfoContext(ctx, "processing event", "topic", topic, "type", event.Type)

		err := handler(ctx, event)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			logger.L().ErrorContext(ctx, "failed to process event", "topic", topic, "error", err)
			return err
		}
		return nil
	}

	return b.next.Subscribe(ctx, topic, instrumentedHandler)
}

func (b *InstrumentedBus) Close() error {
	return b.next.Close()
}
