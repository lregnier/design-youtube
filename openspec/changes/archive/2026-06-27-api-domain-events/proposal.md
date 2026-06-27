## Why

Processing result events (`VideoProcessedEvent`, `VideoFailedEvent`) live in `application/` instead of `domain/video/`, and carry JSON struct tags — a serialization concern that belongs in the infrastructure adapter, not the domain. This splits the event model across two layers and leaks wire format details into business logic.

## What Changes

- Move `VideoProcessedEvent` and `VideoFailedEvent` from `internal/application/events.go` into `internal/domain/video/events.go`; remove all `json:""` struct tags from these types
- Delete `internal/application/events.go`
- Update `ProcessingService.OnProcessed` and `OnFailed` to accept `video.VideoProcessedEvent` and `video.VideoFailedEvent`
- Add adapter-local JSON structs in `sqssubscriber` that handle deserialization and map into the domain event types before calling the service
- Update `processing_service_test.go` to use `video.VideoProcessedEvent` / `video.VideoFailedEvent`

## Capabilities

### New Capabilities

- none

### Modified Capabilities

- `hexagonal-architecture`: requirement for the domain layer ("Domain layer has no external package imports") is already satisfied; this change reinforces the corollary that domain types SHALL NOT carry serialization annotations. Adding a scenario to capture this.

## Impact

- `internal/domain/video/events.go` — add `VideoProcessedEvent`, `VideoFailedEvent` (no JSON tags)
- `internal/application/events.go` — deleted
- `internal/application/processing_service.go` — method signatures change to `video.VideoProcessedEvent` / `video.VideoFailedEvent`
- `internal/application/processing_service_test.go` — construct `video.VideoProcessedEvent` / `video.VideoFailedEvent` instead of `application.*`
- `internal/infrastructure/in/sqssubscriber/subscriber.go` — add local JSON structs; map to domain types before calling service
