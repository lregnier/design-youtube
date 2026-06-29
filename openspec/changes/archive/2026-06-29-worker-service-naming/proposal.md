## Why

The worker's use case `ProcessVideo` and the api's `ProcessingService` do not follow the `*Service` suffix convention used by the other api services (`UploadService`, `CatalogService`). The worker's method is also named `Execute`, a generic verb that conveys no domain meaning, while api service methods use domain-specific names. `ProcessingService` only manages status transitions and its current name doesn't reflect that.

## What Changes

### Worker

- Rename `ProcessVideo` interface → `VideoProcessingService`
- Rename `processVideo` unexported struct → `videoProcessingService`
- Rename `NewProcessVideo` constructor → `NewVideoProcessingService`
- Rename `Execute` method → `Process`
- Update all call sites: `sqssubscriber`, `.mockery.yaml`, generated mocks, tests
- Rename source file `process_video.go` → `video_processing_service.go`

### API

- Rename `ProcessingService` interface → `VideoStatusService`
- Rename `processingService` unexported struct → `videoStatusService`
- Rename `NewProcessingService` constructor → `NewVideoStatusService`
- Rename `HandleVideoProcessingSucceeded` → `MarkReady`, `HandleVideoProcessingFailed` → `MarkFailed`
- Update all call sites: `sqssubscriber`, `.mockery.yaml`, generated mocks, tests
- Rename source file `processing_service.go` → `video_status_service.go`

### READMEs

- Update `backend/worker/README.md` and `backend/api/README.md` to reflect renamed types and methods

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `worker-hexagonal-architecture`: Update requirement wording to reflect `VideoProcessingService` and `Process` method name

## Impact

- Worker and api application layers only; no wire format, domain logic, or infrastructure changes
