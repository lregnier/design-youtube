## 1. Rename infrastructure directories

- [x] 1.1 Move `backend/worker/internal/adapters/inbound/sqsjobs/` → `backend/worker/internal/infrastructure/in/sqsjobs/`; update `package sqsjobs` if needed
- [x] 1.2 Move `backend/worker/internal/adapters/outbound/ffmpeg/` → `backend/worker/internal/infrastructure/out/ffmpeg/`
- [x] 1.3 Move `backend/worker/internal/adapters/outbound/s3storage/` → `backend/worker/internal/infrastructure/out/s3storage/`
- [x] 1.4 Move `backend/worker/internal/adapters/outbound/sqspublisher/` → `backend/worker/internal/infrastructure/out/sqspublisher/`
- [x] 1.5 Delete the now-empty `backend/worker/internal/adapters/` directory

## 2. Merge ports into application

- [x] 2.1 Copy `VideoStorage`, `Transcoder`, and `ResultPublisher` interface definitions from `backend/worker/internal/ports/outbound.go` into `backend/worker/internal/application/process_video.go` (or a new `backend/worker/internal/application/ports.go`)
- [x] 2.2 Delete `backend/worker/internal/ports/`
- [x] 2.3 Update `backend/worker/internal/infrastructure/out/sqspublisher/publisher.go`: replace `ports.ResultPublisher` import/reference with `application.ResultPublisher`
- [x] 2.4 Update `backend/worker/internal/infrastructure/in/sqsjobs/consumer.go`: replace `ports.ResultPublisher` import/reference with `application.ResultPublisher`

## 3. Move domain events to the domain layer

- [x] 3.1 Create `backend/worker/internal/domain/processing/events.go` with `VideoProcessingSucceededEvent{VideoID, ManifestURL, ThumbnailURL string}` and `VideoProcessingFailedEvent{VideoID, Reason string}` — plain Go structs, no JSON tags
- [x] 3.2 Update `backend/worker/internal/infrastructure/out/sqspublisher/publisher.go`: replace `event.VideoProcessed`/`event.VideoFailed` with unexported local wire structs that carry JSON tags; map from domain event fields; remove `internal/event` import
- [x] 3.3 Delete `backend/worker/internal/event/`

## 4. Make ProcessVideo a use case interface

- [x] 4.1 In `backend/worker/internal/application/process_video.go`: define `ProcessVideo` as a Go interface with `Execute(ctx context.Context, job processing.ProcessingJob) error`
- [x] 4.2 Rename the existing `ProcessVideo` struct to unexported `processVideo`; update all method receivers
- [x] 4.3 Update `NewProcessVideo` to return `ProcessVideo` (interface)
- [x] 4.4 Update `backend/worker/internal/infrastructure/in/sqsjobs/consumer.go`: change `processVideo application.ProcessVideo` field type to the interface (no struct change needed if already using the same type name)

## 5. Replace internal mocks with mockery

- [x] 5.1 Create `backend/worker/.mockery.yaml` with output `dir: gen/mocks`, `outpkg: mocks`, `with-expecter: true`; list interfaces `VideoStorage`, `Transcoder`, `ResultPublisher`, `ProcessVideo` from `github.com/lregnier/design-youtube/worker/internal/application`
- [x] 5.2 Run `mockery` from `backend/worker/` to generate `gen/mocks/`
- [x] 5.3 Delete `backend/worker/internal/mocks/`
- [x] 5.4 Update `backend/worker/internal/application/process_video_test.go`: replace `worker/internal/mocks` import with `worker/gen/mocks`
- [x] 5.5 Switch `process_video_test.go` to `package application_test`; add explicit `application.` qualifier to `NewProcessVideo`; keep `buildMasterManifest` test in a separate internal test file or export it if needed — simplest: keep `TestBuildMasterManifest_AllQualities` in `package application` in a new `process_video_internal_test.go`

## 6. Move config into cmd/worker/

- [x] 6.1 Copy `backend/worker/internal/config/config.go` to `backend/worker/cmd/worker/config.go`; change `package config` to `package main`
- [x] 6.2 Update `backend/worker/cmd/worker/main.go`: remove `internal/config` import; replace `config.Load()` with `Load()`
- [x] 6.3 Delete `backend/worker/internal/config/`

## 7. Update main.go import paths

- [x] 7.1 Update `backend/worker/cmd/worker/main.go`: replace all `internal/adapters/inbound/` imports with `internal/infrastructure/in/` and `internal/adapters/outbound/` with `internal/infrastructure/out/`

## 8. Normalize inbound SQS adapter naming

