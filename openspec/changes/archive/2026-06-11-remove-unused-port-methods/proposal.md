## Why

`ports.ObjectStore.PutObject`, `ports.ObjectStore.GetObject`, and `ports.Cache.Delete` are implemented by the `s3store` and `rediscache` adapters but are never called by any use case in `internal/application/`. This violates the existing `hexagonal-architecture` spec requirement that outbound port interfaces be minimal — only methods required by the application layer.

## What Changes

- Remove `PutObject` and `GetObject` from `ports.ObjectStore` (`internal/ports/outbound.go`)
- Remove `Delete` from `ports.Cache` (`internal/ports/outbound.go`)
- Remove the corresponding implementations from `internal/adapters/outbound/s3store/store.go` (`PutObject`, `GetObject`, plus the now-unused `bytes`/`io` imports) and `internal/adapters/outbound/rediscache/cache.go` (`Delete`)
- Regenerate mockery mocks (`internal/gen/mocks/mock_ObjectStore.go`, `internal/gen/mocks/mock_Cache.go`) to drop the removed methods
- No behavior change for existing functionality (upload, confirm chunk, complete upload, get/list videos)
- Worker module is out of scope — its ports are all used

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `hexagonal-architecture`: tighten the "minimal Go interfaces" requirement to explicitly cover this case (no unused methods in outbound port interfaces)

## Impact

- Affected files: `backend/api/internal/ports/outbound.go`, `backend/api/internal/adapters/outbound/s3store/store.go`, `backend/api/internal/adapters/outbound/rediscache/cache.go`, `backend/api/internal/gen/mocks/mock_ObjectStore.go`, `backend/api/internal/gen/mocks/mock_Cache.go`
- No import path changes, no API/HTTP contract changes
