## ADDED Requirements

### Requirement: Composition root contains no HTTP transport construction
`cmd/api/main.go` SHALL NOT construct the HTTP router, register middleware, define routes, or call `http.ListenAndServe` directly. It SHALL only load configuration, construct outbound adapters and use cases, construct the inbound adapters (including the HTTP server via the `internal/adapters/inbound/http` package), and start them. Router construction, middleware registration (logging, recovery, CORS), route registration, strict-handler wiring, and the listen/serve loop SHALL live in `internal/adapters/inbound/http/`.

#### Scenario: main.go delegates router construction and serving to the HTTP adapter
- **WHEN** `cmd/api/main.go` is inspected
- **THEN** it constructs a server via `internal/adapters/inbound/http` and calls its `Start()` method, with no `chi.NewRouter`, middleware registration, route definitions, or `http.ListenAndServe` call present in `main.go`
