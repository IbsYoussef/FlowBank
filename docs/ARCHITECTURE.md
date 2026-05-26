# FlowBank Architecture & Design

## System Overview

FlowBank is an event-driven distributed system that processes banking transactions in real-time. The architecture follows modern microservices patterns with clear separation of concerns.

---

## Local Development Architecture

```
1. PRODUCER (Go)
   ↓
   Publishes JSON transaction events to Kafka topic: "transactions"
   ↓
2. KAFKA / REDPANDA (Message Broker - Docker container)
   ↓
   Distributes events to two independent consumer groups
   ↓
3a. GO CONSUMER                        3b. PYTHON FRAUD DETECTION
    Validates → Processes → DB              Consumes → Scores → DB
   ↓                                        ↓
4. POSTGRESQL (Docker container)
   └── users table
   └── transactions table
   └── fraud_scores table
   ↓
5. REST API (Go) + DASHBOARD
   Serves analytics and fraud data
```

## Production Architecture (AWS)

```
1. PRODUCER (Go - ECR image)
   ↓
   Publishes to Aiven Kafka (managed, hosted in EU)
   ↓
2. AIVEN KAFKA (Managed Kafka 4.1 - SASL_SSL)
   ↓
   ↓──────────────────────────────────────────┐
3a. GO CONSUMER (ECR image)            3b. PYTHON FRAUD DETECTION (ECR image)
    Validates → Processes → RDS             Consumes → Scores → RDS
   ↓                                        ↓
4. AWS RDS POSTGRESQL (Managed, SSL required)
   Hosted in eu-west-2 (London)
   ↓
5. GO API (ECR image) + DASHBOARD
   Exposed on port 80 via Elastic Beanstalk

All containers run on a single AWS Elastic Beanstalk environment (t3.small EC2)
Container images stored in and pulled from AWS ECR
```

---

## Component Details

### 1. Producer Service (`cmd/producer`)

**Purpose**: Generates realistic fake banking transactions

**Responsibilities**:

- Create random transactions (debits/credits) for seed users
- Assign to existing users (Alice and Bob)
- Publish events to Kafka topic every 3 seconds
- Support both local (PLAINTEXT) and production (SASL_SSL) Kafka connections
  **Kafka connection**: Uses `internal/kafka/producer.go` which automatically detects SASL credentials and applies TLS config if present.

### 2. Consumer Service (`cmd/consumer`)

**Purpose**: Processes transaction events from Kafka

**Responsibilities**:

- Consume events from Kafka topic
- Validate transaction data
- Check business rules (overdraft protection for debits)
- Update transaction status (pending → completed/failed)
- Persist to PostgreSQL with atomic balance updates
- Support both local and production Kafka connections

### 3. REST API Service (`cmd/api`)

**Purpose**: Provides HTTP interface for querying data and serving the dashboard

**Endpoints**:

```
GET  /health
     Returns service health status

GET  /api/v1/accounts/{userID}
     Returns user details with current balance and transaction history

GET  /api/v1/analytics
     Returns aggregated dashboard data (totals, charts, recent transactions)

GET  /api/v1/fraud-scores
     Returns recent fraud scores joined with transaction details

GET  /dashboard
     Serves the analytics dashboard HTML
```

### 4. Fraud Detection Service (`fraud-detection/`)

**Purpose**: Scores transactions for fraud risk in real-time

**Language**: Python 3.14 / FastAPI

**Responsibilities**:

- Consume transaction events from the same Kafka topic as the Go consumer
- Apply three fraud detection rules to each transaction
- Write scored results to the `fraud_scores` table
- Expose fraud scores via REST API for dashboard consumption
  **Three Detection Rules**:

| Rule                  | Description                              | Risk Level |
| --------------------- | ---------------------------------------- | ---------- |
| Duplicate Transaction | Same user, same amount within 60 seconds | Medium     |
| High Value            | Transaction above $10,000                | Medium     |
| High Frequency        | More than 5 transactions in 60 seconds   | Medium     |

