## Context

`buildMasterManifest` (`backend/worker/internal/application/process_video.go`) writes the master playlist uploaded to `manifests/{videoId}/master.m3u8`. It currently references each variant playlist as `segments/{videoId}/{quality}/media.m3u8`. HLS players resolve relative URLs against the *manifest's own directory*, so this resolves to `manifests/{videoId}/segments/{videoId}/{quality}/media.m3u8` — which doesn't exist. The actual variant playlists live at the bucket-root `segments/{videoId}/{quality}/media.m3u8`, two directory levels above `manifests/{videoId}/`.

The two prior changes in this session (`make-worker-s3-path-style-configurable`, `make-worker-asset-urls-configurable`) made the *published* `manifestUrl`/`thumbnailUrl` configurable per environment. This bug is unrelated to that config — it's a fixed relationship between two S3 prefixes that's wrong in every environment.

## Goals / Non-Goals

**Goals:**
- Make the master manifest's variant references resolve to the correct, existing S3 objects in both LocalStack and production (CloudFront).
- Fix the path relationship once, independent of `S3_PUBLIC_ENDPOINT_URL` / `CLOUDFRONT_DOMAIN`.

**Non-Goals:**
- Changing variant playlist (`media.m3u8`) generation or `.ts` segment references — these are already correct (same-directory relative paths).
- Changing the S3 key layout (`manifests/{videoId}/...`, `segments/{videoId}/{quality}/...`).
- Reprocessing already-`ready` videos automatically.

## Decisions

- **Use a relative path (`../../segments/{videoId}/{quality}/media.m3u8`) rather than an absolute URL.** Relative paths are resolved by the player against the manifest's own location, so they work unchanged whether the manifest is served from LocalStack (`http://localhost:4566/...`) or CloudFront (`https://{domain}/...`). An absolute URL would require `buildMasterManifest` to know the worker's public base URL, re-coupling this fix to the env-specific config handled by the prior changes.
- **`../../` is correct because both prefixes sit at the same depth under the bucket root.** `manifests/{videoId}/master.m3u8` is two segments deep (`manifests/{videoId}/`); `segments/{videoId}/{quality}/media.m3u8` is reached by going up two levels (`../../`) and back down into `segments/{videoId}/{quality}/`.

## Risks / Trade-offs

- [Already-`ready` videos keep their broken master manifest] → Mitigation: reprocess affected videos (re-enqueue the existing job) after deploying the fix; not automated as part of this change.
- [Future S3 layout changes could silently break the `../../` relationship] → Mitigation: the modified spec requirement now states the resolution invariant explicitly, so future layout changes must consider it.
