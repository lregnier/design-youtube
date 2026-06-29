## ADDED Requirements

### Requirement: README diagrams reflect current architecture
The README files for `backend/api` and `backend/worker` SHALL use terminology that matches the current codebase: port interface names (`EventPublisher`), adapter naming (`sqssubscriber`), and service method names (`HandleVideoProcessingSucceeded`, `HandleVideoProcessingFailed`).

#### Scenario: Worker README uses EventPublisher
- **WHEN** `backend/worker/README.md` is read
- **THEN** the outbound SQS node is labeled `EventPublisher` and the sequence diagram shows a single `Publish(DomainEvent)` call

#### Scenario: API README uses current service method names
- **WHEN** `backend/api/README.md` is read
- **THEN** the processing application box lists `HandleVideoProcessingSucceeded` and `HandleVideoProcessingFailed`, not `ApplyResult`
