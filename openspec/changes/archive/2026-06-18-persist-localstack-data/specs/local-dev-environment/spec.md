## ADDED Requirements

### Requirement: LocalStack state persists across container restarts
The LocalStack container SHALL use a named Docker volume mounted at `/var/lib/localstack` and `PERSISTENCE=1` so that S3 objects, DynamoDB items, and SQS queue definitions survive `docker compose stop` / `docker compose start`.

#### Scenario: Data survives container restart
- **WHEN** a video is uploaded and processed, then `docker compose stop` and `docker compose start` are run
- **THEN** the video record and its S3 assets remain accessible without re-uploading

#### Scenario: Init script does not fail on restart
- **WHEN** LocalStack restores persisted state and the init script runs
- **THEN** the script completes successfully without errors, skipping creation of resources that already exist

### Requirement: Redis state persists across container restarts
The Redis container SHALL use a named Docker volume mounted at `/data` so that cached entries survive container restarts.

#### Scenario: Cache survives restart
- **WHEN** `docker compose stop` and `docker compose start` are run
- **THEN** previously cached video metadata entries are still available in Redis without requiring a DynamoDB round-trip
