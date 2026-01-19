package ratelimit

import (
	"context"
	"time"
)

// Result is the result of a limit check.
type Result struct {
	Allowed   bool
	Remaining int64
	Reset     time.Duration
}

// Limiter determines if an action is allowed.
type Limiter interface {
	// Allow checks if the key is allowed to perform 'cost' operations.
	// period is only relevant for window-based strategies.
	Allow(ctx context.Context, key string, limit int64, period time.Duration) (*Result, error)
}
