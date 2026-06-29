## Context

Both `main.go` files wire adapters with 5–10 lines of AWS SDK option-building per adapter. The actual intent — "create a DynamoDB repository" — is hidden inside that boilerplate. The change extracts each block into a private helper, co-located in an `adapters.go` file in the same `cmd` package.

## Goals / Non-Goals

**Goals:**
- `main()` reads as a flat list of one-liner assignments
- Each adapter's initialization logic is in one focused function
- No changes outside the `cmd` packages

**Non-Goals:**
- Dependency injection framework or container
- Lazy/deferred initialization
- Changing adapter constructor signatures or port interfaces

## Decisions

**`new*` prefix** — `newRepository`, `newStore`, `newPublisher`. Matches Go's unexported private-constructor convention; reads naturally next to the adapters' own exported `New*` constructors. `build*` would also be acceptable but implies more assembly; `init*` clashes with Go's `init()` connotation.

**Helper functions in `main.go`, below `main()`** — Keeps everything in one place; there are only a handful of helpers and they don't justify a second file. `main()` appears first so the entry point is immediately visible; helpers follow beneath it.

**Already-one-line adapters stay inline** — `ffmpeg.NewTranscoder()` (worker) and `rediscache.NewCache(...)` (api) are already single calls with no config branching. Only multi-line initialization blocks get a helper.

**`sqs.Client` not extracted** — In api's `main.go` the SQS client is built twice (once for the publisher, once for the subscriber). Rather than extracting a shared client, each helper builds its own. The duplication is minor and avoids creating shared state between adapters.

**Return port interfaces, not concrete types** — Each helper's return type is the port interface (`application.ObjectStore`, `application.EventPublisher`, etc.), consistent with adapter constructors already returning interfaces.

## Risks / Trade-offs

**Slightly more indirection** → trivial; all functions are in the same package and file, not hidden behind packages or interfaces.
