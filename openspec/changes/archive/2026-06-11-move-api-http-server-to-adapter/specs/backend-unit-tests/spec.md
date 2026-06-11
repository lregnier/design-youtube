## ADDED Requirements

### Requirement: Inbound HTTP adapter handler has unit tests
The `Handler` in `backend/api/internal/adapters/inbound/http/` SHALL have a corresponding `handler_test.go` covering its request/response translation logic for each operation, following the same AAA pattern, mockery v2 mocks, and `Test<Type>_<Method>_<Scenario>` naming conventions as `internal/application/`.

#### Scenario: Handler tests run without infrastructure
- **WHEN** `go test ./internal/adapters/inbound/http/...` is run in `backend/api/`
- **THEN** all tests pass without requiring AWS credentials, a running database, or a real HTTP server, using `internal/mocks`-backed use cases
