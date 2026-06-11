## Why

`backend/api/cmd/api/main.go` currently builds the entire HTTP transport inline: it creates the chi router, registers the Logger/Recoverer/CORS middleware, defines the `/health` route, wires the upload-secret strict middleware, wraps the generated `StrictServerInterface` with `api.NewStrictHandlerWithOptions`, mounts it via `api.HandlerFromMux`, and calls `http.ListenAndServe`. Meanwhile `internal/adapters/inbound/http/` already holds the `Handler` (business logic translation) and the upload-secret middleware — the router/server assembly is the missing piece of that adapter, left stranded in the composition root. This is inconsistent with the hexagonal architecture already enforced for the rest of the api and worker bounded contexts, where `main.go` only wires dependencies and starts the entrypoint.

## What Changes

- Add `httpadapter.NewServer(h *Handler, uploadSecret string, corsAllowedOrigins []string, addr string) *Server` in `backend/api/internal/adapters/inbound/http/`, encapsulating the chi router build (Logger/Recoverer middleware, CORS middleware, `/health` route, upload-secret strict middleware, `api.NewStrictHandlerWithOptions`, `api.HandlerFromMux`) plus a `Start() error` method that wraps `http.ListenAndServe(addr, router)` — mirroring the existing `sqsconsumer.Consumer.Start(ctx)` pattern for the other inbound adapter.
- Update `backend/api/cmd/api/main.go` to construct `httpadapter.NewServer(...)` and call `srv.Start()`. `main.go` keeps only: config loading, AWS config, outbound adapter construction, use case construction, and constructing/starting both inbound adapters (SQS results consumer and HTTP server).
- Make the HTTP listen address configurable via a new `HTTP_ADDR` environment variable, defaulting to `:8080` when unset (preserves current hardcoded behavior).
- No other behavior changes: same routes, same middleware order, same CORS/upload-secret config inputs, same `/health` response.
- Add unit tests for `Handler` (`internal/adapters/inbound/http/handler_test.go`) covering its request/response translation logic for all five operations, using `internal/mocks`-backed use cases and following the same AAA/mockery conventions as `internal/application/`.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- **`hexagonal-architecture`**: add a requirement that the composition root (`main.go`) contains no HTTP transport construction (router, middleware, route registration, listen/serve) — that responsibility belongs to the inbound HTTP adapter package via a `Start()` method
- **`api-runtime-config`**: add a requirement that the HTTP listen address is configurable via the `HTTP_ADDR` environment variable, defaulting to `:8080`
- **`backend-unit-tests`**: add a requirement that the inbound HTTP adapter's `Handler` translation logic has unit tests, following the same AAA/mockery conventions as `internal/application/`

## Impact

- **`backend/api/cmd/api/main.go`**: removes router/middleware/server construction; becomes a pure composition root that constructs and starts both inbound adapters
- **`backend/api/internal/adapters/inbound/http/`**: new file containing the `Server` type (chi router/middleware wiring + `Start()`), moved from `main.go`; new `handler_test.go` with unit tests for `Handler`
- **`backend/api/internal/config/config.go`**: new optional `HTTPAddr` field, read from `HTTP_ADDR`, defaulting to `:8080`
- No external behavior changes — same endpoints, middleware, CORS and upload-secret enforcement, same default port
