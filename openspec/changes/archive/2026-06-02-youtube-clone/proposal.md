## Why

Build a YouTube-like video platform from scratch as a learning and portfolio project, implementing the canonical system design patterns (multipart upload, adaptive bitrate streaming, async transcoding pipeline) documented at hellointerview.com. This demonstrates end-to-end distributed systems skills across a realistic production-grade stack.

## What Changes

- New Go backend exposing REST APIs for video upload, streaming metadata, and catalog browsing
- Upload endpoint protected by a shared secret (API key) — anyone can watch, only holders of the secret can upload
- New React + TypeScript frontend with a video player supporting adaptive bitrate streaming
- New async video processing pipeline (segment splitting, multi-bitrate transcoding, manifest generation)
- New AWS infrastructure provisioned with Terraform (S3, ECS/Fargate, ElastiCache, CloudFront CDN, DynamoDB, SQS)
- New project directory structure: `backend/`, `frontend/`, `infra/`

## Capabilities

### New Capabilities

- `upload-secret`: Shared-secret middleware protecting the upload API — validated via `X-Upload-Secret` request header; secret stored as an environment variable
- `video-upload`: Resumable multipart upload via presigned S3 URLs with per-chunk status tracking
- `video-processing`: Async DAG pipeline — segment splitting, parallel multi-bitrate transcoding (ffmpeg), HLS manifest generation
- `video-streaming`: Serve video metadata and manifest URLs; client-side adaptive bitrate playback via HLS.js
- `video-catalog`: List and retrieve video metadata for browsing the homepage and individual video pages

### Modified Capabilities

## Impact

- New top-level directories: `backend/` (Go), `frontend/` (React/TS), `infra/` (Terraform)
- AWS services: S3 (video storage), ECS Fargate (backend + processing workers), DynamoDB (video metadata), ElastiCache Redis (hot metadata cache), CloudFront (CDN for segments + manifests), ALB (load balancer), SQS (processing queue)
- External dependency: ffmpeg inside worker containers for transcoding
- No existing code is modified — this is a greenfield project
