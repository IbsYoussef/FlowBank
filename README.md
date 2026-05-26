# 🏦 FlowBank

A miniature event-driven transaction processor inspired by Monzo, built with Go, Kafka (RedPanda), Python, and PostgreSQL. Deployed to AWS.

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Python](https://img.shields.io/badge/Python-3776AB?style=for-the-badge&logo=python&logoColor=white)](https://www.python.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Kafka](https://img.shields.io/badge/Kafka-231F20?style=for-the-badge&logo=apache-kafka&logoColor=white)](https://kafka.apache.org/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![AWS](https://img.shields.io/badge/AWS-232F3E?style=for-the-badge&logo=amazon-aws&logoColor=white)](https://aws.amazon.com/)
[![Elastic Beanstalk](https://img.shields.io/badge/Elastic%20Beanstalk-232F3E?style=for-the-badge&logo=amazonwebservices&logoColor=white)](https://aws.amazon.com/)
[![Amazon RDS](https://img.shields.io/badge/Amazon%20RDS-527FFF?style=for-the-badge&logo=amazonrds&logoColor=white)](https://aws.amazon.com/rds/)
[![Amazon ECR](https://img.shields.io/badge/Amazon%20ECR-232F3E?style=for-the-badge&logo=amazonwebservices&logoColor=white)](https://aws.amazon.com/ecr/)
[![JavaScript](https://img.shields.io/badge/JavaScript-F7DF1E?style=for-the-badge&logo=javascript&logoColor=black)](https://developer.mozilla.org/en-US/docs/Web/JavaScript)
[![HTML5](https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white)](https://developer.mozilla.org/en-US/docs/Web/HTML)
[![CSS3](https://img.shields.io/badge/CSS3-1572B6?style=for-the-badge&logo=css3&logoColor=white)](https://developer.mozilla.org/en-US/docs/Web/CSS)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)

---

## 📑 Table of Contents

- [🎯 Project Overview](#-project-overview)
- [🏗️ Architecture](#️-architecture)
- [📁 Project Structure](#-project-structure)
- [🚀 Quick Start](#-quick-start)
- [📊 Analytics Dashboard](#-analytics-dashboard)
- [🛠️ Development Commands](#️-development-commands)
- [📊 Data Models](#-data-models)
- [🔧 Technology Stack](#-technology-stack)
- [📝 Roadmap](#-roadmap)
- [🎓 Learning Goals](#-learning-goals)
- [📚 Resources](#-resources)
- [🤝 Interview Talking Points](#-interview-talking-points)
- [📚 Documentation](#-documentation)
- [📄 License](#-license)
- [👤 Author](#-author)

---

## 🎯 Project Overview

FlowBank demonstrates a distributed system for processing banking transactions using event-driven architecture. Inspired by Monzo's architecture, it showcases modern microservices patterns, real-time data processing, interactive analytics, and fraud detection - all deployed to AWS.

### ✨ Key Features

- ⚡ **Event-Driven Architecture** - Asynchronous transaction processing via Kafka
- 📊 **Real-Time Analytics Dashboard** - Live data visualization
- 🏗️ **Microservices Design** - Independent Go services (Producer, Consumer, API) and a Python fraud detection service
- 🔍 **Real-Time Fraud Detection** - Python FastAPI microservice scoring transactions against three detection rules
- 🔒 **ACID Transactions** - PostgreSQL with row-level locking and atomic updates
- 🐳 **Fully Dockerized** - One-command local deployment with Docker Compose
- 🐍 **Polyglot Architecture** - Go and Python microservices working together on the same Kafka topics
- ☁️ **Cloud Deployed** - Live on AWS Elastic Beanstalk with RDS PostgreSQL and Aiven Kafka
- 📈 **Production Patterns** - Proper error handling, logging, separation of concerns

---

## 🏗️ Architecture

### Local Development Architecture

```
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│   Producer   │────────>│   RedPanda   │────────>│   Consumer   │
│    (Go)      │ Publish │   (Kafka)    │ Consume │    (Go)      │
└──────────────┘         └──────┬───────┘         └──────┬───────┘
                                │                        │
                                │ Consume                ▼
                                │                  ┌──────────────┐
                                └─────────────────>│    Fraud     │
                                                   │   Detection  │
                                                   │   (Python)   │
                                                   └──────┬───────┘
                                                          │
                                                          ▼
┌──────────────┐                                  ┌──────────────┐
│  Dashboard   │<────────── Query ────────────────│  PostgreSQL  │
│ (Go API)     │                                  │  (Storage)   │
└──────────────┘                                  └──────────────┘
```

### Production Architecture (AWS)

```
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│   Producer   │────────>│    Aiven     │────────>│   Consumer   │
│  (Go / ECR)  │ Publish │    Kafka     │ Consume │  (Go / ECR)  │
└──────────────┘         └──────┬───────┘         └──────┬───────┘
                                │                        │
                                │ Consume                ▼
                                │                  ┌──────────────┐
                                └─────────────────>│    Fraud     │
                                                   │  Detection   │
                                                   │ (Python/ECR) │
                                                   └──────┬───────┘
                                                          │
                                                          ▼
┌──────────────┐                                  ┌──────────────┐
│  Dashboard   │<────────── Query ────────────────│   AWS RDS    │
│ (Go API/ECR) │                                  │  PostgreSQL  │
└──────────────┘                                  └──────────────┘

         All services hosted on AWS Elastic Beanstalk (t3.small EC2)
         Container images stored in AWS ECR
```

**Data Flow:**

1. **Producer** generates realistic banking transactions every 3 seconds
2. **Kafka** buffers and distributes events across partitions
3. **Go Consumer** validates, processes, and persists to database with overdraft protection
4. **Python Fraud Detection** consumes the same Kafka topic, scores each transaction against three fraud rules, and writes results to the `fraud_scores` table
5. **API** serves transaction and fraud score data via REST endpoints
6. **Dashboard** visualizes metrics and flagged transactions in real-time with auto-refresh

---

## 📁 Project Structure

```
flowbank/
├── cmd/
│   ├── producer/           # Generates fake banking transactions
│   ├── consumer/           # Processes transactions from Kafka
│   └── api/                # REST API for querying data
├── internal/
│   ├── kafka/              # Kafka producer/consumer helpers + TLS config
│   ├── db/                 # Database connection and queries
│   ├── model/              # Data models (Transaction, User, FraudScore)
│   └── service/            # Business logic layer
├── fraud-detection/        # Python fraud detection microservice
│   ├── app/
│   │   ├── main.py         # FastAPI entry point
│   │   ├── consumer.py     # Kafka consumer logic
│   │   ├── detector.py     # Fraud scoring rules
│   │   ├── database.py     # PostgreSQL connection and writes
│   │   └── models.py       # Pydantic schemas
│   └── requirements.txt
├── deploy/
│   ├── docker-compose.yml  # Local development
│   ├── Dockerfile.*        # Service Dockerfiles
│   └── init.sql            # Database schema
├── docker-compose.yml      # Production (AWS) - uses ECR images
├── web/
│   └── dashboard.html      # Analytics dashboard UI
├── docs/                   # Additional documentation
├── go.mod
├── Makefile
└── README.md
```

---

## 🌐 Live Demo

**[View Live Dashboard →](http://flowbank-prod.eba-fcdmxpas.eu-west-2.elasticbeanstalk.com/dashboard)**

| Endpoint     | URL                                                                                                          |
| ------------ | ------------------------------------------------------------------------------------------------------------ |
| Dashboard    | [/dashboard](http://flowbank-prod.eba-fcdmxpas.eu-west-2.elasticbeanstalk.com/dashboard)                     |
| API Health   | [/health](http://flowbank-prod.eba-fcdmxpas.eu-west-2.elasticbeanstalk.com/health)                           |
| Fraud Scores | [/api/v1/fraud-scores](http://flowbank-prod.eba-fcdmxpas.eu-west-2.elasticbeanstalk.com/api/v1/fraud-scores) |
| Analytics    | [/api/v1/analytics](http://flowbank-prod.eba-fcdmxpas.eu-west-2.elasticbeanstalk.com/api/v1/analytics)       |

---

## 🚀 Quick Start

> **See [docs/QUICKSTART.md](docs/QUICKSTART.md) for the full guide. Start there.**

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
- [Go 1.26+](https://go.dev/dl/) (optional, for local development)
- [Make](https://www.gnu.org/software/make/) (optional, for convenience)

### Running the System

```bash
# Clone the repository
git clone https://github.com/IbsYoussef/FlowBank.git
cd FlowBank

# Start all services
make up

# View logs (optional)
make logs
```

**Local Access Points:**

- 📊 **Dashboard**: http://localhost:8080/dashboard
- 🏥 **API Health**: http://localhost:8080/health
- 💾 **PostgreSQL**: `localhost:5432` (user: `flowbank_user`, pass: `flowbank_pass`, db: `flowbank`)

### Testing the Data Pipeline

#### 1. Check API Health

```bash
curl http://localhost:8080/health
```

#### 2. View Users

```bash
# Coming soon: curl http://localhost:8080/users
```

#### 3. View User Transactions

```bash
# Coming soon: curl http://localhost:8080/users/{user_id}/transactions
```

#### 4. Monitor Kafka Topics

- Open http://localhost:8090
- Navigate to Topics → transactions
- View real-time messages

#### 5. Verify Services

```bash
make ps
# All services should show as "Up" and "healthy"
```

#### 6. Check Database

```bash
make db-shell

# View users
SELECT * FROM users;

# View recent transactions
SELECT * FROM transactions ORDER BY created_at DESC LIMIT 10;

# Check transaction stats
SELECT status, COUNT(*) FROM transactions GROUP BY status;

# View fraud scores
SELECT * FROM fraud_scores ORDER BY scored_at DESC LIMIT 10;
```

#### 7. Monitor Live Transactions

**Option A: Dashboard (Recommended)**

```bash
open http://localhost:8080/dashboard
```

**Option B: Producer Logs**

```bash
make logs-producer
```

**Option C: Consumer Logs**

```bash
make logs-consumer
```

---

## 📊 Analytics Dashboard

![Dashboard Preview](/assets/img.png)

**Local:** http://localhost:8080/dashboard
**Live:** [flowbank-prod.eba-fcdmxpas.eu-west-2.elasticbeanstalk.com/dashboard](http://flowbank-prod.eba-fcdmxpas.eu-west-2.elasticbeanstalk.com/dashboard)

### Features

- 📈 **4 Interactive Charts** - Transaction type distribution, money flow, top merchants, hourly volume
- 🔴 **Live Transaction Feed** - Real-time updates every 5 seconds
- 🚨 **Flagged Transactions Panel** - Suspicious and high-risk transactions with risk badges
- 💰 **Key Metrics** - Total transactions, money in/out, average amount
- 📱 **Responsive Design** - Works on desktop, tablet, mobile

### API Endpoints

| Endpoint                | Method | Description                                  |
| ----------------------- | ------ | -------------------------------------------- |
| `/health`               | GET    | Service health check                         |
| `/api/v1/accounts/{id}` | GET    | User account details with balance            |
| `/api/v1/analytics`     | GET    | Aggregated dashboard analytics               |
| `/api/v1/fraud-scores`  | GET    | Recent fraud scores with transaction details |
| `/dashboard`            | GET    | Analytics dashboard UI                       |

**Example:**

```bash
# Get user account details
curl http://localhost:8080/api/v1/accounts/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11 | jq

# Get analytics data
curl http://localhost:8080/api/v1/analytics | jq

# View dashboard in browser
open http://localhost:8080/dashboard
```

---

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

---

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

### FraudScore

```go
type FraudScore struct {
    TransactionID    string
    RiskScore        string    // "low", "medium", "high"
    Status           string    // "clean", "suspicious", "flagged"
    TriggeredRules   []string  // e.g. ["high_frequency", "high_value"]
    Confidence       float64
    ProcessingTimeMs float64
    ScoredAt         time.Time
}
```

---

## 🔧 Technology Stack

| Component                      | Technology            | Version |
| ------------------------------ | --------------------- | ------- |
| **Language (Backend)**         | Go                    | 1.26    |
| **Language (Fraud Detection)** | Python                | 3.14    |
| **Fraud Detection Framework**  | FastAPI               | 0.136   |
| **Messaging (Local)**          | Kafka / RedPanda      | 7.6.1   |
| **Messaging (Production)**     | Aiven Kafka           | 4.1     |
| **Database (Local)**           | PostgreSQL            | 16      |
| **Database (Production)**      | AWS RDS PostgreSQL    | 16      |
| **Frontend**                   | Chart.js              | 4.4.0   |
| **Containerisation**           | Docker & Compose      | Latest  |
| **Container Registry**         | AWS ECR               | -       |
| **Hosting**                    | AWS Elastic Beanstalk | -       |

---

## 📝 Roadmap

- [x] Project structure and Docker setup
- [x] Implement Kafka producer/consumer
- [x] Build transaction generator (Producer)
- [x] Build transaction processor (Consumer)
- [x] Implement REST API endpoints
- [x] Real-time analytics dashboard
- [x] Add validation and business logic
- [x] Add comprehensive logging
- [x] Fraud detection microservice (Python/FastAPI)
- [x] Real-time fraud scoring dashboard panel
- [x] Cloud deployment (AWS Elastic Beanstalk + RDS + Aiven Kafka)
- [ ] CI/CD pipeline (GitHub Actions → ECR → Elastic Beanstalk)
- [ ] Service status indicators on dashboard
- [ ] Data variety improvements (more merchants, realistic fraud scenarios)
- [ ] Pagination on transaction and fraud feeds
- [ ] Metrics and observability (Prometheus/Grafana)
- [ ] Add unit and integration tests
- [ ] ECS migration for better container orchestration

---

## 🎓 Learning Goals

This project demonstrates understanding of:

- ✅ **Event-Driven Architecture** - Asynchronous communication via message queues
- ✅ **Microservices Design** - Loosely coupled, independently deployable services
- ✅ **Polyglot Systems** - Go and Python microservices sharing the same Kafka topics
- ✅ **Database Transactions** - ACID compliance, row-level locking, atomic updates
- ✅ **RESTful API Design** - Clean endpoints with proper status codes
- ✅ **Fraud Detection** - Rule-based scoring with duplicate, high-value and frequency detection
- ✅ **Data Visualization** - Real-time dashboards with Chart.js
- ✅ **Docker & Containerisation** - Multi-service orchestration
- ✅ **Cloud Deployment** - AWS Elastic Beanstalk, ECR, RDS
- ✅ **Go Best Practices** - Project structure, error handling, concurrency
- ✅ **Python Best Practices** - FastAPI, Pydantic, async patterns
- ✅ **SQL Optimisation** - Aggregated queries, proper indexing

---

## 📚 Resources

- [RedPanda Documentation](https://docs.redpanda.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [FastAPI Documentation](https://fastapi.tiangolo.com/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Aiven Kafka Documentation](https://aiven.io/docs/products/kafka)

---

## 🤝 Interview Talking Points

### Architecture Decisions

- **Why event-driven?** Decouples services, enables async processing, scales horizontally
- **Why Kafka over REST?** Buffering, replay capability, better for high throughput
- **Why PostgreSQL?** ACID compliance crucial for financial data, rich querying
- **Why Go?** Performance, simple concurrency model, excellent for backend services

### Fraud Detection Extension

- **Why a separate Python service?** Demonstrates polyglot microservices - the right tool for the right job. Python's ecosystem for data processing and ML is stronger than Go's
- **Why FastAPI?** Async by default, automatic API docs, Pydantic validation matches the type-safe approach used in Go
- **Why additive architecture?** Zero changes to existing Go services - new consumer on the same Kafka topic. Shows understanding of open/closed principle
- **Three detection rules:** Duplicate transactions, high-value threshold, high-frequency - each independently configurable and extensible

### Cloud Deployment

- **Why Elastic Beanstalk over raw EC2?** Managed platform handles scaling, health checks, rolling deployments. Focus on application not infrastructure
- **Why ECR?** Pre-built images mean the EC2 instance just pulls and runs - no build timeouts or memory pressure on small instances
- **Why Aiven Kafka?** Managed Kafka avoids running resource-heavy Zookeeper + Kafka containers alongside application services on a single instance
- **Why RDS?** Managed PostgreSQL with automatic backups, SSL, and failover. Financial data needs durability guarantees

### Trade-offs & Alternatives

- **Cassandra** for write-heavy workloads with eventual consistency
- **Redis** for caching frequently accessed data
- **WebSockets** for instant dashboard updates (vs 5-second polling)
- **gRPC** for inter-service communication in production

### Production Improvements

- Idempotency keys for exactly-once processing
- Dead-letter queues for failed messages
- Circuit breakers for resilience
- Rate limiting and authentication
- Distributed tracing (Jaeger)
- Metrics and alerting (Prometheus/Grafana)
- ML-based fraud scoring to replace rule-based approach
- CI/CD pipeline (GitHub Actions → ECR → Elastic Beanstalk)
- Service status indicators on dashboard
- Data variety improvements (more merchants, realistic fraud scenarios)
- Pagination on transaction and fraud feeds
- ECS migration for better container orchestration
- UI polish and pagination on dashboard scrolling sections (transaction feed, flagged transactions panel)

---

## 📚 Documentation

**Start here:** [docs/QUICKSTART.md](docs/QUICKSTART.md)

| Document                                                    | Description                                  |
| ----------------------------------------------------------- | -------------------------------------------- |
| [QUICKSTART.md](docs/QUICKSTART.md)                         | Get the system running in 5 minutes          |
| [ARCHITECTURE.md](docs/ARCHITECTURE.md)                     | System design deep dive, local vs production |
| [COMMAND-REFERENCE-LIST.md](docs/COMMAND-REFERENCE-LIST.md) | All make and docker commands                 |
| [DEPENDENCIES.md](docs/DEPENDENCIES.md)                     | Dependency decisions for Go and Python       |
| [IMPROVEMENTS.md](IMPROVEMENTS.md)                          | Planned enhancements and polish items        |

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2026 IbsYoussef

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## 👤 Author

**IbsYoussef**  
🔗 [GitHub](https://github.com/IbsYoussef)

---

<div align="center">

⭐ **Star this repo if you found it helpful!**

[⬆ Back to Top](#-flowbank)

</div>
