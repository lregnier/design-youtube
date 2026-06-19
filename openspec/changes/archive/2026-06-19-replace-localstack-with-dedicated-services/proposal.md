## Why

LocalStack community edition does not persist S3 objects, making local video data ephemeral across container restarts. Replacing LocalStack with three purpose-built services (MinIO, DynamoDB Local, ElasticMQ) gives each service reliable volume persistence and removes a heavy, multi-service emulator in favour of focused, production-analogous tools.

## What Changes

**Infrastructure (docker-compose.yml):**
- Remove the `localstack` service and its `localstack-data` volume
- Add `minio` service (S3-compatible object store) with `minio-data` volume
- Add `minio-init` one-shot container (uses `minio/mc`) to create the bucket and apply CORS on first boot
- Add `dynamodb-local` service (official AWS image) with `dynamodb-data` volume
- Add `dynamodb-init` one-shot container (uses `amazon/aws-cli`) to create the table idempotently
- Add `elasticmq` service (SQS-compatible) with declarative queue config via `docker/elasticmq.conf`
- Remove `LOCALSTACK_ENABLED`, `LOCALSTACK_ENDPOINT`, `S3_USE_PATH_STYLE` env vars from api and worker services
- Add `S3_ENDPOINT_URL`, `S3_PUBLIC_URL`, `DYNAMODB_ENDPOINT_URL`, `SQS_ENDPOINT_URL` env vars for local service discovery

**New config file:**
- `docker/elasticmq.conf` — declares both FIFO queues declaratively (no init script needed)

**Remove:**
- `docker/localstack-init.sh`

**Application code — both `backend/api` and `backend/worker`:**
- Replace `LocalStackEnabled bool` + `LocalStackEndpoint string` config fields with `S3Endpoint string`, `DynamoDBEndpoint string`, `SQSEndpoint string`
- Wire per-service endpoint overrides into AWS SDK client constructors in both `main.go` files
- Rename `LocalStackTransformer` → `EndpointTransformer` and `LocalStackURLBuilder` → `EndpointURLBuilder` to remove the LocalStack coupling from type names
- Change URL strategy condition from `cfg.LocalStackEnabled` to `cfg.S3Endpoint != ""`

**Supersedes:** `persist-localstack-data` — archive that change before applying this one.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `api-runtime-config`: replace `LOCALSTACK_ENABLED` + `LOCALSTACK_ENDPOINT` with `S3_ENDPOINT_URL`, `DYNAMODB_ENDPOINT_URL`, `SQS_ENDPOINT_URL`
- `worker-runtime-config`: same replacement

## Impact

- Affected files: `docker-compose.yml`, `docker/localstack-init.sh` (deleted), `docker/elasticmq.conf` (new), `backend/api/internal/config/config.go`, `backend/api/cmd/api/main.go`, `backend/api/internal/adapters/outbound/s3store/url_transformer.go`, `backend/worker/internal/config/config.go`, `backend/worker/cmd/worker/main.go`, `backend/worker/internal/adapters/outbound/s3storage/url_builder.go`
- No change to port interfaces, application layer, or domain
- Local dev: run `docker compose down -v && docker compose up --build` after applying
