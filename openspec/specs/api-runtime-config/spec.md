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

### Requirement: HTTP listen address is configurable via environment variable
The API SHALL read its HTTP listen address from the `HTTP_ADDR` environment variable. When set, the HTTP server SHALL listen on that address. When unset, the HTTP server SHALL listen on `:8080`.

#### Scenario: Listen address configured
- **WHEN** `HTTP_ADDR` is set to `:9090`
- **THEN** the HTTP server listens on `:9090`

#### Scenario: Listen address unset
- **WHEN** `HTTP_ADDR` is unset
- **THEN** the HTTP server listens on `:8080`

### Requirement: LocalStack mode is configurable via environment variables
The API SHALL read `LOCALSTACK_ENABLED` (bool) to determine whether LocalStack is active, and `LOCALSTACK_ENDPOINT` (string) for the LocalStack URL. When `LOCALSTACK_ENABLED` is `true` or `1`, the API SHALL rewrite presigned upload URLs to use the configured `LOCALSTACK_ENDPOINT` host. When unset or `false`, presigned URLs SHALL be returned unchanged.

#### Scenario: LocalStack enabled
- **WHEN** `LOCALSTACK_ENABLED` is set to `"true"` and `LOCALSTACK_ENDPOINT` is set to `"http://localhost:4566"`
- **THEN** presigned upload URLs have their host replaced with `http://localhost:4566`

#### Scenario: LocalStack disabled
- **WHEN** `LOCALSTACK_ENABLED` is unset or `"false"`
- **THEN** presigned upload URLs are returned unchanged
