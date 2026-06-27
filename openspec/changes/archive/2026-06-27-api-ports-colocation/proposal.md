## Why

The api application layer defined each outbound port interface in its own file (`cache.go`, `event_publisher.go`, `object_store.go`), while the worker already collocates all port interfaces in a single `ports.go`. The inconsistency makes it harder to get an at-a-glance view of what the application layer depends on.

## What Changes

- Merge `cache.go`, `event_publisher.go`, and `object_store.go` into a single `backend/api/internal/application/ports.go`
- Delete the three individual files

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `hexagonal-architecture`: Add requirement that outbound port interfaces are colocated in a single `ports.go` file within the application layer

## Impact

- No API, behavior, or import path changes — all types remain in `package application`
- Worker already follows this convention; api now matches
