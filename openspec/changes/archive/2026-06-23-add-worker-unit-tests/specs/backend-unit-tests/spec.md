## ADDED Requirements

### Requirement: ProcessVideo use case tests cover all execution paths
Every branch in `ProcessVideo.Execute` in `backend/worker/internal/application/` SHALL have at least one corresponding test case. This includes all fatal error paths (returning an error), all non-fatal error paths (publishing a failure event and returning nil), and the non-fatal silent discard paths.

#### Scenario: Duration failure publishes failed event and returns nil
- **WHEN** `Transcoder.Duration` returns an error after a successful download
- **THEN** `ResultPublisher.PublishFailed` is called with a non-empty reason and `Execute` returns nil

#### Scenario: UploadSegments error returns error
- **WHEN** `VideoStorage.UploadSegments` returns an error after all qualities transcode successfully
- **THEN** `Execute` returns a non-nil error wrapping the upstream error

#### Scenario: UploadManifest error returns error
- **WHEN** `VideoStorage.UploadManifest` returns an error after segments are uploaded
- **THEN** `Execute` returns a non-nil error wrapping the upstream error

#### Scenario: Transcode failure for 720p publishes failed event
- **WHEN** `TranscodeHLS` succeeds for 1080p but returns an error for 720p
- **THEN** `ResultPublisher.PublishFailed` is called with a non-empty reason and `Execute` returns nil

#### Scenario: Transcode failure for 360p publishes failed event
- **WHEN** `TranscodeHLS` succeeds for 1080p and 720p but returns an error for 360p
- **THEN** `ResultPublisher.PublishFailed` is called with a non-empty reason and `Execute` returns nil

#### Scenario: UploadThumbnail failure is non-fatal
- **WHEN** `Transcoder.ExtractThumbnail` succeeds but `VideoStorage.UploadThumbnail` returns an error
- **THEN** `ResultPublisher.PublishProcessed` is still called with an empty thumbnail URL and `Execute` returns nil

### Requirement: buildMasterManifest produces correct HLS master playlist
The `buildMasterManifest` helper in `backend/worker/internal/application/` SHALL be covered by at least one unit test that asserts the structure and content of the generated M3U8 playlist.

#### Scenario: Master manifest contains all three quality levels
- **WHEN** `buildMasterManifest` is called with the standard three-quality slice
- **THEN** the returned string starts with `#EXTM3U`, contains three `#EXT-X-STREAM-INF` entries with the correct BANDWIDTH and RESOLUTION values, and includes the relative path to each quality's `media.m3u8`
