## Context

The local dev stack runs LocalStack 3.x (community edition) and Redis 7 via Docker Compose. Neither service has a volume mount today, so all state is stored in the container's writable layer and lost on `docker compose down` or container recreation.

LocalStack 3.x community edition supports persistence via `PERSISTENCE=1` — it snapshots service state to `/var/lib/localstack` between requests and restores it on startup. The init script (`docker/localstack-init.sh`) runs on every startup via the `ready.d` hook; without idempotency guards it will fail when resources already exist after a restore.

## Goals / Non-Goals

**Goals:**
- Data (S3 objects, DynamoDB items, SQS queue definitions) survives `docker compose stop` / `docker compose start`
- Redis cache survives restarts (nice-to-have; avoids cold-cache on restart)
- Init script runs safely whether or not resources already exist

**Non-Goals:**
- Production persistence (this is local dev only)
- Backup or export of LocalStack state
- LocalStack Pro features (persistence works in community edition)

## Decisions

**Named volumes over bind mounts**

Named volumes (`localstack-data`, `redis-data`) are managed by Docker and portable across machines without path assumptions. Bind mounts would work but couple the setup to the host filesystem layout.

**`PERSISTENCE=1` over manual state export**

LocalStack's built-in persistence is the supported path for community edition. Manual state export (e.g. `awslocal dynamodb batch-write-item`) would require additional tooling and wouldn't cover S3 objects.

**Idempotency via existence checks, not `|| true`**

Using `awslocal s3api head-bucket` and `awslocal dynamodb describe-table` before creating resources preserves `set -e` error handling and makes the intent explicit. Using `|| true` would silently swallow real errors (e.g. a malformed table definition).

**Redis persistence**

Redis 7 with default config writes an RDB snapshot periodically. Mounting `/data` is sufficient — no `redis.conf` changes needed. Since Redis is used as a cache here, losing it on restart is not catastrophic, but persistence avoids unnecessary DynamoDB cold reads.

## Risks / Trade-offs

- **Stale LocalStack state**: if the schema changes (new GSI, renamed bucket), the persisted state may conflict. Mitigation: document `docker compose down -v` as the reset command.
- **LocalStack community persistence limitations**: not all services support full fidelity persistence in community edition. S3 objects, DynamoDB items, and SQS queue definitions are well-supported; complex SQS in-flight messages may not be restored. Acceptable for local dev.
- **Disk usage**: video files uploaded to LocalStack S3 will accumulate on the host. Mitigation: periodic `docker compose down -v` to reset.

## Migration Plan

Existing local setups with running containers should run `docker compose down -v` before pulling this change to avoid a volume/state mismatch. After that, `docker compose up` creates the new named volumes and initialises them cleanly.
