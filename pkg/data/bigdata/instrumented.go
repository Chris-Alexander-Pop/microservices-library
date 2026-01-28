package bigdata

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// InstrumentedClient wraps a Client with tracing.
type InstrumentedClient struct {
	next   Client
	tracer trace.Tracer
}

// NewInstrumentedClient creates a new instrumented bigdata client.
func NewInstrumentedClient(next Client) *InstrumentedClient {
	return &InstrumentedClient{
		next:   next,
		tracer: otel.Tracer("pkg/bigdata"),
	}
}

func (c *InstrumentedClient) Query(ctx context.Context, query string, args ...interface{}) (*Result, error) {
	ctx, span := c.tracer.Start(ctx, "bigdata.Query", trace.WithAttributes(
		attribute.String("db.statement", query),
	))
	defer span.End()

	res, err := c.next.Query(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Optional: Log result count
	span.SetAttributes(attribute.Int("db.rows_returned", len(res.Rows)))

	return res, nil
}

func (c *InstrumentedClient) Close() error {
	return c.next.Close()
}
