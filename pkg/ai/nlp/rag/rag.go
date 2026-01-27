// Package rag provides a simple RAG (Retrieval Augmented Generation) orchestrator.
package rag

import (
	"context"
	"strings"

	"github.com/chris-alexander-pop/system-design-library/pkg/ai/nlp/embedding"
	"github.com/chris-alexander-pop/system-design-library/pkg/database/vector"
)

// Orchestrator manages the RAG pipeline.
type Orchestrator struct {
	embedder    embedding.Service
	vectorStore vector.Store
}

// New creates a new RAG orchestrator.
func New(embedder embedding.Service, store vector.Store) *Orchestrator {
	return &Orchestrator{
		embedder:    embedder,
		vectorStore: store,
	}
}

// Ingest adds document text to the knowledge base.
func (o *Orchestrator) Ingest(ctx context.Context, id, text string, metadata map[string]interface{}) error {
	// Simple chunking
	chunks := strings.Split(text, "\n\n")

	for _, chunk := range chunks {
		if strings.TrimSpace(chunk) == "" {
			continue
		}

		vectors, err := o.embedder.Embed(ctx, []string{chunk})
		if err != nil {
			return err
		}

		if metadata == nil {
			metadata = make(map[string]interface{})
		}
		metadata["text"] = chunk

		// In a real system, generate unique ID for chunk
		chunkID := id // simplified, overwrites if multiple chunks

		err = o.vectorStore.Upsert(ctx, chunkID, vectors[0], metadata)
		if err != nil {
			return err
		}
	}

	return nil
}

// Retrieve finds relevant context for a query.
func (o *Orchestrator) Retrieve(ctx context.Context, query string, k int) ([]string, error) {
	vectors, err := o.embedder.Embed(ctx, []string{query})
	if err != nil {
		return nil, err
	}

	results, err := o.vectorStore.Search(ctx, vectors[0], k)
	if err != nil {
		return nil, err
	}

	contexts := make([]string, len(results))
	for i, res := range results {
		if txt, ok := res.Metadata["text"].(string); ok {
			contexts[i] = txt
		}
	}

	return contexts, nil
}
