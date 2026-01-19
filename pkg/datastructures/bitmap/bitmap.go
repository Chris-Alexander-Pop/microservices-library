package bitmap

import (
	"math/bits"
	"sync"
)

// Bitmap implements a dense bitset.
type Bitmap struct {
	data []uint64
	size uint64
	mu   sync.RWMutex
}

// New creates a new Bitmap with the given size.
func New(size uint64) *Bitmap {
	words := (size + 63) / 64
	return &Bitmap{
		data: make([]uint64, words),
		size: size,
	}
}

// Set sets the bit at index i to 1.
func (b *Bitmap) Set(i uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if i >= b.size {
		return
	}
	b.data[i/64] |= 1 << (i % 64)
}

// Clear sets the bit at index i to 0.
func (b *Bitmap) Clear(i uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if i >= b.size {
		return
	}
	b.data[i/64] &^= 1 << (i % 64)
}

// Get returns the bit value at index i.
func (b *Bitmap) Get(i uint64) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if i >= b.size {
		return false
	}
	return (b.data[i/64] & (1 << (i % 64))) != 0
}

// OnesCount returns the number of set bits.
func (b *Bitmap) OnesCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	count := 0
	for _, w := range b.data {
		count += bits.OnesCount64(w)
	}
	return count
}
