## MODIFIED Requirements

### Requirement: S3 path-style access is configurable via environment variable
The worker SHALL use path-style S3 addressing when `S3_ENDPOINT_URL` is set, and virtual-hosted style (the AWS default) when it is unset. The `S3_USE_PATH_STYLE` environment variable is removed.

#### Scenario: Path-style used with custom S3 endpoint
- **WHEN** `S3_ENDPOINT_URL` is set to `"http://minio:9000"`
- **THEN** the worker's S3 client issues requests as `http://minio:9000/bucket/key`

#### Scenario: Virtual-hosted style used in production
- **WHEN** `S3_ENDPOINT_URL` is unset
- **THEN** the worker's S3 client uses virtual-hosted style (`https://bucket.s3.amazonaws.com/key`)

### Requirement: Published asset URLs are configurable via environment variable
The worker SHALL publish `manifestUrl` and `thumbnailUrl` as path-style URLs against `S3_PUBLIC_URL` (`{S3_PUBLIC_URL}/{bucket}/{key}`) when that variable is set. When unset, the worker SHALL publish URLs as `https://{CLOUDFRONT_DOMAIN}/{key}`.

`S3_PUBLIC_URL` is the host-reachable address (e.g. `http://localhost:9000`), distinct from `S3_ENDPOINT_URL` which is the internal Docker-network address used by the SDK client (e.g. `http://minio:9000`).

#### Scenario: Public S3 endpoint configured for local dev
- **WHEN** `S3_PUBLIC_URL` is set to `"http://localhost:9000"`
- **THEN** `manifestUrl` and `thumbnailUrl` are published as `http://localhost:9000/{bucket}/{key}`

#### Scenario: CloudFront used in production
- **WHEN** `S3_PUBLIC_URL` is unset
- **THEN** `manifestUrl` and `thumbnailUrl` are published as `https://{CLOUDFRONT_DOMAIN}/{key}`

## REMOVED Requirements

### Requirement: LocalStack mode is configurable via environment variables (worker)
**Reason**: LocalStack is replaced by dedicated services. `LOCALSTACK_ENABLED` + `LOCALSTACK_ENDPOINT` are superseded by `S3_ENDPOINT_URL`.
**Migration**: Replace `LOCALSTACK_ENABLED=true` + `LOCALSTACK_ENDPOINT=http://localhost:4566` with `S3_ENDPOINT_URL=http://localhost:9000`.

## ADDED Requirements

### Requirement: Per-service AWS endpoint URLs are configurable via environment variables
The worker SHALL read optional endpoint overrides for each AWS service: `S3_ENDPOINT_URL`, `DYNAMODB_ENDPOINT_URL` (unused by the worker but symmetric with api), and `SQS_ENDPOINT_URL`. When set, the corresponding AWS SDK client SHALL use that endpoint.

#### Scenario: Custom S3 endpoint configured
- **WHEN** `S3_ENDPOINT_URL` is set to `"http://minio:9000"`
- **THEN** the worker's S3 client sends all requests to `http://minio:9000`

#### Scenario: Custom SQS endpoint configured
- **WHEN** `SQS_ENDPOINT_URL` is set to `"http://elasticmq:9324"`
- **THEN** the worker's SQS client sends all requests to `http://elasticmq:9324`
