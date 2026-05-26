# FlowBank - Dependencies

## Go Dependencies

### 1. Kafka Client

FlowBank uses `segmentio/kafka-go` for its Go services. Below is also a reference to `confluent-kafka-go/v2` as an alternative worth knowing for interviews.

#### What FlowBank Uses - segmentio/kafka-go

```go
require github.com/segmentio/kafka-go
```

**Why?**

- Pure Go implementation, no C dependencies
- Simple API, easy to understand
- Supports SASL_SSL for production Kafka connections (Aiven)
- Good performance for this use case

**SASL_SSL support**: Custom TLS config in `internal/kafka/tls.go` handles both local (PLAINTEXT) and production (SASL_SSL) connections automatically based on environment variables.

#### Alternative - Confluent's kafka-go v2 (RECOMMENDED for production scale)

```go
require github.com/confluentinc/confluent-kafka-go/v2 v2.3.0
```

**Why?**

- Official Confluent client, most mature
- High performance (uses librdkafka)
- Excellent documentation
- Production-ready
- Works seamlessly with RedPanda

**Alternative**: `github.com/IBM/sarama` (pure Go, no C dependencies)

**Usage Example**:

```go
import (
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

producer, err := kafka.NewProducer(&kafka.ConfigMap{
    "bootstrap.servers": "localhost:19092",
})
```

---

### 2. PostgreSQL Driver - pgx v5 (RECOMMENDED)

```go
require github.com/jackc/pgx/v5 v5.5.0
```

**Why?**

- Fastest PostgreSQL driver for Go
- Native Go implementation
- Better connection pooling via `pgxpool`
- Rich feature set (prepared statements, COPY, etc.)
- Active maintenance
- Supports SSL connections required by AWS RDS

**Alternative**: `github.com/lib/pq` (simpler, but slower)

**Usage Example**:

```go
import (
    "github.com/jackc/pgx/v5/pgxpool"
)

pool, err := pgxpool.New(context.Background(),
    "postgresql://user:pass@localhost:5432/dbname")
```

---

### 3. HTTP Router - Chi or Gorilla Mux (RECOMMENDED)

```go
require github.com/go-chi/chi/v5 v5.0.11
```

**Why Chi?**

- Lightweight, idiomatic
- Great middleware support
- Context-based routing
- No external dependencies

**Alternative**: Standard library `net/http` (fine for MVP - what FlowBank uses)

**Usage Example**:

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

r := chi.NewRouter()
r.Use(middleware.Logger)
r.Get("/users/{id}", getUserHandler)
```

---

### 4. Configuration - Viper (RECOMMENDED)

```go
require github.com/spf13/viper v1.18.2
```

**Why?**

- Handles env vars, files, flags
- Hot reload configuration
- Industry standard

**Alternative**: `github.com/kelseyhightower/envconfig` (simpler)

> FlowBank uses `os.Getenv` directly - Viper would be the next step for a more complex config setup.

---

### 5. Logging - Zap or Zerolog (RECOMMENDED)

```go
require go.uber.org/zap v1.26.0
```

**Why Zap?**

- Extremely fast structured logging
- Used by many production systems
- Rich ecosystem

**Alternative**: `github.com/rs/zerolog` (zero allocation)

**Usage Example**:

```go
import "go.uber.org/zap"

logger, _ := zap.NewProduction()
defer logger.Sync()
logger.Info("transaction processed",
    zap.String("id", txID),
    zap.Int64("amount", amount))
```

> FlowBank uses the standard `log` package - Zap would be the production upgrade.

---

### 6. UUID - Google UUID (RECOMMENDED)

```go
require github.com/google/uuid v1.6.0
```

**Why?**

- Simple, reliable
- Good performance
- Well-maintained

---

### 7. JSON Handling - Standard Library (RECOMMENDED)

Use `encoding/json` from stdlib - it's sufficient for most cases.

**Alternative for performance**: `github.com/goccy/go-json` (3-4x faster)

---

### 8. Error Handling - pkg/errors (OPTIONAL)

```go
require github.com/pkg/errors v0.9.1
```

**Why?**

- Stack trace support
- Error wrapping

**Note**: Go 1.13+ has built-in error wrapping with `fmt.Errorf("%w", err)`

---

## Testing & Development (Go)

### Testify (RECOMMENDED)

```go
require github.com/stretchr/testify v1.8.4
```

**Why?**

- Rich assertion library
- Mocking support
- Suite support

**Usage**:

```go
import "github.com/stretchr/testify/assert"

assert.Equal(t, expected, actual)
assert.NoError(t, err)
```

### Gomock (RECOMMENDED)

```go
require github.com/golang/mock v1.6.0
```

**Why?**

- Official Go mocking framework
- Code generation
- Interface-based

### Testcontainers (ADVANCED)

```go
require github.com/testcontainers/testcontainers-go v0.27.0
```

**Why?**

- Spin up real Kafka, PostgreSQL in tests
- Better than mocks for integration tests

---

## Complete go.mod Reference

```go
module github.com/yourusername/flowbank

go 1.26

require (
    github.com/segmentio/kafka-go latest
    github.com/google/uuid v1.6.0
    github.com/jackc/pgx/v5 v5.5.0
    github.com/spf13/viper v1.18.2
    go.uber.org/zap v1.26.0
)

