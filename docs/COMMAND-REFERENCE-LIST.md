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

## Makefile Shortcuts

| Command         | Equivalent                   | Use Case                     |
| --------------- | ---------------------------- | ---------------------------- |
| `make build`    | `docker compose ... build`   | Compile new code             |
| `make up`       | `docker compose ... up -d`   | Start application            |
| `make down`     | `docker compose ... down`    | Stop containers              |
| `make logs`     | `docker compose ... logs -f` | Check transaction flow       |
| `make clean`    | `docker compose ... down -v` | Full reset (clears database) |
| `make db-shell` | `docker exec ... psql`       | Inspect database             |

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
