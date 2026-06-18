## Why

`S3_PUBLIC_ENDPOINT_URL` encodes two concerns in one variable: whether LocalStack is active and what endpoint to use. A reader has to infer that a non-empty value means LocalStack mode. Splitting into two explicit env vars (`LOCALSTACK_ENABLED` and `LOCALSTACK_ENDPOINT`) makes each concern self-documenting.

## What Changes

- Remove `S3_PUBLIC_ENDPOINT_URL` from both `backend/api` and `backend/worker` config
- Add `LOCALSTACK_ENABLED` (bool) → `Config.LocalStack bool` in both modules
- Add `LOCALSTACK_ENDPOINT` (string) → `Config.LocalStackEndpoint string` in both modules
- Update `cmd/api/main.go` and `cmd/worker/main.go` to branch on `cfg.LocalStack` and pass `cfg.LocalStackEndpoint` to the strategy constructors
- Update `docker-compose.yml` to use the two new env var names

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `api-runtime-config`: replace `S3_PUBLIC_ENDPOINT_URL` with `LOCALSTACK_ENABLED` + `LOCALSTACK_ENDPOINT`
- `worker-runtime-config`: replace `S3_PUBLIC_ENDPOINT_URL` with `LOCALSTACK_ENABLED` + `LOCALSTACK_ENDPOINT`

## Impact

- Affected files: `backend/api/internal/config/config.go`, `backend/api/cmd/api/main.go`, `backend/worker/internal/config/config.go`, `backend/worker/cmd/worker/main.go`, `docker-compose.yml`
- No change to strategy types (`PresignedURLTransformer`, `PublicURLBuilder`) or their implementations
- **BREAKING**: anyone running locally must replace `S3_PUBLIC_ENDPOINT_URL` with `LOCALSTACK_ENABLED=true` + `LOCALSTACK_ENDPOINT=<url>` in their environment
