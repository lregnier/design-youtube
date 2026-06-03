## Why

The API and worker are two distinct bounded contexts with different domain languages, different lifecycles, and different operational concerns — sharing a single Go module blurs those boundaries and couples their deployments. Restructuring `backend/` into `backend/api/` and `backend/worker/` as independent Go modules — within the same monorepo — makes the bounded context separation explicit while keeping everything in one place for development convenience.

## What Changes

- **`backend/` restructured** — split into `backend/api/` (own `go.mod`) and `backend/worker/` (own `go.mod`); each is an independent Go module and bounded context
- **Worker emits integration events** — on completion the worker publishes `VideoProcessed` or `VideoFailed` to a new SQS results queue instead of writing directly to DynamoDB
- **API adds SQS event consumer** — a new background goroutine in the API service polls the results queue and applies events to the `Video` aggregate
- **Worker loses DynamoDB dependency entirely** — its only infrastructure concerns are S3 (read raw, write segments/manifests/thumbnails) and SQS (emit result event)
- **No shared Go code between modules** — each re-implements its own thin infrastructure wrappers; no Shared Kernel
- **New SQS queue** — `video-processing-results.fifo` added to Terraform infra for worker→API event delivery

## Capabilities

### New Capabilities

- `processing-result-event`: Integration event contract (`VideoProcessed`, `VideoFailed`) that the worker emits and the API consumes — defines the message schema shared between bounded contexts
- `api-event-consumer`: API-side SQS consumer that polls the results queue and applies processing outcome to the `Video` aggregate

### Modified Capabilities

- `video-processing`: Worker no longer writes to DynamoDB — it emits a result event instead; processing pipeline itself is unchanged

## Impact

- `backend/` reorganized into `backend/api/` and `backend/worker/` — each with its own `go.mod`, `Dockerfile`, and CI job
- `infra/sqs.tf` adds `video-processing-results.fifo` queue and updates IAM policies
- ECS task definitions updated: worker loses DynamoDB permissions, API gains SQS receive permissions on the results queue
- `docker-compose.yml` updated to reflect new build contexts
- No changes to the OpenAPI spec, frontend, or external API behavior
