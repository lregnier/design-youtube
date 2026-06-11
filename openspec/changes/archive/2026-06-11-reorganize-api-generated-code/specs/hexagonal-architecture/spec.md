## ADDED Requirements

### Requirement: Generated code lives under internal/gen/
All generated code in `backend/api` (oapi-codegen output, mockery mocks) SHALL live under `internal/gen/`, in one subdirectory per generator (`internal/gen/api/`, `internal/gen/mocks/`). Hand-written packages (`internal/domain/`, `internal/application/`, `internal/ports/`, `internal/adapters/`, `internal/config/`) SHALL NOT contain generated code.

#### Scenario: Generated packages are isolated under internal/gen
- **WHEN** `backend/api/internal/` is inspected
- **THEN** the oapi-codegen output is at `internal/gen/api/api.gen.go` and the mockery mocks are at `internal/gen/mocks/`, with no generated files under `internal/domain/`, `internal/application/`, `internal/ports/`, or `internal/adapters/`
