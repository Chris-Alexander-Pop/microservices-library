package memory

import (
	"context"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/concurrency"
	"github.com/chris-alexander-pop/system-design-library/pkg/events"
)

type Event = events.Event

type MemoryBus struct {
	handlers map[string][]events.Handler
	mu       *concurrency.SmartRWMutex
}

func New() *MemoryBus {
	return &MemoryBus{
		handlers: make(map[string][]events.Handler),
		mu:       concurrency.NewSmartRWMutex(concurrency.MutexConfig{Name: "MemoryBus"}),
	}
}

func (m *MemoryBus) Publish(ctx context.Context, topic string, event Event) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	handlers, ok := m.handlers[topic]
	if !ok {
		return nil // No subscribers
	}

	// In memory, we execute synchronous-ish (or spawn goroutines)
	// For simplicity and "async-like" behavior, launch goroutines but wait group could be better for graceful shutdown
	// For this mock, we'll just run them.
	for _, h := range handlers {
		go func(handler events.Handler, evt events.Event) {
			_ = handler(context.Background(), evt) // Detached context for async
		}(h, event)
	}

	return nil
}

func (m *MemoryBus) Subscribe(ctx context.Context, topic string, handler events.Handler) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.handlers[topic] = append(m.handlers[topic], handler)
	return nil
}

func (m *MemoryBus) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers = make(map[string][]events.Handler)
	return nil
}
