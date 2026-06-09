## 1. Root README

- [x] 1.1 Write `README.md` at repo root with project description and local dev quick-start
- [x] 1.2 Add infra-agnostic Mermaid `graph` diagram showing logical components and data flow
- [x] 1.3 Add subprojects table with links to each subproject README
- [x] 1.4 Add CI/CD section summarising the GitHub Actions pipeline

## 2. API README

- [x] 2.1 Write `backend/api/README.md` with service description and environment variable table
- [x] 2.2 Add Mermaid `graph` diagram for hexagonal architecture (inbound adapters → use cases → outbound adapters)
- [x] 2.3 Add API endpoint table grouped by catalog and upload tags
- [x] 2.4 Add Mermaid `sequenceDiagram` for the multipart upload flow
- [x] 2.5 Add Mermaid `stateDiagram-v2` for the video status lifecycle
- [x] 2.6 Add local dev commands section (docker compose, go test, go generate, mockery)

## 3. Worker README

- [x] 3.1 Write `backend/worker/README.md` with service description and environment variable table
- [x] 3.2 Add Mermaid `graph` diagram for hexagonal architecture
- [x] 3.3 Add Mermaid `sequenceDiagram` for the full processing pipeline
- [x] 3.4 Add output quality table (quality, resolution, bitrate)
- [x] 3.5 Add local dev commands section (docker compose, go test, mockery)

## 4. Web README

- [x] 4.1 Replace `frontend/web/README.md` with project-specific content (description, dev commands, tech stack)

## 2b. API README — additional diagram

- [x] 2.7 Add Mermaid `sequenceDiagram` for the get/stream video flow (GET /videos/{id} → cache check → DynamoDB → manifest URL → CDN HLS playback)

## 5. AWS Infrastructure README

- [x] 5.1 Write `infra/aws/README.md` with deployment instructions and variable table
- [x] 5.2 Add Mermaid `graph` diagram showing AWS topology (VPC, subnets, ECS, S3, SQS, DynamoDB, Redis, CloudFront, SSM)
- [x] 5.3 Add AWS resources table
- [x] 5.4 Add state backend configuration snippet and post-deploy SSM step
