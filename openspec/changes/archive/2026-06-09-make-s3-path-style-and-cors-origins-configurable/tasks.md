## 1. Extend API config

- [x] 1.1 Add `S3UsePathStyle bool` to the `Config` struct in `backend/api/internal/config/config.go`, parsed from `S3_USE_PATH_STYLE` via `strconv.ParseBool` (non-required, defaults to false on parse error or empty)
- [x] 1.2 Add `CORSAllowedOrigins string` to the `Config` struct, parsed from `CORS_ALLOWED_ORIGINS` env var
- [x] 1.3 Add `CORS_ALLOWED_ORIGINS` to the required fields map so the API fails fast when unset

## 2. Update main.go wiring

- [x] 2.1 Replace the unconditional `o.UsePathStyle = true` with a conditional: only set `o.UsePathStyle = true` inside the S3 options func when `cfg.S3UsePathStyle` is true; if false, pass no options func (or an empty one)
- [x] 2.2 Replace `AllowedOrigins: []string{"*"}` with `strings.Split(cfg.CORSAllowedOrigins, ",")` — add `"strings"` to imports

## 3. Update docker-compose.yml

- [x] 3.1 Add `S3_USE_PATH_STYLE: "true"` to the `api` service environment block
- [x] 3.2 Add `CORS_ALLOWED_ORIGINS: "*"` to the `api` service environment block
