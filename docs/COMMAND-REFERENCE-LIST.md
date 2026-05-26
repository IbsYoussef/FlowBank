# FlowBank Quick Commands

## Docker Compose Commands

| Command                                                           | Purpose                           |
| ----------------------------------------------------------------- | --------------------------------- |
| `docker compose -f deploy/docker-compose.yml build`               | Build images (after code changes) |
| `docker compose -f deploy/docker-compose.yml up -d`               | Start all services                |
| `docker compose -f deploy/docker-compose.yml up -d --no-deps api` | Restart single service            |
| `docker compose -f deploy/docker-compose.yml ps`                  | Check service status              |
| `docker compose -f deploy/docker-compose.yml logs -f`             | View real-time logs               |
| `docker compose -f deploy/docker-compose.yml down`                | Stop services                     |
| `docker compose -f deploy/docker-compose.yml down -v`             | Clean reset (deletes volumes)     |

---

## Makefile Shortcuts

| Command         | Equivalent                   | Use Case                     |
| --------------- | ---------------------------- | ---------------------------- |
| `make build`    | `docker compose ... build`   | Compile new code             |
| `make up`       | `docker compose ... up -d`   | Start application            |
| `make down`     | `docker compose ... down`    | Stop containers              |
| `make logs`     | `docker compose ... logs -f` | Check transaction flow       |
| `make clean`    | `docker compose ... down -v` | Full reset (clears database) |
| `make db-shell` | `docker exec ... psql`       | Inspect database             |

---

## Common Workflows

**After code changes:**

```bash
make build && make down && make up
```

**Check everything is working:**

```bash
make ps && make logs
```

**Full reset:**

```bash
make clean && make build && make up
```

---

## Fraud Detection Service

| Command                                                            | Purpose                             |
| ------------------------------------------------------------------ | ----------------------------------- |
| `curl http://localhost:8000/health`                                | Check fraud service health directly |
| `curl http://localhost:8080/api/v1/fraud-scores`                   | View fraud scores via Go API        |
| `curl http://localhost:8080/api/v1/fraud-scores?status=flagged`    | View only flagged transactions      |
| `docker compose -f deploy/docker-compose.yml logs fraud-detection` | Fraud detection logs                |

---

## Database Queries

```bash
make db-shell
```

```sql
-- View all users and balances
SELECT user_id, user_name, current_balance FROM users;

-- View recent transactions
SELECT transaction_id, amount, transaction_type, status, created_at
FROM transactions ORDER BY created_at DESC LIMIT 10;

-- Check fraud scores
SELECT transaction_id, risk_score, status, triggered_rules, processing_time_ms
FROM fraud_scores ORDER BY scored_at DESC LIMIT 10;

-- Count by fraud status
SELECT status, COUNT(*) FROM fraud_scores GROUP BY status;

-- Transaction stats
SELECT status, COUNT(*) FROM transactions GROUP BY status;
```

---

## AWS Deployment Commands

```bash
# Deploy to Elastic Beanstalk
eb deploy flowbank-prod

# Check environment status
eb status flowbank-prod

# View deployment events
eb events flowbank-prod

# SSH into EC2 instance
eb ssh flowbank-prod

# Set environment variables
eb setenv KEY=value --region eu-west-2

# View application logs
eb logs flowbank-prod
```

```bash
# Authenticate Docker with ECR
aws ecr get-login-password --region eu-west-2 | docker login --username AWS --password-stdin 084576276551.dkr.ecr.eu-west-2.amazonaws.com

# Build and push image to ECR
docker build -f deploy/Dockerfile.api -t 084576276551.dkr.ecr.eu-west-2.amazonaws.com/flowbank-api:latest .
docker push 084576276551.dkr.ecr.eu-west-2.amazonaws.com/flowbank-api:latest
```
