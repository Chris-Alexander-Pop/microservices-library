package countmin

import (
	"hash/fnv"
	"math"
)

// Sketch is a Count-Min Sketch for frequency estimation.
// It provides approximate counts with probabilistic bounds.
type Sketch struct {
	width  uint
	depth  uint
	count  uint64
	table  [][]uint64
	hasher func(data []byte) uint64
}

// New creates a new Count-Min Sketch.
// epsilon: acceptable error rate (e.g., 0.01)
// delta: confidence (e.g., 0.99 means 1 - 0.01)
func New(epsilon, delta float64) *Sketch {
	width := uint(math.Ceil(math.E / epsilon))
	depth := uint(math.Ceil(math.Log(1 / (1 - delta))))

	table := make([][]uint64, depth)
	for i := range table {
		table[i] = make([]uint64, width)
	}

	return &Sketch{
		width: width,
		depth: depth,
		table: table,
	}
}

// Add adds an item to the sketch.
func (cms *Sketch) Add(data []byte) {
	cms.AddCount(data, 1)
}

// AddCount adds an item with a specific count.
func (cms *Sketch) AddCount(data []byte, count uint64) {
	cms.count += count

	// We need 'depth' pair-wise independent hash functions.
	// For simplicity, we use double hashing to generate 'depth' indices.
	// h(i) = (a + b*i) % width
	h1, h2 := hash(data)

	for i := uint(0); i < cms.depth; i++ {
		idx := (h1 + uint64(i)*h2) % uint64(cms.width)
		cms.table[i][idx] += count
	}
}

// Estimate returns the estimated count for an item.
// It is a point query estimate.
func (cms *Sketch) Estimate(data []byte) uint64 {
	h1, h2 := hash(data)
	minVal := uint64(math.MaxUint64)

	for i := uint(0); i < cms.depth; i++ {
		idx := (h1 + uint64(i)*h2) % uint64(cms.width)
		val := cms.table[i][idx]
		if val < minVal {
			minVal = val
		}
	}
	return minVal
}

// Count returns the total number of items added.
func (cms *Sketch) Count() uint64 {
	return cms.count
}

// Merge merges another sketch into this one.
func (cms *Sketch) Merge(other *Sketch) {
	if cms.width != other.width || cms.depth != other.depth {
		return // Cannot merge incompatible sketches
	}

	cms.count += other.count
	for i := uint(0); i < cms.depth; i++ {
		for j := uint(0); j < cms.width; j++ {
			cms.table[i][j] += other.table[i][j]
		}
	}
}

// Reset clears the sketch.
func (cms *Sketch) Reset() {
	cms.count = 0
	for i := range cms.table {
		for j := range cms.table[i] {
			cms.table[i][j] = 0
		}
	}
}

func hash(data []byte) (uint64, uint64) {
	h := fnv.New64a()
	h.Write(data)
	v1 := h.Sum64()

	h.Write([]byte{0}) // slight variation
	v2 := h.Sum64()
	return v1, v2
}
