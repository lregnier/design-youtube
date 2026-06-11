## 1. Move generated code under internal/gen/

- [x] 1.1 `git mv backend/api/internal/api backend/api/internal/gen/api`
- [x] 1.2 `git mv backend/api/internal/mocks backend/api/internal/gen/mocks`

## 2. Update generator configs

- [x] 2.1 Update `backend/api/api/oapi-codegen.yaml`'s `output:` to `../internal/gen/api/api.gen.go`
- [x] 2.2 Update `backend/api/.mockery.yaml`'s `dir:` to `internal/gen/mocks`

## 3. Update import paths

- [x] 3.1 Update the 4 files importing `github.com/lregnier/design-youtube/api/internal/api` (`internal/adapters/inbound/http/server.go`, `middleware.go`, `handler.go`, `handler_test.go`) to `github.com/lregnier/design-youtube/api/internal/gen/api`
- [x] 3.2 Update the 7 files importing `github.com/lregnier/design-youtube/api/internal/mocks` (`internal/adapters/inbound/http/handler_test.go`, `internal/application/processing/apply_result_test.go`, `internal/application/catalog/list_videos_test.go`, `internal/application/catalog/get_video_test.go`, `internal/application/upload/complete_test.go`, `internal/application/upload/init_test.go`, `internal/application/upload/confirm_chunk_test.go`) to `github.com/lregnier/design-youtube/api/internal/gen/mocks`

## 4. Verify

- [x] 4.1 Regenerate via `go generate ./api/...` (oapi-codegen) and `mockery` (mockery v2) and confirm the diff is empty (only the moved files, no unexpected changes)
- [x] 4.2 `go build ./...`, `go vet ./...`, and `go test ./...` succeed in `backend/api`
- [x] 4.3 Confirm no remaining references to `internal/api"` or `internal/mocks"` (old paths) in `backend/api` via `grep -r`
