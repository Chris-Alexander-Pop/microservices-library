package blob

import (
	"context"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/events"
)

// EventedStore decorates a BlobStore to emit events
type EventedStore struct {
	next BlobStore
	bus  events.Bus
}

func NewEventedStore(next BlobStore, bus events.Bus) *EventedStore {
	return &EventedStore{next: next, bus: bus}
}

func (s *EventedStore) Upload(ctx context.Context, key string, data []byte) error {
	err := s.next.Upload(ctx, key, data)
	if err == nil {
		// Emit Event
		_ = s.bus.Publish(ctx, "blob.uploaded", events.Event{
			ID:        key, // Use blob key as event ID or generate one
			Type:      "blob.uploaded",
			Source:    "pkg/blob",
			Timestamp: time.Now(),
			Payload: map[string]interface{}{
				"key":  key,
				"size": len(data),
			},
		})
	}
	return err
}

func (s *EventedStore) Download(ctx context.Context, key string) ([]byte, error) {
	return s.next.Download(ctx, key)
}

func (s *EventedStore) Delete(ctx context.Context, key string) error {
	err := s.next.Delete(ctx, key)
	if err == nil {
		_ = s.bus.Publish(ctx, "blob.deleted", events.Event{
			ID:        key,
			Type:      "blob.deleted",
			Source:    "pkg/blob",
			Timestamp: time.Now(),
			Payload: map[string]interface{}{
				"key": key,
			},
		})
	}
	return err
}
