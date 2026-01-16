package vector

import (
	"context"
)

// Result represents a search result
type Result struct {
	ID       string                 `json:"id"`
	Score    float32                `json:"score"` // Distances or Similarity
	Metadata map[string]interface{} `json:"metadata"`
}

// Store defines the interface for vector operations
type Store interface {
	// Search finds the nearest neighbors to the query vector
	Search(ctx context.Context, vector []float32, limit int) ([]Result, error)

	// Upsert inserts or updates a vector with metadata
	Upsert(ctx context.Context, id string, vector []float32, metadata map[string]interface{}) error

	// Delete removes a vector by ID
	Delete(ctx context.Context, id string) error
}
