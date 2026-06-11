## Context

`backend/api/internal/api/api.gen.go` is the oapi-codegen output (chi server, strict-server interface, request/response types, models), generated from `backend/api/api/openapi.yaml` via `backend/api/api/oapi-codegen.yaml` (`output: ../internal/api/api.gen.go`, invoked via `go:generate` in `backend/api/api/generate.go`). `backend/api/internal/mocks/` is the mockery output (`VideoRepository`, `ObjectStore`, `Cache`, `Queue` mocks), generated from `.mockery.yaml` (`dir: internal/mocks`). Both are committed to the repo per the `backend-unit-tests` spec. The two share no naming convention, and `internal/api` collides in name with the module-root `api/` directory that holds the spec/config.

## Goals / Non-Goals

**Goals:**
- Group all generated code in `backend/api` under `internal/gen/`, one subdirectory per generator: `internal/gen/api/` (oapi-codegen) and `internal/gen/mocks/` (mockery).
- Keep package names unchanged (`api`, `mocks`) so call sites only need import-path updates, not identifier renames.
- Update the two generator configs so `go generate ./api/...` and `mockery` regenerate into the new locations.

**Non-Goals:**
- No change to `backend/worker` — it has no oapi-codegen output, and its `internal/mocks/` stays put. A follow-up could mirror this convention there later if desired, but isn't part of this change.
- No change to generated content, package names, or the `backend-unit-tests` AAA/naming conventions — only directory location and import paths.

## Decisions

- **Single `internal/gen/` root with one subdirectory per generator**, rather than e.g. co-locating mocks next to the interfaces they mock (`internal/ports/mocks/`, `internal/domain/video/mocks/`). A single `internal/gen/` tree keeps the "generated vs. hand-written" boundary visible at the top level of `internal/` and matches how the two generators are already configured from the module root (one `output`/`dir` setting each) — splitting mocks per-package would mean per-package mockery config blocks for no benefit here, since there's only one mocked package (`internal/ports`) plus `video.VideoRepository`.
- **Keep package names `api` and `mocks` unchanged.** Only the directory (and therefore import path) moves — `internal/api` → `internal/gen/api`, `internal/mocks` → `internal/gen/mocks`. Call sites change their import path but keep using `api.Xxx` / `mocks.NewMockXxx(t)` unchanged, minimizing the diff to import blocks only.
- **`backend/worker` is left as-is.** It only has `internal/mocks/` (no oapi-codegen output), so there's no naming collision to resolve there, and changing it isn't necessary to fix the `backend/api` issue. The `backend-unit-tests` spec update generalizes the mocks-location requirement to permit either layout per module, rather than forcing worker to move too.

## Risks / Trade-offs

- [Stale references to old paths in IDE caches, `go.work` replace directives, or scripts] → Mitigation: grep the whole repo for `internal/api"` and `internal/mocks"` after the move to confirm no leftover references; `go build ./...` and `go vet ./...` will fail loudly on any missed import.
- [Regenerating with `oapi-codegen`/`mockery` after moving configs produces a diff beyond the path change (e.g. version string)] → Mitigation: both tools are configured with `disable-version-string: true` (mockery) / no version banner change expected for oapi-codegen at the same tool version; diff the regenerated output against the moved file to confirm only the package comment/path-relative content is unchanged.
