package resilience

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker_StateTransitions(t *testing.T) {
	cfg := CircuitBreakerConfig{
		Name:             "test-cb",
		FailureThreshold: 2,
		SuccessThreshold: 1,
		Timeout:          100 * time.Millisecond,
	}
	cb := NewCircuitBreaker(cfg)

	ctx := context.Background()
	failErr := errors.New("failure")

	// 1. Initial State: Closed
	if cb.State() != StateClosed {
		t.Errorf("Expected state Closed, got %v", cb.State())
	}

	// 2. Failure 1: Still Closed
	cb.Execute(ctx, func(ctx context.Context) error { return failErr })
	if cb.State() != StateClosed {
		t.Errorf("Expected state Closed, got %v", cb.State())
	}

	// 3. Failure 2: Trip to Open
	cb.Execute(ctx, func(ctx context.Context) error { return failErr })
	if cb.State() != StateOpen {
		t.Errorf("Expected state Open, got %v", cb.State())
	}

	// 4. Request while Open: Fast Fail
	err := cb.Execute(ctx, func(ctx context.Context) error { return nil })
	if !errors.Is(err, ErrCircuitOpen) {
		t.Errorf("Expected ErrCircuitOpen, got %v", err)
	}

	// 5. Wait for Timeout -> Half-Open
	time.Sleep(150 * time.Millisecond)

	// The state transition happens lazily on the next Execute or allowRequest
	// We'll mimic a successful call
	err = cb.Execute(ctx, func(ctx context.Context) error { return nil })
	if err != nil {
		t.Errorf("Expected success in Half-Open, got %v", err)
	}

	// 6. Success -> Closed
	if cb.State() != StateClosed {
		t.Errorf("Expected state Closed after success, got %v", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenFailure(t *testing.T) {
	cfg := CircuitBreakerConfig{
		Name:             "test-cb-fail",
		FailureThreshold: 1,
		SuccessThreshold: 2,
		Timeout:          50 * time.Millisecond,
	}
	cb := NewCircuitBreaker(cfg)
	ctx := context.Background()
	fail := errors.New("fail")

	// Trip to Open
	cb.Execute(ctx, func(ctx context.Context) error { return fail })
	if cb.State() != StateOpen {
		t.Fatalf("Failed to open circuit")
	}

	// Wait for Timeout (Ready for Half-Open)
	time.Sleep(100 * time.Millisecond)

	// Execute failure in Half-Open -> Trip back to Open immediately
	cb.Execute(ctx, func(ctx context.Context) error { return fail })

	if cb.State() != StateOpen {
		t.Errorf("Expected state Open after half-open failure, got %v", cb.State())
	}
}

func TestCircuitBreaker_Reset(t *testing.T) {
	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig("reset-test"))
	cb.setState(StateOpen)
	cb.Reset()
	if cb.State() != StateClosed {
		t.Error("Reset failed to close circuit")
	}
}
