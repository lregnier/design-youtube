## Why

The current backend mixes infrastructure concerns (DynamoDB, S3, Redis, SQS) directly into handler and store packages, making the business logic hard to test in isolation and the architecture hard to reason about. Restructuring around hexagonal architecture (ports and adapters) with DDD tactical patterns makes the codebase a stronger portfolio demonstration of clean architecture and SOLID principles in Go.

## What Changes

- New `internal/domain/` package: `Video` aggregate, value objects (`VideoID`, `VideoStatus`), and the `VideoRepository` port interface — zero external imports
- New `internal/application/` package: one use-case type per operation (`InitUpload`, `ConfirmChunk`, `CompleteUpload`, `GetVideo`, `ListVideos`) depending only on domain types and outbound port interfaces
- New `internal/ports/` package: outbound interface definitions (`VideoRepository`, `ObjectStore`, `Cache`, `Queue`) — the contracts adapters must satisfy
- New `internal/adapters/outbound/` packages: DynamoDB, S3, Redis, and SQS implementations of the port interfaces
- `internal/handler/` becomes `internal/adapters/inbound/http/` — the HTTP inbound adapter wiring oapi-codegen to use cases
- Deleted: `internal/handler/`, `internal/store/`, `internal/middleware/` (middleware moves into the inbound adapter)
- No changes to `backend/api/openapi.yaml`, generated code, or any external behavior

## Capabilities

### New Capabilities

### Modified Capabilities

## Impact

- Affects all files under `backend/internal/` — complete package restructure
- `backend/cmd/api/main.go` and `backend/cmd/worker/main.go` wiring updated to compose adapters into use cases
- No changes to the OpenAPI spec, generated code, Terraform, frontend, or docker-compose
- No breaking changes to the API contract or observable behavior
