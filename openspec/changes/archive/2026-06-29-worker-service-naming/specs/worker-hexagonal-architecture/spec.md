## MODIFIED Requirements

### Requirement: ProcessVideo use case depends only on port interfaces
The `VideoProcessingService` use case in `internal/application/` SHALL accept `VideoStorage`, `Transcoder`, and `EventPublisher` as constructor arguments. It SHALL NOT import any adapter package, AWS SDK, or ffmpeg directly.

#### Scenario: Use case is constructable with mock implementations
- **WHEN** `VideoProcessingService` is instantiated with stub implementations of its three port interfaces
- **THEN** it compiles and can be called without any real infrastructure present
