## MODIFIED Requirements

### Requirement: Outbound port interfaces are defined as minimal Go interfaces
`VideoRepository`, `ObjectStore`, `Cache`, and `Queue` SHALL each be defined as Go interfaces with only the methods required by the application layer. No interface SHALL have more than 6 methods. No interface method SHALL be unused — every method SHALL have at least one call site in `internal/application/`.

#### Scenario: Adapter satisfies port interface
- **WHEN** the DynamoDB adapter is compiled against `VideoRepository`
- **THEN** the Go compiler confirms it satisfies the interface without explicit declaration

#### Scenario: Every port method has a caller
- **WHEN** `internal/ports/outbound.go` is inspected
- **THEN** every method on `VideoRepository`, `ObjectStore`, `Cache`, and `Queue` is called from at least one file under `internal/application/`
