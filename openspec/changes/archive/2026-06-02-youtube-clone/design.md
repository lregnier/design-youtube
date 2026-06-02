## Context

Greenfield portfolio project implementing the YouTube system design pattern from hellointerview.com. The stack is Go backend, React/TypeScript frontend, and Terraform-managed AWS infrastructure. The project targets a single developer running on a personal AWS account, so operational simplicity is weighted heavily alongside correctness of the distributed system patterns.

## Goals / Non-Goals

**Goals:**
- Demonstrate resumable multipart upload via presigned S3 URLs (100MB max, ~10–20 parts)
- Demonstrate async video processing pipeline (ffmpeg segment splitting, multi-bitrate transcoding, HLS manifest generation, auto thumbnail extraction)
- Demonstrate adaptive bitrate streaming to the browser via HLS.js
- Keep AWS costs minimal for a personal portfolio deployment
- Keep the codebase approachable for portfolio review

**Non-Goals:**
- Per-user accounts or social features (comments, likes, subscriptions)
- Search or recommendations
- Production-grade SLA or high-availability deployment
- Real-time analytics or view counts

## Decisions

### 1. Shared secret over JWT auth for upload protection
Use a single `UPLOAD_SECRET` environment variable checked against the `X-Upload-Secret` header on upload-initiating endpoints. Simplifies deployment (no user table, no token lifecycle) while still protecting S3 from anonymous writes. JWT auth can be layered on later without restructuring the upload flow.

Alternatives considered: AWS Cognito (too much managed infra for a personal project), per-user JWT (right pattern but adds scope for zero portfolio benefit at this stage).

### 2. 100MB max video size
Covers 2–3 min of 1080p footage and 5–8 min of screen recordings — enough to demonstrate the full upload and streaming pipeline. Produces 10–20 multipart chunks, which clearly exercises the resumable upload pattern. Keeps CloudFront egress cost to ~$0.03 per full watch across all HLS quality levels.

### 3. DynamoDB over Cassandra for video metadata
The hellointerview design uses Cassandra (partitioned by videoId) for horizontal scale. DynamoDB provides the same partition-key access pattern natively on AWS with zero operational overhead. Same design pattern, no cluster to manage on a personal account.

### 4. HLS over DASH for adaptive bitrate streaming
HLS has near-universal browser support via HLS.js and is the dominant format for on-demand streaming. ffmpeg outputs HLS segments and `.m3u8` manifests directly. DASH offers marginally better compression but requires additional tooling with no meaningful advantage here.

### 5. SQS over Kafka for the processing queue
S3 event notifications trigger an SQS queue; processing workers poll SQS. Kafka provides stronger ordering and replay, but MSK is expensive and unnecessary — each video is an independent unit of work. SQS FIFO queues provide sufficient ordering per video.

### 6. ECS Fargate over Lambda for processing workers
Video transcoding is CPU-intensive and long-running (minutes per video). Lambda's 15-minute timeout and memory ceiling make it unsuitable. Fargate tasks run ffmpeg without time constraints and scale to zero when idle, keeping costs low.

### 7. CloudFront in front of S3 for segments and manifests
HLS segments, manifest files, and thumbnails are served via CloudFront, not directly from S3. Reduces S3 egress costs, adds geographic edge caching, and mirrors the CDN pattern from the reference design. The backend returns CloudFront URLs — the CDN layer is transparent to the client.

### 8. ElastiCache Redis for hot metadata cache
Frequently accessed video metadata (manifest URL, title, thumbnail URL) is cached in Redis with a short TTL to reduce DynamoDB read units. Implements the distributed LRU cache pattern from the reference design with a single-node Redis cluster.

### 9. OpenAPI-first backend with oapi-codegen
The backend API contract is defined in a single `backend/api/openapi.yaml` (OpenAPI 3.0). `oapi-codegen` generates Go server interfaces, request/response types, and a chi-compatible router from that spec. Handlers implement the generated interfaces — the spec is the source of truth, not the implementation. On the frontend, `openapi-typescript` generates TypeScript types from the same spec, keeping the client in sync without manual duplication.

Alternatives considered: hand-written types on both sides (error-prone, diverges over time), gRPC (heavier toolchain, no browser-native support without a gateway).

### 10. Auto-generated thumbnails via ffmpeg
The processing worker extracts a single frame from the video's midpoint using ffmpeg and uploads it to S3 as a JPEG. The CloudFront URL is stored in the video metadata record. No user input required — simplifies the upload flow and keeps the processing pipeline self-contained.

## Risks / Trade-offs

- **Shared secret is not per-user** → If the secret leaks, rotate it via env var update and redeploy. Document this limitation in the README.
- **100MB limit enforced at presigned URL generation** → A client could attempt to upload more by bypassing the backend. Mitigation: set a matching S3 bucket policy `content-length-range` condition on presigned URLs.
- **Single-region deployment** → CloudFront edge distribution mitigates viewer latency, but origin (S3, ECS, DynamoDB) is in one region. Acceptable for a portfolio.
- **Fargate spot task interruptions** → A transcoding task could be interrupted mid-job. Mitigation: SQS message visibility timeout keeps the job re-queued; the worker resumes on retry.
- **DynamoDB on-demand cold start** → Brief latency spikes after idle periods. Mitigation: provisioned capacity if it becomes noticeable during demos.
- **Mid-point thumbnail may not be representative** → For very short clips the midpoint frame is usually fine; for longer videos it may be a scene cut. Acceptable trade-off for a portfolio.
