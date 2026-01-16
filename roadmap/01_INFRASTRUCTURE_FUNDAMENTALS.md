# Infrastructure Fundamentals & Cloud

## Foundational Concepts
- [ ] **Scalability**:
    - Horizontal (Scale Out) vs Vertical (Scale Up) automation.
    - Load Balancing strategies (Round-Robin, Least Connections) implemented in Gateway.
- [ ] **Availability**:
    - Multi-Region / Multi-Zone deployments.
    - Active-Active vs Active-Passive failover configurations.
    - CAP Theorem Trade-offs (Documented decisions for each DB).

## Open Source Messaging (Self-Hosted)
- [ ] **Kafka**:
    - Producer/Consumer using `confluent-kafka-go` or `sarama`.
- [ ] **RabbitMQ**:
    - AMQP 0.9.1 topology management.
- [ ] **NATS**:
    - JetStream advanced patterns (WorkQueue, KV Store).

## AWS (Amazon Web Services)
- [ ] **Computation**: Lambda / Fargate adapters.
- [ ] **Storage**: S3 adapter for `pkg/blob`.
- [ ] **Messaging**: SQS (Queue) & SNS (PubSub).
- [ ] **Auth**: Cognito User Pools.

## GCP (Google Cloud Platform)
- [ ] **Computation**: Cloud Run / Functions.
- [ ] **Storage**: GCS adapter.
- [ ] **Messaging**: Pub/Sub.

## Infrastructure as Code
- [ ] **Terraform**: Modules for every service template.
- [ ] **Kubernetes**: Helm Charts for local/prod disparity.
