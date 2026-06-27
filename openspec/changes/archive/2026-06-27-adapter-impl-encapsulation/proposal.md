## Why

Infrastructure adapter structs are currently exported with constructors returning the concrete pointer type, which leaks implementation details to callers and breaks the hexagonal architecture principle that only port interfaces should cross layer boundaries.

## What Changes

- Unexport all infrastructure adapter implementation structs (e.g., `Repository` → `repository`, `Publisher` → `publisher`, `Subscriber` → `subscriber`, `Store` → `store`, `Transcoder` → `transcoder`)
- Update constructors to return the corresponding port interface instead of `*ConcreteType`

Affected adapters:
- **api** `infrastructure/out/dynamo/repository.go`: `Repository` → `repository`, returns `video.VideoRepository`
- **api** `infrastructure/out/s3/store.go`: `Store` → `store`, returns `application.ObjectStore`
- **api** `infrastructure/out/sqspublisher/publisher.go`: `Publisher` → `publisher`, returns `application.EventPublisher`
- **api** `infrastructure/in/sqssubscriber/subscriber.go`: `Subscriber` → `subscriber` (no port interface — keep returning concrete if no interface exists, or introduce one)
- **worker** `infrastructure/out/sqspublisher/publisher.go`: `Publisher` → `publisher`, returns `application.EventPublisher`
- **worker** `infrastructure/out/s3storage/store.go`: `Store` → `store`, returns `application.VideoStorage`
- **worker** `infrastructure/out/ffmpeg/transcoder.go`: `Transcoder` → `transcoder`, returns `application.Transcoder`
- **worker** `infrastructure/in/sqssubscriber/subscriber.go`: `Subscriber` → `subscriber`, returns nothing new (Start is called directly in main)

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `hexagonal-architecture`: Add requirement that adapter implementation structs are unexported and constructors return port interfaces
- `worker-hexagonal-architecture`: Same requirement for the worker bounded context

## Impact

- `cmd/api/main.go` and `cmd/worker/main.go` callers receive interface types — no behavioral change, but type of local variable changes
- Compile-time interface checks (`var _ Interface = (*impl)(nil)`) remain valid with unexported struct
- No API or wire format changes
