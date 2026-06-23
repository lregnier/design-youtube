## Context

`ProcessVideo.Execute` is the sole use case in the worker application layer. It has a `_test.go` file with 5 tests covering the happy path and 4 error scenarios, but 7 branches remain untested. The `buildMasterManifest` helper is also untested. No mocks need to be added — the three port mocks (`MockVideoStorage`, `MockTranscoder`, `MockResultPublisher`) already exist in `internal/mocks/`.

## Goals / Non-Goals

**Goals:**
- Every branch in `ProcessVideo.Execute` is covered by at least one test
- `buildMasterManifest` is covered by a unit test

**Non-Goals:**
- Adding tests for adapters (`sqsjobs`, `ffmpeg`, `s3storage`, `sqspublisher`)
- Changing any production code
- Restructuring or refactoring existing tests

## Decisions

### Add tests directly to the existing `process_video_test.go`

All new functions go in the existing file rather than a new file. The use case is small and the test file currently has room; splitting would create unnecessary navigation overhead.

### Test `buildMasterManifest` as a package-level function

`buildMasterManifest` is unexported but the test is in the same package (`package application`), so it is directly callable. No export or indirection is needed.

### Cover per-quality transcode failures as separate test cases

The existing `TranscodeFailure` test only exercises the 1080p branch. Because the transcode loop short-circuits on first failure, the 720p and 360p paths are only reachable once 1080p succeeds. Two additional tests are added to confirm the same `PublishFailed` behaviour triggers for downstream qualities.

### `UploadThumbnail` failure is tested as a non-fatal path

`UploadThumbnail` errors are silently discarded (`thumbnailURL, _ = ...`). The test confirms that `PublishProcessed` is still called with an empty thumbnail URL when `UploadThumbnail` returns an error, matching the documented non-fatal intent.

## Risks / Trade-offs

- **Test isolation**: all new tests follow the same mockery v2 / AAA pattern as existing tests — no new risk surface.
