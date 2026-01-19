package audit

import (
	"context"
	"testing"
)

func TestAuditLog(t *testing.T) {
	cfg := Config{
		Enabled: true,
		Redact:  DefaultRedactorConfig(),
	}
	logger := New(cfg)
	ctx := context.Background()

	// Test direct Log
	event := Event{
		EventType: EventTypeLogin,
		Outcome:   OutcomeSuccess,
		ActorID:   "test-user",
	}
	logger.Log(ctx, event)

	// Test Builder
	logger.LogWithBuilder(ctx, EventTypeDataRead).
		Actor("test-user", "user").
		Resource("doc-123", "document").
		Outcome(OutcomeSuccess).
		Send()
}

func TestRedactionInAudit(t *testing.T) {
	cfg := Config{Enabled: true}
	logger := New(cfg) // Uses default redactor

	ctx := context.Background()

	// Should redact credit card
	logger.LogWithBuilder(ctx, EventTypeDataCreate).
		Description("Payment with card 1234-5678-9012-3456").
		Send()
}
