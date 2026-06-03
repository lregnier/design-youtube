## Context

The worker currently has a `cmd/worker/main.go` that contains ~250 lines of mixed concerns: SQS polling, S3 download, ffmpeg execution, S3 upload, manifest generation, thumbnail extraction, and event emission. This mirrors the flat structure the API had before its hexagonal refactor. The `internal/` packages only hold config, event types, and an SQS publisher — no domain model, no application layer, no port abstractions.

The API established the hexagonal pattern for this monorepo. The worker should follow the same structure to be consistent and to demonstrate the pattern applies equally well to a simpler, pipeline-oriented service.

## Goals / Non-Goals

**Goals:**
- Domain layer: `ProcessingJob` value object — zero external imports
- Application layer: `ProcessVideo` use case depending only on port interfaces
- Port interfaces: `VideoStorage` (S3), `Transcoder` (ffmpeg), `ResultPublisher` (SQS results queue)
- Inbound adapter: SQS job consumer decoupled from processing logic
- Outbound adapters: independent implementations of each port
- `cmd/worker/main.go` as a pure composition root
- Identical external behavior

**Non-Goals:**
- Changing the integration event schema
- Adding tests (separate change)
- Changing ffmpeg flags, S3 key structure, or SQS queue configuration

## Decisions

### 1. Package layout mirrors the API

```
internal/
  domain/processing/
    job.go            # ProcessingJob{VideoID, S3Key} value object
  application/
    process_video.go  # ProcessVideo use case
  ports/
    outbound.go       # VideoStorage, Transcoder, ResultPublisher interfaces
  adapters/
    inbound/
      sqsjobs/
        consumer.go   # SQS long-poll loop, parses job, calls use case
    outbound/
      s3storage/
        store.go      # VideoStorage: download raw, upload segments/manifest/thumbnail
      ffmpeg/
        transcoder.go # Transcoder: wraps exec.Command ffmpeg/ffprobe
      sqspublisher/
        publisher.go  # ResultPublisher: sends VideoProcessed/VideoFailed events
  config/
    config.go
  event/
    result.go         # VideoProcessed/VideoFailed types (kept — shared event schema)
```

### 2. Transcoder is a first-class port

Wrapping `exec.Command` behind a `Transcoder` interface is the most valuable abstraction here — it makes the use case testable without ffmpeg installed and decouples the pipeline logic from OS-level process management.

```go
type Transcoder interface {
    Duration(ctx context.Context, inputPath string) (float64, error)
    TranscodeHLS(ctx context.Context, inputPath, outputDir, scale, bitrate string) error
    ExtractThumbnail(ctx context.Context, inputPath, outputPath string, offset float64) error
}
```

Alternatives considered: inlining ffmpeg calls in the use case (untestable), using a third-party Go ffmpeg library (unnecessary dependency for a portfolio).

### 3. VideoStorage covers both read and write

A single `VideoStorage` port handles downloading the raw video and uploading segments, manifests, and thumbnails. This keeps the port count low and reflects that S3 is the single storage concern for the worker.

```go
type VideoStorage interface {
    DownloadRaw(ctx context.Context, videoID string, destPath string) error
    UploadSegments(ctx context.Context, videoID string, segDir string) error
    UploadManifest(ctx context.Context, videoID string, content []byte) (string, error)
    UploadThumbnail(ctx context.Context, videoID string, data []byte) (string, error)
}
```

### 4. ProcessVideo use case owns temp directory lifecycle

The use case creates and cleans up the temp working directory. This keeps the infrastructure management close to the processing logic that needs it, rather than scattering it across adapters.

### 5. event/result.go stays — it is the integration contract

The `VideoProcessed` and `VideoFailed` types define the wire format shared (by documentation) with the API consumer. Keeping them in `internal/event/` makes the contract explicit rather than burying it in the publisher adapter.

## Risks / Trade-offs

- **More files, same logic** → Same trade-off as the API refactor: justified for portfolio consistency and testability.
- **Transcoder interface granularity** → Three methods on one interface vs three separate interfaces. Chosen because all three are always used together in one use case — splitting would be YAGNI.
- **Temp dir in use case** → Slightly outside pure application logic, but the alternative (passing temp dir as a param) leaks infrastructure concerns into the caller. Acceptable trade-off.

## Migration Plan

1. Create new package structure alongside existing code
2. Implement domain, ports, and use case
3. Implement outbound adapters
4. Implement inbound SQS consumer
5. Update `cmd/worker/main.go` to compose new structure
6. Delete old flat code from `internal/queue/` and the body of `main.go`
7. Verify `go build ./...` and `go vet ./...`
