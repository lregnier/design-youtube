## ADDED Requirements

### Requirement: S3 upload completion triggers processing
When S3 receives a CompleteMultipartUpload event for the raw video prefix, it SHALL publish a notification to the SQS processing queue. The processing worker SHALL poll the queue and pick up the event to begin the processing DAG.

#### Scenario: Upload completion enqueues processing job
- **WHEN** S3 receives a CompleteMultipartUpload for a raw video object
- **THEN** an SQS message containing the videoId and S3 key is placed on the processing queue

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
After all segments are uploaded, the worker SHALL generate an HLS master playlist (`.m3u8`) referencing the three quality-level media playlists. The manifest SHALL be uploaded to S3 under `manifests/{videoId}/master.m3u8` and served via CloudFront.

#### Scenario: Master manifest created
- **WHEN** all segment uploads complete successfully
- **THEN** a valid HLS master manifest is present in S3 and accessible via CloudFront

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
If the worker crashes or returns an error without deleting the SQS message, SQS SHALL redeliver the message after the visibility timeout expires. The worker SHALL be idempotent — re-running the pipeline for the same videoId SHALL produce the same output without duplicating S3 objects.

#### Scenario: Worker crash triggers redeliver
- **WHEN** a worker task is interrupted before deleting the SQS message
- **THEN** SQS redelivers the message after the visibility timeout and another worker picks it up
