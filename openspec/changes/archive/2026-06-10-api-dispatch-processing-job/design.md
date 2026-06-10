## Context

The current design intended for S3 to fire an `s3:ObjectCreated:CompleteMultipartUpload` event to the `video-processing.fifo` SQS queue whenever a raw video upload finished. This fails unconditionally: AWS (and LocalStack) prohibit S3 event notifications from targeting FIFO queues. The queue type cannot be changed to standard — FIFO ordering and deduplication are required for the processing pipeline.

The API already holds all job context (videoId, S3 key) at the moment it calls `CompleteMultipartUpload`. The worker's `parseJob` already accepts a raw `{"videoId","s3Key"}` message as a fallback format. The outbound `sqsqueue` adapter and `ports.Queue` interface exist but currently only expose `DeleteMessage`.

## Goals / Non-Goals

**Goals:**
- API sends a processing job message to `video-processing.fifo` immediately after completing the multipart upload
- Worker receives and processes the job without any changes
- LocalStack init script no longer configures S3→SQS notifications (they never worked)

**Non-Goals:**
- Changing the worker's queue type or consumer logic
- Adding retry/dead-letter logic for failed dispatches (out of scope)
- Changing the results queue or results consumer

## Decisions

**API dispatches the job, not S3**
S3 FIFO notification is impossible. The API is the only other system with the complete job context at completion time. Direct dispatch is also simpler: no SNS fanout, no Lambda relay, no standard→FIFO bridge needed.

*Alternative considered*: use a standard SQS queue as the notification target, then bridge to FIFO. Rejected — adds a component with no other benefit; the API already owns this event.

**Reuse existing `ports.Queue` interface, add `SendMessage`**
`Queue` already has `DeleteMessage`. Adding `SendMessage(ctx, body, messageGroupID string) error` keeps the port cohesive and the adapter consistent. The `MessageDeduplicationId` is set to `videoId` — unique per video, safe for content-based dedup.

*Alternative considered*: raw SQS client injected into the use case. Rejected — breaks the hexagonal boundary; adapters exist precisely to hide SDK details.

**Message format: raw job JSON `{"videoId","s3Key"}`**
The worker's `parseJob` already handles this as its fallback path. No worker changes required.

## Risks / Trade-offs

[At-least-once delivery] → If the API call to `SendMessage` fails after `CompleteMultipartUpload` succeeds, the video is stuck in `processing` status with no job in the queue. Mitigation: log the error clearly; a future retry endpoint or admin tool can re-enqueue. This risk existed with S3 notifications too (the notification could also fail).

[FIFO deduplication window] → SQS FIFO deduplicates within a 5-minute window using `MessageDeduplicationId = videoId`. A second complete-upload call for the same video within 5 minutes will silently drop the duplicate message. This is the correct behavior — double-processing a video would be worse.
