## Requirements

### Requirement: Worker domain layer has no external package imports
The `internal/domain/` package SHALL import only the Go standard library. No AWS SDK, ffmpeg, SQS, or HTTP packages SHALL appear in any file under `internal/domain/`.

#### Scenario: Domain compiles without infrastructure dependencies
- **WHEN** `go build ./internal/domain/...` is run in `backend/worker/`
- **THEN** the build succeeds with no imports outside the standard library

### Requirement: ProcessVideo use case depends only on port interfaces
The `ProcessVideo` use case in `internal/application/` SHALL accept `VideoStorage`, `Transcoder`, and `ResultPublisher` as constructor arguments. It SHALL NOT import any adapter package, AWS SDK, or ffmpeg directly.

#### Scenario: Use case is constructable with mock implementations
- **WHEN** `ProcessVideo` is instantiated with stub implementations of its three port interfaces
- **THEN** it compiles and can be called without any real infrastructure present

### Requirement: Transcoder port abstracts all ffmpeg interactions
All ffmpeg and ffprobe invocations SHALL go through the `Transcoder` port interface. The `internal/application/` package SHALL NOT call `exec.Command` directly.

#### Scenario: Transcoder is the only ffmpeg dependency
- **WHEN** the import graph of `internal/application/` is inspected
- **THEN** no import of `os/exec` appears — only the `ports` package

### Requirement: SQS job polling is an inbound adapter
The SQS long-poll loop SHALL live in `internal/adapters/inbound/sqsjobs/`. It SHALL parse the incoming job message and call the `ProcessVideo` use case. It SHALL NOT contain ffmpeg calls, S3 calls, or business logic.

#### Scenario: Inbound adapter delegates to use case
- **WHEN** the SQS consumer receives a valid job message
- **THEN** it calls `ProcessVideo.Execute` and handles success/error — no direct infrastructure calls

### Requirement: cmd/worker/main.go is a pure composition root
`cmd/worker/main.go` SHALL only instantiate adapters, inject them into the use case, inject the use case into the inbound adapter, and start the consumer. It SHALL contain no processing logic, ffmpeg calls, or S3 calls.

#### Scenario: main.go contains no business logic
- **WHEN** `cmd/worker/main.go` is read
- **THEN** it contains only dependency construction and wiring — all logic lives in internal packages
