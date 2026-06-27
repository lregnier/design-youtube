## MODIFIED Requirements

### Requirement: SQS job polling is an inbound adapter
The SQS long-poll loop SHALL live in `internal/infrastructure/in/sqsjobs/`. It SHALL parse the incoming job message and call the `ProcessVideo` use case via its interface. It SHALL NOT contain ffmpeg calls, S3 calls, or business logic.

#### Scenario: Inbound adapter delegates to use case
- **WHEN** the SQS consumer receives a valid job message
- **THEN** it calls `ProcessVideo.Execute` and handles success/error — no direct infrastructure calls

### Requirement: Outbound adapters delegate environment-specific URL building to an injected strategy
The `s3storage` adapter in `backend/worker` SHALL accept a `PublicURLBuilder` (an interface with `AssetURL(bucket, key string) string`) and call it unconditionally — no if/else on `S3_PUBLIC_ENDPOINT_URL` or `CloudFrontDomain` inside the store. The composition root (`cmd/worker/main.go`) SHALL select the concrete implementation (`CloudFrontURLBuilder` for production, `LocalStackURLBuilder` for local dev) based on config.

#### Scenario: Store contains no environment branching
- **WHEN** `backend/worker/internal/infrastructure/out/s3storage/store.go` is inspected
- **THEN** no conditional on `s3PublicEndpointURL` or `cloudfrontDomain` exists — only an unconditional call to the injected `PublicURLBuilder`

#### Scenario: Composition root selects the URL strategy
- **WHEN** `cmd/worker/main.go` is inspected
- **THEN** it constructs either `CloudFrontURLBuilder` or `LocalStackURLBuilder` based on config and injects it into `s3storage.NewStore`

## ADDED Requirements

### Requirement: Worker internal/ contains only three DDD layers
`backend/worker/internal/` SHALL contain exactly three subdirectories: `domain/`, `application/`, and `infrastructure/`. No `ports/`, `event/`, `mocks/`, or `config/` package SHALL exist inside `internal/`.

#### Scenario: internal/ structure is clean
- **WHEN** `backend/worker/internal/` is inspected
- **THEN** only `domain/`, `application/`, and `infrastructure/` directories are present

### Requirement: Port interfaces are defined in the application layer
`VideoStorage`, `Transcoder`, and `ResultPublisher` SHALL be defined as Go interfaces in `internal/application/` alongside `ProcessVideo`. No separate `ports/` package SHALL exist.

#### Scenario: Use case and its ports share a package
- **WHEN** `internal/application/` is inspected
- **THEN** `VideoStorage`, `Transcoder`, `ResultPublisher`, and `ProcessVideo` are all defined in the same package

### Requirement: Worker domain events are defined in the domain layer without JSON tags
All domain events (`VideoProcessingSucceededEvent`, `VideoProcessingFailedEvent`) SHALL be defined in `internal/domain/processing/`. They SHALL be plain Go structs with no `json:""` struct tags. JSON serialization SHALL be handled by adapter-local wire structs in `internal/infrastructure/out/sqspublisher/`.

#### Scenario: Domain event types have no JSON tags
- **WHEN** `internal/domain/processing/` is inspected
- **THEN** no struct field on any event type carries a `json:""` tag

#### Scenario: Publisher uses local wire structs for serialization
- **WHEN** `internal/infrastructure/out/sqspublisher/publisher.go` is inspected
- **THEN** it marshals unexported wire structs (with JSON tags) — not the domain event types directly

### Requirement: ProcessVideo use case is an interface
`ProcessVideo` SHALL be defined as a Go interface in `internal/application/` with an `Execute(ctx context.Context, job processing.ProcessingJob) error` method. The concrete implementation SHALL be an unexported struct. The constructor SHALL return the `ProcessVideo` interface.

#### Scenario: Consumer depends on ProcessVideo interface
- **WHEN** `internal/infrastructure/in/sqsjobs/consumer.go` is inspected
- **THEN** it holds a field of type `application.ProcessVideo` (interface), not the concrete struct

### Requirement: Generated mocks live under gen/mocks/
All mockery-generated mocks for the worker SHALL live under `backend/worker/gen/mocks/`. The `backend/worker/.mockery.yaml` config SHALL specify `dir: gen/mocks`. No hand-written or generated mock SHALL live under `internal/`.

#### Scenario: Mocks are isolated under gen/
- **WHEN** `backend/worker/` is inspected
- **THEN** mock files are at `gen/mocks/` and `internal/` contains no mock files

### Requirement: Worker config lives in cmd/worker/ as package main
The worker `Config` struct and `Load()` function SHALL be defined in `backend/worker/cmd/worker/config.go` as `package main`. `internal/config/` SHALL NOT exist.

#### Scenario: Config is co-located with the entrypoint
- **WHEN** `backend/worker/cmd/worker/` is inspected
- **THEN** `config.go` exists alongside `main.go` in `package main`, and no `internal/config/` directory exists
