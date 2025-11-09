# FlowBank Architecture & Design

## System Overview

FlowBank is an event-driven distributed system that processes banking transactions in real-time. The architecture follows modern microservices patterns with clear separation of concerns.

## Data Flow Pipeline

```
1. PRODUCER (Transaction Generator)
   ↓
   Publishes JSON events to Kafka topic: "transactions"
   ↓
2. KAFKA/REDPANDA (Message Broker)
   ↓
   Distributes events to consumer group
   ↓
3. CONSUMER (Transaction Processor)
   ↓
   Validates → Processes → Persists to DB
   ↓
4. POSTGRESQL (Data Store)
   ↓
   Stores users and transactions
   ↓
5. REST API (Query Interface)
   ↓
   Serves data to clients
```

## Component Details

### 1. Producer Service (`cmd/producer`)

**Purpose**: Generates realistic fake banking transactions

**Responsibilities**:

- Create random transactions (debits/credits)
- Assign to existing users
- Publish events to Kafka topic
- Maintain configurable generation rate

**Event Schema**:

```json
{
  "event_id": "uuid",
  "event_type": "transaction.created",
  "transaction": {
    "id": "uuid",
    "user_id": "uuid",
    "amount": 5000,
    "type": "credit",
    "description": "Salary deposit",
    "merchant_name": "ACME Corp",
    "status": "pending",
    "created_at": "2025-11-06T10:30:00Z"
  },
  "timestamp": "2025-11-06T10:30:00Z"
}
```

### 2. Consumer Service (`cmd/consumer`)

**Purpose**: Processes transaction events from Kafka

**Responsibilities**:

- Consume events from Kafka topic
- Validate transaction data
- Check business rules (e.g., sufficient funds for debits)
- Update transaction status (pending → completed/failed)
- Persist to PostgreSQL
- Handle retries and errors

**Validation Rules**:

- User must exist
- Amount must be > 0
- Type must be "debit" or "credit"
- Debits cannot exceed available balance

### 3. REST API Service (`cmd/api`)

**Purpose**: Provides HTTP interface for querying data

**Endpoints**:

```
GET  /health
     Returns: {"status": "healthy", "service": "flowbank-api"}

GET  /users
     Returns: List of all users

GET  /users/:id
     Returns: User details with current balance

GET  /users/:id/transactions
     Query params: ?limit=20&offset=0
     Returns: Paginated list of user transactions

GET  /transactions/:id
     Returns: Single transaction details
```

### 4. Database Schema

**users table**:

