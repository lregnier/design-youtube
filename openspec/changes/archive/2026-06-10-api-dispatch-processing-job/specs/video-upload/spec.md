## MODIFIED Requirements

### Requirement: Client completes the multipart upload
The backend SHALL expose `POST /videos/{videoId}/upload/complete` (protected by upload secret). The backend SHALL call S3 CompleteMultipartUpload with all collected ETags. On success, the backend SHALL publish a processing job message to the video-processing FIFO queue and update the video metadata status to `processing`.

#### Scenario: Successful completion
- **WHEN** a client posts to complete with all chunk ETags collected
- **THEN** S3 assembles the file, the backend enqueues a processing job, updates the video status to `processing`, and returns 200

#### Scenario: Queue dispatch fails
- **WHEN** S3 CompleteMultipartUpload succeeds but the SQS SendMessage call fails
- **THEN** the backend returns an error to the client; the video status may be `processing` in the store but no job is in the queue
