## Why

`backend/worker/internal/application/process_video_test.go` exists but leaves several error paths in `ProcessVideo.Execute` untested — including `Duration` failure, `UploadSegments` error, `UploadManifest` error, `UploadThumbnail` silent failure, and multiple-quality transcode failures — creating blind spots in a pipeline where a silent error means a video is permanently stuck in processing.

## What Changes

- Add missing `ProcessVideo.Execute` test cases to cover every untested branch in the execution pipeline
- Add a unit test for `buildMasterManifest` (currently no coverage at all)

## Capabilities

### New Capabilities

_None._

### Modified Capabilities

- `backend-unit-tests`: add test scenarios covering the currently untested paths in the worker `ProcessVideo` use case (DurationFailure, UploadSegmentsError, UploadManifestError, UploadThumbnailFailureNonFatal, per-quality TranscodeFailure for 720p and 360p, and buildMasterManifest output)

## Impact

- `backend/worker/internal/application/process_video_test.go` — new test functions added
- No production code changes
- No new dependencies
