## ADDED Requirements

### Requirement: HTTP listen address is configurable via environment variable
The API SHALL read its HTTP listen address from the `HTTP_ADDR` environment variable. When set, the HTTP server SHALL listen on that address. When unset, the HTTP server SHALL listen on `:8080`.

#### Scenario: Listen address configured
- **WHEN** `HTTP_ADDR` is set to `:9090`
- **THEN** the HTTP server listens on `:9090`

#### Scenario: Listen address unset
- **WHEN** `HTTP_ADDR` is unset
- **THEN** the HTTP server listens on `:8080`
