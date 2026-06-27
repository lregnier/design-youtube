## MODIFIED Requirements

### Requirement: Generated code lives under gen/
All generated code in `backend/api` (oapi-codegen output, mockery mocks) SHALL live under `gen/`, in one subdirectory per generator (`gen/api/`, `gen/mocks/`). `gen/` SHALL be a top-level directory alongside `cmd/`, `internal/`, and `openapi/`. Hand-written packages under `internal/` (`internal/domain/`, `internal/application/`, `internal/infrastructure/`) SHALL NOT contain generated code.

#### Scenario: Generated packages are isolated under gen/
- **WHEN** `backend/api/` is inspected
- **THEN** the oapi-codegen output is at `gen/api/api.gen.go` and the mockery mocks are at `gen/mocks/`, with no generated files under `internal/domain/`, `internal/application/`, or `internal/infrastructure/`

## MODIFIED Requirements

### Requirement: OpenAPI spec and codegen config live under openapi/
The OpenAPI specification, oapi-codegen config, and `go:generate` directive for `backend/api` SHALL live at `backend/api/openapi/` (`openapi.yaml`, `oapi-codegen.yaml`, `generate.go`). No directory named `api/` SHALL exist directly under `backend/api/`.

#### Scenario: OpenAPI sources are isolated under openapi/
- **WHEN** `backend/api/` is inspected
- **THEN** `openapi.yaml`, `oapi-codegen.yaml`, and `generate.go` are at `backend/api/openapi/`, and oapi-codegen output is generated at `gen/api/api.gen.go`
