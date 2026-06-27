## 1. Create flat application/ files — catalog slice

- [x] 1.1 Create `backend/api/internal/application/cache.go` (`package application`) — `Cache` interface with `Get` and `Set` methods
- [x] 1.2 Create `backend/api/internal/application/catalog_service.go` — `CatalogService` interface (renamed from `CatalogServicePort`) with `GetVideo` and `ListVideos`; unexported `catalogService` struct with `repo` and `cache` fields; `NewCatalogService(repo, cache) CatalogService` constructor; full method bodies ported from `catalog/service.go`

## 2. Create flat application/ files — upload slice

- [x] 2.1 Create `backend/api/internal/application/object_store.go` — `ObjectStore` interface; value types `MultipartUpload`, `PresignedURL`, `CompletedPart`
- [x] 2.2 Create `backend/api/internal/application/event_publisher.go` — `EventPublisher` interface with `Publish(ctx, video.DomainEvent) error`
- [x] 2.3 Create `backend/api/internal/application/upload_service.go` — `UploadService` interface (renamed from `UploadServicePort`) with `InitUpload`, `ConfirmChunk`, `CompleteUpload`; result types `InitUploadResult`, `ConfirmChunkResult`, `ChunkState`; `MaxFileSize` constant; unexported `uploadService` struct; `NewUploadService(...) UploadService` constructor; full method bodies ported from `upload/service.go`

## 3. Create flat application/ files — processing slice

- [x] 3.1 Create `backend/api/internal/application/events.go` — `VideoProcessedEvent` and `VideoFailedEvent` structs with JSON tags
- [x] 3.2 Create `backend/api/internal/application/processing_service.go` — `ProcessingService` interface (renamed from `ProcessingServicePort`) with `OnProcessed` and `OnFailed`; unexported `processingService` struct; `NewProcessingService(repo) ProcessingService` constructor; full method bodies ported from `processing/service.go`

## 4. Create flat application/ tests

- [x] 4.1 Create `backend/api/internal/application/catalog_service_test.go` (`package application_test`) — port all tests from `catalog/service_test.go`; replace `catalog.NewCatalogService` with `application.NewCatalogService`, `catalog.CatalogService` with `application.CatalogService`
- [x] 4.2 Create `backend/api/internal/application/upload_service_test.go` (`package application_test`) — port all tests from `upload/service_test.go`; replace `upload.*` references with `application.*`; update `mock.AnythingOfType` strings (e.g., `"[]upload.CompletedPart"` → `"[]application.CompletedPart"`)
- [x] 4.3 Create `backend/api/internal/application/processing_service_test.go` (`package application_test`) — port all tests from `processing/service_test.go`; replace `processing.*` references with `application.*`

## 5. Update mockery and regenerate mocks

- [x] 5.1 Update `backend/api/.mockery.yaml` — replace the three sub-package entries with a single `github.com/lregnier/design-youtube/api/internal/application` entry listing: `VideoRepository` (move from `domain/video` entry — keep that), `ObjectStore`, `EventPublisher`, `UploadService`, `Cache`, `CatalogService`, `ProcessingService`
- [x] 5.2 Run `mockery` in `backend/api/` to regenerate all mocks
- [x] 5.3 Delete stale mock files that correspond to old interface names (`mock_UploadServicePort.go`, `mock_CatalogServicePort.go`, `mock_ProcessingServicePort.go`)

## 6. Update infrastructure adapters

- [x] 6.1 Update `infrastructure/out/rediscache/cache.go` — change import from `application/catalog` to `application`; update compile-time check to `var _ application.Cache = (*Cache)(nil)`
- [x] 6.2 Update `infrastructure/out/s3store/store.go` — change import from `application/upload` to `application`; update all `upload.*` references to `application.*`
- [x] 6.3 Update `infrastructure/out/sqspublisher/publisher.go` — change import from `application/upload` to `application`; update compile-time check to `var _ application.EventPublisher = (*Publisher)(nil)`
- [x] 6.4 Update `infrastructure/in/http/handler.go` — change imports from `application/upload`, `application/catalog` to `application`; rename field types and `NewHandler` parameter types from `upload.UploadServicePort`/`catalog.CatalogServicePort` to `application.UploadService`/`application.CatalogService`
- [x] 6.5 Update `infrastructure/in/http/handler_test.go` — update mock type references from `mocks.MockUploadServicePort`/`mocks.MockCatalogServicePort` to `mocks.MockUploadService`/`mocks.MockCatalogService`; update `upload.*` result type references to `application.*`
- [x] 6.6 Update `infrastructure/in/sqsconsumer/consumer.go` — change import from `application/processing` to `application`; update `processing.ProcessingServicePort` to `application.ProcessingService`; update `processing.VideoProcessedEvent`/`processing.VideoFailedEvent` to `application.VideoProcessedEvent`/`application.VideoFailedEvent`

## 7. Update main.go

- [x] 7.1 Update `backend/api/cmd/api/main.go` — replace `application/catalog`, `application/upload`, `application/processing` imports with a single `application` import; update constructor calls to use the `application.` prefix

## 8. Delete old sub-packages and verify

- [x] 8.1 Delete `backend/api/internal/application/upload/`, `backend/api/internal/application/catalog/`, `backend/api/internal/application/processing/` directory trees
- [x] 8.2 Run `go build ./...` and confirm clean
- [x] 8.3 Run `go test ./...` and confirm all tests pass

## 9. Rename SQS infrastructure directories

- [x] 9.1 Rename `infrastructure/out/sqspublisher/` → `infrastructure/out/sqs/`; update package name to `sqs`
- [x] 9.2 Rename `infrastructure/in/sqsconsumer/` → `infrastructure/in/sqs/`; update package name to `sqs`
- [x] 9.3 Update `cmd/api/main.go` — alias AWS SDK as `awssqs`, internal packages as `sqsin`/`sqsout`; `go build ./...` clean

## 10. Adopt Publisher/Subscriber naming throughout

- [x] 10.1 Rename `infrastructure/out/sqsproducer/` → `infrastructure/out/sqspublisher/`; rename `Producer` struct → `Publisher`, `NewProducer` → `NewPublisher`; update package name to `sqspublisher`
- [x] 10.2 Rename `infrastructure/in/sqsconsumer/` → `infrastructure/in/sqssubscriber/`; rename `Consumer` struct → `Subscriber`, `NewConsumer` → `NewSubscriber`; update package name to `sqssubscriber`
- [x] 10.3 Rename `application/event_producer.go` → `event_publisher.go`; rename `EventProducer` interface → `EventPublisher`; update `upload_service.go` field type
- [x] 10.4 Update `.mockery.yaml` (`EventProducer` → `EventPublisher`), regenerate mocks, delete stale `mock_EventProducer.go`
- [x] 10.5 Update `cmd/api/main.go` imports and constructor calls
- [x] 10.6 Update all test files; `go test ./...` clean
