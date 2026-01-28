package metering

import (
	"context"
	"time"
)

// Meter defines the interface for recording usage events.
type Meter interface {
	// RecordUsage ingests a usage event.
	RecordUsage(ctx context.Context, event UsageEvent) error

	// GetUsage retrieves usage events matching the filter.
	GetUsage(ctx context.Context, filter UsageFilter) ([]UsageEvent, error)
}

// Rater defines the interface for calculating costs.
type Rater interface {
	// GetRate returns the price for a specific resource type.
	GetRate(ctx context.Context, resourceType string) (*RateCard, error)

	// CalculateCost estimates the cost for a given usage.
	CalculateCost(ctx context.Context, usage UsageEvent) (float64, error)
}

// UsageEvent represents a single consumption record.
type UsageEvent struct {
	ID           string            `json:"id"`
	TenantID     string            `json:"tenant_id"`
	ResourceType string            `json:"resource_type"` // e.g. "compute.instance.small", "storage.standard"
	ResourceID   string            `json:"resource_id"`
	Quantity     float64           `json:"quantity"` // e.g. hours, GB-months
	Timestamp    time.Time         `json:"timestamp"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// UsageFilter defines criteria for querying usage.
type UsageFilter struct {
	TenantID     string    `json:"tenant_id,omitempty"`
	ResourceType string    `json:"resource_type,omitempty"`
	StartTime    time.Time `json:"start_time,omitempty"`
	EndTime      time.Time `json:"end_time,omitempty"`
}

// RateCard defines the pricing for a resource.
type RateCard struct {
	ResourceType string  `json:"resource_type"`
	PricePerUnit float64 `json:"price_per_unit"`
	Currency     string  `json:"currency"` // e.g. "USD"
	Unit         string  `json:"unit"`     // e.g. "hour", "gb-month"
}

// Config holds configuration for the Metering service.
type Config struct {
	// Driver specifies the Metering backend: "memory", "prometheus", "postgres".
	Driver string `env:"METERING_DRIVER" env-default:"memory"`
}
