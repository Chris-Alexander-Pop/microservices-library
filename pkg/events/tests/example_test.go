package events_test

import (
	"context"
	"fmt"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/events"
	"github.com/chris-alexander-pop/system-design-library/pkg/events/adapters/memory"
)

func Example() {
	// Create an in-memory event bus
	bus := memory.New()
	defer bus.Close()

	ctx := context.Background()

	// Subscribe to events
	handler := func(ctx context.Context, event events.Event) error {
		fmt.Printf("Received: %s\n", event.Type)
		return nil
	}

	_ = bus.Subscribe(ctx, "users", handler)

	// Publish an event
	event := events.Event{
		ID:        "evt-123",
		Type:      "user.created",
		Source:    "user-service",
		Timestamp: time.Now(),
		Payload:   map[string]string{"user_id": "123", "email": "alice@example.com"},
	}

	_ = bus.Publish(ctx, "users", event)

	// Give the handler time to process
	time.Sleep(10 * time.Millisecond)
	// Output: Received: user.created
}

func ExampleEvent() {
	// Create a well-structured event
	event := events.Event{
		ID:        "evt-456",
		Type:      "order.placed",
		Source:    "order-service",
		Timestamp: time.Now().UTC(),
		Payload: map[string]interface{}{
			"order_id": "ord-789",
			"total":    99.99,
			"items":    3,
		},
	}

	fmt.Println(event.Type)
	// Output: order.placed
}
