## Context

`backend/api/cmd/api/main.go` wires outbound adapters, use cases, and inbound adapters, then inline-builds the chi router (middleware, `/health`, strict-handler wrapping, route mounting) and starts `http.ListenAndServe`. `internal/adapters/inbound/http/` already contains `handler.go` (the `StrictServerInterface` implementation) and `middleware.go` (the upload-secret middleware), but not the router/server assembly that ties them together — that logic currently lives in `main.go` instead of the adapter package.

## Goals / Non-Goals

**Goals:**
- Move chi router construction, middleware registration (Logger, Recoverer, CORS), the `/health` route, and strict-handler/mux wiring into `internal/adapters/inbound/http/`.
- Leave `main.go` as a pure composition root: load config, build outbound adapters and use cases, construct the inbound adapters (HTTP router + SQS results consumer), start both.

**Non-Goals:**
- No change to routes, middleware behavior, CORS configuration, upload-secret enforcement, or the `/health` response.
- No change to the SQS results consumer (`sqsconsumer`) — it's already a self-contained adapter.
- No change to the generated `api` package (oapi-codegen output).

## Decisions

- **`Server` struct with `Start() error`, mirroring `sqsconsumer.Consumer.Start(ctx)`.** `NewServer(h *Handler, uploadSecret string, corsAllowedOrigins []string, addr string) *Server` builds the chi router internally and stores it alongside `addr`. `Start()` logs `"listening on " + addr` and calls `http.ListenAndServe(addr, router)`. This makes `main.go` symmetric for both inbound adapters: construct, then `Start` (one via `go`, one blocking) — rather than one adapter returning a handler for `main.go` to feed into `ListenAndServe` directly.
- **Router construction stays as an unexported helper (`newRouter(...) http.Handler`) called by `NewServer`.** This keeps the chi/middleware assembly unit-testable in isolation (e.g. via `httptest.NewServer(srv.Handler())` if ever needed) without `Start` being the only entry point, while `Start` remains the normal path used by `main.go`.
- **Constructor takes the already-constructed `*Handler` and the specific config values it needs (`uploadSecret`, `corsAllowedOrigins`, `addr`)**, mirroring how `httpadapter.NewHandler(...)` and `UploadSecretMiddleware(cfg.UploadSecret)` are already called individually today — just consolidated into one place. Alternative considered: pass the whole `config.Config` into the adapter package — rejected because it would couple the adapter to the config package's shape for values it doesn't otherwise need.
- **HTTP listen address becomes configurable via `HTTP_ADDR`, defaulting to `:8080`.** Previously hardcoded as `":8080"` in `main.go`'s `ListenAndServe` call. Since the address now has to be threaded into `NewServer`, hardcoding it in `main.go` would just move the hardcode rather than remove it. Defaulting (not requiring) `HTTP_ADDR` preserves current behavior for docker-compose/deployments that don't set it. Alternative considered: a numeric `PORT` var (e.g. `"8080"`, prefixed with `:` internally) — rejected in favor of `HTTP_ADDR` since it maps directly to `http.ListenAndServe`'s `addr` argument with no transformation and allows binding a specific host if ever needed.
- **File name `server.go`** in `internal/adapters/inbound/http/`, type name `Server`, constructor `NewServer`, consistent with `NewHandler`/`Handler` naming already in the package.
- **`Handler` unit tests construct real use cases backed by `internal/mocks`-mocked outbound ports, no new interfaces.** `Handler` already holds concrete use-case structs (`upload.InitUpload`, `catalog.GetVideo`, etc.), each constructed from `video.VideoRepository`/`ports.ObjectStore`/`ports.Cache`/`ports.Queue`. Tests build a `Handler` via `NewHandler(...)` with use cases wired to `mocks.NewMockVideoRepository(t)` etc., then call `Handler` methods directly with `api.XxxRequestObject` and assert on the returned `api.XxxResponseObject`. `server.go` (router/middleware wiring) is not unit tested — it has no branching logic of its own.

## Risks / Trade-offs

- [Middleware ordering subtly changes behavior] → Mitigation: move the existing chi setup verbatim (same order: Logger, Recoverer, CORS, then routes), no reordering.
- [`api` package import cycle between adapter and generated code] → Mitigation: `internal/adapters/inbound/http` already imports `internal/api` (for `Handler`/`StrictServerInterface`), so adding `api.NewStrictHandlerWithOptions`/`api.HandlerFromMux` calls there introduces no new dependency direction.
- [`HTTP_ADDR` typo/misconfiguration silently changes the listen address] → Mitigation: default to `:8080` when unset, matching current behavior; no validation beyond what `http.ListenAndServe` itself performs (consistent with how the address was handled before this change).
