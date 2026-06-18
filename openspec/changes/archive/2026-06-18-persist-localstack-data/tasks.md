## 1. Make localstack-init.sh idempotent

- [x] 1.1 Guard `s3 mb` with an existence check: skip bucket creation if `awslocal s3api head-bucket --bucket design-youtube-video-prod` succeeds
- [x] 1.2 Guard `dynamodb create-table` with an existence check: skip table creation if `awslocal dynamodb describe-table --table-name videos` succeeds
- [x] 1.3 Ensure `put-bucket-cors` still runs unconditionally (it is idempotent — safe to re-apply)

## 2. Update docker-compose.yml

- [x] 2.1 Add `PERSISTENCE: "1"` to the LocalStack service `environment` block
- [x] 2.2 Add a `volumes:` entry to the LocalStack service mounting `localstack-data:/var/lib/localstack`
- [x] 2.3 Add a `volumes:` entry to the Redis service mounting `redis-data:/data`
- [x] 2.4 Add a top-level `volumes:` block declaring both `localstack-data:` and `redis-data:` as named volumes

## 3. Verify

- [x] 3.1 Run `docker compose down -v && docker compose up` — confirm init script completes without errors
- [x] 3.2 Upload a video, then run `docker compose stop && docker compose start` — confirm the video record is still present in the UI
- [x] 3.3 Run `docker compose down -v` — confirm volumes are removed cleanly (reset works)
