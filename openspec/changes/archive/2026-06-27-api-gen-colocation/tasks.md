## 1. Move generated code directory

- [x] 1.1 Move `backend/api/internal/gen/` to `backend/api/gen/`

## 2. Update codegen configs

- [x] 2.1 Update `backend/api/.mockery.yaml`: change `dir: internal/gen/mocks` to `dir: gen/mocks`
- [x] 2.2 Update `backend/api/openapi/oapi-codegen.yaml`: change `output: ../internal/gen/api/api.gen.go` to `output: ../gen/api/api.gen.go`

## 3. Update import paths

- [x] 3.1 Update `backend/api/internal/application/catalog_service_test.go`: replace `…/api/internal/gen/mocks` with `…/api/gen/mocks`
- [x] 3.2 Update `backend/api/internal/application/upload_service_test.go`: replace `…/api/internal/gen/mocks` with `…/api/gen/mocks`
- [x] 3.3 Update `backend/api/internal/application/processing_service_test.go`: replace `…/api/internal/gen/mocks` with `…/api/gen/mocks`
- [x] 3.4 Update `backend/api/internal/infrastructure/in/http/handler.go`: replace `…/api/internal/gen/api` with `…/api/gen/api`
- [x] 3.5 Update `backend/api/internal/infrastructure/in/http/middleware.go`: replace `…/api/internal/gen/api` with `…/api/gen/api`
- [x] 3.6 Update `backend/api/internal/infrastructure/in/http/server.go`: replace `…/api/internal/gen/api` with `…/api/gen/api`
- [x] 3.7 Update `backend/api/internal/infrastructure/in/http/handler_test.go`: replace `…/api/internal/gen/mocks` and `…/api/internal/gen/api` with `…/api/gen/mocks` and `…/api/gen/api`

## 4. Update generated file package path comment

- [x] 4.1 Update the package path embedded in `backend/api/gen/api/api.gen.go` if oapi-codegen emitted `internal/gen/api` in any comment or import

## 5. Verify

- [x] 5.1 Run `go build ./...` and confirm clean
- [x] 5.2 Run `go test ./...` and confirm all tests pass
