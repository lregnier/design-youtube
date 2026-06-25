## ADDED Requirements

### Requirement: Worker polls the processing queue for jobs
The processing worker SHALL poll the `video-processing.fifo` SQS queue and begin the processing DAG for each job message it receives.

#### Scenario: Worker picks up the job
- **WHEN** the ECS Fargate worker polls SQS and receives a message
- **THEN** the worker begins the processing pipeline for that videoId

### Requirement: Worker splits video into HLS segments
The worker SHALL use ffmpeg to split the raw video into 6-second HLS segments at three quality levels: 1080p (4000k), 720p (2500k), and 360p (800k). Segments SHALL be uploaded to S3 under a `segments/{videoId}/` prefix.

#### Scenario: Successful segmentation
- **WHEN** the worker processes a valid video file
- **THEN** ffmpeg produces `.ts` segment files for all three quality levels in S3

#### Scenario: Corrupt or unreadable video
- **WHEN** ffmpeg cannot decode the video file
- **THEN** the worker emits a `VideoFailed` integration event to the results queue and deletes the SQS processing message

### Requirement: Worker generates an HLS master manifest
After all segments are uploaded, the worker SHALL generate an HLS master playlist (`.m3u8`) referencing the three quality-level media playlists. The manifest SHALL be uploaded to S3 under `manifests/{videoId}/master.m3u8` and served via CloudFront. Variant playlist references in the master manifest SHALL be relative paths that resolve, against the master manifest's own location, to the variant playlists' actual location at `segments/{videoId}/{quality}/media.m3u8`.

#### Scenario: Master manifest created
- **WHEN** all segment uploads complete successfully
- **THEN** a valid HLS master manifest is present in S3 and accessible via CloudFront

#### Scenario: Variant playlist references resolve correctly
- **WHEN** a player resolves a variant playlist reference from the master manifest relative to the manifest's own URL
- **THEN** the resolved URL points to the existing `segments/{videoId}/{quality}/media.m3u8` object

### Requirement: Worker extracts a thumbnail
The worker SHALL use ffmpeg to extract a single JPEG frame from the video's midpoint (duration / 2). The thumbnail SHALL be uploaded to S3 under `thumbnails/{videoId}/thumb.jpg` and served via CloudFront.

#### Scenario: Thumbnail extracted and stored
- **WHEN** the worker completes processing
- **THEN** a JPEG thumbnail is present at the CloudFront thumbnail URL

### Requirement: Worker updates video metadata on completion
After all pipeline steps succeed, the worker SHALL emit a `VideoProcessed` integration event to `video-processing-results.fifo` containing the CloudFront manifest URL and thumbnail URL. The worker SHALL NOT write to DynamoDB.

#### Scenario: Metadata updated after successful processing
- **WHEN** the full processing DAG completes without error
- **THEN** a `VideoProcessed` event is on the results queue with valid manifestUrl and thumbnailUrl; no DynamoDB write occurs

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

### Requirement: Worker extends visibility timeout during long processing
The worker consumer SHALL periodically extend the SQS message visibility timeout during processing so that long-running transcodes do not cause SQS to redeliver the message while the original task is still running.

#### Scenario: Visibility extended before timeout
- **WHEN** processing is ongoing and the visibility timeout is approaching
- **THEN** the consumer calls `ChangeMessageVisibility` to reset the timeout to 900 seconds

#### Scenario: Heartbeat stops on completion
- **WHEN** processing completes (success or error)
- **THEN** the visibility heartbeat goroutine is stopped and no further `ChangeMessageVisibility` calls are made
