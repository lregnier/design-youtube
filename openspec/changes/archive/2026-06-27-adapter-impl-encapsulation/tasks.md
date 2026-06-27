## 1. api outbound adapters

- [x] 1.1 `dynamo/repository.go`: rename `Repository` → `repository`; update constructor to return `video.VideoRepository`
- [x] 1.2 `s3store/store.go`: rename `Store` → `store`; update constructor to return `application.ObjectStore`
- [x] 1.3 `sqspublisher/publisher.go`: rename `Publisher` → `publisher`; update constructor to return `application.EventPublisher`
- [x] 1.4 `rediscache/cache.go`: rename `Cache` → `cache`; update constructor to return `application.Cache`

## 2. worker outbound adapters

- [x] 2.1 `sqspublisher/publisher.go`: rename `Publisher` → `publisher`; update constructor to return `application.EventPublisher`
- [x] 2.2 `s3storage/store.go`: rename `Store` → `store`; update constructor to return `application.VideoStorage`
- [x] 2.3 `ffmpeg/transcoder.go`: rename `Transcoder` → `transcoder`; update constructor to return `application.Transcoder`

## 3. Verify

- [x] 3.1 Run `go build ./...` from `backend/api/` and confirm clean
- [x] 3.2 Run `go test ./...` from `backend/api/` and confirm all tests pass
- [x] 3.3 Run `go build ./...` from `backend/worker/` and confirm clean
- [x] 3.4 Run `go test ./...` from `backend/worker/` and confirm all tests pass
