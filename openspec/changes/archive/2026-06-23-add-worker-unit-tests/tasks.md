## 1. Add missing ProcessVideo.Execute test cases

- [x] 1.1 In `backend/worker/internal/application/process_video_test.go`: add `TestProcessVideo_Execute_DurationFailure` — mock `DownloadRaw` success, `Duration` returns error; assert `PublishFailed` called with non-empty reason and `Execute` returns nil
- [x] 1.2 Add `TestProcessVideo_Execute_TranscodeFailure720p` — mock download + duration success, `TranscodeHLS` success for 1080p but error for 720p; assert `PublishFailed` called and `Execute` returns nil
- [x] 1.3 Add `TestProcessVideo_Execute_TranscodeFailure360p` — mock download + duration success, `TranscodeHLS` success for 1080p and 720p but error for 360p; assert `PublishFailed` called and `Execute` returns nil
- [x] 1.4 Add `TestProcessVideo_Execute_UploadSegmentsError` — mock full transcode pipeline success, `UploadSegments` returns error; assert `Execute` returns a non-nil error
- [x] 1.5 Add `TestProcessVideo_Execute_UploadManifestError` — mock full transcode + segments upload success, `UploadManifest` returns error; assert `Execute` returns a non-nil error
- [x] 1.6 Add `TestProcessVideo_Execute_UploadThumbnailFailureNonFatal` — mock full transcode + upload pipeline success, `ExtractThumbnail` returns thumbnail bytes, `UploadThumbnail` returns error; assert `PublishProcessed` called with empty thumbnail URL and `Execute` returns nil

## 2. Add buildMasterManifest test

- [x] 2.1 In `backend/worker/internal/application/process_video_test.go`: add `TestBuildMasterManifest_AllQualities` — call `buildMasterManifest` with the package-level `qualities` slice; assert result starts with `#EXTM3U`, contains three `#EXT-X-STREAM-INF` lines with correct BANDWIDTH and RESOLUTION attributes, and contains the correct relative `media.m3u8` paths for each quality

## 3. Add adapter pure-logic tests

- [x] 3.1 Create `backend/worker/internal/adapters/outbound/s3storage/url_builder_test.go`: add `TestCloudFrontURLBuilder_AssetURL` asserting the result is `https://<domain>/<key>` (bucket ignored), and `TestEndpointURLBuilder_AssetURL` asserting the result is `<endpoint>/<bucket>/<key>`
- [x] 3.2 Create `backend/worker/internal/adapters/inbound/sqsjobs/consumer_test.go`: add `TestParseJob_ValidJSON` asserting a well-formed JSON body maps to the correct `ProcessingJob` fields, and `TestParseJob_InvalidJSON` asserting an error is returned for malformed input

## 4. Verify

- [x] 4.1 Run `go test ./internal/...` in `backend/worker/` and confirm all tests pass
