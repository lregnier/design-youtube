## Context

`backend/api/internal/` currently contains four directories: `domain/`, `application/`, `infrastructure/`, and `gen/`. The first three express the hexagonal architecture layers; `gen/` is tooling output (oapi-codegen and mockery). Placing generated code inside `internal/` obscures the layered structure and implies it is architecture, not build artifact.

## Goals / Non-Goals

**Goals:**
- Move `backend/api/internal/gen/` → `backend/api/gen/` so `internal/` contains only the three DDD layers
- Update all import paths, codegen config, and the spec requirement

**Non-Goals:**
- Changing any generated code content or codegen configuration beyond output paths
- Reorganising the structure inside `gen/` (subdirectories `api/` and `mocks/` stay as-is)

## Decisions

### Keep `gen/` at the module root level, not inside `internal/`

`internal/` in Go restricts imports to the parent tree. Removing that restriction for `gen/` has no practical effect — nothing outside `backend/api/` imports the mocks or generated API types, and that is enforced by the module boundary, not by `internal/`. Moving `gen/` out therefore has zero risk and makes `internal/` a clean architectural statement.

Alternative considered: a nested `cmd/api/gen/` or `openapi/gen/`. Rejected — `gen/` is module-wide output (mocks are used by multiple packages), not entrypoint- or spec-specific.

## Risks / Trade-offs

[Regeneration] → After moving, running `go generate ./openapi/...` or `mockery` will write to the new `gen/` path as long as configs are updated first. If configs are updated before the move, an intermediate state exists where `internal/gen/` and `gen/` both exist; the move resolves this.

[Import churn] → Every file importing `…/api/internal/gen/…` needs an import path update. Affected files are all test files and the HTTP infrastructure files — mechanical find-and-replace, no logic changes.
