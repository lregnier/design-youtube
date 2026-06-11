## 1. Rename the directory

- [x] 1.1 `git mv backend/api/api backend/api/openapi`

## 2. Update package name

- [x] 2.1 Change `package api` to `package openapi` in `backend/api/openapi/generate.go`

## 3. Update references

- [x] 3.1 Update `backend/api/README.md`: `api/openapi.yaml` → `openapi/openapi.yaml`, and `go generate ./api/...` → `go generate ./openapi/...`
- [x] 3.2 Confirm `backend/api/openapi/oapi-codegen.yaml`'s `output: ../internal/gen/api/api.gen.go` is still correct relative to the new location (no edit expected)

## 4. Verify

- [x] 4.1 Regenerate via `go generate ./openapi/...`, confirm empty diff in `internal/gen/api/`
- [x] 4.2 `go build ./...`, `go vet ./...`, `go test ./...` succeed
- [x] 4.3 Grep for stale references to `backend/api/api` or `./api/...` (go:generate)
