## ADDED Requirements

### Requirement: Adapter implementation structs are unexported
Infrastructure adapter implementation structs in `backend/worker` SHALL be unexported (lowercase). Constructors SHALL return the port interface type, not a pointer to the concrete struct.

#### Scenario: Constructor returns interface not concrete type
- **WHEN** `NewStore`, `NewTranscoder`, or `NewPublisher` is called in `cmd/worker/main.go`
- **THEN** the returned value is typed as the port interface (`application.VideoStorage`, `application.Transcoder`, `application.EventPublisher`) — not as `*Store`, `*Transcoder`, or `*Publisher`

#### Scenario: Concrete struct is not accessible outside the adapter package
- **WHEN** a file outside `internal/infrastructure/out/s3storage/` attempts to reference `s3storage.Store` directly
- **THEN** the Go compiler rejects it — only `s3storage.NewStore` (returning the interface) is accessible
