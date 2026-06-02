# Backend

Go service with two binaries: an HTTP API server and an SQS-driven video processing worker.

## Module

```
github.com/lregnier/design-youtube/backend
```

## Binaries

| Binary       | Entry point       | Purpose                        |
|--------------|-------------------|--------------------------------|
| `api`        | `cmd/api/`        | HTTP API server (port 8080)    |
| `worker`     | `cmd/worker/`     | SQS polling + ffmpeg pipeline  |

## Package layout

```
api/                    OpenAPI spec + codegen config
  openapi.yaml          Source of truth for the API contract
  oapi-codegen.yaml     Codegen config (chi strict server)
  generate.go           go:generate directive
internal/
  api/                  Generated code — do not edit by hand
  config/               Env var loader (fails fast if any var missing)
  handler/              StrictServerInterface implementation
  middleware/           Upload secret middleware
  store/                DynamoDB + Redis data access
cmd/
  api/                  main.go — wires router, middleware, strict handler
  worker/               main.go — SQS poll loop + ffmpeg pipeline
```

## Architecture note

The upload secret is enforced via `StrictMiddlewareFunc` keyed on `operationID` (`InitUpload`, `ConfirmChunk`, `CompleteUpload`). Read operations (`GetVideos`, `GetVideo`) are public.

## Commands

```bash
# Generate API code from openapi.yaml
go generate ./api/...

# Build both binaries
go build ./...

# Build specific binary
go build ./cmd/api
go build ./cmd/worker

# Vet
go vet ./...
```

## Environment variables (all required)

| Variable           | Description                          |
|--------------------|--------------------------------------|
| `UPLOAD_SECRET`    | Shared secret for upload endpoints   |
| `AWS_REGION`       | AWS region                           |
| `DYNAMODB_TABLE`   | DynamoDB table name (`videos`)       |
| `S3_BUCKET`        | S3 bucket for video storage          |
| `CLOUDFRONT_DOMAIN`| CloudFront distribution domain       |
| `SQS_QUEUE_URL`    | SQS FIFO queue URL                   |
| `REDIS_ADDR`       | Redis address (host:port)            |

## Key constraints

- Do not edit `internal/api/api.gen.go` — regenerate from `openapi.yaml` instead.
- DynamoDB video record status values: `uploading`, `processing`, `ready`, `failed`.
- S3 key prefixes: `raw/{videoId}/original`, `segments/{videoId}/`, `manifests/{videoId}/`, `thumbnails/{videoId}/`.
- CloudFront serves `segments/`, `manifests/`, and `thumbnails/` — never `raw/`.
- Presigned part URLs cap at 10MB per chunk (`content-length-range` condition).
