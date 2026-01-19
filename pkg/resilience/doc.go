/*
Package resilience provides common patterns for building robust, fault-tolerant systems.

This package implements:
  - Circuit Breaker: Prevents cascading failures by stopping requests to failing services.
  - Retry: Automatically retries failed operations with exponential backoff and jitter.

Usage:

	import "github.com/chris-alexander-pop/system-design-library/pkg/resilience"

	// Circuit Breaker
	cb := resilience.NewCircuitBreaker(resilience.DefaultCircuitBreakerConfig("my-service"))

	err := cb.Execute(ctx, func(ctx context.Context) error {
	    return upstream.Call(ctx)
	})

	// Retry
	err := resilience.Retry(ctx, resilience.DefaultRetryConfig(), func(ctx context.Context) error {
	    return upstream.Call(ctx)
	})
*/
package resilience
