## 1. Extend worker config

- [x] 1.1 Add `S3PublicEndpointURL string` to the `Config` struct in `backend/worker/internal/config/config.go`, parsed from `S3_PUBLIC_ENDPOINT_URL` env var (optional, defaults to empty)

## 2. Update s3storage.Store

- [x] 2.1 Add `s3PublicEndpointURL` field to `Store` and accept it in `NewStore` (`backend/worker/internal/adapters/outbound/s3storage/store.go`)
- [x] 2.2 In `UploadManifest` and `UploadThumbnail`, when `s3PublicEndpointURL` is set, build the URL as `fmt.Sprintf("%s/%s/%s", s3PublicEndpointURL, bucket, key)`; when unset, keep the existing `fmt.Sprintf("https://%s/%s", cloudfrontDomain, key)`

## 3. Update main.go wiring

- [x] 3.1 Pass `cfg.S3PublicEndpointURL` into `s3storage.NewStore` in `backend/worker/cmd/worker/main.go`

## 4. Update docker-compose.yml

- [x] 4.1 Add `S3_PUBLIC_ENDPOINT_URL: "http://localhost:4566"` to the `worker` service environment block

## 5. Documentation

- [x] 5.1 Add `S3_PUBLIC_ENDPOINT_URL` row to the Configuration table in `backend/worker/README.md`

## 6. Verify

- [x] 6.1 `go build ./...` and `go vet ./...` pass for `backend/worker`
- [x] 6.2 Rebuild and restart `worker` via `docker compose up --build worker`, re-enqueue the existing video's processing job, and confirm `thumbnailUrl`/`manifestUrl` are now `http://localhost:4566/design-youtube-video-prod/...` and load successfully
