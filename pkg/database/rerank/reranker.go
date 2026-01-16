package rerank

import (
	"context"

	"github.com/chris-alexander-pop/system-design-library/pkg/database/vector"
)

// Reranker defines the interface for reranking search results
type Reranker interface {
	Rerank(ctx context.Context, queryVector []float32, candidates []vector.Result) ([]vector.Result, error)
}

// SimpleScorer implements Reranker by just applying a static multiplier or logic.
// In reality, this would call an LLM or Cross-Encoder model.
// TBI: Implement actual LLM call (e.g. OpenAI/Cohere API or local ONNX model)
type SimpleScorer struct {
	Multiplier float32
}

func (s *SimpleScorer) Rerank(ctx context.Context, queryVector []float32, candidates []vector.Result) ([]vector.Result, error) {
	// Simulate "overengineered" reranking logic
	reranked := make([]vector.Result, len(candidates))
	copy(reranked, candidates)

	for i := range reranked {
		// Boost score arbitrarily for demo
		reranked[i].Score *= s.Multiplier

		if reranked[i].Metadata == nil {
			reranked[i].Metadata = make(map[string]interface{})
		}
		reranked[i].Metadata["reranked"] = true
	}

	return reranked, nil
}
