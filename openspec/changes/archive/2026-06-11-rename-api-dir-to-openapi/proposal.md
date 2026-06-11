## Why

`backend/api/api/` (containing `openapi.yaml`, `oapi-codegen.yaml`, `generate.go`) sits inside the `backend/api` Go module, producing the confusing "api/api" path stutter. Renaming it to `backend/api/openapi/` removes the collision while keeping the directory module-local and discoverable, complementing the recent `internal/gen/` reorganization.

## What Changes

- `git mv backend/api/api backend/api/openapi`
- Update `backend/api/openapi/generate.go`: change `package api` to `package openapi` (directory-local placeholder package, never imported, but should match its directory name by convention)
- `backend/api/openapi/oapi-codegen.yaml`'s `output:` path stays `../internal/gen/api/api.gen.go` (still correct relative to the new directory) — verify no edit needed
- Update `backend/api/README.md`: references to `api/openapi.yaml` → `openapi/openapi.yaml`, and `go generate ./api/...` → `go generate ./openapi/...`
- Regenerate via `go generate ./openapi/...` and confirm the output at `internal/gen/api/api.gen.go` is byte-identical (no behavior or generated-code change)
- Worker module is out of scope (it has no OpenAPI spec)

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `hexagonal-architecture`: the OpenAPI spec/config directory for `backend/api` lives at `internal`-sibling `openapi/` (not `api/`), avoiding the module-name collision

## Impact

- Affected files: `backend/api/api/{openapi.yaml,oapi-codegen.yaml,generate.go}` (renamed to `backend/api/openapi/`), `backend/api/README.md`
- No change to `internal/gen/api/api.gen.go` output, no import path changes (the openapi spec dir is never imported by Go code), no runtime behavior change
