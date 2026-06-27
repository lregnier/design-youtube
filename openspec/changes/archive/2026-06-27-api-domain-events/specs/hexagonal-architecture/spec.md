## MODIFIED Requirements

### Requirement: Domain layer has no external package imports
The `internal/domain/` package SHALL import only the Go standard library. No AWS SDK, Redis, HTTP, or other infrastructure packages SHALL appear in any file under `internal/domain/`. Domain types SHALL NOT carry serialization annotations (e.g., `json:""` struct tags) — serialization is an infrastructure concern and SHALL be handled in adapter-local types.

#### Scenario: Domain compiles without infrastructure dependencies
- **WHEN** `go build ./internal/domain/...` is run in the backend directory
- **THEN** the build succeeds with no imports from `github.com/aws`, `github.com/redis`, or `github.com/go-chi`

#### Scenario: Domain event types have no JSON tags
- **WHEN** `internal/domain/video/events.go` is inspected
- **THEN** no struct field carries a `json:""` tag — all event types are plain Go structs

## ADDED Requirements

### Requirement: All domain events are defined in the domain layer
All events representing things that happened in the domain (`VideoUploadedEvent`, `VideoProcessingSucceededEvent`, `VideoProcessingFailedEvent`) SHALL be defined in `internal/domain/video/`. No event type SHALL be defined in `internal/application/` or any infrastructure package.

#### Scenario: Application package contains no event type definitions
- **WHEN** `internal/application/` is inspected
- **THEN** no file defines a struct whose name ends in `Event`

#### Scenario: Inbound adapter maps wire payload to domain events before calling service
- **WHEN** the SQS subscriber receives a `VideoProcessed` message
- **THEN** it unmarshals the JSON into an adapter-local wire struct, maps the fields to `video.VideoProcessingSucceededEvent`, and calls `ProcessingService.HandleVideoProcessingSucceeded` with the domain type — no JSON tags appear on `video.VideoProcessingSucceededEvent`
