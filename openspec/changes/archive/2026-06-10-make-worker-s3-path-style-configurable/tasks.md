## 1. Extend worker config

- [x] 1.1 Add `S3UsePathStyle bool` to the `Config` struct in `backend/worker/internal/config/config.go`, parsed from `S3_USE_PATH_STYLE` via `strconv.ParseBool` (non-required, defaults to false on parse error or empty)

## 2. Update main.go wiring

- [x] 2.1 In `backend/worker/cmd/worker/main.go`, build an `[]func(*awss3.Options)` slice and append `o.UsePathStyle = true` only when `cfg.S3UsePathStyle` is true; pass it to `awss3.NewFromConfig`

## 3. Update docker-compose.yml

- [x] 3.1 Add `S3_USE_PATH_STYLE: "true"` to the `worker` service environment block

## 4. Documentation

- [x] 4.1 Add `S3_USE_PATH_STYLE` row to the Configuration table in `backend/worker/README.md`

## 5. Verify

- [x] 5.1 `go build ./...` and `go vet ./...` pass for `backend/worker`
- [x] 5.2 Rebuild and restart `worker` via `docker compose up --build worker`, re-run the upload flow, and confirm the video reaches `ready` status (not `failed`)
