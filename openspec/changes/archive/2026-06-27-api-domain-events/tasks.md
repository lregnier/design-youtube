## 1. Extend domain events

- [x] 1.1 Add `VideoProcessedEvent` and `VideoFailedEvent` to `backend/api/internal/domain/video/events.go` — plain Go structs, no JSON tags

## 2. Remove application events

- [x] 2.1 Delete `backend/api/internal/application/events.go`

## 3. Update ProcessingService

- [x] 3.1 Update `backend/api/internal/application/processing_service.go`: change `OnProcessed` and `OnFailed` signatures to accept `video.VideoProcessedEvent` and `video.VideoFailedEvent`; update method bodies accordingly

## 4. Update sqssubscriber adapter

- [x] 4.1 In `backend/api/internal/infrastructure/in/sqssubscriber/subscriber.go`: add unexported wire structs (`videoProcessedMessage`, `videoFailedMessage`) with JSON tags matching the SQS message format
- [x] 4.2 Update `handle` to unmarshal into wire structs and map to `video.VideoProcessedEvent` / `video.VideoFailedEvent` before calling the service

## 5. Regenerate mocks

- [x] 5.1 Run `mockery` from `backend/api/` to regenerate `MockProcessingService` with the updated interface signature

## 6. Update tests

- [x] 6.1 Update `backend/api/internal/application/processing_service_test.go`: replace `application.VideoProcessedEvent` / `application.VideoFailedEvent` with `video.VideoProcessedEvent` / `video.VideoFailedEvent`

## 7. Verify

- [x] 7.1 Run `go build ./...` and confirm clean
- [x] 7.2 Run `go test ./...` and confirm all tests pass

## 8. Rename processing result events

- [x] 8.1 Rename `VideoProcessedEvent` → `VideoProcessingSucceededEvent` and `VideoFailedEvent` → `VideoProcessingFailedEvent` in `backend/api/internal/domain/video/events.go`
- [x] 8.2 Update all references across `processing_service.go`, `sqssubscriber/subscriber.go`, `processing_service_test.go`
- [x] 8.3 Regenerate mocks
- [x] 8.4 Run `go build ./...` and `go test ./...` and confirm clean

## 9. Rename ProcessingService methods

- [x] 9.1 Rename `OnProcessed` → `HandleVideoProcessingSucceeded` and `OnFailed` → `HandleVideoProcessingFailed` in interface and implementation
- [x] 9.2 Update call sites in `sqssubscriber/subscriber.go` and `processing_service_test.go`
- [x] 9.3 Regenerate mocks and verify build + tests
