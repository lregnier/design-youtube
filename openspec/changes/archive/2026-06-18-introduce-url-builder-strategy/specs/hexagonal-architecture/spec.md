## ADDED Requirements

### Requirement: Outbound adapters delegate environment-specific URL building to an injected strategy
The `s3store` adapter in `backend/api` SHALL accept a `PresignedURLTransformer` (an interface with `Transform(presignedURL string) string`) and call it unconditionally — no if/else on `S3_PUBLIC_ENDPOINT_URL` inside the store. The composition root (`cmd/api/main.go`) SHALL select the concrete implementation (`NoOpTransformer` for production, `LocalStackTransformer` for local dev) based on config.

#### Scenario: Store contains no environment branching
- **WHEN** `backend/api/internal/adapters/outbound/s3store/store.go` is inspected
- **THEN** no conditional on `s3PublicEndpointURL` or any endpoint/env string exists — only an unconditional call to the injected `PresignedURLTransformer`

#### Scenario: Composition root selects the URL strategy
- **WHEN** `cmd/api/main.go` is inspected
- **THEN** it constructs either `NoOpTransformer` or `LocalStackTransformer` based on config and injects it into `s3store.NewStore`
