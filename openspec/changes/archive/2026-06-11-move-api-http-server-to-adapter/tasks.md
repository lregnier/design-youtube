## 1. Make the HTTP listen address configurable

- [x] 1.1 Add `HTTPAddr` to `backend/api/internal/config/config.go`, read from `HTTP_ADDR`, defaulting to `:8080` when unset (not in the `required` set)

## 2. Move HTTP server construction into the adapter

- [x] 2.1 In `backend/api/internal/adapters/inbound/http/`, add `server.go` with an unexported `newRouter(h *Handler, uploadSecret string, corsAllowedOrigins []string) http.Handler` that builds the chi router: Logger and Recoverer middleware, CORS middleware, the `/health` route, the upload-secret strict middleware, `api.NewStrictHandlerWithOptions`, and `api.HandlerFromMux` — moved verbatim from `main.go`
- [x] 2.2 In the same file, add a `Server` type with `NewServer(h *Handler, uploadSecret string, corsAllowedOrigins []string, addr string) *Server` (storing `addr` and the router from `newRouter`) and a `Start() error` method that logs `"listening on " + addr` and calls `http.ListenAndServe(addr, router)`
- [x] 2.3 Update `backend/api/cmd/api/main.go` to construct `httpadapter.NewServer(h, cfg.UploadSecret, strings.Split(cfg.CORSAllowedOrigins, ","), cfg.HTTPAddr)` and call `srv.Start()`, removing the inline chi/middleware/route construction and the `nethttp` import if no longer used

## 3. Verify

- [x] 3.1 `go build ./...` and `go vet ./...` succeed in `backend/api`
- [x] 3.2 `go test ./...` passes in `backend/api`
- [x] 3.3 `docker compose up --build api` starts cleanly and `GET /health` returns 200

## 4. Add Handler unit tests

- [x] 4.1 Add `backend/api/internal/adapters/inbound/http/handler_test.go` covering all five `Handler` methods (`GetVideos`, `GetVideo`, `InitUpload`, `ConfirmChunk`, `CompleteUpload`): happy path plus key error-mapping branches (404/400/500), constructing real use cases backed by `internal/mocks`-mocked outbound ports, following the AAA pattern and `Test<Type>_<Method>_<Scenario>` naming
- [x] 4.2 `go test ./...` passes in `backend/api`
