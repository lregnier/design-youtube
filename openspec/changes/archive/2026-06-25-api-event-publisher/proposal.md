## Why

The API's outbound queue port (`Queue.SendMessage`) exposes infrastructure details — raw message body string and group ID — directly to the application layer. The `CompleteUpload` use case manually marshals JSON and constructs queue-specific parameters, coupling domain logic to SQS message format.

## What Changes

- **Replace `Queue` port with `EventPublisher`**: the interface receives a typed domain event object instead of raw strings
- **Introduce `VideoUploadedEvent` domain event**: a struct in the domain layer capturing what happened (`VideoID`, `S3Key`) — no infrastructure concepts
- **Move JSON marshaling and routing to the SQS adapter**: the adapter translates the domain event to an SQS message; the use case is unaware of message format or queue URL
- **Remove `processingJob` struct from use case**: it belongs in the adapter as an infrastructure detail

## Capabilities

### New Capabilities

<!-- No new user-facing capabilities -->

### Modified Capabilities

<!-- No spec-level requirement changes — behavior is identical, this is an internal architectural refactor -->

## Impact

- `backend/api/internal/ports/outbound.go` — replace `Queue` interface with `EventPublisher`
- `backend/api/internal/domain/` — add `VideoUploadedEvent` domain event type
- `backend/api/internal/application/upload/complete.go` — replace `queue.SendMessage(...)` with `publisher.Publish(ctx, VideoUploadedEvent{...})`; remove `processingJob` struct and JSON marshaling
- `backend/api/internal/adapters/outbound/sqsqueue/queue.go` — rename to `sqspublisher`; implement `EventPublisher`; own JSON marshaling and routing
- `backend/api/internal/gen/mocks/mock_Queue.go` — regenerate mock for new interface
- `backend/api/internal/application/upload/complete_test.go` — update to use new mock
- `backend/api/cmd/api/main.go` — update wiring
