## Why

Several architectural changes were applied to the api and worker bounded contexts since the README files were last written. The diagrams and terminology are now out of date and would mislead anyone reading the documentation.

## What Changes

**worker/README.md:**
- Architecture diagram: rename `ResultPublisher` → `EventPublisher` in the SQS outbound node
- Architecture diagram: rename inbound node from "SQS consumer" → "SQS subscriber" (reflects `sqssubscriber` package rename)
- Sequence diagram: replace two separate `PublishProcessed` / `PublishFailed` calls with a single `Publish(DomainEvent)` call, matching the unified `EventPublisher.Publish` interface

**api/README.md:**
- Architecture diagram: rename outbound SQS node label from `Queue, processing jobs` → `EventPublisher` to match the port interface name
- Architecture diagram: rename inbound SQS node label from "SQS consumer" → "SQS subscriber"
- Application box: rename `ApplyResult` → `HandleVideoProcessingSucceeded / HandleVideoProcessingFailed` to reflect the renamed `ProcessingService` methods

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None — documentation only, no spec-level requirements change.

## Impact

- README files only; no code changes
