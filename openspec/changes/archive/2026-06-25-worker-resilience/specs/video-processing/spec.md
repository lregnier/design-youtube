## MODIFIED Requirements

### Requirement: Failed jobs are retried via SQS
If the worker crashes or returns an error without deleting the SQS message, SQS SHALL redeliver the message after the visibility timeout expires. The worker SHALL be idempotent — re-running the pipeline for the same videoId SHALL produce the same output without duplicating S3 objects. After 3 failed receive attempts, SQS SHALL move the message to `video-processing-dlq.fifo` instead of redelivering it.

#### Scenario: Worker crash triggers redeliver
- **WHEN** a worker task is interrupted before deleting the SQS message
- **THEN** SQS redelivers the message after the visibility timeout and another worker picks it up

#### Scenario: Persistent failure routes to DLQ
- **WHEN** a job message has been received and failed 3 times
- **THEN** SQS moves the message to `video-processing-dlq.fifo` and stops redelivering it to the worker

## ADDED Requirements

### Requirement: PublishFailed errors are propagated
If the worker encounters a domain failure (corrupt video, ffprobe error, transcode error) and the call to `PublishFailed` itself fails, the worker SHALL return an error so the SQS message is NOT deleted and the job is retried. The video SHALL NOT be left permanently stuck in `processing` status due to a transient results-queue outage.

#### Scenario: PublishFailed succeeds — message deleted
- **WHEN** a domain failure occurs and `PublishFailed` publishes successfully
- **THEN** the SQS job message is deleted and the video status will be updated to `failed` by the API consumer

#### Scenario: PublishFailed fails — message retained
- **WHEN** a domain failure occurs and `PublishFailed` returns an error
- **THEN** the worker returns the error, the SQS job message is NOT deleted, and SQS redelivers it after the visibility timeout

### Requirement: Worker extends visibility timeout during long processing
The worker consumer SHALL periodically extend the SQS message visibility timeout during processing so that long-running transcodes do not cause SQS to redeliver the message while the original task is still running.

#### Scenario: Visibility extended before timeout
- **WHEN** processing is ongoing and the visibility timeout is approaching
- **THEN** the consumer calls `ChangeMessageVisibility` to reset the timeout to 900 seconds

#### Scenario: Heartbeat stops on completion
- **WHEN** processing completes (success or error)
- **THEN** the visibility heartbeat goroutine is stopped and no further `ChangeMessageVisibility` calls are made
