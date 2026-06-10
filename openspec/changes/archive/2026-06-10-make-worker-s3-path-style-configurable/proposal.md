## Why

The worker's S3 client is constructed with no options, so against LocalStack (`AWS_ENDPOINT_URL=http://localstack:4566`) it issues virtual-hosted-style requests like `http://design-youtube-video-prod.localstack:4566/...`, a hostname that doesn't resolve inside the docker-compose network. `DownloadRaw` fails with a DNS error, the worker reports `PublishFailed("download failed: ...")`, and the video silently ends up `failed` with no logged error. The API solved the same problem for itself via `S3_USE_PATH_STYLE`; the worker needs the same fix.

## What Changes

- Add `S3UsePathStyle bool` to the worker config, parsed from `S3_USE_PATH_STYLE` env var (optional, defaults to false)
- In `cmd/worker/main.go`: apply `UsePathStyle` on the S3 client only when the config flag is true, mirroring the API's `s3Opts` pattern
- In `docker-compose.yml`: set `S3_USE_PATH_STYLE: "true"` for the `worker` service, matching the `api` service
- Document the new env var in `backend/worker/README.md` if a configuration table exists

## Capabilities

### New Capabilities

- `worker-runtime-config`: S3 path-style addressing is configurable via environment variable (mirrors `api-runtime-config`)

### Modified Capabilities

_(none)_

## Impact

- **`backend/worker/internal/config/config.go`**: new `S3UsePathStyle` field
- **`backend/worker/cmd/worker/main.go`**: conditional `UsePathStyle` on the S3 client
- **`docker-compose.yml`**: new env var in the `worker` service block
- **`backend/worker/README.md`**: configuration docs (if applicable)
