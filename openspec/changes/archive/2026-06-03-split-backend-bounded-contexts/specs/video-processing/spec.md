## MODIFIED Requirements

### Requirement: Worker splits video into HLS segments
The worker SHALL use ffmpeg to split the raw video into 6-second HLS segments at three quality levels: 1080p (4000k), 720p (2500k), and 360p (800k). Segments SHALL be uploaded to S3 under a `segments/{videoId}/` prefix.

#### Scenario: Successful segmentation
- **WHEN** the worker processes a valid video file
- **THEN** ffmpeg produces `.ts` segment files for all three quality levels in S3

#### Scenario: Corrupt or unreadable video
- **WHEN** ffmpeg cannot decode the video file
- **THEN** the worker emits a `VideoFailed` integration event to the results queue and deletes the SQS processing message

### Requirement: Worker updates video metadata on completion
After all pipeline steps succeed, the worker SHALL emit a `VideoProcessed` integration event to `video-processing-results.fifo` containing the CloudFront manifest URL and thumbnail URL. The worker SHALL NOT write to DynamoDB.

#### Scenario: Metadata updated after successful processing
- **WHEN** the full processing DAG completes without error
- **THEN** a `VideoProcessed` event is on the results queue with valid manifestUrl and thumbnailUrl; no DynamoDB write occurs
