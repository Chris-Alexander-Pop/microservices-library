# Frontier Tech & Specialized Domains Roadmap

## Web3 & Blockchain
- [ ] **Ethereum**:
    - `go-ethereum` (geth) client wrapper.
    - Smart Contract interaction (ABI binding generation).
    - Wallet management (Key generation, Signing).
- [ ] **Solana**:
    - RPC Client integration.
- [ ] **IPFS**:
    - Content addressing and storage adapters.

## Game Development Backend
- [ ] **Networking**:
    - **UDP Server**: High-performance, low-latency packet handling.
    - **WebSockets**: State synchronization for turn-based games.
- [ ] **Services**:
    - **Matchmaking**: Queue logic based on ELO/Skill.
    - **Leaderboards**: Redis Sorted Sets implementation (`ZADD`, `ZRANGE`).
    - **Inventory**: Atomic item swapping/trading logic.

## Data Engineering & ETL
- [ ] **File Formats**:
    - **Parquet**: Reading/Writing columnar data (`parquet-go`).
    - **Avro / Orc**: Schema-based serialization.
    - **CSV / Excel**: High-performance streaming parsers.
- [ ] **Pipelines**:
    - **DAG Scheduler**: Simple dependency graph execution (A-la Airflow but lightweight).
    - **Stream Processing**: Windowing and aggregation logic for data streams.
