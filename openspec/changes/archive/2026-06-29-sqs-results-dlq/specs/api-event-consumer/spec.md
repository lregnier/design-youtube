## MODIFIED Requirements

### Requirement: Consumer retries on transient errors
If the consumer fails to process a message (e.g. DynamoDB unavailable), it SHALL NOT delete the SQS message. SQS SHALL redeliver the message after the visibility timeout. The consumer SHALL log the error and continue polling. After `maxReceiveCount` (3) failed deliveries, SQS SHALL automatically move the message to `video-processing-results-dlq.fifo` for manual inspection and replay.

#### Scenario: DynamoDB error leaves message in queue
- **WHEN** the consumer receives a message but DynamoDB returns an error during save
- **THEN** the SQS message is not deleted and is redelivered after the visibility timeout

#### Scenario: Persistently failing message is moved to DLQ
- **WHEN** a message has been delivered and failed 3 times
- **THEN** SQS moves the message to `video-processing-results-dlq.fifo` and it is no longer redelivered to the consumer
