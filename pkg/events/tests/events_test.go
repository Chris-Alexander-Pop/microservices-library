package tests

import (
	"context"
	"testing"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/events"
	"github.com/chris-alexander-pop/system-design-library/pkg/events/adapters/memory"
	"github.com/chris-alexander-pop/system-design-library/pkg/test"
)

type EventsTestSuite struct {
	test.Suite
	bus events.Bus
}

func (s *EventsTestSuite) SetupTest() {
	s.Suite.SetupTest()
	s.bus = memory.New()
}

func (s *EventsTestSuite) TestPubSub() {
	received := make(chan events.Event, 1)
	topic := "user.created"

	err := s.bus.Subscribe(s.Ctx, topic, func(ctx context.Context, e events.Event) error {
		received <- e
		return nil
	})
	s.NoError(err)

	evt := events.Event{
		Type:    topic,
		Payload: map[string]string{"user_id": "123"},
	}
	err = s.bus.Publish(s.Ctx, topic, evt)
	s.NoError(err)

	select {
	case e := <-received:
		s.Equal(topic, e.Type)
	case <-time.After(1 * time.Second):
		s.Fail("Timed out waiting for event")
	}
}

func TestEventsSuite(t *testing.T) {
	test.Run(t, new(EventsTestSuite))
}
