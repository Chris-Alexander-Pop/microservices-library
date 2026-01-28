package audit_test

import (
	"context"
	"github.com/chris-alexander-pop/system-design-library/pkg/audit"
	"testing"
)

func TestAuditLog(t *testing.T) {
	cfg := audit.Config{
		Enabled: true,
		Redact:  audit.DefaultRedactorConfig(),
	}
	logger := audit.New(cfg)
	ctx := context.Background()

	// Test direct Log
	event := audit.Event{
		EventType: audit.EventTypeLogin,
		Outcome:   audit.OutcomeSuccess,
		ActorID:   "test-user",
	}
	logger.Log(ctx, event)

	// Test Builder
	logger.LogWithBuilder(ctx, audit.EventTypeDataRead).
		Actor("test-user", "user").
		Resource("doc-123", "document").
		Outcome(audit.OutcomeSuccess).
		Send()
}

func TestRedactionInAudit(t *testing.T) {
	cfg := audit.Config{Enabled: true}
	logger := audit.New(cfg) // Uses default redactor

	ctx := context.Background()

	// Should redact credit card
	logger.LogWithBuilder(ctx, audit.EventTypeDataCreate).
		Description("Payment with card 1234-5678-9012-3456").
		Send()
}
