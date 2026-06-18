## ADDED Requirements

### Requirement: LocalStack mode is configurable via environment variables
The API SHALL read `LOCALSTACK_ENABLED` (bool) to determine whether LocalStack is active, and `LOCALSTACK_ENDPOINT` (string) for the LocalStack URL. When `LOCALSTACK_ENABLED` is `true` or `1`, the API SHALL rewrite presigned upload URLs to use the configured `LOCALSTACK_ENDPOINT` host. When unset or `false`, presigned URLs SHALL be returned unchanged.

#### Scenario: LocalStack enabled
- **WHEN** `LOCALSTACK_ENABLED` is set to `"true"` and `LOCALSTACK_ENDPOINT` is set to `"http://localhost:4566"`
- **THEN** presigned upload URLs have their host replaced with `http://localhost:4566`

#### Scenario: LocalStack disabled
- **WHEN** `LOCALSTACK_ENABLED` is unset or `"false"`
- **THEN** presigned upload URLs are returned unchanged
