## Why

The project has no README files, making it hard to onboard, navigate, or understand how the system fits together at a glance. Adding READMEs with architecture diagrams and sequence diagrams turns the repo into self-documenting reference material.

## What Changes

- Add a root `README.md` with an infrastructure-agnostic system overview and architecture diagram
- Add `backend/api/README.md` with hexagonal architecture diagram, API endpoint table, and sequence diagrams for the upload and result-consumption flows
- Add `backend/worker/README.md` with hexagonal architecture diagram and a sequence diagram for the video processing pipeline
- Add `frontend/web/README.md` replacing the Vite scaffold placeholder with project-specific content
- Add `infra/aws/README.md` with an AWS-specific architecture diagram (VPC, subnets, ECS, S3, SQS, DynamoDB, ElastiCache, CloudFront, SSM) and deployment instructions

## Capabilities

### New Capabilities

- `project-documentation`: README files with Mermaid diagrams covering system architecture (infra-agnostic), per-subproject architecture (hexagonal adapters), backend sequence diagrams (upload flow, processing pipeline, result consumption), and AWS infrastructure topology

### Modified Capabilities

## Impact

- Five new Markdown files added to the repository; no code changes
- `frontend/web/README.md` is replaced (existing file is the Vite scaffold template, not project-specific)
