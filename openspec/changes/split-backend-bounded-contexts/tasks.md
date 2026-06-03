## 1. Infra — Results Queue and IAM

- [ ] 1.1 Add `video-processing-results.fifo` SQS FIFO queue to `infra/sqs.tf` with content-based deduplication and 900s visibility timeout
- [ ] 1.2 Update `infra/iam.tf` worker task role: remove DynamoDB permissions, add `sqs:SendMessage` on the results queue
- [ ] 1.3 Update `infra/iam.tf` API task role: add `sqs:ReceiveMessage`, `sqs:DeleteMessage`, `sqs:ChangeMessageVisibility` on the results queue
- [ ] 1.4 Add `RESULTS_QUEUE_URL` to API ECS task definition environment in `infra/ecs.tf`; remove `DYNAMODB_TABLE` from worker task definition
- [ ] 1.5 Add `results_queue_url` to `infra/outputs.tf`

## 2. Restructure backend/ → backend/api/ and backend/worker/

- [ ] 2.1 Move existing `backend/` contents into `backend/api/` using `git mv`
- [ ] 2.2 Update `go.mod` in `backend/api/` module path to `github.com/lregnier/design-youtube/api`
- [ ] 2.3 Update all internal import paths in `backend/api/` to use the new module path
- [ ] 2.4 Run `go build ./...` in `backend/api/` — zero errors
- [ ] 2.5 Update `docker-compose.yml` build contexts from `./backend` to `./backend/api` and `./backend/worker`
- [ ] 2.6 Update `.github/workflows/ci.yml` working directories and Docker build paths for both services
- [ ] 2.7 Verify `openspec/` design docs reflect the new `backend/api/` and `backend/worker/` layout

## 3. API — Event Consumer

- [ ] 3.1 Add `RESULTS_QUEUE_URL` to `backend/api/internal/config/config.go` — fail fast if unset
- [ ] 3.2 Create `backend/api/internal/application/processing/apply_result.go`: `ApplyProcessingResult` use case with `VideoRepository` dep; handles `VideoProcessed` (calls `MarkReady`) and `VideoFailed` (calls `MarkFailed`); idempotent on already-terminal status
- [ ] 3.3 Create `backend/api/internal/adapters/inbound/sqsconsumer/consumer.go`: SQS long-poll loop that deserializes `eventType`, routes to `ApplyProcessingResult`, deletes message on success, leaves message in queue on error
- [ ] 3.4 Update `backend/api/cmd/api/main.go`: instantiate consumer with results queue URL and `ApplyProcessingResult`; launch as `go consumer.Start(ctx)`
- [ ] 3.5 Run `go build ./...` and `go vet ./...` in `backend/api/` — zero errors

## 4. New Worker Module — backend/worker/

- [ ] 4.1 Create `backend/worker/` directory; initialize `go mod init github.com/lregnier/design-youtube/worker`
- [ ] 4.2 Copy worker processing logic from `backend/api/cmd/worker/` as starting point into `backend/worker/cmd/worker/`
- [ ] 4.3 Remove all DynamoDB imports, config vars, and calls from the worker
- [ ] 4.4 Create `backend/worker/internal/event/result.go`: define `VideoProcessed` and `VideoFailed` structs with `EventType` discriminator field
- [ ] 4.5 Create `backend/worker/internal/queue/publisher.go`: SQS publisher that serializes and sends result events to `video-processing-results.fifo` using `videoId` as message group key
- [ ] 4.6 Replace `markFailed`/`updateReady` DynamoDB calls with `publisher.Emit(VideoFailed{...})` and `publisher.Emit(VideoProcessed{...})`
- [ ] 4.7 Add worker-specific config: `RESULTS_QUEUE_URL`, `S3_BUCKET`, `CLOUDFRONT_DOMAIN`, `SQS_QUEUE_URL`, `AWS_REGION` — remove `DYNAMODB_TABLE`, `REDIS_ADDR`, `UPLOAD_SECRET`
- [ ] 4.8 Write `backend/worker/Dockerfile` (multi-stage build, include ffmpeg)
- [ ] 4.9 Add worker module context to `openspec/changes/split-backend-bounded-contexts/design.md` if needed
- [ ] 4.10 Remove old `backend/api/cmd/worker/` directory (worker now lives in its own module)
- [ ] 4.11 Run `go build ./...` and `go vet ./...` in `backend/worker/` — zero errors

## 5. Update docker-compose for Local Dev

- [ ] 5.1 Update `docker-compose.yml` worker service: point to `./backend/worker`, add `RESULTS_QUEUE_URL`, remove `DYNAMODB_TABLE`
- [ ] 5.2 Update `docker-compose.yml` API service: add `RESULTS_QUEUE_URL`
