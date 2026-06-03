## 1. Domain Layer

- [x] 1.1 Create `internal/domain/video/video.go`: define `VideoID`, `VideoStatus` (uploading/processing/ready/failed), `Chunk`, and `Video` aggregate struct with constructor and methods (`IsReady`, `NextMissingChunk`, `MarkChunkUploaded`)
- [x] 1.2 Create `internal/domain/video/repository.go`: define `VideoRepository` interface (`Save`, `FindByID`, `ListReady`)
- [x] 1.3 Verify `go build ./internal/domain/...` has zero imports outside stdlib

## 2. Port Interfaces

- [x] 2.1 Create `internal/ports/outbound.go`: define `ObjectStore` interface (`CreateMultipartUpload`, `PresignUploadPart`, `CompleteMultipartUpload`, `PutObject`, `GetObject`), `Cache` interface (`Get`, `Set`, `Delete`), and `Queue` interface (`DeleteMessage`)

## 3. Application — Upload Use Cases

- [x] 3.1 Create `internal/application/upload/init.go`: `InitUpload` struct with `VideoRepository` + `ObjectStore` deps; `Execute(ctx, InitUploadCommand) (InitUploadResult, error)` — validates size, creates S3 multipart upload, saves Video, returns presigned URL
- [x] 3.2 Create `internal/application/upload/confirm_chunk.go`: `ConfirmChunk` struct; `Execute` marks chunk uploaded on the Video, saves, returns next presigned URL or done signal
- [x] 3.3 Create `internal/application/upload/complete.go`: `CompleteUpload` struct; `Execute` calls `ObjectStore.CompleteMultipartUpload`, updates Video status to processing
- [x] 3.4 Verify `go build ./internal/application/...` imports only `internal/domain/` and `internal/ports/`

## 4. Application — Catalog Use Cases

- [x] 4.1 Create `internal/application/catalog/get_video.go`: `GetVideo` struct with `VideoRepository` + `Cache` deps; `Execute` checks cache first, falls back to repo, writes cache on miss
- [x] 4.2 Create `internal/application/catalog/list_videos.go`: `ListVideos` struct with `VideoRepository` dep; `Execute` returns all ready videos ordered by uploadedAt desc

## 5. Outbound Adapters

- [x] 5.1 Create `internal/adapters/outbound/dynamo/repository.go`: implement `video.VideoRepository` using AWS SDK v2 DynamoDB — `Save`, `FindByID`, `ListReady` (queries `status-uploadedAt-index` GSI)
- [x] 5.2 Create `internal/adapters/outbound/s3store/store.go`: implement `ports.ObjectStore` using AWS SDK v2 S3 — multipart upload creation, presigned part URLs with `content-length-range`, CompleteMultipartUpload, PutObject, GetObject
- [x] 5.3 Create `internal/adapters/outbound/rediscache/cache.go`: implement `ports.Cache` using go-redis — JSON marshal/unmarshal, 60s TTL on Set
- [x] 5.4 Create `internal/adapters/outbound/sqsqueue/queue.go`: implement `ports.Queue` using AWS SDK v2 SQS — `DeleteMessage`
- [x] 5.5 Verify each adapter satisfies its interface at compile time using blank identifier check: `var _ video.VideoRepository = (*dynamo.Repository)(nil)`

## 6. Inbound HTTP Adapter

- [x] 6.1 Create `internal/adapters/inbound/http/handler.go`: implement `api.StrictServerInterface` — each method constructs a command, calls the relevant use case, maps result to generated response type; no business logic
- [x] 6.2 Move upload secret `StrictMiddlewareFunc` from `internal/middleware/` to `internal/adapters/inbound/http/middleware.go`

## 7. Wire Up and Clean Up

- [x] 7.1 Update `cmd/api/main.go`: instantiate adapters, inject into use cases, inject use cases into HTTP handler, wire strict handler and middleware — delete old imports
- [x] 7.2 Update `cmd/worker/main.go`: replace direct AWS SDK calls with `s3store` and `sqsqueue` adapter instances where applicable; domain types for video status updates
- [x] 7.3 Delete `internal/handler/`, `internal/store/`, `internal/middleware/`
- [x] 7.4 Run `go build ./...` — zero errors
- [x] 7.5 Run `go vet ./...` — zero warnings
- [x] 7.6 Verify import graph: `internal/domain/` and `internal/application/` have no imports from `internal/adapters/`
