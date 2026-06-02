## 1. Project Scaffold

- [x] 1.1 Create top-level directories: `backend/`, `frontend/`, `infra/`
- [x] 1.2 Initialize Go module in `backend/` (`go mod init`)
- [x] 1.3 Scaffold React + TypeScript app in `frontend/` with Vite
- [x] 1.4 Initialize Terraform root module in `infra/` with AWS provider and backend config
- [x] 1.5 Add root `.gitignore` covering Go, Node, and Terraform artifacts

## 2. OpenAPI Contract

- [x] 2.1 Create `backend/api/openapi.yaml` (OpenAPI 3.0) defining all paths: `POST /videos/upload/init`, `POST /videos/{videoId}/upload/confirm-chunk`, `POST /videos/{videoId}/upload/complete`, `GET /videos`, `GET /videos/{videoId}`
- [x] 2.2 Define all request/response schemas: `UploadInitRequest`, `UploadInitResponse`, `ConfirmChunkRequest`, `ConfirmChunkResponse`, `CompleteUploadRequest`, `VideoSummary`, `VideoDetail`
- [x] 2.3 Mark upload endpoints with a custom `x-upload-secret` security scheme in the spec
- [x] 2.4 Validate the spec with `vacuum lint` or `swagger-cli validate` before generating code

## 3. Infra — S3 and CloudFront

- [x] 3.1 Create S3 bucket for video storage with private ACL and versioning disabled
- [x] 3.2 Configure S3 bucket prefixes: `raw/`, `segments/`, `manifests/`, `thumbnails/`
- [x] 3.3 Add S3 bucket CORS policy to allow direct PUT from the browser
- [x] 3.4 Create CloudFront distribution with S3 as origin (OAC — no public S3 access)
- [x] 3.5 Configure CloudFront to serve `segments/`, `manifests/`, and `thumbnails/` prefixes
- [x] 3.6 Output CloudFront domain name as Terraform output

## 4. Infra — DynamoDB, SQS, ElastiCache

- [x] 4.1 Create DynamoDB table `videos` with partition key `videoId` (String), on-demand billing
- [x] 4.2 Add DynamoDB GSI on `status` + `uploadedAt` for catalog list queries
- [x] 4.3 Create SQS FIFO queue `video-processing.fifo` for processing jobs
- [x] 4.4 Configure S3 event notification to send `CompleteMultipartUpload` events to the SQS queue
- [x] 4.5 Create ElastiCache Redis single-node cluster (cache.t3.micro) in a private subnet
- [x] 4.6 Create VPC, subnets, and security groups for ECS and ElastiCache

## 5. Infra — ECS and IAM

- [x] 5.1 Create ECS cluster for backend API and processing worker
- [x] 5.2 Define Fargate task definition for the Go backend API service
- [x] 5.3 Define Fargate task definition for the Go processing worker (with ffmpeg layer)
- [x] 5.4 Create ECS service for the backend API behind an ALB
- [x] 5.5 Create IAM task roles with least-privilege policies (S3 read/write, DynamoDB, SQS, ElastiCache)
- [x] 5.6 Output ALB DNS name as Terraform output

## 6. Backend — Project Structure and Codegen

- [x] 6.1 Add Go dependencies: AWS SDK v2, chi, go-redis, oapi-codegen
- [x] 6.2 Add `oapi-codegen` config (`backend/api/oapi-codegen.yaml`) targeting chi and generating types + server interfaces into `backend/internal/api/`
- [x] 6.3 Add `go generate` directive in `backend/api/` to run codegen from `openapi.yaml`
- [x] 6.4 Run `go generate ./...` and verify generated server interface and types compile
- [x] 6.5 Create config loader that reads env vars (`UPLOAD_SECRET`, `AWS_REGION`, `DYNAMODB_TABLE`, `S3_BUCKET`, `CLOUDFRONT_DOMAIN`, `SQS_QUEUE_URL`, `REDIS_ADDR`) and fails fast if any are missing
- [x] 6.6 Implement upload secret middleware: constant-time compare of `X-Upload-Secret` header, return 401 on mismatch
- [x] 6.7 Wire up chi router using the generated `Handler()` registration function, applying secret middleware to upload routes

## 7. Backend — Video Upload API

