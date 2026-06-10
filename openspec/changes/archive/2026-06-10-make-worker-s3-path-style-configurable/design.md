## Context

The API already exposes `S3_USE_PATH_STYLE` (see `api-runtime-config`) so its S3 client can talk to LocalStack via path-style addressing. The worker's S3 client has no such option, so against LocalStack it issues virtual-hosted-style requests for a hostname that doesn't resolve in the docker-compose network, causing every `DownloadRaw` to fail and the video to end up `failed`.

## Goals / Non-Goals

**Goals:**
- `S3_USE_PATH_STYLE` env var controls path-style S3 access for the worker (default: false), identical semantics to the API's flag
- Local dev docker-compose sets it for the `worker` service so LocalStack works end-to-end

**Non-Goals:**
- Any other S3 client option for the worker
- Changing the API's existing `S3_USE_PATH_STYLE` behavior

## Decisions

**Mirror the API's config field exactly**
Add `S3UsePathStyle bool` to `backend/worker/internal/config/config.go`, parsed the same way as the API's field (optional, `strconv.ParseBool`, defaults to false). Two independent flags (one per service) rather than a shared config package — the services don't currently share config code, and introducing that coupling is out of scope.

**Conditional `s3Opts` in `cmd/worker/main.go`**
Same pattern as `backend/api/cmd/api/main.go`: build an `[]func(*awss3.Options)` slice, append `UsePathStyle = true` only when the flag is set, pass to `awss3.NewFromConfig`.

## Risks / Trade-offs

[Two copies of the same flag/parsing logic] → Acceptable: the api and worker are separate Go modules with no shared internal package today; deduplicating would require introducing a shared module, which is out of scope for this fix.
