## Why

The video processing worker has three silent failure modes that leave videos permanently stuck in `processing` status with no recovery path: `PublishFailed` errors are swallowed, infrastructure errors retry indefinitely with no dead-letter escape, and long transcodes can exceed the SQS visibility timeout causing duplicate concurrent processing.

## What Changes

- **Fix `PublishFailed` error handling**: check the return value in `process_video.go`; if publishing the failure event itself fails, return the error so the job stays in the queue and retries rather than being silently deleted
- **Add a Dead Letter Queue (DLQ)**: provision `video-processing-dlq.fifo` in Terraform and wire it to `video-processing.fifo` with `maxReceiveCount = 3`; messages that fail 3 times are moved to the DLQ instead of looping forever
- **Add a visibility timeout heartbeat**: periodically extend the SQS visibility timeout during processing so long transcodes do not cause SQS to redeliver the message while the original task is still running; the consumer must pass the `ReceiptHandle` to the use case for this purpose

## Capabilities

### New Capabilities

<!-- No new user-facing capabilities -->

### Modified Capabilities

- `video-processing`: Requirements change for retry behaviour, DLQ, and duplicate-processing prevention

## Impact

- `backend/worker/internal/application/process_video.go` — check `PublishFailed` errors; accept `ReceiptHandle` or a heartbeat callback for visibility extension
- `backend/worker/internal/adapters/inbound/sqsjobs/consumer.go` — pass receipt handle / heartbeat function; wire visibility extension goroutine
- `backend/worker/internal/ports/outbound.go` — possibly add a `HeartbeatFunc` type or extend consumer interface
- `infra/aws/sqs.tf` — add DLQ queue resource and `redrive_policy` on `video-processing.fifo`
- `infra/aws/iam.tf` — grant worker IAM permission to call `ChangeMessageVisibility` and `sqs:*` on the DLQ
