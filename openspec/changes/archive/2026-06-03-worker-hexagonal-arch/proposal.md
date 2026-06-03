## Why

The worker bounded context has the right event-based integration contract but its internal structure is flat — all processing logic, infrastructure calls, and SQS polling live in a single `cmd/worker/main.go`. Applying the same hexagonal architecture already established in the API makes the worker's boundaries explicit, its pipeline independently testable, and the codebase consistent in its architectural approach.

## What Changes

- New `internal/domain/processing/` package: `ProcessingJob` value object (`VideoID`, `S3Key`)
- New `internal/application/process_video.go`: `ProcessVideo` use case that orchestrates `VideoStorage`, `Transcoder`, and `ResultPublisher` port interfaces
- New `internal/ports/outbound.go`: `VideoStorage`, `Transcoder`, and `ResultPublisher` interfaces
- New `internal/adapters/inbound/sqsjobs/consumer.go`: SQS long-poll loop as an inbound adapter
- New `internal/adapters/outbound/s3storage/store.go`: S3 `VideoStorage` implementation
- New `internal/adapters/outbound/ffmpeg/transcoder.go`: `Transcoder` implementation wrapping `exec.Command`
- New `internal/adapters/outbound/sqspublisher/publisher.go`: `ResultPublisher` implementation (replaces `internal/queue/publisher.go`)
- `cmd/worker/main.go` reduced to a pure composition root
- Deleted: `internal/queue/publisher.go`, `internal/event/result.go` (moved into appropriate packages)
- No behavior changes: same S3 key structure, same ffmpeg flags, same SQS events emitted

## Capabilities

### New Capabilities

### Modified Capabilities

## Impact

- All files under `backend/worker/internal/` restructured
- `cmd/worker/main.go` simplified — only wires adapters and use case
- No changes to `backend/api/`, `infra/`, `frontend/`, or any integration contracts
- No changes to `openspec/specs/` — this is an internal restructure with no requirement changes
