## ADDED Requirements

### Requirement: OpenAPI spec and codegen config live under openapi/
The OpenAPI specification, oapi-codegen config, and `go:generate` directive for `backend/api` SHALL live at `backend/api/openapi/` (`openapi.yaml`, `oapi-codegen.yaml`, `generate.go`). No directory named `api/` SHALL exist directly under `backend/api/`.

#### Scenario: OpenAPI sources are isolated under openapi/
- **WHEN** `backend/api/` is inspected
- **THEN** `openapi.yaml`, `oapi-codegen.yaml`, and `generate.go` are at `backend/api/openapi/`, and oapi-codegen output continues to be generated at `internal/gen/api/api.gen.go`
