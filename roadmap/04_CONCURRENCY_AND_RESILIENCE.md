# Concurrency & Resilience

## Critical Reliability Features
- [ ] **Rate Limiter**:
    - **Token Bucket**: Burstable traffic control.
    - **Leaky Bucket**: Smooth traffic flow.
    - **Sliding Window**: Precise throttling.
    - Distributed Rate Limiting via Redis Lua scripts.
- [ ] **Circuit Breaker**:
    - State Machine (Open/Closed/Half-Open).
    - Failure threshold detection.
- [ ] **Consistent Hashing**:
    - Hash Ring implementation for distributed load balancing.
    - Virtual Nodes for skew minimization.

## Concurrency Patterns
- [ ] **Worker Pools**: Fixed-size goroutine pools.
- [ ] **Pipelines**: Stream processing stages (Fan-out / Fan-in).
- [ ] **Semaphores**: Limiting concurrent access to resources.
- [ ] **Distributed Locks**: Redlock implementation.

## Optimization
- [ ] **Bloom Filters**: Probabilistic set checking (Database avoidance).
- [ ] **HyperLogLog**: Cardinality estimation.
