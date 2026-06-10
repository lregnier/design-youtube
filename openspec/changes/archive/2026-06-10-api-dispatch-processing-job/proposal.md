## Why

S3 event notifications do not support FIFO SQS queues — this is a hard AWS limitation that affects both LocalStack and production. The current design routes `s3:ObjectCreated:CompleteMultipartUpload` events to `video-processing.fifo`, which always fails with `InvalidArgument`. Videos are uploaded successfully but the worker never receives a job and processing never starts.

## What Changes

- Add `SendMessage` to the `ports.Queue` interface so the API can enqueue jobs
- Implement `SendMessage` in the `sqsqueue` outbound adapter
- Extend `CompleteUpload` use case to dispatch a raw job message to the processing queue after completing the multipart upload
- Wire the queue into `CompleteUpload` in `main.go`
- Remove the S3→SQS event notification configuration from `scripts/localstack-init.sh`
- Remove the now-dead S3 event-notification fallback parsing from the worker's job consumer
- Remove the unused `DeleteMessage` method from `ports.Queue` / `sqsqueue.Queue`
- Update architecture diagrams (root and `backend/api` READMEs) to show the API enqueuing the job directly

## Capabilities

### New Capabilities

- `processing-job-dispatch`: API enqueues a video processing job directly to the FIFO queue when a multipart upload is completed

### Modified Capabilities

- `complete-upload`: After completing the multipart upload in S3, the use case now also publishes a processing job to the queue (new side-effect, same HTTP surface)
- `video-processing`: the worker requirement describing how jobs arrive is updated — the API publishes the job (per `processing-job-dispatch`), and the worker simply polls the queue

## Impact

- **`backend/api/internal/ports/outbound.go`**: `Queue` interface gains `SendMessage`, loses unused `DeleteMessage`
- **`backend/api/internal/adapters/outbound/sqsqueue/queue.go`**: new `SendMessage` implementation; `DeleteMessage` removed
- **`backend/api/internal/application/upload/complete.go`**: `CompleteUpload` struct and constructor gain a `queue` dependency
- **`backend/api/cmd/api/main.go`**: wiring update to pass queue to `CompleteUpload`
- **`scripts/localstack-init.sh`**: S3 notification config block removed; SQS queues retained
- **`backend/worker/internal/adapters/inbound/sqsjobs/consumer.go`**: removed dead S3 event-notification format parsing — jobs always arrive as raw `{"videoId","s3Key"}` JSON
- **`README.md`, `backend/api/README.md`**: diagrams updated to reflect direct job dispatch
