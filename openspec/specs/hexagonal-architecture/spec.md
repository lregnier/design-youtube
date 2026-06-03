## Requirements

### Requirement: Domain layer has no external package imports
The `internal/domain/` package SHALL import only the Go standard library. No AWS SDK, Redis, HTTP, or other infrastructure packages SHALL appear in any file under `internal/domain/`.

#### Scenario: Domain compiles without infrastructure dependencies
- **WHEN** `go build ./internal/domain/...` is run in the backend directory
- **THEN** the build succeeds with no imports from `github.com/aws`, `github.com/redis`, or `github.com/go-chi`

### Requirement: Application layer depends only on domain and port interfaces
The `internal/application/` package SHALL import `internal/domain/` and `internal/ports/` only. It SHALL NOT import any adapter package or AWS SDK directly.

#### Scenario: Use cases are constructable with interface implementations
- **WHEN** a use case struct is instantiated with mock implementations of its port interfaces
- **THEN** it compiles and can be called without any real infrastructure present

### Requirement: Each use case is a distinct type with an Execute method
Each operation (InitUpload, ConfirmChunk, CompleteUpload, GetVideo, ListVideos) SHALL be implemented as a Go struct in `internal/application/` with a single public method that accepts a typed command and returns a typed result.

#### Scenario: InitUpload use case is independently constructable
- **WHEN** `InitUpload` is constructed with a `VideoRepository` and `ObjectStore`
- **THEN** it can execute without knowledge of DynamoDB, S3, or HTTP specifics

### Requirement: Outbound port interfaces are defined as minimal Go interfaces
`VideoRepository`, `ObjectStore`, `Cache`, and `Queue` SHALL each be defined as Go interfaces with only the methods required by the application layer. No interface SHALL have more than 6 methods.

#### Scenario: Adapter satisfies port interface
- **WHEN** the DynamoDB adapter is compiled against `VideoRepository`
- **THEN** the Go compiler confirms it satisfies the interface without explicit declaration

### Requirement: Inbound HTTP adapter contains no business logic
The HTTP handler in `internal/adapters/inbound/http/` SHALL only translate between HTTP request/response types and use-case command/result types. Domain rules (size limits, status transitions, chunk ordering) SHALL live in the application or domain layer, not in the handler.

#### Scenario: Handler delegates to use case
- **WHEN** `InitUpload` HTTP handler is called
- **THEN** it constructs a command, calls the use case `Execute` method, and maps the result — with no direct DynamoDB, S3, or business logic calls

### Requirement: Dependency direction flows inward
Adapters SHALL import domain and application packages. Domain and application packages SHALL NOT import adapter packages. This enforces the hexagonal architecture dependency rule.

#### Scenario: Import graph has no outward dependency
- **WHEN** the import graph of `internal/domain/` and `internal/application/` is inspected
- **THEN** no import path leads to `internal/adapters/` or any infrastructure package
