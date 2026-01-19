package audit

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// InstrumentedAuditor wraps an Auditor with logging and tracing.
type InstrumentedAuditor struct {
	next   Auditor
	tracer trace.Tracer
}

// NewInstrumentedAuditor creates a new instrumented auditor.
func NewInstrumentedAuditor(next Auditor) *InstrumentedAuditor {
	return &InstrumentedAuditor{
		next:   next,
		tracer: otel.Tracer("pkg/audit"),
	}
}

func (a *InstrumentedAuditor) Log(ctx context.Context, event Event) {
	ctx, span := a.tracer.Start(ctx, "audit.Log", trace.WithAttributes(
		attribute.String("event.type", string(event.EventType)),
		attribute.String("event.outcome", string(event.Outcome)),
	))
	defer span.End()

	a.next.Log(ctx, event)
}

func (a *InstrumentedAuditor) LogWithBuilder(ctx context.Context, eventType EventType) *EventBuilder {
	// We delegate directly because the builder eventually calls Log(), which we've instrumented above.
	// However, since LogWithBuilder returns a concrete *EventBuilder tied to the *Logger (likely),
	// and *Logger implements Auditor, extending this wrapper to wrap the builder is complex
	// without returning an interface for the builder.
	// For now, we assume this is mostly a passthrough or simply logs start of build.
	return a.next.LogWithBuilder(ctx, eventType)
}
