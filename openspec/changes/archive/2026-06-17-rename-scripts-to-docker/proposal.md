## Why

`scripts/` exists solely to hold `localstack-init.sh`, a Docker-specific init script for LocalStack. Naming the directory `scripts/` implies general-purpose tooling; renaming it to `docker/` makes the purpose immediately clear and groups all Docker-related supporting files together. `docker-compose.yml` stays at the repo root where Docker Compose expects it by convention.

## What Changes

- `git mv scripts/ docker/`
- Update the LocalStack volume mount in `docker-compose.yml`: `./scripts/localstack-init.sh` → `./docker/localstack-init.sh`

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

(none — no spec-level requirement changes)

## Impact

- Affected files: `scripts/localstack-init.sh` (moved to `docker/`), `docker-compose.yml` (one path updated)
- No behavior change: the script content and compose service definitions are unchanged
