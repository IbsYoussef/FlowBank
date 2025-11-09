# 🏦 FlowBank

A miniature event-driven transaction processor inspired by Monzo, built with Go, Kafka (RedPanda), and PostgreSQL.

## 🎯 Project Overview

FlowBank demonstrates a distributed system for processing banking transactions using event-driven architecture. It showcases:

- **Event Streaming**: Real-time transaction processing via Kafka/RedPanda
- **Microservices**: Three independent Go services (Producer, Consumer, API)
- **Data Pipeline**: Producer → Kafka → Consumer → PostgreSQL → API
- **Docker Orchestration**: Fully containerized development environment

## 🏗️ Architecture

```
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│   Producer   │────────>│   RedPanda   │────────>│   Consumer   │
│ (Fake Txns)  │ publish │   (Kafka)    │ consume │  (Validate)  │
└──────────────┘         └──────────────┘         └──────┬───────┘
                                                         │
                                                         │ persist
                                                         ▼
┌──────────────┐                                  ┌──────────────┐
│  REST API    │ <─────────── read ─────────────> │  PostgreSQL  │
│ (Query Data) │                                  │  (Storage)   │
└──────────────┘                                  └──────────────┘
```

## 📁 Project Structure

```
flowbank/
├── cmd/
│   ├── producer/       # Generates fake banking transactions
│   ├── consumer/       # Processes transactions from Kafka
│   └── api/            # REST API for querying data
├── internal/
│   ├── kafka/          # Kafka producer/consumer helpers
│   ├── db/             # Database connection and queries
│   ├── model/          # Data models (Transaction, User)
│   └── service/        # Business logic layer
├── deploy/
│   ├── docker-compose.yml
│   ├── init.sql        # Database schema
│   └── Dockerfile.*    # Service Dockerfiles
├── go.mod
├── Makefile
└── README.md
```

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.22+ (for local development)
- Make (optional, for convenience)

### Running the System

1. **Start all services:**

   ```bash
   make up
   # or
   docker compose -f deploy/docker-compose.yml up -d
   ```

2. **Check service health:**

   ```bash
   make ps
   ```

3. **View logs:**

   ```bash
   make logs
   # or for specific services
   make logs-producer
   make logs-consumer
   make logs-api
   ```

4. **Access services:**
   - API: http://localhost:8080/health
   - RedPanda Console: http://localhost:8090
   - PostgreSQL: `localhost:5432` (user: `flowbank_user`, db: `flowbank`)

### Testing the Data Pipeline

1. **Check API health:**

   ```bash
   curl http://localhost:8080/health
   ```

2. **View users:**

   ```bash
   # Coming soon: curl http://localhost:8080/users
   ```

3. **View user transactions:**

   ```bash
   # Coming soon: curl http://localhost:8080/users/{user_id}/transactions
   ```

4. **Monitor Kafka topics via RedPanda Console:**

   - Open http://localhost:8090
   - Navigate to Topics → transactions
   - View real-time messages

5. **Query database directly:**
   ```bash
   make db-shell
   # Then run SQL:
   SELECT * FROM users;
   SELECT * FROM transactions ORDER BY created_at DESC LIMIT 10;
   ```

## 🛠️ Development Commands

| Command             | Description                   |
| ------------------- | ----------------------------- |
| `make help`         | Show all available commands   |
| `make deps`         | Install Go dependencies       |
| `make build`        | Build Docker images           |
| `make up`           | Start all services            |
| `make down`         | Stop all services             |
| `make restart`      | Restart services              |
| `make logs`         | Tail logs from all services   |
| `make clean`        | Remove containers and volumes |
| `make db-shell`     | Open PostgreSQL shell         |
| `make kafka-topics` | List Kafka topics             |
| `make test`         | Run Go tests                  |

## 📊 Data Models

### Transaction

```go
type Transaction struct {
    ID           uuid.UUID
    UserID       uuid.UUID
    Amount       int64     // in cents
    Type         string    // "debit" or "credit"
    Description  string
    MerchantName string
    Status       string    // "pending", "completed", "failed"
    CreatedAt    time.Time
}
```

### User

```go
type User struct {
    ID        uuid.UUID
    Email     string
    Name      string
    Balance   int64     // in cents
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## 🔧 Technology Stack

- **Language**: Go 1.22
- **Messaging**: RedPanda (Kafka-compatible)
- **Database**: PostgreSQL 16
- **Containerization**: Docker & Docker Compose
- **Development**: Make, Air (hot reload - coming soon)

## 📝 Roadmap

- [x] Project structure and Docker setup
- [x] Implement Kafka producer/consumer
- [x] Build transaction generator (Producer)
- [x] Build transaction processor (Consumer)
- [x] Implement REST API endpoints
- [x] Add validation and business logic
- [x] Add comprehensive logging
- [ ] Add metrics and observability
- [ ] Add unit and integration tests
- [ ] Add CI/CD pipeline

## 🎓 Learning Goals

This project demonstrates understanding of:

1. **Event-Driven Architecture**: Asynchronous communication via message queues
2. **Microservices**: Loosely coupled, independently deployable services
3. **Distributed Systems**: Service orchestration, data consistency, fault tolerance
4. **Go Best Practices**: Project structure, error handling, concurrency
5. **DevOps**: Containerization, infrastructure as code, observability

## 📚 Resources

- [RedPanda Documentation](https://docs.redpanda.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

## 🤝 Interview Talking Points

- **Why event-driven?** Decouples services, enables async processing, scales horizontally
- **Why RedPanda over Kafka?** Simpler deployment, lower latency, Kafka-compatible
- **Why PostgreSQL?** ACID compliance, robust querying, great for MVP
- **Trade-offs**: Could use Cassandra for write-heavy workloads, Redis for caching
- **Improvements**: Add idempotency, dead-letter queues, circuit breakers, rate limiting

---

**Author**: IbsYoussef

**Purpose**: Backend internship interview project

**Timeline**: 2 weeks to MVP
