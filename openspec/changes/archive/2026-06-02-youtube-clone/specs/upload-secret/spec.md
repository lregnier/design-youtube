## ADDED Requirements

### Requirement: Upload endpoints require a shared secret header
All endpoints that initiate or advance a video upload SHALL require an `X-Upload-Secret` HTTP header. The backend SHALL compare the header value against the `UPLOAD_SECRET` environment variable using a constant-time comparison. Requests with a missing or incorrect secret SHALL be rejected before any S3 or DynamoDB interaction occurs.

#### Scenario: Valid secret allows upload initiation
- **WHEN** a client sends `POST /videos/upload/init` with a correct `X-Upload-Secret` header
- **THEN** the server returns 200 with the presigned URL data

#### Scenario: Missing secret is rejected
- **WHEN** a client sends `POST /videos/upload/init` without an `X-Upload-Secret` header
- **THEN** the server returns 401 Unauthorized

#### Scenario: Wrong secret is rejected
- **WHEN** a client sends `POST /videos/upload/init` with an incorrect `X-Upload-Secret` value
- **THEN** the server returns 401 Unauthorized

#### Scenario: Read endpoints are public
- **WHEN** a client sends `GET /videos` or `GET /videos/{videoId}` without any secret header
- **THEN** the server returns the requested data normally

### Requirement: Secret is configured via environment variable
The `UPLOAD_SECRET` value SHALL be read from the `UPLOAD_SECRET` environment variable at server startup. The server SHALL refuse to start if the variable is missing or empty.

#### Scenario: Server starts with secret configured
- **WHEN** the server starts with `UPLOAD_SECRET` set to a non-empty value
- **THEN** the server initializes successfully and begins accepting requests

#### Scenario: Server refuses to start without secret
- **WHEN** the server starts with `UPLOAD_SECRET` unset or empty
- **THEN** the server logs a fatal error and exits with a non-zero status code
