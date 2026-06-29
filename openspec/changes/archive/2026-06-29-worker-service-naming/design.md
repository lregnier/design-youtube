## Context

The api uses `*Service` suffix for all use case interfaces and domain-specific method names. The worker's `ProcessVideo` + `Execute` predates this convention. The api's `ProcessingService` has an ambiguous name — it doesn't convey that it only handles status transitions. This change normalizes both bounded contexts.

## Goals / Non-Goals

**Goals:**
- Consistent `*Service` suffix on all use case interfaces in both bounded contexts
- Domain-specific method names (`Process`, `MarkReady`, `MarkFailed`) instead of generic verbs
- `VideoStatusService` name that reflects the actual responsibility (status management)

**Non-Goals:**
- Changing method signatures, parameters, or logic
- Renaming port interfaces (`VideoStorage`, `Transcoder`, `EventPublisher`)
- Moving status transition responsibility across services (UploadService keeps its own transitions)

## Decisions

**`VideoProcessingService` (worker)** — describes the capability (video processing), consistent with other `*Service` names in the codebase. `Process` as the method name matches the domain verb and replaces the generic `Execute`.

**`VideoStatusService` (api)** — the service only transitions video status (`uploading → processing → ready/failed`). The name makes the narrow responsibility explicit. `MarkReady`/`MarkFailed` are domain verbs that match what the service does (they mirror `video.MarkReady()` / `video.MarkFailed()` on the domain model).

**UploadService keeps its status transitions** — `InitUpload` (sets `uploading`) and `CompleteUpload` (sets `processing`) are inseparable from their upload mechanics (multipart S3, SQS dispatch). Centralizing them into `VideoStatusService` would introduce a service-to-service dependency with no real benefit.
