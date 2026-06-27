## Why

`internal/config/` is a package with a single consumer: `cmd/api/main.go`. Giving it its own internal package adds indirection (an import, a package boundary, a separate directory) for no benefit — config loading is an entrypoint concern and belongs with the entrypoint.

## What Changes

- Move `backend/api/internal/config/config.go` → `backend/api/cmd/api/config.go` as `package main`
- Remove the `import` of `internal/config` from `main.go`; all config types are now in the same package
- Delete `backend/api/internal/config/` directory

## Capabilities

### New Capabilities

- none

### Modified Capabilities

- none

## Impact

- `backend/api/internal/config/config.go` — deleted
- `backend/api/cmd/api/config.go` — new file, same content, `package main`
- `backend/api/cmd/api/main.go` — remove `config` import, drop `cfg.` prefix changes (types now in same package, no prefix needed for the package — but field access stays the same)
