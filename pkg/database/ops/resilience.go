package ops

import (
	"context"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
)

// WithRetry executes the operation with exponential backoff retries.
// Useful for transient network errors or db connection glitches.
func WithRetry(ctx context.Context, attempts int, backoff time.Duration, op func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				backoff *= 2 // Exponential backoff
			}
		}

		err = op()
		if err == nil {
			return nil
		}

		// If context is canceled, stop retrying
		if ctx.Err() != nil {
			return ctx.Err()
		}
		// In a real system, we would check if err is "Retryable" (e.g. 5xx, timeout, lock contention)
		// For now, we assume simplistic retry for any error unless explicitly wrapped as Permanent.
	}
	return errors.Wrap(err, "max retries exceeded")
}
