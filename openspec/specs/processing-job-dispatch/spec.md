## Requirements

### Requirement: API dispatches processing job on upload completion
When a multipart upload is successfully completed, the API SHALL publish a job message to the video-processing FIFO queue. The message body SHALL be JSON with `videoId` and `s3Key` fields. The `MessageGroupId` SHALL be set to the videoId. The `MessageDeduplicationId` SHALL be set to the videoId.

#### Scenario: Job enqueued after successful complete
- **WHEN** the client calls `POST /videos/{videoId}/upload/complete` and `CompleteMultipartUpload` succeeds in S3
- **THEN** a message `{"videoId":"<id>","s3Key":"raw/<id>/original"}` is sent to the processing FIFO queue and the video status is updated to `processing`

#### Scenario: S3 completion fails — no job enqueued
- **WHEN** the S3 `CompleteMultipartUpload` call returns an error
- **THEN** no message is sent to the queue and the error is returned to the client

#### Scenario: Queue send fails after S3 completion
- **WHEN** S3 `CompleteMultipartUpload` succeeds but the SQS `SendMessage` call fails
- **THEN** the error is returned to the client and the video remains in `processing` status in the repository (no rollback of S3 or DynamoDB)
