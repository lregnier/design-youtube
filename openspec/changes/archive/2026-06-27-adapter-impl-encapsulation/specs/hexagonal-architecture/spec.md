## ADDED Requirements

### Requirement: Adapter implementation structs are unexported
Infrastructure adapter implementation structs SHALL be unexported (lowercase). Constructors SHALL return the port interface type, not a pointer to the concrete struct. This ensures callers in the composition root depend only on the interface, not the implementation.

#### Scenario: Constructor returns interface not concrete type
- **WHEN** `NewRepository`, `NewStore`, or `NewPublisher` is called in `cmd/api/main.go`
- **THEN** the returned value is typed as the port interface (`video.VideoRepository`, `application.ObjectStore`, `application.EventPublisher`) — not as `*Repository`, `*Store`, or `*Publisher`

#### Scenario: Concrete struct is not accessible outside the adapter package
- **WHEN** a file outside `internal/infrastructure/out/dynamo/` attempts to reference `dynamo.Repository` directly
- **THEN** the Go compiler rejects it — only `dynamo.NewRepository` (returning the interface) is accessible
