## ADDED Requirements

### Requirement: Published asset URLs are configurable via environment variable
The worker SHALL read a public S3 endpoint override from the `S3_PUBLIC_ENDPOINT_URL` environment variable. When set, the worker SHALL publish `manifestUrl` and `thumbnailUrl` as path-style URLs against that endpoint (`{S3_PUBLIC_ENDPOINT_URL}/{bucket}/{key}`). When unset, the worker SHALL publish these URLs as `https://{CLOUDFRONT_DOMAIN}/{key}` (the AWS default for production).

#### Scenario: Public endpoint configured for LocalStack
- **WHEN** `S3_PUBLIC_ENDPOINT_URL` is set to `http://localhost:4566`
- **THEN** `manifestUrl` and `thumbnailUrl` are published as `http://localhost:4566/{bucket}/{key}`

#### Scenario: Public endpoint unset for production
- **WHEN** `S3_PUBLIC_ENDPOINT_URL` is unset
- **THEN** `manifestUrl` and `thumbnailUrl` are published as `https://{CLOUDFRONT_DOMAIN}/{key}`
