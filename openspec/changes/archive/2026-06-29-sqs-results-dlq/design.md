## Context

`video-processing.fifo` already has a DLQ (`video-processing-dlq.fifo`) with `maxReceiveCount=3` configured via Terraform. `video-processing-results.fifo` has no equivalent, so a message that fails processing on every delivery loops forever. The fix is a new DLQ queue and a `redrive_policy` on the results queue — two resource additions to `sqs.tf`.

## Goals / Non-Goals

**Goals:**
- Bound retries on the results queue to 3 deliveries
- Give operators a recoverable queue (DLQ) to inspect and replay failed messages
- Match the pattern already in place for the job queue

**Non-Goals:**
- Application-level retry logic — SQS handles this entirely
- Alerting or monitoring on the DLQ (separate concern)
- Changing `maxReceiveCount` on the existing job queue DLQ

## Decisions

**Separate DLQ per queue** — `video-processing-results-dlq.fifo` is a distinct resource rather than reusing `video-processing-dlq.fifo`. Mixing messages from the job queue and results queue in one DLQ makes it harder to diagnose failures and replay selectively.

**`maxReceiveCount=3`** — matches the job queue. Fast idempotent operations (DynamoDB writes) that fail 3 times in a row indicate a systemic problem, not a transient blip.

**FIFO DLQ** — AWS requires the DLQ for a FIFO queue to also be FIFO.

## Risks / Trade-offs

**Messages move to DLQ silently** → operators must monitor DLQ depth via CloudWatch to catch persistent failures. Acceptable for now; alerting can be added later.
