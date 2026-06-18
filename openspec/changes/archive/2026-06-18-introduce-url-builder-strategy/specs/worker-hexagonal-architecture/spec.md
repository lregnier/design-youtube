## ADDED Requirements

### Requirement: Outbound adapters delegate environment-specific URL building to an injected strategy
The `s3storage` adapter in `backend/worker` SHALL accept a `PublicURLBuilder` (an interface with `AssetURL(bucket, key string) string`) and call it unconditionally — no if/else on `S3_PUBLIC_ENDPOINT_URL` or `CloudFrontDomain` inside the store. The composition root (`cmd/worker/main.go`) SHALL select the concrete implementation (`CloudFrontURLBuilder` for production, `LocalStackURLBuilder` for local dev) based on config.

#### Scenario: Store contains no environment branching
- **WHEN** `backend/worker/internal/adapters/outbound/s3storage/store.go` is inspected
- **THEN** no conditional on `s3PublicEndpointURL` or `cloudfrontDomain` exists — only an unconditional call to the injected `PublicURLBuilder`

#### Scenario: Composition root selects the URL strategy
- **WHEN** `cmd/worker/main.go` is inspected
- **THEN** it constructs either `CloudFrontURLBuilder` or `LocalStackURLBuilder` based on config and injects it into `s3storage.NewStore`