- [x] 8.1 Rename `backend/worker/internal/infrastructure/in/sqsjobs/` → `sqsconsumer/`; change `package sqsjobs` → `package sqsconsumer`
- [x] 8.2 Rename `consumer.go` and `consumer_test.go` file references (package declaration only; file names stay)
- [x] 8.3 Update `backend/worker/cmd/worker/main.go`: replace `sqsjobs` import path and `sqsjobs.NewConsumer` with `sqsconsumer.NewConsumer`

## 9. Align inbound SQS adapter name with api pattern

- [x] 9.1 Rename `sqsconsumer/` → `sqssubscriber/`; change `package sqsconsumer` → `package sqssubscriber`
- [x] 9.2 Rename files `consumer.go` → `subscriber.go` and `consumer_test.go` → `subscriber_test.go`; rename type `Consumer` → `Subscriber` and constructor `NewConsumer` → `NewSubscriber`
- [x] 9.3 Update `backend/worker/cmd/worker/main.go`: replace `sqsconsumer` import and `sqsconsumer.NewConsumer` with `sqssubscriber.NewSubscriber`

## 10. Use domain event types in ResultPublisher

- [x] 10.1 Update `ResultPublisher` in `backend/worker/internal/application/ports.go`: change signatures to `PublishProcessed(ctx, evt processing.VideoProcessingSucceededEvent) error` and `PublishFailed(ctx, evt processing.VideoProcessingFailedEvent) error`
- [x] 10.2 Update call sites in `backend/worker/internal/application/process_video.go`: construct and pass event structs instead of raw strings
- [x] 10.3 Update `backend/worker/internal/infrastructure/in/sqssubscriber/subscriber.go`: pass `processing.VideoProcessingFailedEvent` to `PublishFailed`
- [x] 10.4 Update `backend/worker/internal/infrastructure/out/sqspublisher/publisher.go`: accept event structs; map fields to wire structs
- [x] 10.5 Regenerate mocks: run `mockery` from `backend/worker/`
- [x] 10.6 Update `backend/worker/internal/application/process_video_test.go`: adjust `PublishProcessed` / `PublishFailed` mock expectations to pass event structs

## 11. Rename ResultPublisher → EventPublisher

- [x] 11.1 Rename `ResultPublisher` → `EventPublisher` in `backend/worker/internal/application/ports.go`
- [x] 11.2 Update all references in `backend/worker/internal/application/process_video.go`
- [x] 11.3 Update `backend/worker/internal/infrastructure/out/sqspublisher/publisher.go`: `application.ResultPublisher` → `application.EventPublisher`
- [x] 11.4 Update `backend/worker/internal/infrastructure/in/sqssubscriber/subscriber.go`: `application.ResultPublisher` → `application.EventPublisher`
- [x] 11.5 Update `backend/worker/.mockery.yaml`: rename `ResultPublisher` → `EventPublisher`
- [x] 11.6 Regenerate mocks: run `mockery` from `backend/worker/`; delete `gen/mocks/mock_ResultPublisher.go`
- [x] 11.7 Update `backend/worker/internal/application/process_video_test.go`: `mocks.NewMockResultPublisher` → `mocks.NewMockEventPublisher`

## 12. Align EventPublisher with api single-method interface

- [x] 12.1 Add `DomainEvent` interface with private `domainEvent()` marker to `backend/worker/internal/domain/processing/events.go`; implement it on `VideoProcessingSucceededEvent` and `VideoProcessingFailedEvent`
- [x] 12.2 Update `EventPublisher` in `backend/worker/internal/application/ports.go`: replace two-method signature with `Publish(ctx context.Context, event processing.DomainEvent) error`
- [x] 12.3 Update `backend/worker/internal/application/process_video.go`: call `publisher.Publish(ctx, processing.VideoProcessingSucceededEvent{...})`
- [x] 12.4 Update `backend/worker/internal/infrastructure/in/sqssubscriber/subscriber.go`: call `publisher.Publish(ctx, processing.VideoProcessingFailedEvent{...})`
- [x] 12.5 Update `backend/worker/internal/infrastructure/out/sqspublisher/publisher.go`: replace `PublishProcessed`/`PublishFailed` with a single `Publish` method that type-switches on `processing.DomainEvent`
- [x] 12.6 Regenerate mocks: run `mockery` from `backend/worker/`
- [x] 12.7 Update `backend/worker/internal/application/process_video_test.go`: adjust mock expectations to use `Publish` with event struct

## 13. Verify

- [x] 13.1 Run `go build ./...` from `backend/worker/` and confirm clean
- [x] 13.2 Run `go test ./...` from `backend/worker/` and confirm all tests pass
