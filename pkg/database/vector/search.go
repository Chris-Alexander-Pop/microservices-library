package vector

import (
	"container/heap"
	"context"
	"runtime"

	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// SearchFunc defines the signature for searching a single shard
type SearchFunc func(ctx context.Context, shardID string, vector []float32, limit int) ([]Result, error)

// ResultHeap implements heap.Interface for a Min-Heap of Results based on Score.
// We use a Min-Heap to maintain the Top-K HIGHEST scores.
// The root of the heap will be the item with the LOWEST score among the top K.
// If we find an item with score > root, we pop root and push new item.
type ResultHeap []Result

func (h ResultHeap) Len() int           { return len(h) }
func (h ResultHeap) Less(i, j int) bool { return h[i].Score < h[j].Score } // Min-Heap based on Score
func (h ResultHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *ResultHeap) Push(x interface{}) {
	*h = append(*h, x.(Result))
}

func (h *ResultHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// ScatterGatherSearch executes search across all shards concurrently
// Implements bounded concurrency and Heap-based aggregation.
func ScatterGatherSearch(
	ctx context.Context,
	vector []float32,
	limit int,
	shardIDs []string,
	searchFn SearchFunc,
) ([]Result, error) {
	// 1. Concurrency Control
	// Limit concurrency to NumCPU * 2 to prevent explosion
	maxConcurrency := int64(runtime.NumCPU() * 2)
	sem := semaphore.NewWeighted(maxConcurrency)

	resultsChan := make(chan []Result, len(shardIDs))
	g, ctx := errgroup.WithContext(ctx)

	// 2. Scatter
	for _, id := range shardIDs {
		shardID := id
		if err := sem.Acquire(ctx, 1); err != nil {
			return nil, err
		}

		g.Go(func() error {
			defer sem.Release(1)

			// Fail-safe
			defer func() {
				if r := recover(); r != nil {
					// In real app, log metrics.Counter("panic.search", 1)
				}
			}()

			res, err := searchFn(ctx, shardID, vector, limit)
			if err != nil {
				return errors.Wrap(err, "search failed on shard "+shardID)
			}
			resultsChan <- res
			return nil
		})
	}

	// 3. Wait
	if err := g.Wait(); err != nil {
		return nil, err
	}
	close(resultsChan)

	// 4. Gather with Min-Heap (Maintain Top-K)
	h := &ResultHeap{}
	heap.Init(h)

	for shardResults := range resultsChan {
		for _, res := range shardResults {
			if h.Len() < limit {
				heap.Push(h, res)
			} else {
				// If new score is better than the worst of our top-k
				if res.Score > (*h)[0].Score {
					heap.Pop(h)
					heap.Push(h, res)
				}
			}
		}
	}

	// 5. Extract and Sort (Heap logic leaves them unordered mostly)
	// Pop all gives us smallest to largest of the top K.
	finalResults := make([]Result, h.Len())
	for i := h.Len() - 1; i >= 0; i-- {
		finalResults[i] = heap.Pop(h).(Result)
	}
	// ResultHeap is Min-Heap.
	// Pop element 1: Smallest.
	// Pop element 2: Second smallest.
	// ...
	// Pop element K: Largest.
	// We filled `finalResults` from back to front (i=Len-1 down to 0).
	// So:
	// finalResults[Len-1] = Smallest
	// finalResults[0] = Largest
	// Thus finalResults is sorted Descending (Highest score first). Correct.

	return finalResults, nil
}
