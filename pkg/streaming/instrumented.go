package streaming

import (
	"context"

	"github.com/chris-alexander-pop/system-design-library/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// InstrumentedClient wraps a Client with logging and tracing.
type InstrumentedClient struct {
	next   Client
	tracer trace.Tracer
}

// NewInstrumentedClient creates a new InstrumentedClient.
func NewInstrumentedClient(next Client) *InstrumentedClient {
	return &InstrumentedClient{
		next:   next,
		tracer: otel.Tracer("pkg/streaming"),
	}
}

func (c *InstrumentedClient) PutRecord(ctx context.Context, streamName string, partitionKey string, data []byte) error {
	ctx, span := c.tracer.Start(ctx, "streaming.PutRecord", trace.WithAttributes(
		attribute.String("stream.name", streamName),
		attribute.String("partition.key", partitionKey),
		attribute.Int("data.size", len(data)),
	))
	defer span.End()

	logger.L().InfoContext(ctx, "putting record to stream", "stream", streamName, "partition_key", partitionKey)

	err := c.next.PutRecord(ctx, streamName, partitionKey, data)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.L().ErrorContext(ctx, "failed to put record", "stream", streamName, "error", err)
	}
	return err
}

func (c *InstrumentedClient) Close() error {
	return c.next.Close()
}