Risk scoring:

- 0 rules triggered → LOW / clean
- 1 rule triggered → MEDIUM / suspicious
- 2+ rules triggered → HIGH / flagged
  **Kafka connection**: Uses SASL_SSL in production with base64-encoded CA certificate passed as environment variable.

### 4. Database Schema

---

### users table

```sql
user_id         UUID PRIMARY KEY
email           VARCHAR(100) UNIQUE NOT NULL
user_name       VARCHAR(100) NOT NULL
current_balance BIGINT NOT NULL DEFAULT 0  -- in cents
created_at      TIMESTAMP WITH TIME ZONE
updated_at      TIMESTAMP WITH TIME ZONE   -- auto-updated via trigger
```

### transactions table

```sql
transaction_id  UUID PRIMARY KEY
user_id         UUID NOT NULL REFERENCES users(user_id)
amount          BIGINT NOT NULL             -- in cents
transaction_type VARCHAR(20) NOT NULL       -- 'debit' or 'credit'
description     VARCHAR(255)
merchant_name   VARCHAR(100)
status          VARCHAR(20) NOT NULL        -- 'COMPLETED' or 'FAILED'
created_at      TIMESTAMP WITH TIME ZONE
```

### fraud_scores table

```sql
fraud_score_id    UUID PRIMARY KEY DEFAULT gen_random_uuid()
transaction_id    UUID NOT NULL REFERENCES transactions(transaction_id)
risk_score        VARCHAR(10) NOT NULL      -- 'low', 'medium', 'high'
status            VARCHAR(20) NOT NULL      -- 'clean', 'suspicious', 'flagged'
triggered_rules   TEXT[] NOT NULL DEFAULT '{}'
confidence        DECIMAL(3,2) NOT NULL
processing_time_ms DECIMAL(10,2)
scored_at         TIMESTAMP WITH TIME ZONE
```

**Indexes**:

- `idx_transactions_user_id_created_at` - fraud frequency queries
- `idx_transactions_user_id_amount` - duplicate detection queries
- `idx_fraud_scores_transaction_id` - join performance
- `idx_fraud_scores_status` - dashboard filtering
- `idx_fraud_scores_scored_at` - ordering

---

---

## Infrastructure

### Local (Docker Compose)

| Service          | Image                           | Port                              |
| ---------------- | ------------------------------- | --------------------------------- |
| RedPanda (Kafka) | confluentinc/cp-kafka:7.6.1     | 9092 (internal), 19092 (external) |
| Zookeeper        | confluentinc/cp-zookeeper:7.6.1 | 2181                              |
| PostgreSQL       | postgres:16-alpine              | 5432                              |
| Producer         | Built from Dockerfile.producer  | -                                 |
| Consumer         | Built from Dockerfile.consumer  | -                                 |
| API              | Built from Dockerfile.api       | 8080                              |
| Fraud Detection  | Built from Dockerfile.fraud     | 8000                              |

### Production (AWS)

| Component          | Service                          | Notes                                                     |
| ------------------ | -------------------------------- | --------------------------------------------------------- |
| Compute            | AWS Elastic Beanstalk (t3.small) | Runs all containers via Docker Compose                    |
| Container Registry | AWS ECR                          | 4 repositories (api, producer, consumer, fraud-detection) |
| Database           | AWS RDS PostgreSQL 16            | eu-west-2, db.t3.micro, SSL required                      |
| Message Broker     | Aiven Kafka 4.1                  | Free tier, EU region, SASL_SSL                            |
| Instance type      | t3.small (2GB RAM)               | Required to run 4 containers simultaneously               |

---

## Design Patterns & Best Practices

### Event-Driven Architecture

- **Decoupling**: Services communicate via events, not direct calls
- **Scalability**: Each service can scale independently
- **Extensibility**: New consumers (fraud detection) added without modifying existing services
- **Resilience**: Kafka provides buffering and retry mechanisms

