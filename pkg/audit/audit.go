package audit

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/logger"
)

// Config configures the audit package.
type Config struct {
	// Enabled toggles audit logging on/off.
	Enabled bool `env:"AUDIT_ENABLED" env-default:"true"`

	// RedactConfig for handling sensitive data.
	Redact RedactorConfig
}

// Auditor defines the interface for audit logging.
type Auditor interface {
	// Log records an audit event.
	Log(ctx context.Context, event Event)

	// LogWithBuilder starts a fluent event builder.
	LogWithBuilder(ctx context.Context, eventType EventType) *EventBuilder
}

// Logger handles audit logging.
type Logger struct {
	log      *slog.Logger
	redactor *Redactor
	config   Config
}

// New creates a new audit logger.
func New(cfg Config) *Logger {
	// If Redactor is needed, it should be configured
	var r *Redactor
	if cfg.Redact.Replacement != "" || len(cfg.Redact.CustomPatterns) > 0 {
		r = NewRedactor(cfg.Redact)
	} else {
		// Default redactor
		r = NewRedactor(DefaultRedactorConfig())
	}

	return &Logger{
		log:      logger.L(),
		redactor: r,
		config:   cfg,
	}
}

// Log records an audit event.
func (l *Logger) Log(ctx context.Context, event Event) {
	if !l.config.Enabled {
		return
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	// Redact sensitive data
	if l.redactor != nil {
		event = l.redactor.RedactEvent(event)
	}

	// Convert to JSON for structured logging
	// In a real system, this might send to an async queue (Kafka/SQS)
	// For this library, we log structured to stdout
	data, _ := json.Marshal(event)

	l.log.InfoContext(ctx, "audit",
		"event", string(data),
		"event_type", string(event.EventType),
		"outcome", string(event.Outcome),
		"actor_id", event.ActorID,
		"target_id", event.TargetID,
	)
}

// LogWithBuilder provides a fluent interface for building audit events.
func (l *Logger) LogWithBuilder(ctx context.Context, eventType EventType) *EventBuilder {
	return &EventBuilder{
		logger: l,
		ctx:    ctx,
		event: Event{
			Timestamp: time.Now().UTC(),
			EventType: eventType,
			Outcome:   OutcomeSuccess,
		},
	}
}

// EventType categorizes audit events.
type EventType string

const (
	// Authentication events
	EventTypeLogin          EventType = "auth.login"
	EventTypeLogout         EventType = "auth.logout"
	EventTypeLoginFailed    EventType = "auth.login_failed"
	EventTypeMFAEnabled     EventType = "auth.mfa_enabled"
	EventTypeMFADisabled    EventType = "auth.mfa_disabled"
	EventTypePasswordChange EventType = "auth.password_change"

	// Authorization events
	EventTypeAccessGranted EventType = "authz.access_granted"
	EventTypeAccessDenied  EventType = "authz.access_denied"

	// Data events
	EventTypeDataCreate EventType = "data.create"
	EventTypeDataRead   EventType = "data.read"
	EventTypeDataUpdate EventType = "data.update"
	EventTypeDataDelete EventType = "data.delete"
	EventTypeDataExport EventType = "data.export"

	// Admin events
	EventTypeConfigChange EventType = "admin.config_change"
	EventTypeUserCreate   EventType = "admin.user_create"
	EventTypeUserModify   EventType = "admin.user_modify"
	EventTypeUserDelete   EventType = "admin.user_delete"
	EventTypeRoleChange   EventType = "admin.role_change"

	// Security events
	EventTypeSecurityAlert      EventType = "security.alert"
	EventTypeRateLimited        EventType = "security.rate_limited"
	EventTypeSuspiciousActivity EventType = "security.suspicious"
)

// Outcome indicates the result of an operation.
type Outcome string

const (
	OutcomeSuccess Outcome = "success"
	OutcomeFailure Outcome = "failure"
	OutcomeUnknown Outcome = "unknown"
)

// Event represents a structured audit event.
type Event struct {
	// Required fields
	Timestamp time.Time `json:"timestamp"`
	EventType EventType `json:"event_type"`
	Outcome   Outcome   `json:"outcome"`

	// Actor information
	ActorID        string `json:"actor_id,omitempty"`
	ActorType      string `json:"actor_type,omitempty"` // user, service, system
	ActorIP        string `json:"actor_ip,omitempty"`
	ActorUserAgent string `json:"actor_user_agent,omitempty"`

	// Target information
	TargetID   string `json:"target_id,omitempty"`
	TargetType string `json:"target_type,omitempty"`

	// Resource information
	ResourceID   string `json:"resource_id,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`

	// Operation details
	Action      string `json:"action,omitempty"`
	Description string `json:"description,omitempty"`

	// Additional context
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Request details
	RequestID     string `json:"request_id,omitempty"`
	SessionID     string `json:"session_id,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`

	// Error details (for failures)
	ErrorCode    string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// EventBuilder provides a fluent interface for building audit events.
type EventBuilder struct {
	logger *Logger
	ctx    context.Context
	event  Event
}

func (b *EventBuilder) Actor(id, actorType string) *EventBuilder {
	b.event.ActorID = id
	b.event.ActorType = actorType
	return b
}

func (b *EventBuilder) ActorIP(ip string) *EventBuilder {
	b.event.ActorIP = ip
	return b
}

func (b *EventBuilder) Target(id, targetType string) *EventBuilder {
	b.event.TargetID = id
	b.event.TargetType = targetType
	return b
}

func (b *EventBuilder) Resource(id, resourceType string) *EventBuilder {
	b.event.ResourceID = id
	b.event.ResourceType = resourceType
	return b
}

func (b *EventBuilder) Action(action string) *EventBuilder {
	b.event.Action = action
	return b
}

func (b *EventBuilder) Description(desc string) *EventBuilder {
	b.event.Description = desc
	return b
}

func (b *EventBuilder) Outcome(outcome Outcome) *EventBuilder {
	b.event.Outcome = outcome
	return b
}

func (b *EventBuilder) Error(code, message string) *EventBuilder {
	b.event.Outcome = OutcomeFailure
	b.event.ErrorCode = code
	b.event.ErrorMessage = message
	return b
}

func (b *EventBuilder) Metadata(key string, value interface{}) *EventBuilder {
	if b.event.Metadata == nil {
		b.event.Metadata = make(map[string]interface{})
	}
	b.event.Metadata[key] = value
	return b
}

func (b *EventBuilder) RequestID(id string) *EventBuilder {
	b.event.RequestID = id
	return b
}

func (b *EventBuilder) Send() {
	b.logger.Log(b.ctx, b.event)
}
