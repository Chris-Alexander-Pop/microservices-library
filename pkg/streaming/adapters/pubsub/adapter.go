package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub/v2"
)

// Adapter provides streaming interface for GCP Pub/Sub.
type Adapter struct {
	client    *pubsub.Client
	projectID string
}

// New creates a new Pub/Sub streaming adapter.
func New(ctx context.Context, projectID string) (*Adapter, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &Adapter{client: client, projectID: projectID}, nil
}

// topicFullName returns the fully qualified topic name.
func (a *Adapter) topicFullName(topicName string) string {
	return fmt.Sprintf("projects/%s/topics/%s", a.projectID, topicName)
}

// PutRecord publishes a record to the specified topic.
func (a *Adapter) PutRecord(ctx context.Context, topicName, partitionKey string, data []byte) error {
	// Get publisher for the topic
	publisher := a.client.Publisher(a.topicFullName(topicName))
	defer publisher.Stop()

	// Enable ordering if partition key is provided
	if partitionKey != "" {
		publisher.EnableMessageOrdering = true
	}

	// Publish message with ordering key
	res := publisher.Publish(ctx, &pubsub.Message{
		Data:        data,
		OrderingKey: partitionKey,
	})
	_, err := res.Get(ctx)
	return err
}

// Close closes the Pub/Sub client.
func (a *Adapter) Close() error {
	return a.client.Close()
}
