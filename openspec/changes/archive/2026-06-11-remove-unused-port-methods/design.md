## Context

`ports.ObjectStore` and `ports.Cache` (`backend/api/internal/ports/outbound.go`) each declare one method with no caller in `internal/application/`:
- `ObjectStore.PutObject(ctx, key, data, contentType) error` — implemented in `s3store.Store`, calls `s3.PutObject`
- `ObjectStore.GetObject(ctx, key) ([]byte, error)` — implemented in `s3store.Store`, calls `s3.GetObject` + `io.ReadAll`
- `Cache.Delete(ctx, key) error` — implemented in `rediscache.Cache`, calls `redis.Del`

Confirmed via grep across `internal/application/`, `internal/adapters/inbound/`, `cmd/`, and all `_test.go` files: zero call sites.

## Goals / Non-Goals

**Goals:**
- Remove the 3 unused methods from the port interfaces, their adapter implementations, and the generated mocks
- Keep `internal/adapters/outbound/s3store` and `rediscache` packages compiling cleanly (drop now-unused imports)
- Tighten the `hexagonal-architecture` "minimal interfaces" requirement so this stays enforced

**Non-Goals:**
- No change to `Cache.Get`/`Set`, `ObjectStore.CreateMultipartUpload`/`PresignUploadPart`/`CompleteMultipartUpload`, or `Queue.SendMessage` — all are actively used
- No worker module changes
- No change to cache TTL/eviction behavior (Redis keys still expire via the existing `ttl` constant; `Delete` was never relied upon for invalidation)

## Decisions

- **Remove rather than keep "for future use"**: per the existing spec, port interfaces SHALL contain only methods required by the application layer. If cache invalidation or direct object get/put becomes necessary later, the method can be re-added alongside its first caller (YAGNI).
- **Regenerate mocks via `mockery` rather than hand-edit**: keeps `internal/gen/mocks/` consistent with its generator config, per the `backend-unit-tests` spec requirement that mocks are mockery-generated.
- **Drop `bytes`/`io` imports from `s3store/store.go`**: both were only used by `PutObject`/`GetObject`; leaving them would fail `go vet`/build.

## Risks / Trade-offs

- **Removing a port method is technically a breaking change to the interface** → mitigated: both interfaces are internal (`internal/ports`), no external consumers; build/vet/test after the change confirms nothing else depended on them.
- **Mock regeneration could produce unrelated diff noise** → mitigated: inspect `git diff` after running `mockery` to confirm only the removed methods disappear.
