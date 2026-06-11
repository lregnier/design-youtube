## Why

HLS playback returns a 404 for every variant playlist. The master manifest (`manifests/{videoId}/master.m3u8`) references variant playlists with the relative path `segments/{videoId}/{quality}/media.m3u8`, but players resolve relative URLs against the master manifest's own directory (`manifests/{videoId}/`), producing `manifests/{videoId}/segments/{videoId}/{quality}/media.m3u8` — a path that doesn't exist. The actual variant playlists live at the bucket-root `segments/{videoId}/{quality}/media.m3u8`. This breaks playback in both LocalStack and production (CloudFront), since it's a path-relationship bug, not an environment/URL-config issue.

## What Changes

- In `buildMasterManifest` (`backend/worker/internal/application/process_video.go`), change the variant reference template from `segments/%s/%s/media.m3u8` to `../../segments/%s/%s/media.m3u8` so it correctly resolves from `manifests/{videoId}/` back to the bucket-root `segments/{videoId}/{quality}/media.m3u8`

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- **`video-processing`**: tighten "Worker generates an HLS master manifest" to require that variant playlist references resolve correctly relative to the master manifest's own location, closing the gap that allowed this bug

## Impact

- **`backend/worker/internal/application/process_video.go`**: one-line format string fix in `buildMasterManifest`
- Any already-`ready` videos retain their broken master manifest until reprocessed
