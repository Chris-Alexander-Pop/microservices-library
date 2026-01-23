package vector

import (
	"context"
)

// Config holds configuration for a vector database.
type Config struct {
	// Driver specifies the vector backend: "pinecone", "memory".
	Driver string `env:"VECTOR_DRIVER" env-default:"pinecone"`

	// Host is the vector service endpoint.
	Host string `env:"VECTOR_HOST"`

	// APIKey is the authentication key.
	APIKey string `env:"VECTOR_API_KEY"`

	// Environment is the vector service environment (e.g., "us-west1-gcp").
	Environment string `env:"VECTOR_ENVIRONMENT"`

	// ProjectID is the project identifier.
	ProjectID string `env:"VECTOR_PROJECT_ID"`

	// IndexName is the name of the vector index.
	IndexName string `env:"VECTOR_INDEX_NAME"`

	// Dimension is the size of the vectors (e.g., 1536 for OpenAI embeddings).
	Dimension int `env:"VECTOR_DIMENSION" env-default:"1536"`
}

// Result represents a search result.
type Result struct {
	ID       string                 `json:"id"`
	Score    float32                `json:"score"` // Distance or Similarity
	Metadata map[string]interface{} `json:"metadata"`
}

// Store defines the interface for vector operations.
type Store interface {
	// Search finds the nearest neighbors to the query vector.
	Search(ctx context.Context, vector []float32, limit int) ([]Result, error)

	// Upsert inserts or updates a vector with metadata.
	Upsert(ctx context.Context, id string, vector []float32, metadata map[string]interface{}) error

	// Delete removes a vector by ID.
	Delete(ctx context.Context, id string) error

	// Close releases resources.
	Close() error
}
