# New Templates & Packages

## Reference Architectures
*Production-ready implementations of common "System Design" problems.*
- [ ] **URL Shortener Service**:
    - High write throughput, Base62 encoding, unique ID generation.
- [ ] **Chat Service**:
    - WebSocket management, Pub/Sub fan-out, Message persistence (Cassandra).
- [ ] **Notification System**:
    - Priority Queues, Rate Limiting per user, Provider abstraction (Email/SMS).
- [ ] **Feed Service**:
    - Fan-out on Write vs Fan-out on Read patterns.

## Service Types
- [ ] **Serverless Function**: AWS Lambda / GCP Cloud Function handler adapters.
- [ ] **BFF (Backend For Frontend)**: Aggregation layer for specific UIs.
- [ ] **CLI Tool**: Cobra-based framework for ops tools.

## Domain Packages
- [ ] **Payments**: Stripe/PayPal intent abstractions.
- [ ] **Email**: Templating and sending (SendGrid/SES).
- [ ] **GeoLocation**: H3/S2 library helpers.