- [x] 7.1 Implement generated `PostVideosUploadInit` handler: validate file size ≤ 100MB, create S3 multipart upload, write DynamoDB record with status `uploading`, return `UploadInitResponse`
- [x] 7.2 Add S3 presigned URL generation with `content-length-range` condition capped at 10MB per part
- [x] 7.3 Implement generated `PostVideosVideoIdUploadConfirmChunk` handler: mark chunk uploaded in DynamoDB, return presigned URL for next part (or none if last)
- [x] 7.4 Implement generated `PostVideosVideoIdUploadComplete` handler: call S3 CompleteMultipartUpload with all ETags, update DynamoDB status to `processing`
- [x] 7.5 Implement resume logic in init handler: if videoId exists with status `uploading`, return current chunk state and presigned URL for first missing part

## 8. Backend — Video Metadata API

- [x] 8.1 Implement generated `GetVideosVideoId` handler: check Redis cache, fall back to DynamoDB, return `VideoDetail` (or 404)
- [x] 8.2 Implement Redis cache write on DynamoDB miss with 60-second TTL
- [x] 8.3 Implement generated `GetVideos` handler: query DynamoDB GSI for status=`ready` ordered by uploadedAt desc, return array of `VideoSummary`

## 9. Processing Worker

- [x] 9.1 Create separate Go binary in `backend/cmd/worker/` that polls SQS
- [x] 9.2 Implement SQS polling loop with visibility timeout extension during processing
- [x] 9.3 Download raw video from S3 to a temp directory
- [x] 9.4 Run ffmpeg to produce HLS segments at 1080p (4000k), 720p (2500k), and 360p (800k) with 6-second segment duration
- [x] 9.5 Upload all segment `.ts` files and media playlists to S3 under `segments/{videoId}/`
- [x] 9.6 Generate and upload HLS master manifest to `manifests/{videoId}/master.m3u8`
- [x] 9.7 Run ffmpeg to extract JPEG thumbnail at video midpoint, upload to `thumbnails/{videoId}/thumb.jpg`
- [x] 9.8 Update DynamoDB record: status `ready`, CloudFront manifest URL, CloudFront thumbnail URL
- [x] 9.9 Delete SQS message on success; on ffmpeg failure set status `failed` and delete message
- [x] 9.10 Ensure worker is idempotent (safe to re-run for the same videoId)

## 10. Frontend — TypeScript Client Generation

- [x] 10.1 Add `openapi-typescript` dev dependency to `frontend/`
- [x] 10.2 Add npm script `generate:api` that runs `openapi-typescript ../backend/api/openapi.yaml -o src/api/types.gen.ts`
- [x] 10.3 Run `generate:api` and verify the generated types match the expected request/response shapes
- [x] 10.4 Create `src/api/client.ts` wrapping `fetch` with base URL config and typed request/response helpers using the generated types

## 11. Frontend — Homepage (Catalog)

- [x] 11.1 Set up React Router with routes: `/` (homepage), `/videos/:videoId` (player), `/upload` (upload form)
- [x] 11.2 Create typed `getVideos()` API function using generated `VideoSummary` type
- [x] 11.3 Build `VideoCard` component: thumbnail image, title, formatted upload date
- [x] 11.4 Build `HomePage` component: fetch catalog on mount, render grid of `VideoCard` components, show empty state when list is empty

## 12. Frontend — Video Player Page

- [x] 12.1 Add HLS.js dependency
- [x] 12.2 Create typed `getVideo(videoId)` API function using generated `VideoDetail` type
- [x] 12.3 Build `VideoPlayer` component: initialize HLS.js with manifest URL, attach to `<video>` element, show thumbnail as poster
- [x] 12.4 Build `VideoPage` component: fetch metadata on mount, render `VideoPlayer` when status is `ready`, show loading/processing state otherwise

## 13. Frontend — Upload Page

- [x] 13.1 Build `UploadForm` component: file input (accept video/*), title field, description field, secret field, submit button
- [x] 13.2 Add client-side file size validation (reject > 100MB before any network call)
- [x] 13.3 Implement multipart upload orchestration using generated `UploadInitResponse` and `ConfirmChunkResponse` types: call init, loop through chunks calling confirm-chunk, call complete
- [x] 13.4 Display per-chunk progress bar during upload
- [x] 13.5 Handle 401 response from backend with "invalid secret" error message
- [x] 13.6 Redirect to homepage on successful upload completion

## 14. Docker and CI

- [x] 14.1 Write `Dockerfile` for Go backend API (multi-stage build)
- [x] 14.2 Write `Dockerfile` for processing worker (multi-stage build, include ffmpeg)
- [x] 14.3 Write `Dockerfile` for React frontend (build stage + nginx serve stage)
- [x] 14.4 Add `docker-compose.yml` for local development (backend, worker, Redis, LocalStack for S3/SQS/DynamoDB)
- [x] 14.5 Add GitHub Actions workflow: lint + test on push, build and push Docker images to ECR on merge to main
