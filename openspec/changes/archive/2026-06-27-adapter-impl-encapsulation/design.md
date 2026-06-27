## Context

Both bounded contexts follow hexagonal architecture where port interfaces are defined in the `application` layer and implemented by infrastructure adapters. However, adapter structs are currently exported (`Repository`, `Publisher`, etc.) and constructors return concrete pointer types (`*Repository`, `*Publisher`). This means `main.go` holds references to concrete types rather than interfaces, and any package that imports the adapter package can instantiate or type-assert to the concrete struct.

## Goals / Non-Goals

**Goals:**
- Unexport all adapter implementation structs
- Make constructors return the port interface type
- Keep compile-time interface satisfaction checks in place

**Non-Goals:**
- Introducing new port interfaces where none exist (e.g., `Subscriber` in main is used by calling `.Start()` directly — out of scope)
- Changing any behavior or wire format
- Modifying the port interface definitions themselves

## Decisions

**Unexported struct + interface return type**
Each adapter follows the same mechanical pattern:
```go
// Before
type Publisher struct { ... }
func NewPublisher(...) *Publisher { ... }

// After
type publisher struct { ... }
func NewPublisher(...) application.EventPublisher { return &publisher{...} }
```
The compile-time guard `var _ Interface = (*impl)(nil)` continues to work with unexported structs and acts as documentation that the struct satisfies the interface.

**Subscriber structs have no port interface — leave exported**
`sqssubscriber.Subscriber` in both api and worker is called directly via `.Start(ctx)` in `main.go`. There is no port interface for inbound adapters (they are drivers, not driven). These remain exported until a `Driver` interface pattern is introduced separately.

## Risks / Trade-offs

- **Low risk** — purely mechanical rename; no logic changes
- Callers in `main.go` that previously relied on the concrete type for additional methods beyond the interface will break — but inspection confirms none do
