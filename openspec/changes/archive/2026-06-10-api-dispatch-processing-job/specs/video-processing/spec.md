## REMOVED Requirements

### Requirement: S3 upload completion triggers processing
**Reason**: S3 event notifications cannot target FIFO SQS queues. The API now publishes the job to `video-processing.fifo` directly after `CompleteMultipartUpload` succeeds (see `processing-job-dispatch`).
**Migration**: No client-facing change. The worker continues to poll `video-processing.fifo`; see the new "Worker polls the processing queue for jobs" requirement.

## ADDED Requirements

### Requirement: Worker polls the processing queue for jobs
The processing worker SHALL poll the `video-processing.fifo` SQS queue and begin the processing DAG for each job message it receives.

#### Scenario: Worker picks up the job
- **WHEN** the ECS Fargate worker polls SQS and receives a message
- **THEN** the worker begins the processing pipeline for that videoId
