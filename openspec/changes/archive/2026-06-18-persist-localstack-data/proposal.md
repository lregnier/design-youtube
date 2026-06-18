## Why

LocalStack (S3, DynamoDB, SQS) and Redis lose all state when their containers stop. Every restart requires re-uploading videos and re-processing them from scratch, which makes the local development loop painful.

## What Changes

- Add a named Docker volume (`localstack-data`) mounted at `/var/lib/localstack` and set `PERSISTENCE=1` in the LocalStack service definition in `docker-compose.yml`
- Add a named Docker volume (`redis-data`) mounted at `/data` in the Redis service definition
- Declare both volumes in the top-level `volumes:` block of `docker-compose.yml`
- Make `docker/localstack-init.sh` idempotent:
  - Guard `s3 mb` with an existence check (skip if bucket already exists)
  - Guard `dynamodb create-table` with an existence check (skip if table already exists)
  - `sqs create-queue` is already idempotent — no change needed

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

(none — this is an infrastructure/dev-environment change with no spec-level behavior change)

## Impact

- Affected files: `docker-compose.yml`, `docker/localstack-init.sh`
- No application code changes, no API changes
- **BREAKING** (local dev only): existing LocalStack containers must be stopped and volumes removed (`docker compose down -v`) before restarting to avoid stale state from a previous non-persistent setup
