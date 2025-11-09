# FlowBank Quick Start Guide

This guide will help you get FlowBank up and running in under 5 minutes.

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
# Expected: go version go1.22.x or higher
```

## 🚀 Step 1: Clone and Navigate

```bash
# If you haven't already, navigate to the project
cd flowbank

# Verify structure
ls -la
# You should see: cmd/, internal/, deploy/, Makefile, README.md, etc.
```

## 🏗️ Step 2: Start the Infrastructure

Start all services with a single command:

```bash
make up
```

This will:

1. Start RedPanda (Kafka) on ports 19092, 8090
2. Start PostgreSQL on port 5432
3. Build and start Producer service
4. Build and start Consumer service
5. Build and start API service

**Expected output**:

```
✅ FlowBank services started!
   - API:              http://localhost:8080/health
   - RedPanda Console: http://localhost:8090
   - PostgreSQL:       localhost:5432
```

## 🔍 Step 3: Verify Services

Check that all services are running:

```bash
make ps
```

**Expected output**:

```
NAME                    IMAGE                  STATUS
flowbank-api            flowbank-api           Up 30 seconds
flowbank-consumer       flowbank-consumer      Up 30 seconds
flowbank-producer       flowbank-producer      Up 30 seconds
flowbank-postgres       postgres:16-alpine     Up 35 seconds (healthy)
flowbank-redpanda       redpanda:v24.2.9       Up 35 seconds (healthy)
flowbank-console        console:v2.7.2         Up 30 seconds
```

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

-- Exit
\q
```

### Test 3: Monitor Kafka Topics

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

### Test 4: Watch Service Logs

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
Producer heartbeat: 2025-11-06T10:30:00Z
Producer heartbeat: 2025-11-06T10:30:05Z
```

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
```

Watch as transactions flow from Producer → Kafka → Consumer → Database!

## 🎯 Step 6: Test API Endpoints (Coming Soon)

Once you implement the API endpoints, test them:

```bash
# Get all users
curl http://localhost:8080/users | jq

# Get specific user
curl http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000 | jq

# Get user transactions
curl http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/transactions | jq
```

## 🛑 Stopping Services

When done, stop all services:

```bash
make down
```

To clean up everything including volumes:

```bash
make clean
```

## 🐛 Troubleshooting

### Services won't start

```bash
# Check Docker is running
docker ps

# Check for port conflicts
lsof -i :8080  # API
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

### Build errors

```bash
# Clean and rebuild
make clean
make build
make up
```

## 📝 Next Steps

Now that your environment is running:

1. **Implement Producer Logic** - Generate realistic fake transactions
2. **Implement Consumer Logic** - Process and validate transactions
3. **Implement API Endpoints** - Query users and transactions
4. **Add Tests** - Unit and integration tests
5. **Add Observability** - Logging, metrics, monitoring

## 🎓 Learning Resources

- Review [ARCHITECTURE.md](./ARCHITECTURE.md) for system design details
- Check [README.md](./README.md) for project overview
- Explore RedPanda Console: http://localhost:8090
- Read the code in `cmd/`, `internal/` directories

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

**Need help?** Check the [README.md](./README.md) or [ARCHITECTURE.md](./ARCHITECTURE.md) for more details.

Good luck with your internship interview! 🚀
