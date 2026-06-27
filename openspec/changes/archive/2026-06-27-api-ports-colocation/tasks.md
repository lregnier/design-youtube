## 1. Consolidate port interfaces

- [x] 1.1 Create `backend/api/internal/application/ports.go` with `Cache`, `EventPublisher`, `ObjectStore` and associated value types
- [x] 1.2 Delete `cache.go`, `event_publisher.go`, and `object_store.go`

## 2. Verify

- [x] 2.1 Run `go build ./...` from `backend/api/` and confirm clean
- [x] 2.2 Run `go test ./...` from `backend/api/` and confirm all tests pass
