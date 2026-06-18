## MODIFIED Requirements

### Requirement: Published asset URLs are configurable via environment variable
The worker SHALL read `LOCALSTACK_ENABLED` (bool) and `LOCALSTACK_ENDPOINT` (string) from the environment. When `LOCALSTACK_ENABLED` is `true` or `1`, the worker SHALL publish `manifestUrl` and `thumbnailUrl` as path-style URLs against the configured `LOCALSTACK_ENDPOINT` (`{LOCALSTACK_ENDPOINT}/{bucket}/{key}`). When unset or `false`, the worker SHALL publish these URLs as `https://{CLOUDFRONT_DOMAIN}/{key}`.

#### Scenario: LocalStack enabled
- **WHEN** `LOCALSTACK_ENABLED` is set to `"true"` and `LOCALSTACK_ENDPOINT` is set to `"http://localhost:4566"`
- **THEN** `manifestUrl` and `thumbnailUrl` are published as `http://localhost:4566/{bucket}/{key}`

#### Scenario: LocalStack disabled for production
- **WHEN** `LOCALSTACK_ENABLED` is unset or `"false"`
- **THEN** `manifestUrl` and `thumbnailUrl` are published as `https://{CLOUDFRONT_DOMAIN}/{key}`

## REMOVED Requirements

### Requirement: Published asset URLs are configurable via environment variable (S3_PUBLIC_ENDPOINT_URL)
**Reason**: Replaced by explicit `LOCALSTACK_ENABLED` + `LOCALSTACK_ENDPOINT` variables.
**Migration**: Replace `S3_PUBLIC_ENDPOINT_URL=http://localhost:4566` with `LOCALSTACK_ENABLED=true` and `LOCALSTACK_ENDPOINT=http://localhost:4566`.
