## Context

The backend currently lives as a single Go module in `backend/` with two binaries (`cmd/api` and `cmd/worker`) and shared internal packages. While the hexagonal architecture refactor established clean layer boundaries, both binaries still share the same source tree, the same `go.mod`, and the same CI pipeline. The worker also writes directly to DynamoDB — crossing the bounded context boundary by owning state that belongs to the API context.

The goal is to make the bounded context separation structural: separate repos, separate modules, event-based integration only.

## Goals / Non-Goals

**Goals:**
- `backend/` becomes `design-youtube-api/` — the API bounded context owns `Video` aggregate and all DynamoDB writes
- New `design-youtube-worker/` repo — pure processing pipeline, no domain model, no DynamoDB
- Worker communicates outcomes via integration events on a new SQS results queue
- API has a background event consumer that applies worker outcomes to the `Video` aggregate
- No shared Go code between the two repos — each re-implements its own thin infrastructure wrappers
- Terraform updated with the new results queue and adjusted IAM policies

**Non-Goals:**
- Changing the OpenAPI spec or any external API behavior
- Adding a message broker (SNS, Kafka) — SQS FIFO is sufficient for this point-to-point integration
- Shared domain library / Shared Kernel — explicitly rejected; each context owns its own model

## Decisions

### 1. Integration event via SQS FIFO, not SNS fanout

Worker publishes to `video-processing-results.fifo`; API has a single consumer. No fanout needed — there is exactly one consumer of processing results. SQS FIFO gives ordering per `videoId` as message group key and at-least-once delivery with deduplication.

Alternatives considered: SNS → SQS (unnecessary indirection for one consumer); direct HTTP callback from worker to API (couples worker to API's network address and availability).

### 2. Integration event schema is the only shared contract

The message payload is the only thing the two repos need to agree on:

```json
// VideoProcessed
{
  "eventType": "VideoProcessed",
  "videoId": "...",
  "manifestUrl": "https://cdn.example.com/manifests/.../master.m3u8",
  "thumbnailUrl": "https://cdn.example.com/thumbnails/.../thumb.jpg"
}

// VideoFailed
{
  "eventType": "VideoFailed",
  "videoId": "...",
  "reason": "ffmpeg decode error"
}
```

This schema is documented in `processing-result-event` spec. Both repos evolve independently as long as they honour this contract.

### 3. Worker has no domain model — just a job struct

The worker doesn't know what a `Video` aggregate is. It receives a job (`videoId`, `s3Key`), does the pipeline, emits a result. Its internal model:

```go
type ProcessingJob    struct { VideoID, S3Key string }
type ProcessingResult struct { VideoID, ManifestURL, ThumbnailURL, Reason string; Success bool }
```

No aggregates, no repositories, no status transitions. This keeps the worker simple and aligned with its bounded context.

### 4. API event consumer runs as a background goroutine

The API service starts an SQS long-poll loop in a goroutine alongside the HTTP server. On receiving a `VideoProcessed` event it loads the `Video` aggregate via `VideoRepository`, calls `MarkReady(manifestURL, thumbnailURL)`, and saves. On `VideoFailed` it calls `MarkFailed()`. The consumer uses the same `VideoRepository` port and DynamoDB adapter already in place — no new infrastructure.

Alternatives considered: separate ECS task for the consumer (more operational complexity for minimal benefit at this scale).

### 5. No shared Go module — independent infrastructure wrappers

Both repos will have their own S3 and SQS client setup (~30 lines each). This is deliberate: avoiding a shared library keeps the bounded contexts independent and avoids the Shared Kernel anti-pattern. The duplication is trivial.

### 6. Monorepo with independent Go modules per bounded context

Both services live under `backend/` in this repo, but as separate Go modules:

```
backend/
  api/      ← go.mod: github.com/lregnier/design-youtube/api
  worker/   ← go.mod: github.com/lregnier/design-youtube/worker
frontend/
infra/
```

Each module has its own `go.mod`, `Dockerfile`, and CI job. They share no Go code. This gives bounded context independence (separate dependency graphs, separate build/deploy pipelines) while keeping everything in one repo for development convenience — no need to clone multiple repos to work on the full system.

The Terraform in `infra/` manages both services — the new results queue and updated IAM policies are added there.

## Risks / Trade-offs

- **At-least-once delivery** → API consumer must be idempotent: applying `VideoProcessed` to an already-ready video is a no-op.
- **Event ordering** → Using `videoId` as SQS FIFO message group key ensures ordered delivery per video. Two events for the same video won't interleave.
- **Worker repo not in this openspec** → The `design-youtube-worker` repo is new and external; its tasks are tracked here for coordination but its code lives elsewhere.
- **Rename of `backend/`** → Git history for `backend/` is preserved; `backend/api/` and `backend/worker/` are created from the existing content via `git mv`.

## Migration Plan

1. Add `video-processing-results.fifo` SQS queue to Terraform; update IAM policies
2. Move `backend/` content into `backend/api/`; update `go.mod` module path and all imports
3. Create `backend/worker/` as a new Go module; copy worker logic; strip DynamoDB; add SQS result emission
4. Add API event consumer goroutine; wire into `backend/api/cmd/api/main.go`
5. Update ECS task definitions, `docker-compose.yml`, and CI workflows for new paths
6. Update OpenSpec design docs to reflect new `backend/api/` and `backend/worker/` structure
