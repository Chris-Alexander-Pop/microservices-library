/*
Package audit provides structured audit logging for compliance and security.

This package includes:
  - Structured audit events (SIEM-ready)
  - Event types for common operations
  - PII redaction utilities

Usage:

	import "github.com/chris-alexander-pop/system-design-library/pkg/audit"

	cfg := audit.Config{Enabled: true}
	auditor := audit.New(cfg)

	auditor.LogWithBuilder(ctx, audit.EventTypeLogin).
		Actor("user-123", "user").
		Outcome(audit.OutcomeSuccess).
		Send()
*/
package audit
