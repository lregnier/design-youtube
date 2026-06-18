## 1. backend/api — update Config

- [x] 1.1 In `backend/api/internal/config/config.go`: replace `S3PublicEndpointURL string` field with `LocalStack bool` and `LocalStackEndpoint string`
- [x] 1.2 In `Load()`: read `LOCALSTACK_ENABLED` via `strconv.ParseBool` into `LocalStack`; read `LOCALSTACK_ENDPOINT` into `LocalStackEndpoint`; remove `S3PublicEndpointURL` assignment
- [x] 1.3 In `Load()`: add validation — return an error if `LocalStack == true` and `LocalStackEndpoint == ""`

## 2. backend/api — update main.go

- [x] 2.1 In `cmd/api/main.go`: replace the `cfg.S3PublicEndpointURL != ""` check with `cfg.LocalStack`; pass `cfg.LocalStackEndpoint` to `s3store.NewLocalStackTransformer`

## 3. backend/worker — update Config

- [x] 3.1 In `backend/worker/internal/config/config.go`: replace `S3PublicEndpointURL string` field with `LocalStack bool` and `LocalStackEndpoint string`
- [x] 3.2 In `Load()`: read `LOCALSTACK_ENABLED` via `strconv.ParseBool` into `LocalStack`; read `LOCALSTACK_ENDPOINT` into `LocalStackEndpoint`; remove `S3PublicEndpointURL` assignment
- [x] 3.3 In `Load()`: add validation — return an error if `LocalStack == true` and `LocalStackEndpoint == ""`

## 4. backend/worker — update main.go

- [x] 4.1 In `cmd/worker/main.go`: replace the `cfg.S3PublicEndpointURL != ""` check with `cfg.LocalStack`; pass `cfg.LocalStackEndpoint` to `s3storage.NewLocalStackURLBuilder`

## 5. docker-compose.yml

- [x] 5.1 Replace `S3_PUBLIC_ENDPOINT_URL: http://localhost:4566` with `LOCALSTACK_ENABLED: "true"` and `LOCALSTACK_ENDPOINT: http://localhost:4566` in both the `api` and `worker` service definitions

## 6. Verify

- [x] 6.1 `go build ./...` and `go vet ./...` pass in both `backend/api` and `backend/worker`
- [x] 6.2 `go test ./...` passes in both modules
- [x] 6.3 Grep for remaining `S3PublicEndpointURL` or `S3_PUBLIC_ENDPOINT_URL` references in config and main files (should be zero)
