## Context

The local dev stack currently uses LocalStack as an all-in-one AWS emulator. LocalStack community edition reliably emulates DynamoDB and SQS but does not persist S3 objects. Replacing it with three purpose-built services gives each AWS service a dedicated, production-analogous emulator with reliable volume persistence:

- **MinIO** — S3-compatible object store; production-quality persistence, supports presigned URLs and CORS
- **amazon/dynamodb-local** — official AWS DynamoDB emulator; persists to a local SQLite file
- **ElasticMQ** — lightweight SQS-compatible queue server; queues defined declaratively in a config file

The application currently uses a single `AWS_ENDPOINT_URL` env var (picked up automatically by the Go AWS SDK) to redirect all service calls to LocalStack. With separate services on different ports, this must become per-service endpoint overrides injected at client construction time.

## Goals / Non-Goals

**Goals:**
- Full volume persistence for all three data stores
- No behaviour change for production (endpoint overrides are only set when env vars are present)
- Clean type names (remove "LocalStack" from adapter types)
- Minimal application code change — only config and composition root

**Non-Goals:**
- MinIO console access from the app (the console on port 9001 is for manual inspection only)
- SQS message persistence across ElasticMQ restarts (queue definitions persist via config; in-flight messages do not — acceptable for local dev)
- Changing port interfaces, application layer, or domain

## Decisions

**Per-service endpoint config fields, not a single override**

With three separate services, a single `AWS_ENDPOINT_URL` no longer makes sense. Each service gets its own optional env var (`S3_ENDPOINT_URL`, `DYNAMODB_ENDPOINT_URL`, `SQS_ENDPOINT_URL`). When unset, the AWS SDK uses real AWS endpoints — production requires no change.

**Go AWS SDK v2 per-service endpoint injection**

Each service client is constructed with an `Options` function that sets `BaseEndpoint` when the config field is non-empty:
```go
// S3 (api main.go)
s3Opts := []func(*awss3.Options){}
if cfg.S3Endpoint != "" {
    s3Opts = append(s3Opts, func(o *awss3.Options) {
        o.BaseEndpoint = aws.String(cfg.S3Endpoint)
    })
}

// DynamoDB (api main.go)
dynamoOpts := []func(*dynamodb.Options){}
if cfg.DynamoDBEndpoint != "" {
    dynamoOpts = append(dynamoOpts, func(o *dynamodb.Options) {
        o.BaseEndpoint = aws.String(cfg.DynamoDBEndpoint)
    })
}
```

**Two S3 config fields: `S3Endpoint` (SDK) and `S3PublicURL` (browser-facing)**

`S3_ENDPOINT_URL` → `S3Endpoint` is the internal Docker-network address injected into the SDK client (`http://minio:9000`). `S3_PUBLIC_URL` → `S3PublicURL` is the host-reachable address used when building URLs returned to the browser (`http://localhost:9000`). This mirrors the old split between `AWS_ENDPOINT_URL` and `LOCALSTACK_ENDPOINT`.

The `PresignedURLTransformer` (api) and `PublicURLBuilder` (worker) strategies are selected at startup. The condition changes from `cfg.LocalStackEnabled` to `cfg.S3PublicURL != ""`. The transformer/builder receives `cfg.S3PublicURL`.

**Rename `LocalStackTransformer` → `EndpointTransformer`, `LocalStackURLBuilder` → `EndpointURLBuilder`**

These types rewrite URLs based on a custom endpoint — they are not LocalStack-specific. Renaming decouples them from a tool name.

**ElasticMQ queue definitions in `elasticmq.conf`, not an init script**

ElasticMQ reads its queue configuration from a Hocon config file at startup. This is idempotent by design and eliminates the need for an init container for SQS.

**MinIO init via `minio/mc` one-shot container**

MinIO has no built-in bucket/CORS init mechanism. A one-shot `minio-init` container using the `minio/mc` CLI creates the bucket and sets CORS after MinIO is healthy. The `mc mb --ignore-existing` flag makes it idempotent.

**DynamoDB init via `amazon/aws-cli` one-shot container**

DynamoDB Local has no built-in table init. A one-shot `dynamodb-init` container runs the `aws dynamodb create-table` command with `--no-fail-on-already-existing` (or guarded with `|| true`) after DynamoDB Local is healthy.

**Remove `S3_USE_PATH_STYLE`**

MinIO requires path-style addressing (`http://minio:9000/bucket/key`) and it is always enabled for local dev. Rather than a configurable flag, path-style is hardcoded when `S3Endpoint` is set. In production (no `S3Endpoint`), virtual-hosted style is used. This removes a config variable that was only ever set to `true` locally.

**Remove `docker/localstack-init.sh`**

Replaced entirely by the per-service init approach above.

## Risks / Trade-offs

- **More services in compose** → `docker compose up` starts 6 services (was 4). Each has a healthcheck; startup time increases slightly.
- **ElasticMQ SQS URL format** → Queue URLs become `http://elasticmq:9324/000000000000/video-processing.fifo`. The account ID (`000000000000`) and region must match what the app sends. ElasticMQ is configured with `accountId = 000000000000` and `region = us-east-1` to match existing queue URL config.
- **MinIO CORS for presigned URLs** → MinIO supports CORS via `mc anonymous` or bucket policies. The init container must set this correctly or browser PUT requests for uploads will fail.
- **`aws` import in DynamoDB init container** → `amazon/aws-cli` needs `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and `AWS_DEFAULT_REGION` set (any non-empty values work for local).
