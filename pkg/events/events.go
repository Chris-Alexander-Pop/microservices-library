package events

import (
	"context"

	"github.com/nats-io/nats.go"
)

type EventStore interface {
	Publish(ctx context.Context, topic string, data []byte) error
	Subscribe(ctx context.Context, topic string, handler func([]byte) error) error
	Close()
}

type NatsEventStore struct {
	conn *nats.Conn
}

func NewNats(url string) (*NatsEventStore, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{conn: nc}, nil
}

func (n *NatsEventStore) Publish(ctx context.Context, topic string, data []byte) error {
	return n.conn.Publish(topic, data)
}

func (n *NatsEventStore) Subscribe(ctx context.Context, topic string, handler func([]byte) error) error {
	_, err := n.conn.Subscribe(topic, func(msg *nats.Msg) {
		_ = handler(msg.Data)
	})
	return err
}

func (n *NatsEventStore) Close() {
	n.conn.Close()
}
