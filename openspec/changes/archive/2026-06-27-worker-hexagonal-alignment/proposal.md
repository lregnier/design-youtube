## Why

The worker bounded context diverges from the architectural patterns established in the api: it uses `adapters/inbound/` and `adapters/outbound/` instead of `infrastructure/in/` and `infrastructure/out/`, keeps port interfaces in a separate `ports/` package, stores events in `internal/event/` with JSON tags, keeps mocks and config inside `internal/`, and exposes the `ProcessVideo` use case as a concrete struct rather than an interface. Aligning the worker closes the conceptual gap between the two contexts and makes the codebase consistent.

## What Changes

- **BREAKING** Rename `internal/adapters/inbound/` → `internal/infrastructure/in/` and `internal/adapters/outbound/` → `internal/infrastructure/out/`
- Merge `internal/ports/outbound.go` port interfaces (`VideoStorage`, `Transcoder`, `ResultPublisher`) into `internal/application/` alongside the use case
- Move `internal/event/result.go` events (`VideoProcessed`, `VideoFailed`) into `internal/domain/processing/events.go`; remove all JSON tags; rename to `VideoProcessingSucceededEvent` and `VideoProcessingFailedEvent` to match api naming
- Add adapter-local wire structs in `sqspublisher` for JSON serialization
- Move `internal/mocks/` → `gen/mocks/`; add `.mockery.yaml`; regenerate mocks
- Move `internal/config/config.go` → `cmd/worker/config.go` as `package main`
- Make `ProcessVideo` an interface with an unexported `processVideo` implementation struct; constructor returns the interface

## Capabilities

### New Capabilities

- none

### Modified Capabilities

- `worker-hexagonal-architecture`: all path-based requirements referencing `internal/adapters/` update to `internal/infrastructure/`; new requirements for domain events and use case interface pattern added

## Impact

- `backend/worker/internal/adapters/` — deleted; contents moved to `internal/infrastructure/`
- `backend/worker/internal/ports/` — deleted; interfaces merged into `internal/application/`
- `backend/worker/internal/event/` — deleted; events moved to `internal/domain/processing/`
- `backend/worker/internal/mocks/` — deleted; replaced by `gen/mocks/` via mockery
- `backend/worker/internal/config/` — deleted; moved to `cmd/worker/config.go`
- `backend/worker/internal/application/process_video.go` — use case becomes interface
- `backend/worker/cmd/worker/main.go` — import paths updated, config qualifier dropped
- `backend/worker/.mockery.yaml` — new file
