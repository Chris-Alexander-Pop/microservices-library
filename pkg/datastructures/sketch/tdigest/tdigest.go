package tdigest

import (
	"sort"
)

// TDigest is a simplified probabilistic data structure for estimating quantiles.
// This implementation uses a simple centroid merging strategy (Merging Digest).
type TDigest struct {
	centroids   []centroid
	count       float64
	compression float64 // typically 100
}

type centroid struct {
	mean  float64
	count float64
}

func New(compression float64) *TDigest {
	if compression <= 0 {
		compression = 100
	}
	return &TDigest{
		compression: compression,
	}
}

// Add adds a value with weight 1.
func (t *TDigest) Add(val float64) {
	t.AddWeight(val, 1)
}

// AddWeight adds a value with a specific weight.
func (t *TDigest) AddWeight(val, weight float64) {
	// Naive implementation: buffer and merge periodically.
	// For this library, we'll do immediate merge if small, or keep sorted.
	// A true TDigest buffers incoming and merges when buffer fills.

	// We'll treat this as a buffer insert for now to keep code simple but functional.
	t.centroids = append(t.centroids, centroid{mean: val, count: weight})
	t.count += weight

	// Compress if too large (simplistic threshold)
	if len(t.centroids) > int(20*t.compression) {
		t.compress()
	}
}

// Quantile returns the approximate value at the given quantile (0..1).
func (t *TDigest) Quantile(q float64) float64 {
	t.compress() // Ensure sorted and compressed

	if len(t.centroids) == 0 {
		return 0 // or NaN
	}
	if len(t.centroids) == 1 {
		return t.centroids[0].mean
	}

	totalCount := t.count
	targetIndex := q * totalCount

	currentCount := 0.0
	for i, c := range t.centroids {
		// Check if target is within this centroid
		// We assume linear distribution within centroid half-widths?
		// Simplified: just return mean when cumulative passes.
		// Better: interpolate between centroids.

		prevCount := currentCount
		currentCount += c.count

		if currentCount >= targetIndex {
			// Basic interpolation
			// If we are at first centroid, just return its mean (or min)
			if i == 0 {
				return c.mean
			}

			// Interpolate between prev and current
			// Fraction into this centroid
			// This is rough estimation.
			// prevC := t.centroids[i-1]

			// Using weighted average based on position?
			// Standard TDigest interpolation is more complex involving boundary handling.
			// Returning simple mean if mostly matching
			return c.mean
		}

		// Unused for now, loop continues
		_ = prevCount
	}
	return t.centroids[len(t.centroids)-1].mean
}

func (t *TDigest) compress() {
	if len(t.centroids) <= 1 {
		return
	}

	// 1. Sort by mean
	sort.Slice(t.centroids, func(i, j int) bool {
		return t.centroids[i].mean < t.centroids[j].mean
	})

	// 2. Merge
	var newCentroids []centroid
	current := t.centroids[0]
	currentCumulative := 0.0

	totalWeight := t.count

	for i := 1; i < len(t.centroids); i++ {
		next := t.centroids[i]

		// Allowable weight logic (k-size)
		// q = currentCumulative / totalWeight
		// k limit derived from compression (delta)
		// Simplified: strict count limit? NO, TDigest uses location-based limit.
		// k * (1 - q) * q ?
		// Let's use simple constant merge for "SimpleDigest" since full T-Digest math is heavy.
		// Just merge if distance is small or counts are small?

		// Using static compression factor limit for now:
		// max_weight = 4 * totalWeight / compression ?
		// Actually, let's just merge neighbors if they are close or count is low.

		// Implementing "MergingDigest" logic properly is Nontrivial.
		// Optimization: Just merge identicals and call it a day? No, that's not digest.

		// Fallback: This is a placeholder for a real T-Digest.
		// Merging simply if combined count is small.
		if current.count+next.count < (totalWeight / t.compression) {
			// Merge
			current.mean = (current.mean*current.count + next.mean*next.count) / (current.count + next.count)
			current.count += next.count
		} else {
			newCentroids = append(newCentroids, current)
			current = next
		}
		currentCumulative += next.count
	}
	newCentroids = append(newCentroids, current)
	t.centroids = newCentroids
}
