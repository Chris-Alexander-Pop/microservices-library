package adaptive

import (
	"math"
	"sync"
	"time"
)

// Limiter implements an adaptive concurrency limiter based on TCP Vegas-style latency analysis.
// It estimates the queuing delay and adjusts the concurrency limit to maintain optimal throughput and low latency.
type Limiter struct {
	limit    float64
	minLimit float64
	maxLimit float64

	// Metrics
	inflight int
	mu       sync.Mutex

	// Sample window
	minRTT time.Duration
}

func New(minLimit, maxLimit float64) *Limiter {
	return &Limiter{
		limit:    minLimit,
		minLimit: minLimit,
		maxLimit: maxLimit,
		minRTT:   math.MaxInt64,
	}
}

// Acquire tries to acquire a concurrency token.
// Returns true if allowed, false if rejected (shed load).
func (l *Limiter) Acquire() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if float64(l.inflight) >= l.limit {
		return false
	}
	l.inflight++
	return true
}

// Release releases the token and updates the limiter with sample RTT.
func (l *Limiter) Release(rtt time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.inflight--
	if l.inflight < 0 {
		l.inflight = 0
	}

	// Update MinRTT (over long horizon, reset occasionally in real impl)
	if rtt < l.minRTT {
		l.minRTT = rtt
	}

	// Add sample
	// For simplicity, we adjust on every Request or small batches.
	// Vegas Formula:
	// queueSize = inflight * (1 - minRTT / rtt)
	// alpha = 3, beta = 6 (standard TCP params roughly)

	// We dampen updates
	queueDelay := rtt - l.minRTT
	if queueDelay < 0 {
		queueDelay = 0
	}

	// Derived queue size approximation
	queueSize := float64(l.inflight+1) * (float64(queueDelay) / float64(rtt))

	alpha := 3.0
	beta := 6.0

	if queueSize < alpha {
		// Latency is low, increase limit
		l.limit += 0.1
	} else if queueSize > beta {
		// Latency is high, decrease limit
		l.limit -= 0.1
	}

	// Enforce bounds
	if l.limit < l.minLimit {
		l.limit = l.minLimit
	} else if l.limit > l.maxLimit {
		l.limit = l.maxLimit
	}
}

// Limit returns current limit.
func (l *Limiter) Limit() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.limit
}
