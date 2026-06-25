## 1. Define the domain event and port

- [x] 1.1 In `backend/api/internal/domain/`, add `VideoUploadedEvent` struct with `VideoID` and `S3Key` fields, and a `DomainEvent` interface it satisfies
- [x] 1.2 In `backend/api/internal/ports/outbound.go`, replace the `Queue` interface with `EventPublisher` having a single `Publish(ctx context.Context, event domain.DomainEvent) error` method

## 2. Update the SQS adapter

- [x] 2.1 Rename `backend/api/internal/adapters/outbound/sqsqueue/` to `sqspublisher/` and update the package name
- [x] 2.2 Rewrite the adapter to implement `EventPublisher`: type-switch on the event, marshal to JSON, send to SQS; return an error for unknown event types

## 3. Update the application layer

- [x] 3.1 In `complete.go`, replace the `queue ports.Queue` dependency with `publisher ports.EventPublisher`; remove the `processingJob` struct and JSON marshaling; call `publisher.Publish(ctx, domain.VideoUploadedEvent{...})`

## 4. Update tests and mocks

- [x] 4.1 Regenerate the mock for `EventPublisher` (replace `mock_Queue.go` with `mock_EventPublisher.go`) using mockery
- [x] 4.2 Update `complete_test.go` to use the new mock and `VideoUploadedEvent`

## 5. Update wiring

- [x] 5.1 In `cmd/api/main.go`, update the import path and constructor call to use the renamed `sqspublisher` package
