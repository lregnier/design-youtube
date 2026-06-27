## Why

`internal/` currently mixes DDD architecture layers (`domain/`, `application/`, `infrastructure/`) with generated tooling output (`gen/`). Moving `gen/` out makes `internal/` a clean, readable expression of the hexagonal architecture with exactly three layers.

## What Changes

- Move `backend/api/internal/gen/` → `backend/api/gen/`
- Update all Go import paths from `…/api/internal/gen/…` to `…/api/gen/…`
- Update `.mockery.yaml` output path
- Update `openapi/generate.go` `go:generate` directive
- Update the hexagonal-architecture spec requirement that mandates `internal/gen/`

## Capabilities

### New Capabilities

- none

### Modified Capabilities

- `hexagonal-architecture`: requirement "Generated code lives under internal/gen/" changes to "Generated code lives under gen/" (top-level alongside cmd/, internal/, openapi/)

## Impact

- `backend/api/internal/gen/` — deleted
- `backend/api/gen/` — new location for all generated code
- All `*_test.go` files importing `…/api/internal/gen/mocks` — import path updated
- `backend/api/internal/infrastructure/in/http/handler.go`, `server.go`, `middleware.go` — import path updated
- `backend/api/.mockery.yaml` — output path updated
- `backend/api/openapi/generate.go` — `go:generate` path updated
- `openspec/specs/hexagonal-architecture/spec.md` — requirement updated
