# feat/fraud-detection

## What is this branch?

This branch extends FlowBank with a Python-based fraud detection microservice. If you're skimming this, here's the one-liner: **transactions flow through Kafka, this service picks them up, scores them for fraud risk, and surfaces the results on the existing dashboard in real time.**

No existing Go services are modified. This is purely additive.

---

## Why this feature?

FlowBank already demonstrates event-driven architecture, microservices design, and real-time data processing. Fraud detection was the natural next layer - it's a real problem every fintech system solves, and it adds a second language (Python) to the stack, which demonstrates the ability to work in polyglot distributed systems.

The decision to build it as a separate Python microservice rather than extending the existing Go services was deliberate. It keeps services loosely coupled, shows that the Kafka topic can support multiple independent consumers, and adds a FastAPI project to the portfolio that is directly relevant to fintech engineering roles.

---

## What this service does

The fraud detection service consumes transaction events from the same Kafka topic as the existing Go consumer. For each transaction it applies three detection rules:

- **Duplicate transactions** - flags the same account making the same amount transaction within 60 seconds
- **High value threshold** - flags any transaction above $10,000
- **High frequency** - flags any account producing more than 5 transactions within a 60-second window

Each transaction is scored as low, medium, or high risk and given a status of clean, suspicious, or flagged depending on how many rules trigger. Results are written to a new `fraud_scores` table in the shared PostgreSQL instance and surfaced on the existing dashboard.

---

## Architecture

```
Kafka Topic (transactions)
        │
        ├──> Go Consumer (existing FlowBank service, unchanged)
        │
        └──> Python Fraud Detection Consumer (this service)
                  │
                  ▼
          Score Transaction
          ├─ Rule 1: Duplicate check
          ├─ Rule 2: High value check
          └─ Rule 3: High frequency check
                  │
                  ▼
          Write to fraud_scores table
                  │
                  ▼
          Dashboard reads and displays risk scores
```

### Folder structure

```
fraud-detection/
├── app/
│   ├── main.py         # FastAPI entry point, routes, startup/shutdown
│   ├── models.py       # Pydantic schemas for transactions and fraud scores
│   ├── detector.py     # Fraud scoring rules, single responsibility
│   ├── consumer.py     # Kafka consumer logic only
│   └── database.py     # PostgreSQL connection and queries
├── requirements.txt
└── .env.example
```

Each file has a single responsibility. This makes the service easy to test, easy to extend, and easy to reason about.

---

## Tech stack

| Layer         | Technology                | Notes                                     |
| ------------- | ------------------------- | ----------------------------------------- |
| Language      | Python 3.12               |                                           |
| Framework     | FastAPI                   | Async, type-safe, auto-generates API docs |
| Message queue | Kafka via confluent-kafka | Same topic as existing Go services        |
| Database      | PostgreSQL via asyncpg    | Shared instance, new fraud_scores table   |
| Validation    | Pydantic                  | Request/response schemas                  |
| Config        | pydantic-settings         | Reads from .env file                      |
| Container     | Docker                    | Added to existing docker-compose.yml      |

---

## Build plan

| Phase    | Goal                                                                         |
| -------- | ---------------------------------------------------------------------------- |
| Days 1-2 | FastAPI fundamentals, project setup, running app with basic endpoints        |
| Days 3-4 | Kafka consumer connected and receiving live transaction events               |
| Days 5-6 | Fraud detection rules implemented, scores written to PostgreSQL              |
| Day 7    | Docker integration, full system running with one command                     |
| Day 8    | Dashboard updated, deployed, demo seeded with clean and flagged transactions |

---

## Dashboard changes

The existing transaction feed will get two additions:

- A risk score column and status badge (green/amber/red) on each transaction row
- A separate flagged transactions panel showing only high-risk entries

No new URL or deployment is needed. The same dashboard endpoint displays the updated data.

---

## Current status

- [x] Project setup, virtual environment, dependencies
- [x] FastAPI app running with root and health endpoints
- [x] Pydantic models
- [x] Kafka consumer
- [x] Fraud detection logic
- [x] PostgreSQL integration
- [x] Docker integration
- [x] Dashboard updates
- [x] Deployment

Live URL: http://flowbank-prod.eba-fcdmxpas.eu-west-2.elasticbeanstalk.com/dashboard
