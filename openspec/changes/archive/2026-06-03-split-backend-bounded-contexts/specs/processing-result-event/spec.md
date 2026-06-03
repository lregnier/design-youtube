## ADDED Requirements

### Requirement: Worker emits a VideoProcessed event on success
After successfully completing the full processing pipeline (segmentation, transcoding, manifest generation, thumbnail extraction), the worker SHALL publish a `VideoProcessed` message to the `video-processing-results.fifo` SQS queue. The message SHALL include: `eventType` ("VideoProcessed"), `videoId`, `manifestUrl` (CloudFront URL), and `thumbnailUrl` (CloudFront URL).

#### Scenario: Successful processing emits VideoProcessed
- **WHEN** the worker completes the full pipeline without error
- **THEN** a `VideoProcessed` JSON message is placed on `video-processing-results.fifo` with the correct videoId, manifestUrl, and thumbnailUrl

### Requirement: Worker emits a VideoFailed event on pipeline failure
If ffmpeg fails to decode or process the video, the worker SHALL publish a `VideoFailed` message to the `video-processing-results.fifo` SQS queue. The message SHALL include: `eventType` ("VideoFailed"), `videoId`, and `reason` (error description). The worker SHALL NOT write to DynamoDB.

#### Scenario: ffmpeg failure emits VideoFailed
- **WHEN** ffmpeg cannot decode or process the raw video file
- **THEN** a `VideoFailed` JSON message is placed on `video-processing-results.fifo` with the videoId and reason

#### Scenario: Worker has no DynamoDB dependency
- **WHEN** the worker service starts and processes jobs
- **THEN** it makes no DynamoDB API calls under any circumstances

### Requirement: Integration event schema is versioned
Both events SHALL include a top-level `eventType` string field that consumers can use to discriminate message type without attempting full deserialization. The schema SHALL be:

```json
// VideoProcessed
{
  "eventType": "VideoProcessed",
  "videoId": "<uuid>",
  "manifestUrl": "<cloudfront-url>",
  "thumbnailUrl": "<cloudfront-url>"
}

// VideoFailed
{
  "eventType": "VideoFailed",
  "videoId": "<uuid>",
  "reason": "<error description>"
}
```

#### Scenario: Consumer discriminates event type via eventType field
- **WHEN** the API consumer receives a message from the results queue
- **THEN** it reads `eventType` first and routes to the appropriate handler without full deserialization errors
