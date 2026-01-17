package vectorclock

import (
	"bytes"
	"encoding/gob"
	"sync"
)

// Clock represents a Vector Clock.
type Clock map[string]uint64

// Ordering represents the relationship between two clocks.
type Ordering int

const (
	Equal Ordering = iota
	Before
	After
	Concurrent
)

// New creates a new Vector Clock.
func New() Clock {
	return make(Clock)
}

// Increment increments the logical clock for the given node.
func (vc Clock) Increment(node string) {
	vc[node]++
}

// Merge merges another clock into this one.
func (vc Clock) Merge(other Clock) {
	for node, time := range other {
		if vc[node] < time {
			vc[node] = time
		}
	}
}

// Copy returns a deep copy of the clock.
func (vc Clock) Copy() Clock {
	newVC := make(Clock, len(vc))
	for k, v := range vc {
		newVC[k] = v
	}
	return newVC
}

// Compare returns the ordering of this clock relative to another.
func (vc Clock) Compare(other Clock) Ordering {
	var isBefore, isAfter bool

	// Check all keys in vc
	for node, time := range vc {
		otherTime, exists := other[node]
		if !exists {
			// If missing in other, vc has newer info (implying After)
			// Unless other has keys vc doesn't have.
			isAfter = true
		} else {
			if time > otherTime {
				isAfter = true
			} else if time < otherTime {
				isBefore = true
			}
		}
	}

	// Check keys in other that are missing in vc
	for node := range other {
		if _, exists := vc[node]; !exists {
			isBefore = true
		}
	}

	if isBefore && isAfter {
		return Concurrent
	}
	if isBefore {
		return Before
	}
	if isAfter {
		return After
	}
	return Equal
}

// Bytes returns a serialized byte slice.
func (vc Clock) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(vc); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// FromBytes deserializes bytes into a Clock.
func FromBytes(data []byte) (Clock, error) {
	var vc Clock
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&vc); err != nil {
		return nil, err
	}
	return vc, nil
}

// ThreadSafeClock wraps a Clock with a mutex.
type ThreadSafeClock struct {
	clock Clock
	mu    sync.RWMutex
}

func NewThreadSafe() *ThreadSafeClock {
	return &ThreadSafeClock{
		clock: New(),
	}
}

func (ts *ThreadSafeClock) Increment(node string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.clock.Increment(node)
}

func (ts *ThreadSafeClock) Merge(other Clock) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.clock.Merge(other)
}

func (ts *ThreadSafeClock) Compare(other Clock) Ordering {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.clock.Compare(other)
}

func (ts *ThreadSafeClock) Snapshot() Clock {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.clock.Copy()
}
