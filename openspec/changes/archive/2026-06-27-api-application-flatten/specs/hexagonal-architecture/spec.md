## MODIFIED Requirements

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
