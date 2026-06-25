## Why

The application layer (`process_video.go`) currently calls `PublishFailed` directly on domain failures, meaning it owns retry semantics — deciding when a job is permanently failed. This couples infrastructure concerns (SQS retry count) into the domain/application layer, breaking hexagonal boundaries.

## What Changes

- **Remove `PublishFailed` calls from `process_video.go`**: domain failures (download, ffprobe, transcode) now return plain errors; the use case no longer decides whether a failure is permanent
- **Consumer reads `ApproximateReceiveCount`**: the SQS consumer checks the receive count on each incoming message and passes it to the retry decision
- **Consumer calls `PublishFailed` at max retries**: on the final attempt the consumer publishes the failure event itself, then deletes the message; on earlier attempts it logs and does not delete (SQS redelivers)
- **Inject `ResultPublisher` into the consumer**: the consumer needs access to the outbound publisher port to call `PublishFailed`

## Capabilities

### New Capabilities

<!-- No new user-facing capabilities -->

### Modified Capabilities

- `video-processing`: Requirement changes — retry behaviour and failure-event publishing now owned by the consumer, not the use case

## Impact

- `backend/worker/internal/application/process_video.go` — remove all `PublishFailed` call sites; return plain errors for domain failures
- `backend/worker/internal/application/process_video_test.go` — update tests: domain failure cases now expect an error returned (not `PublishFailed` called)
- `backend/worker/internal/adapters/inbound/sqsjobs/consumer.go` — request `ApproximateReceiveCount` attribute; inject `ResultPublisher`; call `PublishFailed` when receive count hits max
- `backend/worker/cmd/worker/main.go` — wire `ResultPublisher` into `Consumer`
