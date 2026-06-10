## ADDED Requirements

### Requirement: CORS allowed origins are configurable via environment variable
The API SHALL read allowed CORS origins from the `CORS_ALLOWED_ORIGINS` environment variable. The value SHALL be a comma-separated list of origin strings. The API SHALL refuse to start if `CORS_ALLOWED_ORIGINS` is not set.

#### Scenario: Single origin configured
- **WHEN** `CORS_ALLOWED_ORIGINS` is set to `https://app.example.com`
- **THEN** the API allows cross-origin requests only from `https://app.example.com`

#### Scenario: Wildcard configured for local dev
- **WHEN** `CORS_ALLOWED_ORIGINS` is set to `*`
- **THEN** the API allows cross-origin requests from any origin

#### Scenario: Variable not set
- **WHEN** `CORS_ALLOWED_ORIGINS` is unset at startup
- **THEN** the API fails to start with a clear error message

### Requirement: S3 path-style access is configurable via environment variable
The API SHALL read the S3 path-style flag from the `S3_USE_PATH_STYLE` environment variable. When the value is `true` or `1`, the S3 client SHALL use path-style addressing. When unset or any other value, the S3 client SHALL use virtual-hosted style (the AWS default).

#### Scenario: Path-style enabled for LocalStack
- **WHEN** `S3_USE_PATH_STYLE` is set to `"true"`
- **THEN** the S3 client issues requests as `http://host/bucket/key`

#### Scenario: Path-style disabled for production
- **WHEN** `S3_USE_PATH_STYLE` is unset or `"false"`
- **THEN** the S3 client uses virtual-hosted style (`https://bucket.s3.amazonaws.com/key`)
