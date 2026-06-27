## Requirements

### Requirement: Domain layer has no external package imports
The `internal/domain/` package SHALL import only the Go standard library. No AWS SDK, Redis, HTTP, or other infrastructure packages SHALL appear in any file under `internal/domain/`.

#### Scenario: Domain compiles without infrastructure dependencies
- **WHEN** `go build ./internal/domain/...` is run in the backend directory
- **THEN** the build succeeds with no imports from `github.com/aws`, `github.com/redis`, or `github.com/go-chi`

### Requirement: Application layer depends only on domain and port interfaces
The `internal/application/` package SHALL be a single flat package (no subdirectories). It SHALL import `internal/domain/` only. It SHALL NOT import any adapter package or AWS SDK directly. Outbound port interfaces (`ObjectStore`, `EventPublisher`, `Cache`) SHALL be defined inside `internal/application/` alongside the service interfaces that consume them.

#### Scenario: Use cases are constructable with interface implementations
- **WHEN** a service struct is instantiated with mock implementations of its port interfaces
- **THEN** it compiles and can be called without any real infrastructure present

### Requirement: Each use case is a distinct type with an Execute method
Each application service (`UploadService`, `CatalogService`, `ProcessingService`) SHALL be defined as a Go interface in `internal/application/` with named methods for each operation. The concrete implementation SHALL be an unexported struct. Constructors SHALL return the interface type.

#### Scenario: UploadService is independently constructable via its interface
- **WHEN** `NewUploadService` is called with implementations of `VideoRepository`, `ObjectStore`, `EventPublisher`, and a bucket name
- **THEN** it returns an `UploadService` interface value that can execute without knowledge of DynamoDB, S3, or HTTP specifics

### Requirement: Outbound port interfaces are defined as minimal Go interfaces
`VideoRepository`, `ObjectStore`, `Cache`, and `Queue` SHALL each be defined as Go interfaces with only the methods required by the application layer. No interface SHALL have more than 6 methods. No interface method SHALL be unused — every method SHALL have at least one call site in `internal/application/`.

#### Scenario: Adapter satisfies port interface
- **WHEN** the DynamoDB adapter is compiled against `VideoRepository`
- **THEN** the Go compiler confirms it satisfies the interface without explicit declaration

#### Scenario: Every port method has a caller
- **WHEN** `internal/ports/outbound.go` is inspected
- **THEN** every method on `VideoRepository`, `ObjectStore`, `Cache`, and `Queue` is called from at least one file under `internal/application/`

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

### Requirement: Composition root contains no HTTP transport construction
`cmd/api/main.go` SHALL NOT construct the HTTP router, register middleware, define routes, or call `http.ListenAndServe` directly. It SHALL only load configuration, construct outbound adapters and use cases, construct the inbound adapters (including the HTTP server via the `internal/adapters/inbound/http` package), and start them. Router construction, middleware registration (logging, recovery, CORS), route registration, strict-handler wiring, and the listen/serve loop SHALL live in `internal/adapters/inbound/http/`.

#### Scenario: main.go delegates router construction and serving to the HTTP adapter
- **WHEN** `cmd/api/main.go` is inspected
- **THEN** it constructs a server via `internal/adapters/inbound/http` and calls its `Start()` method, with no `chi.NewRouter`, middleware registration, route definitions, or `http.ListenAndServe` call present in `main.go`

### Requirement: Generated code lives under gen/
All generated code in `backend/api` (oapi-codegen output, mockery mocks) SHALL live under `gen/`, in one subdirectory per generator (`gen/api/`, `gen/mocks/`). `gen/` SHALL be a top-level directory alongside `cmd/`, `internal/`, and `openapi/`. Hand-written packages under `internal/` (`internal/domain/`, `internal/application/`, `internal/infrastructure/`) SHALL NOT contain generated code.

#### Scenario: Generated packages are isolated under gen/
- **WHEN** `backend/api/` is inspected
- **THEN** the oapi-codegen output is at `gen/api/api.gen.go` and the mockery mocks are at `gen/mocks/`, with no generated files under `internal/domain/`, `internal/application/`, or `internal/infrastructure/`

### Requirement: OpenAPI spec and codegen config live under openapi/
The OpenAPI specification, oapi-codegen config, and `go:generate` directive for `backend/api` SHALL live at `backend/api/openapi/` (`openapi.yaml`, `oapi-codegen.yaml`, `generate.go`). No directory named `api/` SHALL exist directly under `backend/api/`.

#### Scenario: OpenAPI sources are isolated under openapi/
- **WHEN** `backend/api/` is inspected
- **THEN** `openapi.yaml`, `oapi-codegen.yaml`, and `generate.go` are at `backend/api/openapi/`, and oapi-codegen output is generated at `gen/api/api.gen.go`

### Requirement: Outbound adapters delegate environment-specific URL building to an injected strategy
The `s3store` adapter in `backend/api` SHALL accept a `PresignedURLTransformer` (an interface with `Transform(presignedURL string) string`) and call it unconditionally — no if/else on `S3_PUBLIC_ENDPOINT_URL` inside the store. The composition root (`cmd/api/main.go`) SHALL select the concrete implementation (`NoOpTransformer` for production, `LocalStackTransformer` for local dev) based on config.

#### Scenario: Store contains no environment branching
- **WHEN** `backend/api/internal/adapters/outbound/s3store/store.go` is inspected
- **THEN** no conditional on `s3PublicEndpointURL` or any endpoint/env string exists — only an unconditional call to the injected `PresignedURLTransformer`

#### Scenario: Composition root selects the URL strategy
- **WHEN** `cmd/api/main.go` is inspected
- **THEN** it constructs either `NoOpTransformer` or `LocalStackTransformer` based on config and injects it into `s3store.NewStore`
