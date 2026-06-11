## Context

`backend/api/api/` holds the OpenAPI spec (`openapi.yaml`), the oapi-codegen config (`oapi-codegen.yaml`), and a placeholder `generate.go` carrying the `//go:generate` directive. The `backend/api` Go module is itself named `api`, so this path reads as "api/api". The recent `internal/gen/` reorganization moved generated output out of `internal/api/`, but left this source-side `api/` directory untouched.

## Goals / Non-Goals

**Goals:**
- Rename `backend/api/api/` to `backend/api/openapi/`, removing the path stutter
- Keep the directory module-local (not moved to repo root) since it's only relevant to `backend/api`
- Update the `go:generate` directive's placeholder package name to match its new directory (`openapi`)
- Update README references and confirm regeneration is a no-op

**Non-Goals:**
- No change to `internal/gen/api/api.gen.go` content or location
- No change to the OpenAPI spec content itself
- No worker module changes (it has no OpenAPI spec)

## Decisions

- **Rename to `openapi/` rather than relocate to repo root**: the spec is specific to `backend/api`'s HTTP API; keeping it module-local matches where `oapi-codegen.yaml`'s relative `output:` path and `go:generate` resolve from.
- **Update `package api` → `package openapi` in `generate.go`**: this package is never imported (it only exists to host the `//go:generate` directive), but matching package name to directory name avoids `go vet`/tooling confusion and follows Go convention.
- **Leave `oapi-codegen.yaml`'s `output: ../internal/gen/api/api.gen.go` and `package: api` unchanged**: both are relative to/describe the *generated output* location (`internal/gen/api/`), which doesn't move.

## Risks / Trade-offs

- **Stale references** → grep the repo for `api/openapi.yaml`, `./api/...` (go:generate), and `backend/api/api` after the move; confirm only the README needed updates.
- **Regeneration drift** → run `go generate ./openapi/...` and `mockery` after the move, confirm `git status` shows no diff in `internal/gen/`.
