## 1. Remove unused methods from port interfaces

- [x] 1.1 Remove `PutObject` and `GetObject` from `ports.ObjectStore` in `internal/ports/outbound.go`
- [x] 1.2 Remove `Delete` from `ports.Cache` in `internal/ports/outbound.go`

## 2. Remove unused adapter implementations

- [x] 2.1 Remove `Store.PutObject` and `Store.GetObject` from `internal/adapters/outbound/s3store/store.go`, and drop the now-unused `bytes`/`io` imports
- [x] 2.2 Remove `Cache.Delete` from `internal/adapters/outbound/rediscache/cache.go`

## 3. Regenerate mocks

- [x] 3.1 Run `mockery` to regenerate `internal/gen/mocks/mock_ObjectStore.go` and `internal/gen/mocks/mock_Cache.go`, confirm only the removed methods disappear

## 4. Verify

- [x] 4.1 `go build ./...`, `go vet ./...`, `go test ./...` succeed
- [x] 4.2 Grep for any remaining references to `PutObject`, `GetObject` (on `ObjectStore`/`Store`), or `Cache.Delete`