### Additive Extension Pattern

The fraud detection service was added without modifying any existing Go services. It simply subscribes to the same Kafka topic as a new consumer group. This demonstrates the open/closed principle at the architectural level.

### Polyglot Microservices

Go is used for the transaction processing pipeline (performance, concurrency). Python is used for fraud detection (FastAPI async support, better data processing ecosystem, extensible to ML).

### TLS Configuration

Both Go and Python Kafka clients support dual-mode connection: PLAINTEXT for local development (no env vars set) and SASL_SSL for production (env vars present). The CA certificate is stored as a base64-encoded environment variable and decoded at startup.

---

## Trade-offs & Decisions

| Decision               | Chosen            | Alternative  | Reasoning                                   |
| ---------------------- | ----------------- | ------------ | ------------------------------------------- |
| Message Broker (local) | RedPanda          | Kafka        | Simpler for local dev, Kafka-compatible     |
| Message Broker (prod)  | Aiven Kafka       | AWS MSK      | Free tier available, no Zookeeper overhead  |
| Database               | PostgreSQL        | Cassandra    | ACID needed for financial data              |
| Language (pipeline)    | Go                | Java/Node    | Performance, concurrency, simple deployment |
| Language (fraud)       | Python            | Go           | Async FastAPI, better ML extensibility      |
| Orchestration (local)  | Docker Compose    | Kubernetes   | Sufficient for local dev                    |
| Hosting                | Elastic Beanstalk | ECS/EKS      | Managed platform, simpler for MVP           |
| Container Registry     | ECR               | Docker Hub   | Integrated with AWS IAM, private            |
| API Framework (Python) | FastAPI           | Flask/Django | Async, auto-docs, Pydantic validation       |

---

## Scalability Considerations

### Current Limitations

- Single EC2 instance running all containers
- Two seed users (Alice and Bob) generate high-frequency fraud flags
- 3-second consumer delay to avoid race conditions between Go consumer and fraud detection

### How to Scale

1. **Horizontal Scaling**: Run multiple consumer instances (Kafka consumer groups handle partition assignment)
2. **Separate Instances**: Move fraud detection to its own EC2 instance or ECS task
3. **Database**: Read replicas for analytics queries, connection pooling
4. **Kafka**: Increase partitions for higher throughput
5. **ECS/EKS**: Migrate to proper container orchestration for production

---

5. **ECS/EKS**: Migrate to proper container orchestration for production

## Future Enhancements

### Short Term

- CI/CD pipeline (GitHub Actions → ECR → Elastic Beanstalk)
- Service status indicators on dashboard
- More realistic transaction data (more users, merchant categories)
- Pagination on transaction and fraud feeds

### Medium Term

- Prometheus metrics and Grafana dashboards
- Unit and integration tests
- Reduced consumer delay via retry mechanism

### Long Term

- ECS migration for better container orchestration
- ML-based fraud scoring to replace rule-based approach
- Authentication and rate limiting
- Dead-letter queues for failed messages

---

## Interview Discussion Points

- **Q**: Why event-driven vs REST between services?
  **A**: Asynchronous processing, better scalability, decoupling, resilience. The fraud detection service was added without changing any existing service.
- **Q**: What happens if the Consumer crashes?
  **A**: Kafka retains messages. The consumer can resume from the last committed offset - no messages lost.
- **Q**: How do you ensure exactly-once processing?
  **A**: Currently at-least-once. Production improvement would add idempotency keys and database constraints.
- **Q**: Why Python for fraud detection rather than Go?
  **A**: Python's data processing ecosystem is stronger and the service is easily extensible to ML-based scoring. FastAPI's async model integrates cleanly with Kafka consumption.
- **Q**: How would you handle high traffic?
  **A**: Scale consumers horizontally, add caching, increase Kafka partitions, move to ECS for better orchestration.
