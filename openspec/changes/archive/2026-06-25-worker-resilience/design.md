## Context

The worker consumer (`sqsjobs/consumer.go`) polls `video-processing.fifo` with a 900 s visibility timeout. `process_video.go` handles domain failures by calling `publisher.PublishFailed` and returning `nil` (message deleted), and infrastructure failures by returning an error (message stays, retried). Two silent failure modes exist: (1) `PublishFailed` return values are ignored — if the results queue is unavailable the job is deleted and the video is permanently stuck; (2) there is no DLQ, so a persistently-failing job loops forever. A third issue: 900 s may not cover all transcodes, causing SQS to redeliver a message while the original task is still running.

## Goals / Non-Goals

**Goals:**
- `PublishFailed` failures surface as errors, keeping the job in the queue for retry
- Persistently-failing jobs are moved to a DLQ after 3 attempts instead of looping forever
- Long-running processing jobs extend their own visibility timeout so SQS does not redeliver prematurely

**Non-Goals:**
- DLQ alerting or operator tooling (out of scope for now)
- Automatic re-processing from the DLQ (manual re-drive via AWS console is sufficient)
- Retry backoff between attempts (SQS FIFO does not support per-message delay; visibility timeout provides natural spacing)

## Decisions

**Check `PublishFailed` errors and return them**
If `PublishFailed` fails, the correct behaviour is to leave the message in the queue so SQS redelivers it after the visibility timeout. The DLQ ensures it doesn't loop indefinitely. This is a one-line fix at each call site.

**DLQ: FIFO queue with `maxReceiveCount = 3`**
`video-processing.fifo` requires a FIFO DLQ — AWS mandates matching queue types for redrive policies. `maxReceiveCount = 3` gives three attempts before the message is quarantined. Three is enough to survive transient failures (network blip, S3 throttle) while catching persistent issues quickly.

*Alternative considered*: higher count (e.g. 5). Rejected — more retries mean longer time-to-detect for hard failures; 3 strikes the right balance.

**Heartbeat: goroutine extending visibility every 5 minutes**
The consumer spawns a goroutine alongside `processVideo.Execute` that calls `ChangeMessageVisibility` every 5 minutes with a fresh 900 s timeout. The goroutine is stopped (via `context.CancelFunc`) when `Execute` returns. This approach requires no changes to the `ProcessVideo` use case interface — the heartbeat lives entirely in the consumer layer.

*Alternative considered*: pass a heartbeat callback into `Execute`. Rejected — it couples the application layer to SQS infrastructure details, breaking the hexagonal boundary.

**IAM: grant `sqs:ChangeMessageVisibility` and DLQ read permissions to the worker task role**
The worker ECS task role already has `sqs:ReceiveMessage`, `sqs:DeleteMessage`. `ChangeMessageVisibility` must be added explicitly. The DLQ itself does not need to be consumed by the worker — no additional policy for that.

## Risks / Trade-offs

[Heartbeat goroutine leak] If `processVideo.Execute` hangs indefinitely the heartbeat goroutine also hangs → Mitigated by passing the consumer's poll `ctx` to the heartbeat; when the ECS task is stopped the context is cancelled.

[DLQ messages require manual re-drive] Messages in the DLQ are not automatically retried → Acceptable; persistent failures need human investigation before re-processing anyway.

[3 retries may re-transcode a long video multiple times] Infrastructure retries re-run the full pipeline → Acceptable; S3 uploads are idempotent (same key, same content overwrites cleanly).
