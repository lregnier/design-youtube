## ADDED Requirements

### Requirement: Outbound port interfaces are colocated in ports.go
All outbound port interfaces in `internal/application/` SHALL be defined in a single `ports.go` file. Individual files per interface (e.g., `cache.go`, `event_publisher.go`, `object_store.go`) SHALL NOT exist.

#### Scenario: All port interfaces are in ports.go
- **WHEN** `backend/api/internal/application/` is inspected
- **THEN** `Cache`, `EventPublisher`, and `ObjectStore` are all defined in `ports.go` with no separate per-interface files