- **Q**: Why Elastic Beanstalk rather than ECS?
  **A**: Faster to deploy for an MVP. ECS would be the natural next step for production - the ECR setup already supports it.

---

## Security Considerations (Future)

- [ ] Authentication/Authorization (JWT)
- [ ] Rate limiting
- [ ] Input validation and sanitization
- [ ] TLS for Kafka and PostgreSQL (already implemented for production via SASL_SSL)
- [ ] Secrets management (Vault)
- [ ] API key management

## Monitoring & Alerting (Future)

- [ ] Prometheus metrics
- [ ] Grafana dashboards
- [ ] Log aggregation (ELK stack)
- [ ] Health check endpoints
- [ ] Dead letter queues for failed messages

## Testing Strategy

### Unit Tests

- Business logic in `internal/service`
- Database queries
- Kafka producer/consumer wrappers
- Fraud detection rules in `fraud-detection/app/detector.py`

### Integration Tests

- End-to-end: Producer → Kafka → Consumer → DB → API
- End-to-end: Producer → Kafka → Fraud Detection → DB → API
- Database migrations
- API endpoints

### Load Testing

- Use k6 or Locust
- Simulate thousands of transactions/second
- Measure latency and throughput
- Verify sub-100ms p95 fraud scoring latency

## Development Workflow

1. Start infrastructure: `make up`
2. View logs: `make logs`
3. Make code changes
4. Rebuild: `make build && make restart`
5. Test: `curl http://localhost:8080/health`
6. Query DB: `make db-shell`
7. Monitor Kafka: http://localhost:8090
8. Check fraud scores: `curl http://localhost:8080/api/v1/fraud-scores`

## Future Enhancements

### Phase 2

- [ ] Add comprehensive validation
- [ ] Add unit tests
- [ ] Add integration tests
- [ ] Implement idempotency
- [ ] CI/CD pipeline (GitHub Actions → ECR → Elastic Beanstalk)

### Phase 3 (Production Ready)

- [ ] Add authentication/authorization
- [ ] Implement rate limiting
- [ ] Add caching layer (Redis)
- [ ] Add monitoring and alerting (Prometheus/Grafana)
- [ ] Implement circuit breakers
- [ ] Add dead letter queues
- [ ] ECS migration for better container orchestration
- [ ] ML-based fraud scoring

## Interview Discussion Points

### Architecture Questions

- **Q**: Why event-driven vs REST between services?
  **A**: Asynchronous processing, better scalability, decoupling, resilience. The fraud detection service was added without modifying any existing service.
- **Q**: What happens if Consumer crashes?
  **A**: Kafka retains messages, consumer can resume from last committed offset.
- **Q**: How do you ensure exactly-once processing?
  **A**: Currently at-least-once. Production improvement would add idempotency keys, database constraints, and Kafka transactions.
- **Q**: How would you handle high traffic?
  **A**: Scale consumers horizontally, add caching, optimise queries, rate limiting, increase Kafka partitions.
- **Q**: What's the CAP theorem trade-off here?
  **A**: Prioritising Consistency + Partition Tolerance (PostgreSQL). Eventual consistency for reads is acceptable.
- **Q**: Why Python for fraud detection rather than Go?
  **A**: Python's data processing ecosystem is stronger and the service is easily extensible to ML-based scoring. FastAPI's async model integrates cleanly with Kafka consumption.
- **Q**: Why Elastic Beanstalk rather than ECS?
  **A**: Faster to deploy for an MVP. ECS would be the natural next step - the ECR setup already supports it.

## References

- [Designing Data-Intensive Applications](https://dataintensive.net/)
- [Microservices Patterns](https://microservices.io/patterns/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Kafka: The Definitive Guide](https://www.confluent.io/resources/kafka-the-definitive-guide/)
- [FastAPI Documentation](https://fastapi.tiangolo.com/)
