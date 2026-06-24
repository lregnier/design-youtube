## Context

`GET /videos` currently calls `repo.ListReady()`, which queries a DynamoDB GSI (`status-uploadedAt-index`) filtering by `status = "ready"`. The `VideoSummary` response schema already includes a `status` field, so no API contract changes are needed. `VideoCard` always renders as a `<Link>` and assumes `thumbnailUrl` is present. `HomePage` fetches once on mount with no refresh.

## Goals / Non-Goals

**Goals:**
- Videos appear on the home page immediately after upload completes (in `processing` state)
- The home page auto-refreshes until all processing videos resolve
- `VideoCard` communicates status clearly without being clickable for non-ready videos

**Non-Goals:**
- WebSockets or server-sent events for real-time push (polling is sufficient for this scale)
- Showing videos in `uploading` state (the upload flow is still in progress)
- Persisting polling state across page reloads
- Status filtering or sorting controls on the home page

## Decisions

### Replace `ListReady` with `List` (Scan excluding `uploading`)

The existing `ListReady` uses a GSI query for a single status value. To return multiple statuses (`processing`, `ready`, `failed`) we have three options:

1. **Multiple GSI queries** (one per status, merge and sort in Go) — correct but verbose
2. **DynamoDB Scan with FilterExpression** — simpler; scans all items but acceptable at this project's data volume
3. **New GSI on a different partition key** — over-engineered for this scale

**Decision:** DynamoDB Scan with `FilterExpression: status <> :uploading`, then sort the results by `uploadedAt` descending in Go. Simple, no infrastructure change, correct for the dataset size.

The `VideoRepository` interface drops `ListReady` and adds `List` with the same signature. The mock is regenerated via mockery.

### Poll every 5 seconds, stop when no processing videos remain

Alternatives considered:
- **Manual refresh button** — requires user action, poor UX for a background job
- **WebSocket / SSE push** — accurate but adds backend infrastructure complexity
- **Constant polling regardless of state** — wastes requests when nothing is processing

**Decision:** `HomePage` starts a 5-second `setInterval` when the video list contains at least one `processing` entry. The interval is cleared (and not started) when all videos are `ready` or `failed`. The interval is also cleared on component unmount.

### `VideoCard` renders a `<div>` for non-ready videos instead of `<Link>`

Non-ready videos must not be navigable (the video page requires a manifest URL). Wrapping in a `<div>` instead of `<Link>` avoids cursor/focus confusion without extra click-handler logic.

- `processing`: thumbnail area replaced with a dark placeholder and centered "Processing…" text
- `failed`: thumbnail area replaced with a dark placeholder and centered "Failed" text in red

## Risks / Trade-offs

- **Scan cost** → acceptable at this project's scale; switch to multi-query if the table grows significantly
- **5-second poll** → adds load proportional to open browser tabs with processing videos; acceptable for a single-user tool
- **`thumbnailUrl` may be empty for processing/failed** → `VideoCard` must not render a broken `<img>` tag for non-ready videos; the placeholder div handles this
