## ADDED Requirements

### Requirement: S3 path-style access is configurable via environment variable
The worker SHALL read the S3 path-style flag from the `S3_USE_PATH_STYLE` environment variable. When the value is `true` or `1`, the S3 client SHALL use path-style addressing. When unset or any other value, the S3 client SHALL use virtual-hosted style (the AWS default).

#### Scenario: Path-style enabled for LocalStack
- **WHEN** `S3_USE_PATH_STYLE` is set to `"true"`
- **THEN** the worker's S3 client issues requests as `http://host/bucket/key`

#### Scenario: Path-style disabled for production
- **WHEN** `S3_USE_PATH_STYLE` is unset or `"false"`
- **THEN** the worker's S3 client uses virtual-hosted style (`https://bucket.s3.amazonaws.com/key`)