require (
    github.com/stretchr/testify v1.8.4
)
```

---

## Installation Commands

```bash
# Core dependencies
go get github.com/segmentio/kafka-go
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/go-chi/chi/v5
go get github.com/google/uuid
go get github.com/spf13/viper
go get go.uber.org/zap

# Test dependencies
go get github.com/stretchr/testify

# Tidy up
go mod tidy
```

---

## Alternative Stack (Simpler for MVP)

```go
module github.com/yourusername/flowbank

go 1.26

require (
    github.com/IBM/sarama v1.42.1    // Pure Go Kafka
    github.com/lib/pq v1.10.9        // Simple PostgreSQL
    github.com/google/uuid v1.6.0    // UUID generation
)
```

**Pros**: Fewer dependencies, pure Go, easier to understand
**Cons**: Slower performance, less features

---

## Recommended for Interviews

1. **Start Simple**: Use stdlib as much as possible
2. **Add Kafka Client**: `segmentio/kafka-go` or `confluent-kafka-go/v2`
3. **Add PostgreSQL**: `pgx/v5` (noticeable performance boost)
4. **Add UUID**: `google/uuid` (already using in models)
5. **Add Testing**: `testify` (shows you care about testing)

**Skip for now** (add later):

- Viper (use `os.Getenv` for MVP)
- Zap (use `log` stdlib for MVP)
- Chi (use `net/http` for MVP)

---

## Package Selection Criteria

1. **Maintenance**: Last commit date, active issues
2. **Performance**: Benchmarks, production usage
3. **Documentation**: Examples, godoc quality
4. **Community**: Stars, forks, used by whom
5. **Dependencies**: Avoid dependency hell

## Where to Find Packages

- **Awesome Go**: https://github.com/avelino/awesome-go
- **Go Packages**: https://pkg.go.dev/
- **GitHub Topics**: https://github.com/topics/golang

---

## Performance Tips

1. **Connection Pooling**: Always use connection pools (pgxpool, kafka producer reuse)
2. **Batch Operations**: Batch inserts/reads when possible
3. **Prepared Statements**: Use for repeated queries
4. **Context Timeouts**: Always set deadlines
5. **Graceful Shutdown**: Handle SIGTERM properly

---

## Python Dependencies (Fraud Detection Service)

### FastAPI

```
fastapi[standard]==0.136.1
```

**Why?**

- Async by default - integrates cleanly with Kafka consumption
- Automatic API documentation at `/docs`
- Built on Pydantic for request/response validation
- The `[standard]` extra includes uvicorn and the FastAPI CLI

### Pydantic / pydantic-settings

```
pydantic==2.13.4
pydantic-settings==2.14.1
```

**Why?**

- Type-safe data validation for Kafka message deserialization
- `pydantic-settings` reads environment variables cleanly into typed config classes
- Pydantic v2 only - v1 support was dropped

### asyncpg

```
asyncpg==0.31.0
```

**Why?**

- Async PostgreSQL driver - required for FastAPI's async model
- Fastest PostgreSQL driver for Python
- Connection pooling support
- Handles SSL connections for AWS RDS

### confluent-kafka

```
confluent-kafka==2.14.0
```

**Why?**

- Official Confluent Kafka client for Python
- Uses librdkafka underneath - most mature and battle-tested
- Supports SASL_SSL required for Aiven Kafka

**SASL_SSL support**: The CA certificate is passed as a base64-encoded environment variable, decoded at startup, and written to a temp file for the confluent-kafka client.

### Complete requirements.txt

```
# FastAPI
fastapi[standard]==0.136.1

# Pydantic dependencies
pydantic==2.13.4
pydantic-settings==2.14.1

# PostgreSQL driver
asyncpg==0.31.0

# Kafka
confluent-kafka==2.14.0
```

---

## Interview Talking Points

**Q**: Why confluent-kafka-go over sarama for Go?

**A**: "While sarama is pure Go and easier to deploy, confluent-kafka-go uses librdkafka underneath which is battle-tested and performs better. For a transaction processing system, performance and reliability are critical. FlowBank uses segmentio/kafka-go which is a good middle ground - pure Go but well maintained."

**Q**: Why pgx over lib/pq?

**A**: "pgx is faster (2-3x in benchmarks), has better connection pooling, and supports PostgreSQL-specific features like COPY. It's actively maintained and widely used in production. For financial transactions where we process thousands of events per second, the performance gain matters."

**Q**: Why confluent-kafka for Python over kafka-python?

**A**: "confluent-kafka uses librdkafka which is the most battle-tested Kafka client across all languages. For a production-grade fraud detection service, reliability matters more than a pure-Python implementation."

**Q**: Why asyncpg over psycopg3 for Python?

**A**: "asyncpg has a cleaner async API and better performance. psycopg3 is a valid alternative but asyncpg is the established choice for FastAPI + PostgreSQL."

---

## Next Steps

1. Review this document
2. Install Go dependencies: `go mod tidy`
3. Install Python dependencies: `pip install -r fraud-detection/requirements.txt`
4. Test with real Kafka and PostgreSQL via `make up`
