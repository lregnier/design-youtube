## 1. API main.go

- [x] 1.1 Add `newRepository(cfg *Config, awsCfg aws.Config) video.VideoRepository` below `main()` in `backend/api/cmd/api/main.go`
- [x] 1.2 Add `newStore(cfg *Config, awsCfg aws.Config) application.ObjectStore`
- [x] 1.3 Add `newPublisher(cfg *Config, awsCfg aws.Config) application.EventPublisher`
- [x] 1.4 Replace inline adapter blocks in `main()` with one-liner calls to the helpers; leave `cache` inline (already one line)
- [x] 1.5 Build and confirm clean: `go build ./...` from `backend/api/`

## 2. Worker main.go

- [x] 2.1 Add `newStore(cfg *Config, awsCfg aws.Config) application.VideoStorage` below `main()` in `backend/worker/cmd/worker/main.go`
- [x] 2.2 Add `newPublisher(cfg *Config, awsCfg aws.Config) application.EventPublisher`
- [x] 2.3 Replace inline adapter blocks in `main()` with one-liner calls to the helpers
- [x] 2.4 Build and confirm clean: `go build ./...` from `backend/worker/`
