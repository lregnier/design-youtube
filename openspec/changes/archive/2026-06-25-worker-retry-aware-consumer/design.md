## Context

After the `worker-resilience` change, `process_video.go` calls `PublishFailed` at each domain failure site (download, ffprobe, transcode). This means the application layer decides when a failure is permanent — it embeds retry semantics (implicitly: "first failure = permanent"). The SQS consumer's `ApproximateReceiveCount` is unused. The hexagonal boundary is violated: the use case is coupled to the SQS retry model.

## Goals / Non-Goals

**Goals:**
- `process_video.go` returns plain errors for all failures — no `PublishFailed` call sites
- The SQS consumer reads `ApproximateReceiveCount` and is the single place that decides permanent failure
- `PublishFailed` is called by the consumer on the final attempt; earlier attempts log and let SQS redeliver
- The application layer is fully decoupled from SQS retry mechanics

**Non-Goals:**
- Changing the DLQ configuration or `maxReceiveCount` (stays at 3)
- Adding per-failure-type retry strategies (all domain failures treated equally)
- Backoff between retries (SQS FIFO does not support per-message delay)

## Decisions

**Consumer injects and owns `ResultPublisher`**
The consumer receives `ResultPublisher` as a constructor argument. On the final attempt it calls `publisher.PublishFailed` directly. This keeps the outbound adapter dependency in the inbound adapter layer — both are infrastructure, and the consumer already calls `processVideo.Execute` which is an application-layer dependency. Injecting a second outbound port here is consistent with how the application layer is wired in `main.go`.

*Alternative considered*: pass the receive count into `Execute` and let the use case decide. Rejected — this re-introduces SQS semantics into the application layer.

*Alternative considered*: a separate DLQ consumer that calls `PublishFailed` when it reads from the DLQ. Rejected — adds a second consumer process for a concern that's naturally handled at receive time; also increases time-to-detect since the message must reach the DLQ first.

**`ApproximateReceiveCount` as the retry signal**
SQS populates `ApproximateReceiveCount` as a message attribute on every `ReceiveMessage` call (when explicitly requested). The consumer compares it against `maxReceiveCount` (constant: 3, matching the DLQ redrive policy). No state management needed.

**On final attempt: call `PublishFailed`, then delete the message**
If `PublishFailed` itself fails on the final attempt, the consumer returns an error (message not deleted, goes to DLQ). The DLQ then acts as the final safety net rather than the primary failure-routing mechanism.

**`process_video.go` returns domain errors directly**
Each domain failure site becomes `return fmt.Errorf("download failed: %w", err)` etc. No `ResultPublisher` in the use case at all — the port dependency is removed from `ProcessVideo`.

## Risks / Trade-offs

[Video stays in `processing` status during retries] On earlier attempts the video is not marked failed — users see "processing" while SQS retries. → Acceptable: this is the correct UX for transient failures that will resolve.

[Final-attempt `PublishFailed` failure sends to DLQ without `failed` status] If `PublishFailed` fails on attempt 3, the message goes to DLQ and the video is stuck in `processing` forever. → Same risk as before; DLQ requires manual re-drive regardless.

[`ApproximateReceiveCount` is approximate] AWS documents this as best-effort. In rare cases the count may skip values. → At `maxReceiveCount = 3` this means at worst one extra retry before `PublishFailed` is called — acceptable.
