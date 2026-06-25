## MODIFIED Requirements

### Requirement: Failed jobs are retried via SQS
If the worker crashes or returns an error without deleting the SQS message, SQS SHALL redeliver the message after the visibility timeout expires. The worker SHALL be idempotent — re-running the pipeline for the same videoId SHALL produce the same output without duplicating S3 objects. After 3 failed receive attempts, SQS SHALL move the message to `video-processing-dlq.fifo` instead of redelivering it.

The SQS consumer SHALL own the retry decision. The consumer SHALL read `ApproximateReceiveCount` from the incoming message and compare it against `maxReceiveCount` (3). On domain failure:
- If `ApproximateReceiveCount < maxReceiveCount`: the consumer SHALL log the failure and return without deleting the message, allowing SQS to redeliver.
- If `ApproximateReceiveCount >= maxReceiveCount`: the consumer SHALL call `PublishFailed` to emit a `VideoFailed` event, then delete the message.

The `ProcessVideo` use case SHALL NOT call `PublishFailed` — it SHALL return plain errors for all domain failures.

#### Scenario: Worker crash triggers redeliver
- **WHEN** a worker task is interrupted before deleting the SQS message
- **THEN** SQS redelivers the message after the visibility timeout and another worker picks it up

#### Scenario: Persistent failure routes to DLQ
- **WHEN** a job message has been received and failed 3 times
- **THEN** SQS moves the message to `video-processing-dlq.fifo` and stops redelivering it to the worker

#### Scenario: Transient failure — early attempt, no failure event
- **WHEN** a domain failure occurs and `ApproximateReceiveCount` is below `maxReceiveCount`
- **THEN** the consumer logs the error, does NOT delete the message, and does NOT emit a `VideoFailed` event; SQS redelivers after the visibility timeout

#### Scenario: Persistent failure — final attempt, failure event emitted
- **WHEN** a domain failure occurs and `ApproximateReceiveCount` equals `maxReceiveCount`
- **THEN** the consumer emits a `VideoFailed` event via `PublishFailed` and deletes the SQS message
