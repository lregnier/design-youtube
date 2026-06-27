## Context

The API application layer currently has three sub-packages: `application/upload`, `application/catalog`, `application/processing`. Each sub-package contains a `ports.go` (interfaces) and `service.go` (implementation). The sub-package split was motivated by avoiding a shared `ports/` import cycle, but it introduced unnecessary package-boundary overhead: callers must import three separate paths, mockery must target three packages, and the `Port` suffix on every interface name is visual noise. Since all three slices belong to the same bounded context and domain, a single flat `application/` package is the right scope.

## Goals / Non-Goals

**Goals:**
- Single `application/` package — one import path for all application-layer types
- Service interfaces named without `Port` suffix (`UploadService`, `CatalogService`, `ProcessingService`)
- Interface and implementation co-located in one file per service
- Outbound port interfaces each in their own file (`object_store.go`, `event_publisher.go`, `cache.go`)
- Implementation structs unexported; constructors return the interface type
- All callers updated; no subdirectory packages remain

**Non-Goals:**
- Changing any business logic or method signatures
- Modifying the domain layer
- Changing infrastructure adapter internals beyond import path updates

## Decisions

### 1. Interface named `UploadService`, struct named `uploadService`

Go allows an exported interface and an unexported struct with near-identical names in the same package. Callers see only `UploadService` (the interface); the struct is an implementation detail. Constructors return `UploadService`, ensuring callers are always programming against the interface.

Alternative considered: keep struct exported as `UploadServiceImpl`. Rejected — the `Impl` suffix is a Java-ism and adds no value in Go where the interface is the primary abstraction.

### 2. One file per service, interface + struct together

Each service file (e.g., `catalog_service.go`) holds the `CatalogService` interface, the `catalogService` struct, its constructor, and all methods. Outbound ports go in separate files because they are consumed by infrastructure adapters independently and grouping them with a specific service would misrepresent ownership.

Alternative considered: keep `ports.go` and rename it. Rejected — the original reason for splitting was clarity of inbound vs outbound; with a flat package the distinction is captured by file naming instead.

### 3. Value types and event types stay in `application/`

`MultipartUpload`, `PresignedURL`, `CompletedPart`, `VideoProcessedEvent`, `VideoFailedEvent` remain in the `application` package. They are application-layer DTOs, not domain objects, so they belong here rather than in `domain/video`.

### 4. Flat file layout

```
application/
  catalog_service.go    — CatalogService interface + catalogService struct + methods
  cache.go              — Cache interface
  upload_service.go     — UploadService interface + uploadService struct + methods + result types
  object_store.go       — ObjectStore interface + MultipartUpload, PresignedURL, CompletedPart
  event_publisher.go    — EventPublisher interface
  processing_service.go — ProcessingService interface + processingService struct + methods
  events.go             — VideoProcessedEvent, VideoFailedEvent
```

### 5. mockery.yaml consolidated to one package

`backend/api/.mockery.yaml` moves all interface targets under a single `github.com/lregnier/design-youtube/api/internal/application` entry. Generated mocks are regenerated.

## Risks / Trade-offs

[Import cycle in tests] → Eliminated. With a single package, service tests are no longer external (`package application_test`) unless we choose to keep them that way. Since generated mocks import `application` and the tests are now also in `application`, we must keep tests as `package application_test`. This is identical to the pattern already in place for upload and processing tests.

[Larger package surface] → All types are visible across the package without import qualification. Acceptable: the application layer is small and its types are intentionally shared across slices (e.g., the handler uses both `UploadService` and `CatalogService`).

## Migration Plan

1. Create new flat files under `application/` (`package application`)
2. Update `.mockery.yaml` and regenerate mocks
3. Update all infrastructure adapters to new import path and renamed interfaces
4. Update `main.go`
5. Update all test files
6. Delete the three subdirectory packages
7. Run `go build ./...` and `go test ./...`
