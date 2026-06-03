## ADDED Requirements

### Requirement: API service polls the results queue for processing events
The API service SHALL run a background SQS long-poll loop against `video-processing-results.fifo` alongside the HTTP server. The consumer SHALL start when the API process starts and run until the process terminates.

#### Scenario: Consumer starts with the API process
- **WHEN** the API ECS task starts
- **THEN** both the HTTP server and the SQS results consumer are running within the same process

### Requirement: VideoProcessed event marks the Video aggregate as ready
When the consumer receives a `VideoProcessed` event, it SHALL load the `Video` aggregate by `videoId` via `VideoRepository`, call `MarkReady(manifestURL, thumbnailURL)`, and save the updated aggregate. The consumer SHALL delete the SQS message after a successful save.

#### Scenario: VideoProcessed transitions video to ready
- **WHEN** the consumer receives a valid `VideoProcessed` message
- **THEN** the Video aggregate status is `ready`, manifestUrl and thumbnailUrl are set, and the SQS message is deleted

#### Scenario: VideoProcessed is idempotent on already-ready video
- **WHEN** the consumer receives a `VideoProcessed` message for a video already in `ready` status
- **THEN** the aggregate is saved with the same values and the message is deleted without error

### Requirement: VideoFailed event marks the Video aggregate as failed
When the consumer receives a `VideoFailed` event, it SHALL load the `Video` aggregate, call `MarkFailed()`, and save. The consumer SHALL delete the SQS message after a successful save.

#### Scenario: VideoFailed transitions video to failed
- **WHEN** the consumer receives a valid `VideoFailed` message
- **THEN** the Video aggregate status is `failed` and the SQS message is deleted

### Requirement: Consumer retries on transient errors
If the consumer fails to process a message (e.g. DynamoDB unavailable), it SHALL NOT delete the SQS message. SQS SHALL redeliver the message after the visibility timeout. The consumer SHALL log the error and continue polling.

#### Scenario: DynamoDB error leaves message in queue
- **WHEN** the consumer receives a message but DynamoDB returns an error during save
- **THEN** the SQS message is not deleted and is redelivered after the visibility timeout
