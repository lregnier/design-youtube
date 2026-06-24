## Why

After completing an upload, the user is redirected to the home page but sees nothing — their video only appears once it finishes transcoding (minutes later), with no indication that processing is underway. The home page should reflect the full lifecycle: a video should appear immediately after upload in a "processing" state, then become clickable once ready.

## What Changes

- `GET /videos` changes from returning only `ready` videos to returning all videos with status `processing`, `ready`, or `failed` (excluding `uploading`, which is still in-progress upload)
- `VideoRepository` gains a `List` method (replaces `ListReady`) that returns all non-uploading videos ordered by upload timestamp descending
- `VideoCard` renders differently per status: clickable with thumbnail for `ready`; non-clickable with a "Processing…" indicator for `processing`; non-clickable with a "Failed" badge for `failed`
- `HomePage` polls `GET /videos` every 5 seconds while at least one video is in `processing` state, and stops polling when none remain

## Capabilities

### New Capabilities

_None._

### Modified Capabilities

- `video-catalog`: `GET /videos` must include `processing` and `failed` videos, not only `ready`; `VideoCard` must show status-aware UI; homepage must auto-refresh during processing

## Impact

- `backend/api/internal/domain/video/repository.go` — add `List` to `VideoRepository` interface
- `backend/api/internal/adapters/outbound/dynamo/repository.go` — implement `List` (Scan excluding `uploading`, sorted by uploadedAt desc)
- `backend/api/internal/application/catalog/list_videos.go` — call `repo.List()` instead of `repo.ListReady()`
- `backend/api/internal/gen/mocks/mock_VideoRepository.go` — regenerate mock
- `backend/api/internal/application/catalog/list_videos_test.go` — update test
- `frontend/web/src/components/VideoCard.tsx` — status-aware rendering
- `frontend/web/src/pages/HomePage.tsx` — add polling while processing
