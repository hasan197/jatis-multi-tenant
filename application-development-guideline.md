# Application Development Guidelines

This document outlines the technical standards and requirements for developing applications. Adherence to these guidelines ensures consistency, security, and scalability across all deliverables.

---

## 1. Non-Functional Requirements

### 1.1 Performance

#### Throughput & Latency

* **Mandatory**:
  * Handle high volumes of inbound/outbound messages without degradation.
  * Ensure end-to-end processing time ≤100ms under normal conditions.
* **Recommended**:
  * Maintain performance during peak traffic.

#### Resource Utilization

* **Mandatory**:
  * Optimize CPU/memory usage within acceptable thresholds.
  * Prevent memory leaks and excessive resource consumption.

### 1.2 Scalability

#### Horizontal Scalability

* **Mandatory**:
  * Design stateless services for horizontal scaling.
  * Support scaling via instance addition without major reconfiguration.

#### Load Balancing

* **Mandatory**:
  * Configure load balancers to reroute traffic from unhealthy instances.
* **Recommended**:
  * Implement load balancing for even traffic distribution.

### 1.3 Security

#### Data Security

* **Mandatory**:
  * All data security practices **must** align with the [OWASP Top 10](https://owasp.org/www-project-top-ten/) and [OWASP Application Security Verification Standard (ASVS)](https://owasp.org/www-project-application-security-verification-standard/).
  * Encrypt sensitive data (API keys, user info, session tokens) **at rest and in transit** using industry-standard protocols (e.g., AES-256 for storage, TLS 1.3 for communication).
  * Use HTTPS with TLS 1.2+ for **all** inter-service communication.
  * Follow the [OWASP Cryptographic Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cryptographic_Storage_Cheat_Sheet.html) for encryption and key management.
  * Validate and sanitize inputs/outputs per the [OWASP Input Validation Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Input_Validation_Cheat_Sheet.html).
  * Never log sensitive data (credit card numbers, passwords, tokens) – adhere to [OWASP Logging Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html).

#### Access Control

* **Mandatory**:
  * Restrict access to configs, logs, Redis, and databases to authorized personnel.

#### Session Security

* **Mandatory**:
  * Generate secure session IDs to prevent hijacking.
  * Exclude session data from logs/errors.

### 1.4 Reliability & Availability

#### Uptime

* **Mandatory**:
  * Achieve ≥99.9% availability.
  * Implement redundancy for critical components.

#### Fault Tolerance

* **Mandatory**:
  * Gracefully handle external dependency failures (e.g., RabbitMQ, Redis).
  * Auto-reconnect for persistent connections.

#### Data Integrity

* **Mandatory**:
  * Ensure exactly-once message processing.
  * Use ACID-compliant transactions where applicable.
  * Maintain consistent session data with user interactions.

### 1.5 Maintainability

#### Code Quality

* **Mandatory**:
  * Follow clean architecture principles (modularity, separation of concerns).
  * Adhere to coding standards (naming, documentation).

#### Documentation

* **Mandatory**:
  * Provide updated technical design docs, API specs, and user guides.

#### Testing

* **Mandatory**:
  * Implement automated unit, integration, and end-to-end tests.
  * Target ≥80% test coverage.
* **Recommended**:
  * Use CI/CD pipelines for automated testing.

### 1.6 Observability

#### Logging

* **Mandatory**:
  * Use structured JSON logging with timestamps, log levels, and Request IDs.
  * Exclude sensitive data from logs.

#### Monitoring

* **Mandatory**:
  * Expose Prometheus metrics at `/metrics`.
  * Track message processing times, error rates, and RAG metrics.

#### Tracing

* **Mandatory**:
  * Propagate Request IDs across services/logs.
* **Recommended**:
  * Implement distributed tracing.


---

## 2. Technology Stack

### 2.1 GoLang (v1.21)

#### Libraries

* **Web Framework**: `echo`
* **Config**: `viper`
* **Logging**: `logrus` + `lumberjack`
* **Messaging**: `go-rabbitmq`
* **HTTP Client**: `resty`
* **Cron**: `robfig`
* **PostgreSQL**: `pgx`, `pgxutil`
* **CLI**: `cobra`
* **Redis**: `go-redis`

#### Project Structure

```
project-root/
├── cmd/                  # Entry points for your applications (e.g., CLI, API server)
│   ├── root.go
│   ├── serve.go
│   └── version.go   
├── internal/             # Private application code (not importable by others)
│   ├── modules/          # Feature-based modules (e.g., "user", "product")
│   │   └── user/         # Example: User module)
│   │       ├── domain/   # Core business logic
│   │       │   ├── entity.go
│   │       │   └── repository.go  # Interface definitions (e.g., `UserRepository`)
│   │       ├── usecase/  # Business logic/use cases (implements domain interfaces)
│   │       │   └── user_usecase.go
│   │       ├── delivery/ # Delivery layer (HTTP/gRPC/messaging)
│   │       │   ├── http/
│   │       │   │   ├── handler.go
│   │       │   │   └── router.go
│   │       │   └── messaging/ 
│   │       │       └── rabbitmq/
│   │       │           └── user_created_handler.go  
│   │       └── repository/ # Data layer implementations (e.g., DB, cache)
│   │           └── postgresql/  # PostgreSQL-specific implementation
│   │               └── user_repository.go
│   └── config/           # Configuration management
│       └── config.go
├── pkg/                  # Public reusable code
│       ├── infrastructure/
│       │   ├── config/
│       │   ├── logging/
│       │   ├── database/
│       │   ├── cache/
│       │   ├── messaging/
│       │   └── httpclient/
│       ├── delivery/
│       │   └── middleware/
│       └── utils/
│           ├── util1.go
│           └── util2.go            
├── configs/
│   └── config.yaml
├── api/                  # API contracts (OpenAPI/Swagger, Protobuf)
├── scripts/              # Deployment/utility scripts
├── tests/                # Integration/e2e tests
└── main.go               # Main application entry  
├── go.mod
├── go.sum
└── README.md       
```

### 2.2 Python (≥3.8)

* **Web Framework**: `fastapi`
* **Model**: `pydantic`
* **Config**: `pydantic-settings`
* **Logging**: `structlog`
* **Messaging**: `aio-pika`
* **HTTP Client**: `aio-http`
* **Cron**: `apscheduler`
* **PostgreSQL**: `psycopg`, `psycopg-binary`, `psycopg-pool`
* **CLI**: `argparse`
* **Redis**: `redis`

### 2.3 Client-Side TypeScript (v4+)

#### Libraries

* **Framework**: `react@18`
* **Tooling**: `vite`
* **UI**: `material-ui`
* **State/Caching**: `react-query`, `zustand`
* **HTTP**: `axios`
* **Routing**: `react-router-dom`
* **Charts**: `react-chartjs-2`
* **Testing**: `vitest`, `react-testing-library`, `msw`

### 2.4 Server-Side NodeJS/TypeScript (v16.3+)

#### Libraries

* **Web Framework**: `express`
* **HTTP Client**: `got`
* **Config**: `dotenv`
* **Logging**: `winston`


---

## 3. Compliance

* **Deliverables**:
  * Source code adhering to specified libraries/versions.
  * Documentation (technical design, API specs).
  * Test reports with coverage metrics.
* **Reviews**:
  * Code and architecture reviews will validate adherence to guidelines.
* **Updates**:
  * Notify the architecture team of any proposed library/version changes.