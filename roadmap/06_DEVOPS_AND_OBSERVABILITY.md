# DevOps, CI/CD & Observability Roadmap

## Documentation & Spec Generation
- [ ] **OpenAPI (Swagger)**:
    - `swag`: CLI annotation parser to generate docs from Go comments.
    - `oapi-codegen`: Generate Go Types and Fiber/Echo/Chi servers FROM the YAML spec (Design First approach).
- [ ] **AsyncAPI**: Generation for Event Driven Architectures.
- [ ] **Protobuf**: Documentation generation (`protoc-gen-doc`).

## CI/CD Pipelines (GitHub Actions / GitLab CI)
- [ ] **Reusable Workflows**:
    - Build & Test (Go)
    - Lint (golangci-lint)
    - Security Scan (Trivy, Govulncheck)
    - Docker Build & Push
- [ ] **Release Automation**:
    - Semantic Release (Auto-tagging based on commits)
    - Goreleaser configuration for CLI tools
- [ ] **GitOps**:
    - ArgoCD manifest generation
    - FluxCD integration

## Observability Stack
- [ ] **Logging**:
    - **Loki**: Log aggregation setup
    - **Elasticsearch/Fluentd/Kibana (EFK)**: Stack generation
- [ ] **Metrics**:
    - **Prometheus**: ServiceMonitor generation for K8s
    - **Grafana**: Pre-built Dashboards (Golang Runtime, HTTP, Postgres)
- [ ] **Tracing**:
    - **Jaeger**: All-in-one setup
    - **Tempo**: Trace storage integration
- [ ] **Profiling**:
    - **Pprof**: Continuous profiling integration (Pyroscope or Parca)

## Infrastructure as Code (IaC)
- [ ] **Terraform Modules**:
    - AWS RDS / ElastiCache setup
    - GKE / EKS Cluster setup
- [ ] **Pulumi**: Go-based infrastructure definition templates
- [ ] **Ansible**: Playbooks for bare-metal deployment
