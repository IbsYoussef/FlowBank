# 🏦 FlowBank

A miniature event-driven transaction processor inspired by Monzo, built with Go, Kafka (RedPanda), and PostgreSQL.

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Kafka](https://img.shields.io/badge/Kafka-231F20?style=for-the-badge&logo=apache-kafka&logoColor=white)](https://kafka.apache.org/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
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
- [📚 Additional Documentation](#-additional-documentation)
- [📄 License](#-license)
- [👤 Author](#-author)

---

## 🎯 Project Overview

FlowBank demonstrates a distributed system for processing banking transactions using event-driven architecture. Inspired by Monzo's architecture, it showcases modern microservices patterns, real-time data processing, and interactive analytics.

### ✨ Key Features

- ⚡ **Event-Driven Architecture** - Asynchronous transaction processing via Kafka
- 📊 **Real-Time Analytics Dashboard** - Live visualization with Chart.js
- 🏗️ **Microservices Design** - Three independent Go services (Producer, Consumer, API)
- 🔒 **ACID Transactions** - PostgreSQL with row-level locking and atomic updates
- 🐳 **Fully Dockerized** - One-command deployment with Docker Compose
- 📈 **Production Patterns** - Proper error handling, logging, separation of concerns

---

## 🏗️ Architecture

```
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│   Producer   │────────>│     Kafka    │────────>│   Consumer   │
│ (Generates)  │ Publish │   (Stream)   │ Consume │ (Validates)  │
└──────────────┘         └──────────────┘         └──────┬───────┘
                                                         │
                                                         ▼
┌──────────────┐                                  ┌──────────────┐
│  Dashboard   │<────────── Query ────────────────│  PostgreSQL  │
│ (Analytics)  │                                  │  (Storage)   │
└──────────────┘                                  └──────────────┘
```

**Data Flow:**
1. **Producer** generates realistic banking transactions every 3 seconds
2. **Kafka** buffers and distributes events across partitions
3. **Consumer** validates, processes, and persists to database with overdraft protection
4. **API** serves transaction data via REST endpoints
5. **Dashboard** visualizes metrics in real-time with auto-refresh

---

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
├── web/
│   └── dashboard.html  # Analytics dashboard UI
├── docs/               # Additional documentation
├── go.mod
├── Makefile
└── README.md
```

---

## 🚀 Quick Start

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
- [Go 1.25+](https://go.dev/dl/) (optional, for local development)
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

**Access Points:**
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
```

#### 7. Monitor Live Transactions

**Option A: Dashboard (Recommended)**
```bash
open http://localhost:8080/dashboard
# Watch the live transaction feed update every 5 seconds
```

**Option B: Producer Logs**
```bash
make logs-producer
# See transactions being generated every 3 seconds
```

**Option C: Consumer Logs**
```bash
make logs-consumer
# See transactions being processed and marked as COMPLETED/FAILED
```

---

## 📊 Analytics Dashboard

![Dashboard Preview](/assets/img.png)

**Access at:** http://localhost:8080/dashboard

### Features

- 📈 **4 Interactive Charts**
   - Transaction Type Distribution (Pie Chart)
   - Money Flow Analysis (Bar Chart)
   - Top 5 Merchants (Horizontal Bar)
   - Hourly Transaction Volume (Line Chart)
- 🔴 **Live Transaction Feed** - Real-time updates every 5 seconds
- 💰 **Key Metrics** - Total transactions, money in/out, average amount
- 📱 **Responsive Design** - Works on desktop, tablet, mobile

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Service health check |
| `/api/v1/accounts/{id}` | GET | User account details with balance |
| `/api/v1/analytics` | GET | Aggregated dashboard analytics |
| `/dashboard` | GET | Analytics dashboard UI |

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

---

## 🔧 Technology Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| **Language** | Go | 1.25 |
| **Messaging** | Kafka (Confluent) | 7.6.1 |
| **Database** | PostgreSQL | 16 |
| **Frontend** | Chart.js | 4.4.0 |
| **Containerization** | Docker & Compose | Latest |

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
- [ ] Add metrics and observability (Prometheus/Grafana)
- [ ] Add unit and integration tests
- [ ] Add CI/CD pipeline
- [ ] Cloud deployment (AWS ECS)

---

## 🎓 Learning Goals

This project demonstrates understanding of:

- ✅ **Event-Driven Architecture** - Asynchronous communication via message queues
- ✅ **Microservices Design** - Loosely coupled, independently deployable services
- ✅ **Database Transactions** - ACID compliance, row-level locking, atomic updates
- ✅ **RESTful API Design** - Clean endpoints with proper status codes
- ✅ **Data Visualization** - Real-time dashboards with Chart.js
- ✅ **Docker & Containerization** - Multi-service orchestration
- ✅ **Go Best Practices** - Project structure, error handling, concurrency
- ✅ **SQL Optimization** - Aggregated queries, proper indexing

---

## 📚 Resources

- [RedPanda Documentation](https://docs.redpanda.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

---

## 🤝 Interview Talking Points

### Architecture Decisions

- **Why event-driven?** Decouples services, enables async processing, scales horizontally
- **Why Kafka over REST?** Buffering, replay capability, better for high throughput
- **Why PostgreSQL?** ACID compliance crucial for financial data, rich querying
- **Why Go?** Performance, simple concurrency model, excellent for backend services

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

---

## 📚 Additional Documentation

- [Architecture Deep Dive](docs/ARCHITECTURE.md)
- [Quick Start Guide](docs/QUICKSTART.md)
- [Command Reference](docs/COMMAND-REFERENCE-LIST.md)
- [Dependencies Guide](docs/DEPENDENCIES.md)

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