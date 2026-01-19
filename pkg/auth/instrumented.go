package auth

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// InstrumentedVerifier wraps a Verifier with tracing.
type InstrumentedVerifier struct {
	next   Verifier
	tracer trace.Tracer
}

// NewInstrumentedVerifier creates a new instrumented verifier.
func NewInstrumentedVerifier(next Verifier) *InstrumentedVerifier {
	return &InstrumentedVerifier{
		next:   next,
		tracer: otel.Tracer("pkg/auth"),
	}
}

func (v *InstrumentedVerifier) Verify(ctx context.Context, token string) (*Claims, error) {
	ctx, span := v.tracer.Start(ctx, "auth.Verify")
	defer span.End()

	claims, err := v.next.Verify(ctx, token)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("auth.subject", claims.Subject),
		attribute.String("auth.role", claims.Role),
	)

	return claims, nil
}
