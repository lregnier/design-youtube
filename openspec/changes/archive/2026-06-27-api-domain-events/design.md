## Context

There are two event definitions in the codebase:

- `domain/video/events.go`: `DomainEvent` interface + `VideoUploadedEvent` — outbound events emitted by the domain, correctly free of JSON tags
- `application/events.go`: `VideoProcessedEvent` + `VideoFailedEvent` — inbound processing results consumed by `ProcessingService`, incorrectly carrying `json:""` struct tags

The `sqssubscriber` currently unmarshals SQS message bodies directly into `application.VideoProcessedEvent` / `application.VideoFailedEvent`, relying on the JSON tags embedded in those types. This makes application-layer types responsible for wire format — a hexagonal architecture violation.

## Goals / Non-Goals

**Goals:**
- Consolidate all domain events into `domain/video/events.go`, tag-free
- Push JSON deserialization into the `sqssubscriber` adapter via local wire structs
- `ProcessingService` methods accept pure domain types with no serialization knowledge

**Non-Goals:**
- Changing the SQS message format or any wire protocol
- Adding `EventType` field to domain events (it is a transport concern)

## Decisions

### Domain events carry no struct tags

`json:""` tags are a serialization contract with a specific wire format. Domain types represent business concepts, not wire payloads. If the wire format changes (field rename, new transport), only the adapter changes — domain types are unaffected.

Alternative considered: keep JSON tags on domain types for convenience. Rejected — it couples the domain to a specific serialization format and makes the domain layer harder to reuse across transports.

### Adapter-local wire structs in sqssubscriber

`sqssubscriber` adds unexported structs (`videoProcessedMessage`, `videoFailedMessage`) with JSON tags matching the SQS message format. After unmarshaling, it maps to `video.VideoProcessedEvent` / `video.VideoFailedEvent` before calling the service. This keeps the mapping collocated with the transport.

### ProcessingService methods take domain event types

`OnProcessed(ctx, video.VideoProcessedEvent)` and `OnFailed(ctx, video.VideoFailedEvent)` — the service operates on domain vocabulary, not wire vocabulary. The adapter is responsible for the translation.

### EventType field stays in the adapter

`EventType` is used only for routing in `sqssubscriber` (switch on `"VideoProcessed"` / `"VideoFailed"`). It is a transport routing key, not a domain concept. It belongs in the adapter's wire struct, not in the domain event.

## Risks / Trade-offs

[Test updates] → `processing_service_test.go` constructs event values directly. Switching from `application.VideoProcessedEvent` to `video.VideoProcessedEvent` requires updating field construction in tests — mechanical, no logic changes.

[Mock regeneration] → `ProcessingService` interface signature changes, so the mockery-generated `MockProcessingService` will be stale. Running `mockery` regenerates it; no manual mock editing needed.
