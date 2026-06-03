## Context

The backend was implemented as a flat package structure: `internal/handler/`, `internal/store/`, and `internal/middleware/`. Business logic, infrastructure calls, and HTTP concerns are mixed together — the handler directly constructs S3 presigned URLs and DynamoDB records. This makes the core logic untestable without real AWS infrastructure and obscures the architecture to a code reviewer.

The refactor introduces a deliberate blend of three architectural influences:

- **DDD tactical patterns** (Vernon, "Implementing Domain-Driven Design"): `Video` aggregate, `VideoID` value object, `VideoRepository` interface, and ubiquitous language naming (`InitUpload`, `ConfirmChunk`). These live in the domain layer.
- **Clean Architecture** (Martin, "Clean Architecture"): the four-ring dependency model — Entities → Use Cases → Interface Adapters → Frameworks & Drivers. Our `domain/`, `application/`, `ports/`, and `adapters/` map directly to these rings. The strict inward dependency rule is enforced: inner packages never import outer ones.
- **Hexagonal Architecture** (Cockburn): the ports-and-adapters split that makes infrastructure swappable. Clean Architecture absorbs this pattern — our `ports/` and `adapters/` are its concrete expression.

The result is closer to Clean Architecture in structure (explicit use-case layer, strict dependency rule) with DDD vocabulary and tactical patterns inside. Vernon's layered DDD is a valid alternative but puts less emphasis on the use-case layer; for this portfolio the Clean Architecture framing makes the intent more legible. Everything is kept idiomatic Go: small focused interfaces, no unnecessary wrapping, packages named by what they do rather than abstract layer labels.

## Goals / Non-Goals

**Goals:**
- Domain layer with zero infrastructure imports — `Video` aggregate and `VideoRepository` interface compile with no external dependencies
- Application layer that orchestrates use cases by depending on port interfaces, not concrete implementations
- Outbound adapters (DynamoDB, S3, Redis, SQS) that satisfy port interfaces and contain all AWS SDK code
- Inbound HTTP adapter that wires oapi-codegen strict server to use cases
- Identical external behavior — same API, same DynamoDB schema, same S3 key structure
- Go-idiomatic: interfaces defined where they are consumed (application layer), not in a separate `interfaces/` package

**Non-Goals:**
- Domain events or event sourcing
- CQRS — reads and writes use the same repository interface
- Changing the OpenAPI spec, generated code, Terraform, or frontend
- Adding tests (separate concern, separate change)

## Decisions

### 1. Architectural lineage: Clean Architecture over Vernon's DDD layering

Vernon's DDD describes four layers (UI → Application → Domain → Infrastructure) but treats the Application layer as thin orchestration with little emphasis on explicit use-case types. Martin's Clean Architecture gives the Use Cases ring equal status to Entities and enforces the dependency rule formally. This design follows Clean Architecture's structure because it makes the intent more legible in code — a reader can immediately see which ring a package belongs to.

DDD tactical patterns (aggregate, value object, repository interface) are applied inside the domain ring. The two schools are complementary: DDD tells you *what* to put in the domain; Clean Architecture tells you *how to structure the layers around it*.

| Clean Architecture ring | Our package         |
|-------------------------|---------------------|
| Entities                | `domain/video/`     |
| Use Cases               | `application/`      |
| Interface Adapters      | `ports/`, `adapters/inbound/` |
| Frameworks & Drivers    | `adapters/outbound/` |

### 2. Package layout

```
internal/
  domain/
    video/
      video.go        # Video aggregate + value types (VideoID, VideoStatus, Chunk)
      repository.go   # VideoRepository interface
  application/
    upload/
      init.go         # InitUpload use case
      confirm_chunk.go
      complete.go
    catalog/
      get_video.go    # GetVideo use case
      list_videos.go
  ports/
    outbound.go       # ObjectStore, Cache, Queue interfaces
  adapters/
    inbound/
      http/
        handler.go    # StrictServerInterface implementation, wires use cases
        middleware.go # Upload secret StrictMiddlewareFunc
    outbound/
      dynamo/         # VideoRepository implementation
      s3store/        # ObjectStore implementation
      rediscache/     # Cache implementation
      sqsqueue/       # Queue implementation
```

Rationale: naming by domain concept (`upload/`, `catalog/`) rather than CRUD layer. Keeps related use cases co-located.

### 2. VideoRepository interface lives in the domain package

Go's convention is "accept interfaces, return structs" and "define interfaces where they are used." The `VideoRepository` interface is consumed by the domain and application layers, so it lives in `domain/video/`. The DynamoDB adapter imports `domain/video` and implements it — dependency flows inward, never outward.

Alternatives considered: defining all interfaces in a top-level `ports/` package. Rejected — splits the interface from the type it operates on, which is un-Go-like.

### 3. ObjectStore, Cache, and Queue interfaces in `ports/outbound.go`

These are infrastructure concerns with no domain semantics. They are thin, method-count-minimal interfaces. Keeping them in a single `ports/outbound.go` avoids creating a package per interface while still making the contracts explicit.

### 4. Use cases as structs with a single Execute (or named) method

```go
type InitUpload struct {
    repo  video.VideoRepository
    store ports.ObjectStore
}

func (uc InitUpload) Execute(ctx context.Context, cmd InitUploadCommand) (InitUploadResult, error)
```

Each use case takes its dependencies at construction time (dependency injection via constructor), receives a typed command, returns a typed result. No global state. This keeps each use case independently testable by swapping in mock implementations of the interfaces.

Alternatives considered: plain functions. Functions work but make dependency injection awkward and don't group related config. Structs with a single public method are idiomatic for this pattern in Go.

### 5. Inbound HTTP adapter delegates entirely to use cases

The HTTP handler (`adapters/inbound/http/handler.go`) implements `api.StrictServerInterface`. Each method constructs a command from the request object, calls the relevant use case, and maps the result to a response type. No business logic lives in the handler.

### 6. Upload secret middleware stays in the inbound adapter

The secret check is an HTTP transport concern, not a domain rule. It stays in `adapters/inbound/http/middleware.go` as a `StrictMiddlewareFunc`, exactly as before — just relocated.

## Risks / Trade-offs

- **More files, same logic** → The refactor adds packages without adding features. Justified for portfolio legibility; would need justification in a production codebase.
- **Interface proliferation** → Go's `io.Reader` philosophy warns against interfaces for their own sake. Mitigation: keep interfaces to what's actually needed — `VideoRepository` has ~5 methods, `ObjectStore` has ~3.
- **Worker binary needs updating** → `cmd/worker/main.go` directly calls AWS SDK today. After the refactor it should compose adapters the same way `cmd/api/main.go` does. Risk of missing this. Mitigation: covered explicitly in tasks.

## Migration Plan

1. Create new package structure alongside existing packages
2. Implement domain types and port interfaces
3. Implement use cases against interfaces
4. Implement outbound adapters (can test against real infra or LocalStack)
5. Implement inbound HTTP adapter
6. Update `cmd/api/main.go` and `cmd/worker/main.go` to compose new structure
7. Delete old `internal/handler/`, `internal/store/`, `internal/middleware/`
8. Verify `go build ./...` and `go vet ./...` pass
