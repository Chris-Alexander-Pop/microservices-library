package analytics

import (
	"context"

	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
)

var (
	// ErrCounterNotFound is returned when operating on a non-existent counter.
	ErrCounterNotFound = errors.NotFound("counter not found", nil)
)

// Config holds configuration for the Analytics tracker.
type Config struct {
	// Precision for HyperLogLog (4-18). Default 14.
	Precision uint8 `env:"ANALYTICS_PRECISION" env-default:"14"`
}

// Tracker defines the interface for tracking unique events/elements.
type Tracker interface {
	// Add records an element for the given counter name.
	Add(ctx context.Context, counter string, element string) error

	// Count returns the estimated unique count for the given counter.
	Count(ctx context.Context, counter string) (uint64, error)

	// Reset clears a specific counter.
	Reset(ctx context.Context, counter string) error

	// Multi-counter operations could be added here
}
