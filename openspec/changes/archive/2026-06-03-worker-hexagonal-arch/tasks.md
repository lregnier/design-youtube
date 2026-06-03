## 1. Domain Layer

- [x] 1.1 Create `internal/domain/processing/job.go`: define `ProcessingJob` struct with `VideoID` and `S3Key` string fields
- [x] 1.2 Verify `go build ./internal/domain/...` has zero imports outside stdlib

## 2. Port Interfaces

- [x] 2.1 Create `internal/ports/outbound.go`: define `VideoStorage` interface (`DownloadRaw`, `UploadSegments`, `UploadManifest`, `UploadThumbnail`), `Transcoder` interface (`Duration`, `TranscodeHLS`, `ExtractThumbnail`), and `ResultPublisher` interface (`PublishProcessed`, `PublishFailed`)

## 3. Application — ProcessVideo Use Case

- [x] 3.1 Create `internal/application/process_video.go`: `ProcessVideo` struct with `VideoStorage`, `Transcoder`, `ResultPublisher`, `CloudFrontDomain` deps; `Execute(ctx, ProcessingJob) error` — creates temp dir, calls storage.DownloadRaw, calls transcoder for each quality level, calls storage.UploadSegments + UploadManifest + UploadThumbnail, calls publisher.PublishProcessed on success or PublishFailed on error
- [x] 3.2 Verify `go build ./internal/application/...` imports only `internal/domain/`, `internal/ports/`, `internal/event/`, and stdlib

## 4. Outbound Adapters

- [x] 4.1 Create `internal/adapters/outbound/s3storage/store.go`: implement `ports.VideoStorage` using AWS SDK v2 S3 — `DownloadRaw` (streams to temp file), `UploadSegments` (walks dir, PutObject each .ts and .m3u8), `UploadManifest` (PutObject master.m3u8, returns CloudFront URL), `UploadThumbnail` (PutObject thumb.jpg, returns CloudFront URL); add compile-time interface check `var _ ports.VideoStorage = (*Store)(nil)`
- [x] 4.2 Create `internal/adapters/outbound/ffmpeg/transcoder.go`: implement `ports.Transcoder` — `Duration` (runs ffprobe, parses float), `TranscodeHLS` (runs ffmpeg with scale/bitrate/hls_time flags), `ExtractThumbnail` (runs ffmpeg -ss -frames:v 1); add compile-time interface check
- [x] 4.3 Create `internal/adapters/outbound/sqspublisher/publisher.go`: implement `ports.ResultPublisher` — `PublishProcessed` emits `event.VideoProcessed`, `PublishFailed` emits `event.VideoFailed`, both using `videoId` as FIFO message group key; add compile-time interface check
- [x] 4.4 Delete `internal/queue/publisher.go` (replaced by sqspublisher)

## 5. Inbound Adapter — SQS Job Consumer

- [x] 5.1 Create `internal/adapters/inbound/sqsjobs/consumer.go`: SQS long-poll loop; parses S3 event notification or raw `processingJob` JSON into `domain.ProcessingJob`; calls `processVideo.Execute`; deletes SQS message on success; leaves message on error for retry

## 6. Wire Up and Clean Up

- [x] 6.1 Rewrite `cmd/worker/main.go` as a pure composition root: instantiate AWS config, s3storage adapter, ffmpeg adapter, sqspublisher adapter, ProcessVideo use case, sqsjobs consumer; call `consumer.Start(ctx)`
- [x] 6.2 Delete the old processing logic from the previous `cmd/worker/main.go` body (all helper functions: `downloadS3`, `uploadBytes`, `uploadDir`, `videoDuration`, `transcode`, `extractThumbnail`, `buildMasterManifest`, `emitFailed`)
- [x] 6.3 Run `go build ./...` in `backend/worker/` — zero errors
- [x] 6.4 Run `go vet ./...` in `backend/worker/` — zero warnings
- [x] 6.5 Verify `internal/domain/` and `internal/application/` have no imports from `internal/adapters/`
