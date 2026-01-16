# Enterprise Architecture Patterns Roadmap

## Domain-Driven Design (DDD)
- [ ] **Aggregate Root Base**:
    - Helper struct for tracking domain events (`RecordEvent`, `PullEvents`).
- [ ] **Value Objects**:
    - Generic `Money`, `Address`, `Email` value objects with validation.
- [ ] **Specification Pattern**:
    - `Spec` interface for composable business rules (`And`, `Or`, `Not`).
    - SQL generation from specifications.
- [ ] **Repository Pattern**:
    - Generic `Repository[T]` interface with `Find`, `Save`, `Delete`.

## CQRS (Command Query Responsibility Segregation)
- [ ] **Command Bus**:
    - In-memory synchronous bus (`Dispatch(cmd)`).
    - Middleware support (Logging, Validation, Transactions).
- [ ] **Query Bus**:
    - Caching decorator for queries.

## Event Sourcing
- [ ] **Event Store**:
    - Append-only log interface.
    - Optimistic concurrency control (Version checks).
- [ ] **Snapshots**:
    - Strategy for snapshotting aggregates every N events.
- [ ] **Projections**:
    - Async workers to project events into Read Models (SQL/Mongo).

## Transactions
- [ ] **Unit of Work**:
    - Interface to manage transactions across multiple repositories.
- [ ] **Saga Pattern**:
    - Orchestration-based Saga manager (State machine).
    - Choreography helpers (Event headers for correlation).
- [ ] **Outbox Pattern**:
    - Guaranteed delivery of domain events to the message bus.
