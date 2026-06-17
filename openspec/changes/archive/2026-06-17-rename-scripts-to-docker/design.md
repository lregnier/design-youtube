## Context

`scripts/` contains one file: `localstack-init.sh`, mounted into the LocalStack container via `docker-compose.yml`. The rename is purely cosmetic — no content changes, no new dependencies, no migration concerns.

## Goals / Non-Goals

**Goals:**
- Rename `scripts/` → `docker/` to reflect that its contents are Docker-specific
- Update the volume mount path in `docker-compose.yml` to match

**Non-Goals:**
- Moving `docker-compose.yml` (stays at root)
- Changing the script content or compose service definitions
