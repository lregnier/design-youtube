## Context

The worker's `internal/application/ports.go` defines all outbound port interfaces (`VideoStorage`, `Transcoder`, `EventPublisher`) in one place. The api had the same interfaces spread across three files. Consolidating them mirrors the worker and follows the single-file port convention.

## Goals / Non-Goals

**Goals:**
- Single `ports.go` file in `backend/api/internal/application/` containing `Cache`, `EventPublisher`, `ObjectStore`, and the value types used by `ObjectStore` (`MultipartUpload`, `PresignedURL`, `CompletedPart`)

**Non-Goals:**
- Changing any interface definitions or method signatures
- Moving service interfaces (`UploadService`, `CatalogService`, `ProcessingService`)

## Decisions

**One file, no subdirectories** — consistent with the worker pattern and the existing requirement that `internal/application/` is a flat package.
