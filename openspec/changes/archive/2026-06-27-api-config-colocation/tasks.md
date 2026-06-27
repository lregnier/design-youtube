## 1. Move config into cmd/api/

- [x] 1.1 Copy `backend/api/internal/config/config.go` to `backend/api/cmd/api/config.go`; change `package config` to `package main`
- [x] 1.2 Remove the `internal/config` import from `backend/api/cmd/api/main.go`; update any `config.Config` or `config.Load()` references to just `Config` and `Load()` (same package, no qualifier needed)
- [x] 1.3 Delete `backend/api/internal/config/` directory
- [x] 1.4 Run `go build ./...` and confirm clean
