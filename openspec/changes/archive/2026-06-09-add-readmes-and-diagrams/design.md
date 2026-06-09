## Context

The repo has no README files. There is one auto-generated placeholder at `frontend/web/README.md` (Vite scaffold boilerplate) that needs replacing. All other READMEs are missing. The system uses a multi-bounded-context backend (api + worker), a React frontend, and Terraform-managed AWS infrastructure — each with distinct architecture and configuration worth documenting independently.

## Goals / Non-Goals

**Goals:**
- One root README with an infrastructure-agnostic system diagram (shows logical components and data flow, not cloud provider specifics)
- Per-subproject READMEs at `backend/api/`, `backend/worker/`, `frontend/web/`, and `infra/aws/`
- Mermaid diagrams throughout (GitHub renders them natively, no external tool needed)
- Sequence diagrams for the two key backend flows: upload orchestration (client → API → S3) and video processing pipeline (SQS → worker → S3 → results)
- AWS README includes provider-specific topology: VPC, public/private subnets, ECS Fargate, S3, SQS FIFO queues, DynamoDB, ElastiCache Redis, CloudFront, SSM Parameter Store
- Hexagonal architecture diagrams for api and worker (inbound adapters → application → outbound adapters)
- Environment variable tables and quick-start commands for each subproject

**Non-Goals:**
- No code changes — documentation only
- No external diagram hosting (Mermaid in Markdown is sufficient)
- No API reference docs (the OpenAPI spec at `backend/api/api/openapi.yaml` is the authoritative source)

## Decisions

### Mermaid over static images
Static images require a separate authoring and export step and go stale. Mermaid diagrams live in the Markdown file, are version-controlled, and render on GitHub without any tooling. Alternative (draw.io, Lucidchart) considered and rejected — too much friction to maintain.

### Root diagram is infra-agnostic
The root README describes logical system topology (client, API, worker, object store, queue, cache, CDN) without naming AWS. This keeps it readable regardless of where the system runs and avoids duplicating the AWS-specific detail that belongs in `infra/aws/README.md`.

### One spec, one capability (`project-documentation`)
All five READMEs serve the same concern: discoverability and onboarding. Splitting into five capabilities would create unnecessary overhead with no benefit — the requirements are simple and the implementation is mechanical.

### Replace `frontend/web/README.md` entirely
The existing file is Vite scaffold boilerplate with no project-specific content. Replacing it is not a breaking change.

## Risks / Trade-offs

- **Mermaid syntax errors** → Only render on push to GitHub; local preview requires a Mermaid-aware editor. Mitigation: keep diagrams simple and test with GitHub's preview before merging.
- **Diagram drift** → Architecture evolves and diagrams become stale. Mitigation: diagrams are co-located with the code they describe; reviewers can catch drift in PRs. No automation needed at this scale.
