## ADDED Requirements

### Requirement: Backend returns video metadata including manifest URL
The backend SHALL expose `GET /videos/{videoId}`. For videos with status `ready`, the response SHALL include: videoId, title, description, status, CloudFront manifest URL, CloudFront thumbnail URL, and upload timestamp. For videos not yet ready, the response SHALL include status only.

#### Scenario: Ready video returns full metadata
- **WHEN** a client requests `GET /videos/{videoId}` for a video with status `ready`
- **THEN** the server returns 200 with all metadata fields including manifest URL and thumbnail URL

#### Scenario: Processing video returns partial metadata
- **WHEN** a client requests `GET /videos/{videoId}` for a video with status `processing`
- **THEN** the server returns 200 with status `processing` and no manifest URL

#### Scenario: Unknown video returns 404
- **WHEN** a client requests `GET /videos/{videoId}` for a non-existent videoId
- **THEN** the server returns 404

### Requirement: Metadata is served from Redis cache with DynamoDB fallback
The backend SHALL check Redis for the video metadata before querying DynamoDB. Cache entries SHALL have a TTL of 60 seconds. On a cache miss, the backend SHALL fetch from DynamoDB and populate the cache.

#### Scenario: Cache hit serves metadata without DynamoDB query
- **WHEN** a video metadata record exists in Redis
- **THEN** the backend returns the cached record and does not query DynamoDB

#### Scenario: Cache miss falls back to DynamoDB
- **WHEN** a video metadata record is not in Redis
- **THEN** the backend queries DynamoDB, returns the result, and writes it to Redis

### Requirement: Frontend player streams video using HLS.js with adaptive bitrate
The frontend video page SHALL use HLS.js to load the master manifest URL and play the video. HLS.js SHALL automatically select the appropriate bitrate segment based on the client's current network conditions. The player SHALL support play, pause, and seek.

#### Scenario: Player loads and starts playback
- **WHEN** a user navigates to a video page for a ready video
- **THEN** HLS.js fetches the master manifest, selects an initial quality level, and begins buffering and playing

#### Scenario: Player adapts quality on bandwidth change
- **WHEN** the client's available bandwidth drops significantly during playback
- **THEN** HLS.js switches to a lower bitrate segment without interrupting playback

#### Scenario: Player displays thumbnail before playback starts
- **WHEN** a user navigates to a video page before pressing play
- **THEN** the thumbnail image is displayed as the video poster
