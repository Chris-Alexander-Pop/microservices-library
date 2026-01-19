package resilience

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetry_Success(t *testing.T) {
	calls := 0
	err := Retry(context.Background(), DefaultRetryConfig(), func(ctx context.Context) error {
		calls++
		if calls < 3 {
			return errors.New("temp fail")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected success, got %v", err)
	}
	if calls != 3 {
		t.Errorf("Expected 3 calls, got %d", calls)
	}
}

func TestRetry_MaxAttempts(t *testing.T) {
	cfg := DefaultRetryConfig()
	cfg.MaxAttempts = 3
	cfg.InitialBackoff = 1 * time.Millisecond // Fast test

	calls := 0
	failErr := errors.New("steady fail")

	err := Retry(context.Background(), cfg, func(ctx context.Context) error {
		calls++
		return failErr
	})

	if err != failErr {
		t.Errorf("Expected failErr, got %v", err)
	}
	if calls != 3 {
		t.Errorf("Expected 3 calls, got %d", calls)
	}
}

func TestRetry_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := DefaultRetryConfig()
	cfg.InitialBackoff = 100 * time.Millisecond

	// Cancel immediately
	cancel()

	err := Retry(ctx, cfg, func(ctx context.Context) error {
		return errors.New("should act on context")
	})

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected ContextCanceled, got %v", err)
	}
}