```sql
- id (UUID, primary key)
- email (VARCHAR, unique)
- name (VARCHAR)
- balance (BIGINT) -- in cents
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

**transactions table**:

```sql
- id (UUID, primary key)
- user_id (UUID, foreign key)
- amount (BIGINT) -- in cents
- type (VARCHAR) -- 'debit' or 'credit'
- description (TEXT)
- merchant_name (VARCHAR)
- status (VARCHAR) -- 'pending', 'completed', 'failed'
- created_at (TIMESTAMP)
```

**Indexes**:

- `idx_transactions_user_id` on transactions(user_id)
- `idx_transactions_created_at` on transactions(created_at DESC)
- `idx_transactions_user_created` on transactions(user_id, created_at DESC)

**Trigger**: Automatically updates user balance when transaction status changes to 'completed'

## Infrastructure Components

### RedPanda (Kafka)

**Configuration**:

- Topic: `transactions`
- Partitions: 3 (for parallelism)
- Replication: 1 (dev environment)

**Why RedPanda over Kafka?**:

- Simpler deployment (single binary)
- Lower resource consumption
- Kafka-compatible API
- Better for development/MVP

### PostgreSQL

**Why PostgreSQL?**:

- ACID compliance for financial data
- Rich querying capabilities
- Excellent tooling and documentation
- Triggers for automatic balance updates
- Great for MVP and can scale with partitioning

**Alternative (Cassandra)**: Better for write-heavy workloads, eventual consistency acceptable

## Design Patterns & Best Practices

### 1. Event-Driven Architecture

- **Decoupling**: Services communicate via events, not direct calls
- **Scalability**: Each service can scale independently
- **Resilience**: Kafka provides buffering and retry mechanisms

### 2. Database Per Service (Conceptual)

- Each service owns its data model
- API service reads from PostgreSQL
- Could be split into separate DBs for true microservices

### 3. Idempotency (Future Enhancement)

- Use event_id to prevent duplicate processing
- Critical for financial transactions

### 4. Circuit Breaker (Future Enhancement)

- Protect services from cascading failures
- Implement timeouts and fallbacks

### 5. Observability (Future Enhancement)

- Structured logging
- Metrics (Prometheus)
- Distributed tracing (Jaeger)

## Trade-offs & Decisions

| Decision       | Chosen         | Alternative | Reasoning                                   |
| -------------- | -------------- | ----------- | ------------------------------------------- |
| Message Broker | RedPanda       | Kafka       | Simpler for MVP, Kafka-compatible           |
| Database       | PostgreSQL     | Cassandra   | ACID needed, easier querying                |
| Language       | Go             | Java/Node   | Performance, concurrency, simple deployment |
| Orchestration  | Docker Compose | Kubernetes  | Sufficient for local dev, easier to demo    |
| API Framework  | stdlib http    | Gin/Echo    | Keep it simple, show fundamentals           |

## Scalability Considerations

### Current Limitations (MVP):

- Single instance of each service
- Local Docker Compose deployment
- No load balancing
- No caching layer

### How to Scale:

1. **Horizontal Scaling**: Run multiple consumer instances (Kafka consumer groups)
2. **Database**: Read replicas, connection pooling, partitioning
3. **Caching**: Add Redis for frequently accessed data
4. **Load Balancing**: Add nginx/HAProxy for API
5. **Kubernetes**: Deploy to K8s for production orchestration

## Security Considerations (Future)

- [ ] Authentication/Authorization (JWT)
- [ ] Rate limiting
- [ ] Input validation and sanitization
- [ ] TLS for Kafka and PostgreSQL
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

### Integration Tests

- End-to-end: Producer → Kafka → Consumer → DB → API
- Database migrations
- API endpoints

### Load Testing

- Use k6 or Locust
- Simulate thousands of transactions/second
- Measure latency and throughput

## Development Workflow

1. Start infrastructure: `make up`
2. View logs: `make logs`
3. Make code changes
4. Rebuild: `make build && make restart`
5. Test: `curl http://localhost:8080/health`
6. Query DB: `make db-shell`
7. Monitor Kafka: http://localhost:8090

## Future Enhancements

### Phase 2 (Week 3-4)

- [ ] Implement full API endpoints
- [ ] Add comprehensive validation
- [ ] Add unit tests
- [ ] Add integration tests
- [ ] Implement idempotency

### Phase 3 (Production Ready)

- [ ] Add authentication/authorization
- [ ] Implement rate limiting
- [ ] Add caching layer (Redis)
- [ ] Set up CI/CD pipeline
- [ ] Deploy to cloud (AWS/GCP)
- [ ] Add monitoring and alerting
- [ ] Implement circuit breakers
- [ ] Add dead letter queues

## Interview Discussion Points

### Architecture Questions:

- **Q**: Why event-driven vs REST between services?
  **A**: Asynchronous processing, better scalability, decoupling, resilience

- **Q**: What happens if Consumer crashes?
  **A**: Kafka retains messages, consumer can resume from last committed offset

- **Q**: How do you ensure exactly-once processing?
  **A**: Idempotency keys, database constraints, Kafka transactions

- **Q**: How would you handle high traffic?
  **A**: Scale consumers horizontally, add caching, optimize queries, rate limiting

- **Q**: What's the CAP theorem trade-off here?
  **A**: Prioritizing Consistency + Partition Tolerance (PostgreSQL), eventual consistency for reads acceptable

## References

- [Designing Data-Intensive Applications](https://dataintensive.net/)
- [Microservices Patterns](https://microservices.io/patterns/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Kafka: The Definitive Guide](https://www.confluent.io/resources/kafka-the-definitive-guide/)
