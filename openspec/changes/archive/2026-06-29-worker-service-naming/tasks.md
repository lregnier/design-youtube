## Worker: Rename ProcessVideo → VideoProcessingService

### 1. Rename use case file and types

- [x] 1.1 Rename `process_video.go` → `video_processing_service.go`
- [x] 1.2 In `video_processing_service.go`: rename interface `ProcessVideo` → `VideoProcessingService`, struct `processVideo` → `videoProcessingService`, constructor `NewProcessVideo` → `NewVideoProcessingService`, method `Execute` → `Process`

### 2. Rename test files

- [x] 2.1 Rename `process_video_test.go` → `video_processing_service_test.go`; update all `NewProcessVideo` → `NewVideoProcessingService` and `uc.Execute` → `uc.Process`
- [x] 2.2 Rename `process_video_internal_test.go` → `video_processing_service_internal_test.go`

### 3. Update call sites

- [x] 3.1 Update `backend/worker/internal/infrastructure/in/sqssubscriber/subscriber.go`: field type `application.ProcessVideo` → `application.VideoProcessingService`; method call `pv.Execute` → `pv.Process`; `NewSubscriber` parameter type
- [x] 3.2 Update `backend/worker/cmd/worker/main.go`: `NewProcessVideo` → `NewVideoProcessingService`

### 4. Update mocks

- [x] 4.1 Update `backend/worker/.mockery.yaml`: `ProcessVideo` → `VideoProcessingService`
- [x] 4.2 Run `mockery` from `backend/worker/`; delete `gen/mocks/mock_ProcessVideo.go`
- [x] 4.3 Update `backend/worker/internal/application/video_processing_service_test.go`: `mocks.NewMockProcessVideo` → `mocks.NewMockVideoProcessingService`

### 5. Verify worker

- [x] 5.1 Run `go build ./...` from `backend/worker/` and confirm clean
- [x] 5.2 Run `go test ./...` from `backend/worker/` and confirm all tests pass

---

## API: Rename ProcessingService → VideoStatusService

### 6. Rename use case file and types

- [x] 6.1 Rename `processing_service.go` → `video_status_service.go`
- [x] 6.2 In `video_status_service.go`: rename interface `ProcessingService` → `VideoStatusService`, struct `processingService` → `videoStatusService`, constructor `NewProcessingService` → `NewVideoStatusService`, method `HandleVideoProcessingSucceeded` → `MarkReady`, method `HandleVideoProcessingFailed` → `MarkFailed`

### 7. Rename test file

- [x] 7.1 Rename `processing_service_test.go` → `video_status_service_test.go`; update all `NewProcessingService` → `NewVideoStatusService`, `HandleVideoProcessingSucceeded` → `MarkReady`, `HandleVideoProcessingFailed` → `MarkFailed`, and test function names accordingly

### 8. Update call sites

- [x] 8.1 Update `backend/api/internal/infrastructure/in/sqssubscriber/subscriber.go`: field type `application.ProcessingService` → `application.VideoStatusService`; method calls `svc.HandleVideoProcessingSucceeded` → `svc.MarkReady`, `svc.HandleVideoProcessingFailed` → `svc.MarkFailed`; `NewSubscriber` parameter type
- [x] 8.2 Update `backend/api/cmd/api/main.go`: `NewProcessingService` → `NewVideoStatusService`

### 9. Update mocks

- [x] 9.1 Update `backend/api/.mockery.yaml`: `ProcessingService` → `VideoStatusService`
- [x] 9.2 Run `mockery` from `backend/api/`; delete `gen/mocks/mock_ProcessingService.go`
- [x] 9.3 Update any test files that reference `mocks.NewMockProcessingService` → `mocks.NewMockVideoStatusService`

### 10. Verify api

- [x] 10.1 Run `go build ./...` from `backend/api/` and confirm clean
- [x] 10.2 Run `go test ./...` from `backend/api/` and confirm all tests pass

---

## README: Sync terminology

### 11. Update READMEs

- [x] 11.1 Update `backend/worker/README.md`: replace `ProcessVideo` with `VideoProcessingService` and `Execute` with `Process` in diagrams and text
- [x] 11.2 Update `backend/api/README.md`: replace `ProcessingService` with `VideoStatusService` and `HandleVideoProcessingSucceeded`/`HandleVideoProcessingFailed` with `MarkReady`/`MarkFailed` in diagrams and text
