## 1. Rename the directory

- [x] 1.1 `git mv scripts docker`

## 2. Update docker-compose.yml

- [x] 2.1 Update the LocalStack volume mount from `./scripts/localstack-init.sh` to `./docker/localstack-init.sh`

## 3. Verify

- [x] 3.1 Grep for any remaining references to `scripts/` in the repo (excluding git history and openspec)
- [x] 3.2 Confirm `docker compose config` parses cleanly (validates the compose file)
