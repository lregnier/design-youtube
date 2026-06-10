## Why

The worker publishes `manifestUrl`/`thumbnailUrl` as `https://{CLOUDFRONT_DOMAIN}/{key}`. In local dev `CLOUDFRONT_DOMAIN=localhost`, but nothing serves HTTPS on `localhost:443` — these URLs return connection-refused in the browser (confirmed via curl: `https://localhost/thumbnails/.../thumb.jpg` fails, while the LocalStack path-style equivalent `http://localhost:4566/design-youtube-video-prod/thumbnails/.../thumb.jpg` returns 200). The API solved the same class of problem for presigned upload URLs via `S3_PUBLIC_ENDPOINT_URL`; the worker needs the same escape hatch for the asset URLs it publishes.

## What Changes

- Add `S3PublicEndpointURL string` to the worker config, parsed from `S3_PUBLIC_ENDPOINT_URL` env var (optional, defaults to empty)
- In `s3storage.Store`: when `s3PublicEndpointURL` is set, build `manifestUrl`/`thumbnailUrl` as `{s3PublicEndpointURL}/{bucket}/{key}` (path-style); when unset, keep the existing `https://{cloudfrontDomain}/{key}` (CloudFront) behavior
- In `cmd/worker/main.go`: pass `cfg.S3PublicEndpointURL` into `s3storage.NewStore`
- In `docker-compose.yml`: set `S3_PUBLIC_ENDPOINT_URL: "http://localhost:4566"` for the `worker` service, matching the `api` service
- Document the new env var in `backend/worker/README.md`

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `worker-runtime-config`: add "Published asset URLs are configurable via environment variable" requirement

## Impact

- **`backend/worker/internal/config/config.go`**: new `S3PublicEndpointURL` field
- **`backend/worker/internal/adapters/outbound/s3storage/store.go`**: conditional URL construction in `UploadManifest`/`UploadThumbnail`
- **`backend/worker/cmd/worker/main.go`**: pass new config field to `s3storage.NewStore`
- **`docker-compose.yml`**: new env var in the `worker` service block
- **`backend/worker/README.md`**: configuration docs

## Note

This only affects `manifestUrl`/`thumbnailUrl` for *new* processing results. The existing `ready` video in the local stack has stale `https://localhost/...` URLs from before this fix and will need reprocessing to pick up corrected URLs (a verification step, not a code task).
