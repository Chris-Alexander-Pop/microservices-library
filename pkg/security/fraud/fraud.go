package fraud

import (
	"context"
	"time"
)

// Config configures the fraud detection system.
type Config struct {
	// Provider specifies which fraud provider to use (memory, maxmind, etc.).
	Provider string `env:"SECURITY_FRAUD_PROVIDER" env-default:"memory"`
}

// Evaluation represents the result of a fraud check.
type Evaluation struct {
	RiskScore float64           `json:"risk_score"` // 0.0 to 1.0 (1.0 = high risk)
	Action    string            `json:"action"`     // allow, block, review
	Reasons   []string          `json:"reasons"`
	Metadata  map[string]string `json:"metadata"`
	CheckID   string            `json:"check_id"`
	Timestamp time.Time         `json:"timestamp"`
}

// UserEvent represents an action taken by a user that needs evaluation.
type UserEvent struct {
	UserID    string            `json:"user_id"`
	IPAddress string            `json:"ip_address"`
	UserAgent string            `json:"user_agent"`
	Action    string            `json:"action"` // login, purchase, signup
	Amount    float64           `json:"amount,omitempty"`
	Currency  string            `json:"currency,omitempty"`
	Metadata  map[string]string `json:"metadata"`
}

// Detector defines the interface for fraud detection.
type Detector interface {
	Score(ctx context.Context, event UserEvent) (*Evaluation, error)
}
