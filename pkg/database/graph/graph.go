package graph

import (
	"context"
)

// Vertex represents a node in the graph.
type Vertex struct {
	ID         string
	Label      string
	Properties map[string]interface{}
}

// Edge represents a relationship between two vertices.
type Edge struct {
	ID         string
	Label      string
	FromID     string
	ToID       string
	Properties map[string]interface{}
}

// Config holds configuration for the graph database.
type Config struct {
	Driver   string `env:"GRAPH_DRIVER" env-default:"memory"`
	Host     string `env:"GRAPH_HOST"`
	Port     int    `env:"GRAPH_PORT"`
	User     string `env:"GRAPH_USER"`
	Password string `env:"GRAPH_PASSWORD"`
	Database string `env:"GRAPH_DATABASE"`
}

// Interface defines the methods for a graph database.
type Interface interface {
	// AddVertex adds a vertex to the graph.
	AddVertex(ctx context.Context, v *Vertex) error

	// AddEdge adds an edge to the graph.
	AddEdge(ctx context.Context, e *Edge) error

	// GetVertex retrieves a vertex by ID.
	GetVertex(ctx context.Context, id string) (*Vertex, error)

	// GetNeighbors retrieves neighbor vertices connected by an edge with the given label.
	GetNeighbors(ctx context.Context, vertexID string, edgeLabel string, direction string) ([]*Vertex, error)

	// Query executes a query (Gremlin, Cypher, etc.).
	Query(ctx context.Context, query string, args map[string]interface{}) (interface{}, error)

	// Close releases resources.
	Close() error
}
