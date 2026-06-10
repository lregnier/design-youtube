## Why

`UsePathStyle` and `AllowedOrigins: ["*"]` are hardcoded for local dev but wrong for production — path-style S3 access is deprecated by AWS, and a wildcard CORS origin is a security risk. Both values need to come from environment variables so each environment can configure them correctly.

## What Changes

- Add `S3UsePathStyle bool` to the API config, parsed from `S3_USE_PATH_STYLE` env var (optional, defaults to false)
- Add `CORSAllowedOrigins string` to the API config, parsed from `CORS_ALLOWED_ORIGINS` env var (required)
- In `main.go`: apply `UsePathStyle` only when the config flag is true; split `CORSAllowedOrigins` by comma and pass as `AllowedOrigins`
- In `docker-compose.yml`: set `S3_USE_PATH_STYLE: "true"` and `CORS_ALLOWED_ORIGINS: "*"` for local dev

## Capabilities

### New Capabilities

_(none — this is purely a configurability improvement, no new user-facing behaviour)_

### Modified Capabilities

_(no spec-level requirement changes — existing behaviour is preserved, just made environment-driven)_

## Impact

- **`backend/api/internal/config/config.go`**: two new fields
- **`backend/api/cmd/api/main.go`**: conditional `UsePathStyle`, dynamic `AllowedOrigins`
- **`docker-compose.yml`**: two new env vars in the `api` service block
