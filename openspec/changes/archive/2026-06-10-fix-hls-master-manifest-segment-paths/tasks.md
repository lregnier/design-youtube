## 1. Fix master manifest path generation

- [x] 1.1 In `backend/worker/internal/application/process_video.go`, change `buildMasterManifest`'s variant reference template from `segments/%s/%s/media.m3u8` to `../../segments/%s/%s/media.m3u8`

## 2. Verify

- [x] 2.1 Rebuild the worker and re-enqueue the existing `ready` video's processing job
- [x] 2.2 Confirm the regenerated `manifests/{videoId}/master.m3u8` references variants as `../../segments/{videoId}/{quality}/media.m3u8`
- [x] 2.3 Confirm `http://localhost:4566/design-youtube-video-prod/segments/{videoId}/1080p/media.m3u8` returns 200 and the video plays in the browser
