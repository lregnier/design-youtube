## 1. Tooling Setup — backend/api

- [x] 1.1 Add `github.com/stretchr/testify` and `github.com/vektra/mockery/v2` to `backend/api/go.mod` (`go get`)
- [x] 1.2 Create `backend/api/.mockery.yaml`: configure mockery to generate mocks for `video.VideoRepository`, `ports.ObjectStore`, `ports.Cache` into `internal/mocks/`
- [x] 1.3 Run `mockery` in `backend/api/` and verify `internal/mocks/` contains `MockVideoRepository`, `MockObjectStore`, `MockCache`

## 2. Tooling Setup — backend/worker

- [x] 2.1 Add `github.com/stretchr/testify` to `backend/worker/go.mod` (`go get`)
- [x] 2.2 Create `backend/worker/.mockery.yaml`: configure mockery to generate mocks for `ports.VideoStorage`, `ports.Transcoder`, `ports.ResultPublisher` into `internal/mocks/`
- [x] 2.3 Run `mockery` in `backend/worker/` and verify `internal/mocks/` contains `MockVideoStorage`, `MockTranscoder`, `MockResultPublisher`

## 3. API Tests — Upload Use Cases

- [x] 3.1 Create `backend/api/internal/application/upload/init_test.go`:
  - `TestInitUpload_Execute_ValidRequest` — happy path: creates multipart upload, saves video, returns presigned URL
  - `TestInitUpload_Execute_FileTooLarge` — fileSize > 100MB returns error, no repo/store calls
  - `TestInitUpload_Execute_ResumeExistingUpload` — videoId provided with status `uploading`, returns next missing chunk URL
  - `TestInitUpload_Execute_RepoSaveError` — repo.Save returns error, propagated

- [x] 3.2 Create `backend/api/internal/application/upload/confirm_chunk_test.go`:
  - `TestConfirmChunk_Execute_MorePartsRemaining` — marks chunk uploaded, returns next presigned URL, done=false
  - `TestConfirmChunk_Execute_LastChunk` — marks final chunk, done=true, no presigned URL
  - `TestConfirmChunk_Execute_VideoNotFound` — repo returns nil, error returned

- [x] 3.3 Create `backend/api/internal/application/upload/complete_test.go`:
  - `TestCompleteUpload_Execute_Success` — calls CompleteMultipartUpload, updates status to processing
  - `TestCompleteUpload_Execute_VideoNotFound` — repo returns nil, error returned
  - `TestCompleteUpload_Execute_S3Error` — S3 CompleteMultipartUpload fails, error propagated

## 4. API Tests — Catalog Use Cases

- [x] 4.1 Create `backend/api/internal/application/catalog/get_video_test.go`:
  - `TestGetVideo_Execute_CacheHit` — cache returns data, repo.FindByID never called
  - `TestGetVideo_Execute_CacheMissPopulatesCache` — cache miss, repo called, result written to cache
  - `TestGetVideo_Execute_VideoNotFound` — repo returns nil, nil returned with no error

- [x] 4.2 Create `backend/api/internal/application/catalog/list_videos_test.go`:
  - `TestListVideos_Execute_ReturnsList` — repo returns ready videos, all returned
  - `TestListVideos_Execute_EmptyList` — repo returns empty slice, empty slice returned

## 5. API Tests — Processing Use Case

- [x] 5.1 Create `backend/api/internal/application/processing/apply_result_test.go`:
  - `TestApplyProcessingResult_OnProcessed_Success` — loads video, calls MarkReady, saves
  - `TestApplyProcessingResult_OnProcessed_IdempotentOnReadyVideo` — video already ready, MarkReady called again, saves without error
  - `TestApplyProcessingResult_OnProcessed_VideoNotFound` — repo returns nil, error returned
  - `TestApplyProcessingResult_OnFailed_Success` — loads video, calls MarkFailed, saves
  - `TestApplyProcessingResult_OnFailed_VideoNotFound` — repo returns nil, error returned

## 6. Worker Tests — ProcessVideo Use Case

- [x] 6.1 Create `backend/worker/internal/application/process_video_test.go`:
  - `TestProcessVideo_Execute_SuccessfulPipeline` — all steps succeed, PublishProcessed called with correct URLs
  - `TestProcessVideo_Execute_DownloadFailure` — DownloadRaw fails, PublishFailed called, no transcode attempted
  - `TestProcessVideo_Execute_TranscodeFailure` — TranscodeHLS fails for one quality, PublishFailed called
  - `TestProcessVideo_Execute_ThumbnailFailureNonFatal` — ExtractThumbnail fails, PublishProcessed still called (thumbnail URL may be empty)
  - `TestProcessVideo_Execute_PublishProcessedError` — PublishProcessed returns error, error propagated

## 7. CI Update

- [x] 7.1 Add `go test ./...` step to the `api` job in `.github/workflows/ci.yml` after `go build ./...`
- [x] 7.2 Add `go test ./...` step to the `worker` job in `.github/workflows/ci.yml` after `go build ./...`

## 8. Verify

- [x] 8.1 Run `go test ./internal/application/...` in `backend/api/` — all tests pass
- [x] 8.2 Run `go test ./internal/application/...` in `backend/worker/` — all tests pass
