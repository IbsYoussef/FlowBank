# FlowBank - Recommended Go Dependencies (2025)

## 🎯 Core Dependencies

### 1. Kafka Client - **Confluent's kafka-go v2** (RECOMMENDED)

```go
// go.mod
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

### 2. PostgreSQL Driver - **pgx v5** (RECOMMENDED)

```go
// go.mod
require github.com/jackc/pgx/v5 v5.5.0
```

**Why?**

- Fastest PostgreSQL driver for Go
- Native Go implementation
- Better connection pooling
- Rich feature set (prepared statements, COPY, etc.)
- Active maintenance

**Alternative**: `github.com/lib/pq` (simpler, but slower)

**Usage Example**:

```go
import (
    "github.com/jackc/pgx/v5/pgxpool"
)

pool, err := pgxpool.New(context.Background(),
    "postgresql://user:pass@localhost:5432/dbname")
```

### 3. HTTP Router - **Chi** or **Gorilla Mux** (RECOMMENDED)

```go
// go.mod
require github.com/go-chi/chi/v5 v5.0.11
```

**Why Chi?**

- Lightweight, idiomatic
- Great middleware support
- Context-based routing
- No external dependencies

**Alternative**: Standard library `net/http` (fine for MVP!)

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

### 4. Configuration - **Viper** (RECOMMENDED)

```go
// go.mod
require github.com/spf13/viper v1.18.2
```

**Why?**

- Handles env vars, files, flags
- Hot reload configuration
- Industry standard

**Alternative**: `github.com/kelseyhightower/envconfig` (simpler)

### 5. Logging - **Zap** or **Zerolog** (RECOMMENDED)

```go
// go.mod
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

## 🛠️ Utility Dependencies

### 6. UUID - **Google UUID** (RECOMMENDED)

```go
// go.mod
require github.com/google/uuid v1.6.0
```

**Why?**

- Simple, reliable
- Good performance
- Well-maintained

### 7. JSON Handling - **Standard Library** (RECOMMENDED)

Use `encoding/json` from stdlib - it's sufficient for most cases.

**Alternative for performance**: `github.com/goccy/go-json` (3-4x faster)

### 8. Error Handling - **pkg/errors** (OPTIONAL)

```go
// go.mod
require github.com/pkg/errors v0.9.1
```

**Why?**

- Stack trace support
- Error wrapping

**Note**: Go 1.13+ has built-in error wrapping with `fmt.Errorf("%w", err)`

## 🧪 Testing & Development

### 9. Testing - **Testify** (RECOMMENDED)

```go
// go.mod
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

### 10. Mocking - **Gomock** (RECOMMENDED)

```go
// go.mod
require github.com/golang/mock v1.6.0
```

**Why?**

- Official Go mocking framework
- Code generation
- Interface-based

### 11. Integration Testing - **Testcontainers** (ADVANCED)

```go
// go.mod
require github.com/testcontainers/testcontainers-go v0.27.0
```

**Why?**

- Spin up real Kafka, PostgreSQL in tests
- Better than mocks for integration tests

## 📦 Complete go.mod for FlowBank

Here's the recommended `go.mod` for your project:

```go
module github.com/yourusername/flowbank

go 1.22

require (
	github.com/confluentinc/confluent-kafka-go/v2 v2.3.0
	github.com/go-chi/chi/v5 v5.0.11
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.5.0
	github.com/spf13/viper v1.18.2
	go.uber.org/zap v1.26.0
)

require (
	// Test dependencies
	github.com/stretchr/testify v1.8.4
)
```

## 🎯 Installation Commands

Install all dependencies:

```bash
# Core dependencies
go get github.com/confluentinc/confluent-kafka-go/v2/kafka
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

## 🔄 Alternative Stack (Simpler for MVP)

If you want to minimize dependencies for the MVP:

```go
module github.com/yourusername/flowbank

go 1.22

require (
	github.com/IBM/sarama v1.42.1              // Pure Go Kafka
	github.com/lib/pq v1.10.9                  // Simple PostgreSQL
	github.com/google/uuid v1.6.0              // UUID generation
)
```

**Pros**:

- Fewer dependencies
- Pure Go (no C dependencies)
- Easier to understand

**Cons**:

- Slower performance
- Less features

## 🚀 Recommended for Interviews

For your 2-week timeline, I recommend:

1. **Start Simple**: Use stdlib as much as possible
2. **Add Kafka Client**: `confluent-kafka-go/v2` (must-have)
3. **Add PostgreSQL**: `pgx/v5` (noticeable performance boost)
4. **Add UUID**: `google/uuid` (already using in models)
5. **Add Testing**: `testify` (shows you care about testing)

**Skip for now** (add later):

- Viper (use `os.Getenv` for MVP)
- Zap (use `log` stdlib for MVP)
- Chi (use `net/http` for MVP)

## 📝 Package Selection Criteria

When choosing packages, consider:

1. **Maintenance**: Last commit date, active issues
2. **Performance**: Benchmarks, production usage
3. **Documentation**: Examples, godoc quality
4. **Community**: Stars, forks, used by whom
5. **Dependencies**: Avoid dependency hell

## 🔍 Where to Find Packages

- **Awesome Go**: https://github.com/avelino/awesome-go
- **Go Packages**: https://pkg.go.dev/
- **GitHub Topics**: https://github.com/topics/golang

## ⚡ Performance Tips

1. **Connection Pooling**: Always use connection pools (pgxpool, kafka producer reuse)
2. **Batch Operations**: Batch inserts/reads when possible
3. **Prepared Statements**: Use for repeated queries
4. **Context Timeouts**: Always set deadlines
5. **Graceful Shutdown**: Handle SIGTERM properly

## 🎓 Interview Talking Points

**Q**: Why confluent-kafka-go over sarama?

**A**: "While sarama is pure Go and easier to deploy, confluent-kafka-go uses librdkafka underneath which is battle-tested and performs better. For a transaction processing system, performance and reliability are critical. However, for a quick MVP, sarama would also work fine."

**Q**: Why pgx over lib/pq?

**A**: "pgx is faster (2-3x in benchmarks), has better connection pooling, and supports PostgreSQL-specific features like COPY. It's actively maintained and widely used in production. For financial transactions where we process thousands of events per second, the performance gain matters."

---

**Next Steps**:

1. Review this document
2. Install dependencies: `go mod tidy`
3. Start implementing producer/consumer logic
4. Test with real Kafka and PostgreSQL

Good luck! 🚀
