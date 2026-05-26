# FlowBank Quick Start Guide

Get FlowBank up and running in under 5 minutes.

> **Live Demo**: [flowbank-prod.eba-fcdmxpas.eu-west-2.elasticbeanstalk.com/dashboard](http://flowbank-prod.eba-fcdmxpas.eu-west-2.elasticbeanstalk.com/dashboard)

---

## ✅ Prerequisites Check

Before starting, ensure you have:

```bash
# Check Docker
docker --version
# Expected: Docker version 20.x or higher

# Check Docker Compose
docker compose version
# Expected: Docker Compose version v2.x or higher

# Check Go (optional, for local development)
go version
# Expected: go version go1.26.x or higher
```

---

## 🚀 Step 1: Clone and Navigate

```bash
# Clone the repository
git clone https://github.com/IbsYoussef/FlowBank.git
cd FlowBank

# Verify structure
ls -la
# You should see: cmd/, internal/, deploy/, fraud-detection/, Makefile, README.md, etc.
```

---

## 🏗️ Step 2: Start the Infrastructure

Start all services with a single command:

```bash
make up
```

This will:

1. Start RedPanda (Kafka) on ports 19092, 9092
2. Start Zookeeper on port 2181
3. Start PostgreSQL on port 5432
4. Build and start Producer service
5. Build and start Consumer service
6. Build and start API service on port 8080
7. Build and start Fraud Detection service on port 8000

**Expected output**:

```
✅ FlowBank services started!
   - API:              http://localhost:8080/health
   - RedPanda Console: http://localhost:8090
   - PostgreSQL:       localhost:5432
```

---

## 🔍 Step 3: Verify Services

Check that all services are running:

```bash
make ps
```

**Expected output**:

```
NAME                          IMAGE                      STATUS
flowbank-api                  flowbank-api               Up 30 seconds
flowbank-consumer             flowbank-consumer          Up 30 seconds
flowbank-producer             flowbank-producer          Up 30 seconds
flowbank-fraud-detection      flowbank-fraud-detection   Up 30 seconds
flowbank-postgres             postgres:16-alpine         Up 35 seconds (healthy)
flowbank-redpanda             cp-kafka:7.6.1             Up 35 seconds (healthy)
flowbank-zookeeper            cp-zookeeper:7.6.1         Up 35 seconds (healthy)
```

---

## 🧪 Step 4: Test the Pipeline

### Test 1: API Health Check

```bash
curl http://localhost:8080/health
```

**Expected**:

```json
{ "status": "healthy", "service": "flowbank-api" }
```

### Test 2: Check Database

```bash
make db-shell
```

Then run:

```sql
-- View seed users
SELECT * FROM users;

-- View recent transactions (may be empty initially)
SELECT * FROM transactions ORDER BY created_at DESC LIMIT 5;

-- View fraud scores
SELECT transaction_id, risk_score, status, triggered_rules
FROM fraud_scores ORDER BY scored_at DESC LIMIT 5;

-- Exit
\q
```

### Test 3: Check Fraud Scores API

```bash
curl http://localhost:8080/api/v1/fraud-scores | jq
```

### Test 4: Monitor Kafka Topics

**Option A: Command Line**

```bash
make kafka-topics
```

**Expected**: Should show `transactions` topic

**Option B: Web UI**

1. Open browser: http://localhost:8090
2. Navigate to "Topics"
3. Click on "transactions"
4. You should see messages flowing in real-time!

### Test 5: Watch Service Logs

View logs from all services:

```bash
make logs
```

Or watch specific services:

```bash
# Producer logs
make logs-producer

# Consumer logs
make logs-consumer

# API logs
make logs-api
```

**Expected producer output**:

```
PUBLISHED: TxID abc-123 | User a0eebc99... | Type debit | Amount 5000 cents
PUBLISHED: TxID def-456 | User b1eebc99... | Type credit | Amount 12500 cents
```

---

## 📊 Step 5: Visualize the Data Flow

Open three terminal windows:

**Terminal 1 - Producer logs**:

```bash
make logs-producer
```

**Terminal 2 - Consumer logs**:

```bash
make logs-consumer
```

**Terminal 3 - Database queries**:

```bash
make db-shell
# Then repeatedly run:
SELECT COUNT(*) FROM transactions;
SELECT COUNT(*) FROM fraud_scores;
```

Watch as transactions flow from Producer → Kafka → Consumer → Database, and fraud scores appear in real-time!

---

## 🎯 Step 6: Test API Endpoints

```bash
# Get user account details
curl http://localhost:8080/api/v1/accounts/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11 | jq

# Get analytics data
curl http://localhost:8080/api/v1/analytics | jq

# Get fraud scores
curl http://localhost:8080/api/v1/fraud-scores | jq

# View dashboard in browser
open http://localhost:8080/dashboard
```

```bash
# Coming soon: curl http://localhost:8080/users | jq
# Coming soon: curl http://localhost:8080/users/{user_id}/transactions | jq
```

---

## 🛑 Stopping Services

When done, stop all services:

```bash
make down
```

To clean up everything including volumes:

```bash
make clean
```

---

## 🐛 Troubleshooting

### Services won't start

```bash
# Check Docker is running
docker ps

# Check for port conflicts
lsof -i :8080  # API
lsof -i :8000  # Fraud Detection
lsof -i :19092 # Kafka
lsof -i :5432  # PostgreSQL

# View detailed logs
docker compose -f deploy/docker-compose.yml logs
```

### Can't connect to database

```bash
# Verify PostgreSQL is healthy
docker ps | grep postgres

# Check connection
docker exec -it flowbank-postgres psql -U flowbank_user -d flowbank -c "SELECT 1;"
```

### Kafka not receiving messages

```bash
# Check RedPanda status
docker exec -it flowbank-redpanda rpk cluster health

# List topics
docker exec -it flowbank-redpanda rpk topic list

# Consume messages manually
docker exec -it flowbank-redpanda rpk topic consume transactions
```

### Fraud scores not appearing

The fraud detection service has a 3-second delay to avoid a race condition with the Go consumer. Wait a few seconds after transactions appear in the feed before checking fraud scores.

### Build errors

```bash
# Clean and rebuild
make clean
make build
make up
```

---

## 📝 Next Steps

Now that your environment is running:

1. **Explore the Dashboard** - Open http://localhost:8080/dashboard
2. **Review the Architecture** - See [ARCHITECTURE.md](./ARCHITECTURE.md)
3. **Add Tests** - Unit and integration tests
4. **Add Observability** - Logging, metrics, monitoring
5. **Try the Live Demo** - [flowbank-prod.eba-fcdmxpas.eu-west-2.elasticbeanstalk.com/dashboard](http://flowbank-prod.eba-fcdmxpas.eu-west-2.elasticbeanstalk.com/dashboard)

---

## 🎓 Learning Resources

- Review [ARCHITECTURE.md](./ARCHITECTURE.md) for system design details
- Check [README.md](../README.md) for project overview
- Explore RedPanda Console: http://localhost:8090
- Read the code in `cmd/`, `internal/`, `fraud-detection/` directories

---

## ✨ Quick Commands Reference

| Command             | Description                   |
| ------------------- | ----------------------------- |
| `make up`           | Start all services            |
| `make down`         | Stop all services             |
| `make restart`      | Restart services              |
| `make logs`         | View all logs                 |
| `make ps`           | Show running containers       |
| `make db-shell`     | PostgreSQL shell              |
| `make kafka-topics` | List Kafka topics             |
| `make clean`        | Remove all containers/volumes |

---

**Need help?** Check the [README.md](../README.md) or [ARCHITECTURE.md](./ARCHITECTURE.md) for more details.
