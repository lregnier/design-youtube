## Context

The API `CompleteUpload` use case directly constructs a JSON-encoded SQS message body and calls `queue.SendMessage(body, groupID)`. This means the use case owns the message format (`processingJob` struct, JSON marshaling) and knows SQS-specific parameters (message group ID). The `Queue` port is a thin wrapper around `SendMessage` with no domain semantics. The result is infrastructure concerns leaking into the application layer.

## Goals / Non-Goals

**Goals:**
- `EventPublisher` port speaks domain language: `Publish(ctx, event)`
- `VideoUploadedEvent` is the only domain event the API publishes today — defined in the domain layer
- JSON marshaling, queue URL selection, and message group ID construction move to the SQS adapter
- No behavior change — the same SQS message is sent, just assembled in the right layer

**Non-Goals:**
- Supporting multiple queues from a single adapter instance (only one queue today)
- Event sourcing or an event bus — this is a simple outbound port, not an event store
- Changing the worker-side message format

## Decisions

**`EventPublisher` interface with `Publish(ctx, DomainEvent) error`**
A single method accepting an interface type keeps the port minimal. The adapter type-switches on the concrete event type to determine routing and serialization. New event types require only an adapter change, not a port change.

*Alternative considered*: one method per event type (`PublishVideoUploaded`, `PublishVideoDeleted`). Rejected — the port grows with every new event; the adapter already needs to know the type to route it.

**`VideoUploadedEvent` in the domain layer**
Domain events describe things that happened in the domain — they belong next to the domain types that produce them, not in the application or infrastructure layer.

**`processingJob` JSON struct moves to the SQS adapter**
It's a wire format detail, not a domain concept. The adapter owns the translation from domain event to SQS message body.

**Rename `sqsqueue` package to `sqspublisher`**
Reflects the new role — it's not a generic queue wrapper, it's the event publisher adapter.

## Risks / Trade-offs

[Type switch in adapter is not exhaustive at compile time] If a new event type is added but the adapter isn't updated, the publish silently does nothing. → Mitigate with a default case that returns an error: `return fmt.Errorf("unknown event type: %T", event)`.
