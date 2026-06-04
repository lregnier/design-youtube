## Why

The hexagonal architecture established in both bounded contexts makes the application layer fully testable in isolation — every use case depends on interfaces, not concrete infrastructure. Adding unit tests now validates the business logic, documents expected behaviour, and makes future refactors safe. Tests also strengthen the portfolio by demonstrating disciplined engineering practice alongside the architectural patterns.

## What Changes

- **mockery v2** configured via `.mockery.yaml` in `backend/api/` and `backend/worker/` — auto-generates mocks for all port interfaces
- **testify** added as a test dependency in both modules
- **API test files** — one per use case in `internal/application/`:
  - `upload/init_test.go` (InitUpload)
  - `upload/confirm_chunk_test.go` (ConfirmChunk)
  - `upload/complete_test.go` (CompleteUpload)
  - `catalog/get_video_test.go` (GetVideo)
  - `catalog/list_videos_test.go` (ListVideos)
  - `processing/apply_result_test.go` (ApplyProcessingResult)
- **Worker test file** — `internal/application/process_video_test.go` (ProcessVideo)
- All tests follow **AAA pattern** with `// Arrange`, `// Act`, `// Assert` comments
- All tests use generated mocks, never real infrastructure

## Capabilities

### New Capabilities

### Modified Capabilities

## Impact

- New files only — no existing code modified
- `backend/api/` and `backend/worker/` gain `.mockery.yaml`, `go.sum` updates, and generated `mocks/` directories
- CI workflow already runs `go build ./...` — add `go test ./...` step for both modules
