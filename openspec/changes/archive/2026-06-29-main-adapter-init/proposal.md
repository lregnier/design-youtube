## Why

`main()` in both `api` and `worker` contains verbose inline adapter initialization blocks that mix AWS SDK option-building with adapter construction. The wiring of services is buried in boilerplate, making each `main()` hard to scan at a glance.

## What Changes

- Extract each adapter's initialization into a private `new*` helper function in a new `adapters.go` file alongside `main.go`, for both `backend/api/cmd/api/` and `backend/worker/cmd/worker/`
- `main()` becomes a flat sequence of one-liner assignments: `repo := newRepository(cfg, awsCfg)`
- Helpers extracted in **api**: `newRepository`, `newStore`, `newCache`, `newPublisher`
- Helpers extracted in **worker**: `newStore`, `newPublisher`
- `ffmpeg.NewTranscoder()` stays inline — it's already one line with no config branching

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None.

## Impact

- `backend/api/cmd/api/main.go` — simplified; new `backend/api/cmd/api/adapters.go`
- `backend/worker/cmd/worker/main.go` — simplified; new `backend/worker/cmd/worker/adapters.go`
- No changes to application, domain, or infrastructure packages
