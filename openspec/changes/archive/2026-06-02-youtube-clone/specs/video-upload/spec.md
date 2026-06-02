## ADDED Requirements

### Requirement: Client initiates a multipart upload and receives presigned URLs
The backend SHALL expose `POST /videos/upload/init` (protected by upload secret). The request body SHALL include video title, description, total file size, and total chunk count. The backend SHALL reject requests where `fileSize` exceeds 104857600 bytes (100MB). On success, the backend SHALL create a multipart upload in S3, create a video metadata record in DynamoDB with status `uploading`, and return the videoId, S3 uploadId, and a presigned URL for the first chunk.

#### Scenario: Valid upload initiation under size limit
- **WHEN** a client posts a valid init request with `fileSize` ≤ 100MB
- **THEN** the server returns 200 with `videoId`, `uploadId`, and a presigned part URL for part 1

#### Scenario: File exceeds 100MB limit
- **WHEN** a client posts an init request with `fileSize` > 104857600
- **THEN** the server returns 400 with an error message stating the size limit

### Requirement: Client uploads chunks directly to S3 via presigned URLs
Each chunk SHALL be uploaded by the client directly to S3 using the presigned URL returned by the backend. Chunk size SHALL be between 5MB and 10MB except for the final chunk which MAY be smaller. S3 returns an ETag per successfully uploaded part.

#### Scenario: Successful chunk upload
- **WHEN** a client PUTs a chunk to the presigned S3 URL
- **THEN** S3 returns 200 with an `ETag` header for that part

### Requirement: Client confirms each uploaded chunk with the backend
The backend SHALL expose `POST /videos/{videoId}/upload/confirm-chunk`. The request body SHALL include `partNumber` and `eTag`. The backend SHALL mark that chunk as uploaded in the DynamoDB metadata record and return a presigned URL for the next chunk (if any remain).

#### Scenario: Chunk confirmed, more parts remain
- **WHEN** a client confirms a chunk and there are remaining parts
- **THEN** the server marks the chunk uploaded and returns a presigned URL for the next part

#### Scenario: Final chunk confirmed
- **WHEN** a client confirms the last chunk
- **THEN** the server marks all chunks uploaded and returns no next presigned URL

### Requirement: Client completes the multipart upload
The backend SHALL expose `POST /videos/{videoId}/upload/complete` (protected by upload secret). The backend SHALL call S3 CompleteMultipartUpload with all collected ETags. On success, the video metadata status SHALL be updated to `processing`.

#### Scenario: Successful completion
- **WHEN** a client posts to complete with all chunk ETags collected
- **THEN** S3 assembles the file, the backend updates status to `processing`, and returns 200

### Requirement: Upload is resumable
If a client calls `POST /videos/upload/init` with an existing `videoId` that has status `uploading`, the backend SHALL return the current chunk status and a presigned URL for the first not-yet-uploaded chunk, allowing the client to resume without re-uploading completed chunks.

#### Scenario: Resume after interruption
- **WHEN** a client re-initiates an upload for a videoId with some chunks already uploaded
- **THEN** the server returns the already-uploaded chunk list and a presigned URL starting from the first missing chunk
