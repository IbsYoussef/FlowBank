.PHONY: help build up down logs clean test deps

# Default target
help:
	@echo "FlowBank - Event-Driven Transaction Processor"
	@echo ""
	@echo "Available commands:"
	@echo "  make deps         - Install Go dependencies"
	@echo "  make build        - Build all Docker images"
	@echo "  make up           - Start all services"
	@echo "  make down         - Stop all services"
	@echo "  make restart      - Restart all services"
	@echo "  make logs         - Show logs from all services"
	@echo "  make logs-api     - Show API service logs"
	@echo "  make logs-producer - Show producer service logs"
	@echo "  make logs-consumer - Show consumer service logs"
	@echo "  make ps           - Show running services"
	@echo "  make clean        - Remove all containers and volumes"
	@echo "  make db-shell     - Open PostgreSQL shell"
	@echo "  make kafka-topics - List Kafka topics"
	@echo "  make test         - Run Go tests"
	@echo "  make fmt          - Format Go code"
	@echo "  make lint         - Run Go linter"

# Install dependencies
deps:
	go mod download
	go mod tidy

# Build Docker images
build:
	docker compose -f deploy/docker-compose.yml build

# Start all services
up:
	docker compose -f deploy/docker-compose.yml up -d
	@echo ""
	@echo "✅ FlowBank services started!"
	@echo "   - API:              http://localhost:8080/health"
	@echo "   - RedPanda Console: http://localhost:8090"
	@echo "   - PostgreSQL:       localhost:5432"
	@echo ""

# Stop all services
down:
	docker compose -f deploy/docker-compose.yml down

# Restart services
restart: down up

# View logs
logs:
	docker compose -f deploy/docker-compose.yml logs -f

logs-api:
	docker compose -f deploy/docker-compose.yml logs -f api

logs-producer:
	docker compose -f deploy/docker-compose.yml logs -f producer

logs-consumer:
	docker compose -f deploy/docker-compose.yml logs -f consumer

# Show running services
ps:
	docker compose -f deploy/docker-compose.yml ps

# Clean everything
clean:
	docker compose -f deploy/docker-compose.yml down -v
	docker system prune -f

# Database shell
db-shell:
	docker exec -it flowbank-postgres psql -U flowbank_user -d flowbank

# Kafka topic management
kafka-topics:
	docker exec -it flowbank-redpanda rpk topic list

kafka-create-topic:
	docker exec -it flowbank-redpanda rpk topic create transactions --partitions 3 --replicas 1

kafka-consume:
	docker exec -it flowbank-redpanda rpk topic consume transactions --format json

# Development helpers
test:
	go test -v ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

# Quick development workflow
dev: deps build up logs