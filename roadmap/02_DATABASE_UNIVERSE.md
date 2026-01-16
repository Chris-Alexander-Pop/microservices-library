# The Database Universe

## Core Concepts
- [ ] **Sharding & Partitioning**:
    - Logic for application-level sharding (e.g. by UserID).
    - Database-level partitioning (Range, Hash, List).
- [ ] **Replication**:
    - Read Replica lag handling in `pkg/database`.
    - Leader Election for failover.

## SQL & Relational
- [ ] **Postgres**: Advanced features (Partitions, RLS).
- [ ] **MySQL/MariaDB**: Deadlock handling.
- [ ] **SQLite**: Embedded for edge services.

## NoSQL & Document
- [ ] **MongoDB**: Aggregation pipelines.
- [ ] **Cassandra/Scylla**: Wide-column store patterns.
- [ ] **DynamoDB**: Single Table Design helpers.

## Vector & AI
- [ ] **pgvector**: Postgres as a Vector DB.
- [ ] **Pinecone/Milvus**: Specialized Search.

## Caching
- [ ] **Strategies**:
    - Cache-Aside pattern helper (`GetOrSet`).
    - Write-Through / Write-Back interfaces.
- [ ] **Redis**: Distributed Cache.
- [ ] **Ristretto**: In-Memory (Local) Cache.
