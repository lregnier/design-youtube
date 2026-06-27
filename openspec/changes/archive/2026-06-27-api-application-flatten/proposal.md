## Why

The current application layer splits each domain slice into its own subdirectory package (`upload/`, `catalog/`, `processing/`). This adds package-boundary boilerplate (import paths, `Port` suffix on interfaces, separate `ports.go` files) without providing meaningful isolation — all three slices operate on the same domain and share the same bounded context. Flattening them into a single `application/` package removes accidental complexity and aligns with Go's preference for small, cohesive packages.

## What Changes

- Merge `application/upload/`, `application/catalog/`, and `application/processing/` into a single `application/` package
- Rename service interfaces: drop the `Port` suffix (`UploadServicePort` → `UploadService`, `CatalogServicePort` → `CatalogService`, `ProcessingServicePort` → `ProcessingService`)
- Co-locate each service interface with its implementation struct in one file (e.g., `catalog_service.go` holds both `CatalogService` interface and unexported `catalogService` struct)
- Make implementation structs unexported; constructors return the interface type
- Give each outbound port interface its own file (`cache.go`, `object_store.go`, `event_publisher.go`)
- Keep value types (`MultipartUpload`, `PresignedURL`, `CompletedPart`) and event types (`VideoProcessedEvent`, `VideoFailedEvent`) in the `application` package
- **BREAKING**: import paths change from `application/catalog`, `application/upload`, `application/processing` to `application` everywhere (infrastructure adapters, `main.go`, `mockery.yaml`, tests)

## Capabilities

### New Capabilities

- none

### Modified Capabilities

- `hexagonal-architecture`: internal package layout of the application layer changes (subdirs → flat); no requirement-level behavior changes, only structural

## Impact

- `backend/api/internal/application/` — all files restructured
- `backend/api/internal/infrastructure/in/http/` — import paths updated
- `backend/api/internal/infrastructure/in/sqsconsumer/` — import paths updated
- `backend/api/internal/infrastructure/out/rediscache/`, `s3store/`, `sqspublisher/` — compile-time interface checks updated
- `backend/api/cmd/api/main.go` — import paths updated
- `backend/api/.mockery.yaml` — consolidated to single `application` package entry
- `backend/api/internal/gen/mocks/` — regenerated
- All test files importing application sub-packages updated
