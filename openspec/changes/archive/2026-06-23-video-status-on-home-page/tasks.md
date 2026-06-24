## 1. Backend — repository

- [x] 1.1 In `backend/api/internal/domain/video/repository.go`: rename `ListReady` to `List` in the `VideoRepository` interface (same signature: `List(ctx context.Context) ([]*Video, error)`)
- [x] 1.2 In `backend/api/internal/adapters/outbound/dynamo/repository.go`: rename `ListReady` to `List` and replace the GSI query with a DynamoDB Scan that uses `FilterExpression: "#st <> :uploading"` and sorts results by `uploadedAt` descending in Go after unmarshalling

## 2. Backend — application and mock

- [x] 2.1 In `backend/api/internal/application/catalog/list_videos.go`: change `uc.repo.ListReady(ctx)` to `uc.repo.List(ctx)`
- [x] 2.2 In `backend/api/internal/gen/mocks/mock_VideoRepository.go`: rename the generated `MockVideoRepository.ListReady` method to `List` (update method name, `_ListReady_Call` struct, and all references within the file)
- [x] 2.3 In `backend/api/internal/application/catalog/list_videos_test.go`: update the mock call from `repo.EXPECT().ListReady(...)` to `repo.EXPECT().List(...)`

## 3. Frontend — VideoCard

- [x] 3.1 In `frontend/web/src/components/VideoCard.tsx`: for `status === "ready"`, keep the current `<Link>` wrapper with thumbnail `<img>`; for `status === "processing"`, render a `<div>` (not a Link) with a dark `#111` placeholder replacing the thumbnail area containing centered "Processing…" text in `#aaa`; for `status === "failed"`, render a `<div>` with the same placeholder containing centered "Failed" text in `#c00`
- [x] 3.2 In `frontend/web/src/components/VideoCard.tsx`: add a CSS spinner above the "Processing…" label for `status === "processing"` — a 20px rotating circle via `@keyframes _vc-spin` injected with a `<style>` tag, border-based with a solid top segment in `#aaa`

## 4. Frontend — HomePage polling

- [x] 4.1 In `frontend/web/src/pages/HomePage.tsx`: extract the fetch into a `loadVideos` function; after each successful fetch, if any returned video has `status === "processing"` start (or keep) a `setInterval` of 5000 ms that calls `loadVideos`; if none are processing, clear the interval; clear the interval on unmount via `useEffect` cleanup

## 5. Verify

- [x] 5.1 Run `go test ./internal/...` in `backend/api/` and confirm all tests pass
- [ ] 5.2 Start the app and upload a video; confirm it appears on the home page as "Processing…" immediately after upload, then transitions to a clickable ready card once the worker finishes
