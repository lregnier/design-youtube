## Context

`backend/api/internal/config/config.go` defines a `Config` struct and a `Load()` function that reads environment variables. It is imported only by `cmd/api/main.go`. No other package in the module uses it.

## Goals / Non-Goals

**Goals:**
- Eliminate the unnecessary `internal/config` package by moving its single file into `cmd/api/` as `package main`
- Remove the corresponding import from `main.go`

**Non-Goals:**
- Changing any config fields, defaults, or loading logic
- Making config available to other packages

## Decisions

### Move to `package main`, not a sub-package of `cmd/api/`

`cmd/api/config.go` uses `package main` — same package as `main.go`. No import needed; types and functions are directly accessible. This is the idiomatic Go approach when config is only needed by the entrypoint.

Alternative considered: `cmd/api/config/config.go` as a separate package. Rejected — it just recreates the same problem with a different path.

## Risks / Trade-offs

[Testability] → Config loading is currently in its own package, which could theoretically be tested in isolation. In practice there are no tests for it, and entrypoint config is rarely unit-tested. No regression risk.
