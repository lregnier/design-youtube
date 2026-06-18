## Why

Both `backend/api` and `backend/worker` branch on `S3_PUBLIC_ENDPOINT_URL != ""` inside their S3 adapter methods to decide how to build public-facing URLs — either rewriting a presigned URL's host for LocalStack, or choosing between a LocalStack path-style URL and a CloudFront URL. This env-specific branching lives inside the store's core logic, coupling S3 operations to deployment knowledge. Replacing it with a strategy pattern moves that concern to a single decision point at startup (composition root), leaving the store focused solely on S3 operations.

## What Changes

**`backend/api` — presigned URL transformation:**
- Introduce `PresignedURLTransformer` interface (`Transform(url string) string`) in `internal/adapters/outbound/s3store/`
- Implement `NoOpTransformer` (production: returns URL unchanged) and `LocalStackTransformer` (replaces host with `S3_PUBLIC_ENDPOINT_URL`)
- `Store` accepts a `PresignedURLTransformer` instead of `s3PublicEndpointURL string`; `PresignUploadPart` calls `s.transformer.Transform(url)` unconditionally
- `cmd/api/main.go` selects the implementation based on config: `LocalStackTransformer` if `S3PublicEndpointURL != ""`, `NoOpTransformer` otherwise

**`backend/worker` — public asset URL building:**
- Introduce `PublicURLBuilder` interface (`AssetURL(bucket, key string) string`) in `internal/adapters/outbound/s3storage/`
- Implement `CloudFrontURLBuilder` (production: `https://<domain>/<key>`) and `LocalStackURLBuilder` (dev: `<endpoint>/<bucket>/<key>`)
- `Store` accepts a `PublicURLBuilder` instead of `cloudfrontDomain` + `s3PublicEndpointURL`; `assetURL` calls `s.urlBuilder.AssetURL(s.bucket, key)` unconditionally
- `cmd/worker/main.go` selects the implementation based on config

**Both modules:**
- Remove `s3PublicEndpointURL` and (for worker) `cloudfrontDomain` fields from the `Store` struct
- Delete `rewriteHost` helper (api) and the `if/else` in `assetURL` (worker)
- Update unit tests to inject stub implementations

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `hexagonal-architecture`: outbound adapters SHALL delegate environment-specific URL building to an injected strategy; the store itself SHALL contain no branching on deployment environment

## Impact

- Affected files: `backend/api/internal/adapters/outbound/s3store/store.go`, `backend/api/cmd/api/main.go`, `backend/worker/internal/adapters/outbound/s3storage/store.go`, `backend/worker/cmd/worker/main.go`
- New files: `backend/api/internal/adapters/outbound/s3store/url_transformer.go`, `backend/worker/internal/adapters/outbound/s3storage/url_builder.go`
- No change to port interfaces, domain, or application layer
- No behavior change for existing functionality
