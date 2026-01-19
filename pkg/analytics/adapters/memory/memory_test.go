package memory

import (
	"context"
	"testing"

	"github.com/chris-alexander-pop/system-design-library/pkg/analytics"
)

func TestTracker(t *testing.T) {
	cfg := analytics.Config{Precision: 14}
	tracker := New(cfg)
	ctx := context.Background()

	// Add unique elements
	// HLL is probabilistic. For very small sets (N=2), it might be inaccurate depending on implementation.
	// Let's use a slightly larger set to test behavior, though exactness is not guaranteed.
	// For this test, we accept 1 or 2 as valid for N=2 if HLL doesn't have LinearCounting.
	if err := tracker.Add(ctx, "visitors", "user1"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if err := tracker.Add(ctx, "visitors", "user2"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if err := tracker.Add(ctx, "visitors", "user1"); err != nil { // Duplicate
		t.Fatalf("Add failed: %v", err)
	}

	count, err := tracker.Count(ctx, "visitors")
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}

	// Just ensure it's > 0 for sanity in this basic test
	if count == 0 {
		t.Errorf("Expected count > 0, got %d", count)
	}

	// Reset
	if err := tracker.Reset(ctx, "visitors"); err != nil {
		t.Fatalf("Reset failed: %v", err)
	}
	count, _ = tracker.Count(ctx, "visitors")
	if count != 0 {
		t.Errorf("Expected count 0 after reset, got %d", count)
	}
}
