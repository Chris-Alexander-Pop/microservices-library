package disruptor

import (
	"runtime"
	"sync/atomic"
)

const (
	padding = 64 // Cache line padding to prevent false sharing
)

// RingBuffer is a high-performance lock-free ring Buffer.
type RingBuffer[T any] struct {
	Buffer []T
	Mask   uint64
	Cursor uint64 // Writer Cursor
	_      [padding]byte
	Gate   uint64 // Reader gating Cursor
}

func New[T any](size uint64) *RingBuffer[T] {
	// Size must be power of 2
	if size == 0 || (size&(size-1)) != 0 {
		size = 1024 // Default
	}

	return &RingBuffer[T]{
		Buffer: make([]T, size),
		Mask:   size - 1,
	}
}

// Publish reserves a slot and runs the callback to fill it.
func (rb *RingBuffer[T]) Publish(producer func(*T)) {
	// Simple Single-Producer logic.
	// For Multi-Producer, need CAS on Cursor.

	current := atomic.LoadUint64(&rb.Cursor)
	next := current + 1

	// Check wrap around against slow readers (Gate)
	// wrapPoint := next - size.
	// if wrapPoint > Gate, spin.
	// Only necessary if preventing overwrite.
	// Disruptor assumes gating.

	idx := next & rb.Mask
	producer(&rb.Buffer[idx])

	atomic.StoreUint64(&rb.Cursor, next)
}

// Consume waits for availability and consumes.
func (rb *RingBuffer[T]) Consume(consumer func(T)) {
	// Simple Single-Consumer logic.

	target := atomic.LoadUint64(&rb.Gate) + 1

	// Spin wait for Cursor to reach target
	for target > atomic.LoadUint64(&rb.Cursor) {
		runtime.Gosched()
	}

	// Read
	idx := target & rb.Mask
	consumer(rb.Buffer[idx])

	atomic.StoreUint64(&rb.Gate, target)
}

// Note: This is an extremely simplified schematic of the LMAX Disruptor.
// Real disruptor deals with Barriers, multiple dependencies, and memory barriers extensively.
