# design-youtube

YouTube-like video platform built as a learning and portfolio project. Implements the canonical system design patterns: resumable multipart upload, async HLS transcoding pipeline, and adaptive bitrate streaming.

## Stack

| Layer    | Tech                                      |
|----------|-------------------------------------------|
| Backend  | Go, chi, oapi-codegen, AWS SDK v2         |
| Frontend | React 19, TypeScript, Vite, HLS.js        |
| Infra    | Terraform, AWS (ECS Fargate, S3, DynamoDB, SQS, ElastiCache, CloudFront) |

## Repo layout

```
backend/    Go API server + processing worker
frontend/   React/TS SPA
infra/      Terraform (AWS)
openspec/   Specs, change history, and archived proposals
```

Each subdirectory has its own `AGENTS.md` with domain-specific context. Read the relevant one before working in that area.

## Key conventions

- The OpenAPI contract lives at `backend/api/openapi.yaml` — it is the source of truth for the API. Do not hand-edit `backend/internal/api/api.gen.go`.
- Run `go generate ./api/...` in `backend/` after changing `openapi.yaml`.
- Run `npm run generate:api` in `frontend/` after changing `openapi.yaml`.
- Upload endpoints require `X-Upload-Secret` header. Read endpoints are public.
- Max video upload size: 100MB.

## Local development

```bash
docker compose up   # starts API, worker, Redis, LocalStack
```

See each subdirectory's `AGENTS.md` for build and test commands.
