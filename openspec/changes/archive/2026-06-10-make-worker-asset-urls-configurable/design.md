## Context

The worker's `s3storage.Store` builds `manifestUrl`/`thumbnailUrl` as `https://{cloudfrontDomain}/{key}`, intended for production where `CLOUDFRONT_DOMAIN` is a real CloudFront distribution fronting the S3 bucket. In local dev, `CLOUDFRONT_DOMAIN=localhost` doesn't serve anything, so these URLs are unreachable from the browser. The API faced the analogous problem for presigned upload URLs and solved it with `S3_PUBLIC_ENDPOINT_URL`, which rewrites the scheme+host of an S3-generated URL to a browser-accessible LocalStack endpoint.

## Goals / Non-Goals

**Goals:**
- `S3_PUBLIC_ENDPOINT_URL` env var, when set, makes the worker publish path-style LocalStack URLs for `manifestUrl`/`thumbnailUrl`
- When unset (production), behavior is unchanged — URLs are still `https://{cloudfrontDomain}/{key}`
- Local dev docker-compose sets it for the `worker` service

**Non-Goals:**
- Rewriting URLs of *already-processed* videos (existing `ready` records keep their stored URL until reprocessed)
- Changing the API's existing `S3_PUBLIC_ENDPOINT_URL` / `rewriteHost` logic for presigned URLs

## Decisions

**Branch on presence of `s3PublicEndpointURL`, not a `rewriteHost`-style rewrite**
The API's `rewriteHost` rewrites an existing S3-generated URL's scheme+host. The worker doesn't have an S3-generated URL to rewrite — it constructs the URL itself from `cloudfrontDomain` + key. So instead of generating a CloudFront URL and rewriting it, the worker picks which template to use upfront:
- `s3PublicEndpointURL` set → `{s3PublicEndpointURL}/{bucket}/{key}` (path-style, e.g. `http://localhost:4566/design-youtube-video-prod/thumbnails/{id}/thumb.jpg`)
- unset → `https://{cloudfrontDomain}/{key}` (existing CloudFront behavior)

**`bucket` is already a `Store` field**
`NewStore(client, bucket, cloudfrontDomain)` already carries `bucket`, so no new field is needed beyond `s3PublicEndpointURL`.

**Same env var name as the API (`S3_PUBLIC_ENDPOINT_URL`)**
Keeps the local-dev mental model consistent: "this is the browser-reachable LocalStack endpoint", used the same way by both services even though the underlying code paths differ (rewrite vs. construct).

## Risks / Trade-offs

[Stale URLs on already-processed videos] → The one `ready` video in the local stack keeps its old `https://localhost/...` URLs until reprocessed. Mitigation: documented as a verification step (re-enqueue the job, as done for the previous path-style fix), not a code concern.
