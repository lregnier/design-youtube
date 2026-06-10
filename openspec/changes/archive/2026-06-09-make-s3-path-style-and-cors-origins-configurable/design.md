## Context

After adding CORS middleware and S3 path-style access for local dev, both values are hardcoded. `AllowedOrigins: ["*"]` is a security concern in production. `UsePathStyle: true` is deprecated by AWS for production buckets. Both need to be driven by environment variables.

## Goals / Non-Goals

**Goals:**
- `S3_USE_PATH_STYLE` env var controls path-style S3 access (default: false)
- `CORS_ALLOWED_ORIGINS` env var controls allowed origins (required, comma-separated)
- Local dev docker-compose sets both for LocalStack compatibility

**Non-Goals:**
- Per-route CORS configuration
- CORS credentials, max-age, or expose-headers configuration
- Any other S3 client options

## Decisions

**`S3UsePathStyle` as bool, not string**
`strconv.ParseBool` handles `"true"/"false"/"1"/"0"` — cleaner than a string comparison in main.go. The field is optional with a false default, so an unset var means virtual-hosted style (correct for prod).

**`CORSAllowedOrigins` as a single comma-separated string in config, split at wire-up**
Keeps config.go simple (one field, one env var). The split happens in main.go where the chi cors handler is constructed. This matches the common pattern for list-typed env vars.

**`CORS_ALLOWED_ORIGINS` is required**
The API must not silently allow all origins in production due to a missing env var. Requiring it forces an explicit decision per environment. Local dev sets `"*"` deliberately.

## Risks / Trade-offs

[Misconfigured origins in prod] → If an operator sets an incorrect domain, CORS preflight requests fail silently for end users. Mitigation: document the expected format in the env var table in the API README.

[Path-style deprecation] → AWS has paused enforcement but may re-enable it. `S3_USE_PATH_STYLE` must remain unset (or `false`) in production from the start to avoid future breakage.
