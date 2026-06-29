## ADDED Requirements

### Requirement: Root README
The root `README.md` SHALL provide an infrastructure-agnostic overview of the system.

It SHALL include:
- A one-paragraph description of what the platform does
- A Mermaid `graph` diagram showing logical components (client, API, worker, object store, queue, cache, CDN) and the data flow between them — without naming any cloud provider
- A subprojects table linking to each subproject README
- A local development quick-start section (Docker Compose command, service URLs, upload secret)
- A CI/CD section summarising the GitHub Actions pipeline

#### Scenario: Reader understands system at a glance
- **WHEN** a developer opens the repository root on GitHub
- **THEN** they can see a system diagram, understand the high-level data flow, and navigate to any subproject README within one click

#### Scenario: Diagram is infra-agnostic
- **WHEN** the root diagram is inspected
- **THEN** it contains no AWS-specific service names (no S3, SQS, DynamoDB, ECS, CloudFront, ElastiCache)

---

### Requirement: API README
`backend/api/README.md` SHALL document the REST API service.

It SHALL include:
- A Mermaid `graph` diagram showing the hexagonal architecture: inbound adapters (HTTP handler, SQS consumer), application use cases (upload, catalog, processing), and outbound adapters (DynamoDB, S3, Redis, SQS)
- An API endpoint table grouped by tag (catalog, upload) with method, path, and description columns
- A Mermaid `sequenceDiagram` for the multipart upload flow: client → InitUpload → presigned PUT to S3 → ConfirmChunk loop → CompleteUpload
- A Mermaid `stateDiagram-v2` for the video status lifecycle (uploading → processing → ready/failed)
- An environment variables table
- Local dev commands (docker compose, go test, go generate, mockery)

#### Scenario: Upload flow is clear
- **WHEN** a developer reads the API README
- **THEN** the sequence diagram shows every HTTP call and the role of the presigned URL in the upload flow

#### Scenario: Get/stream flow is clear
- **WHEN** a developer reads the API README
- **THEN** a sequence diagram shows the full get/stream flow: client calls GET /videos/{id}, API reads from cache or DynamoDB and returns the manifest URL, client fetches the HLS manifest and segments from the CDN

#### Scenario: Architecture layers are visible
- **WHEN** a developer reads the API README
- **THEN** the hexagonal diagram clearly shows the separation between inbound adapters, application use cases, and outbound adapters

---

### Requirement: Worker README
`backend/worker/README.md` SHALL document the video processing worker.

It SHALL include:
- A Mermaid `graph` diagram showing the hexagonal architecture: inbound adapter (SQS consumer), application use case (ProcessVideo), and outbound adapters (S3 VideoStorage, FFmpeg Transcoder, SQS ResultPublisher)
- A Mermaid `sequenceDiagram` for the full processing pipeline: SQS job → download raw → ffprobe duration → transcode loop (1080p, 720p, 360p) → upload segments → upload manifest → extract thumbnail → publish result
- An output quality table (quality, resolution, bitrate)
- An environment variables table
- Local dev commands (docker compose, go test, mockery)

#### Scenario: Processing pipeline is clear
- **WHEN** a developer reads the worker README
- **THEN** the sequence diagram shows every step from receiving the SQS message to publishing the result, including the transcoding loop across all three qualities

---

### Requirement: Web README
`frontend/web/README.md` SHALL replace the Vite scaffold boilerplate with project-specific content.

It SHALL include:
- Brief description of the web client's role in the platform
- Development commands (npm install, npm run dev, npm run build)
- Docker Compose usage note
- Tech stack list (React, TypeScript, Vite, nginx)

#### Scenario: Placeholder is gone
- **WHEN** `frontend/web/README.md` is read
- **THEN** it contains no Vite scaffold boilerplate text (no "React + TypeScript + Vite", no ESLint configuration instructions)

---

### Requirement: AWS Infrastructure README
`infra/aws/README.md` SHALL document the Terraform-managed AWS infrastructure.

It SHALL include:
- A Mermaid `graph` diagram showing: VPC with public subnets (ALB) and private subnets (ECS API, ECS Worker, ElastiCache Redis), S3, DynamoDB, two SQS FIFO queues, CloudFront (with OAC → S3), and SSM Parameter Store
- An AWS resources table (resource type, name/purpose)
- Deployment instructions: `terraform init`, `terraform plan`, `terraform apply`, and the SSM post-deploy step to set the upload secret
- State backend configuration snippet
- Terraform variables table

#### Scenario: AWS topology is visible
- **WHEN** a developer reads the infra README
- **THEN** the diagram shows which resources are in the VPC, which are public vs private, and how traffic flows from the internet to the ECS services and CloudFront to S3

#### Scenario: Deployment steps are complete
- **WHEN** a developer follows the deployment instructions
- **THEN** they can provision the full infrastructure and set the upload secret without consulting any other document

---

### Requirement: README diagrams reflect current architecture
The README files for `backend/api` and `backend/worker` SHALL use terminology that matches the current codebase: port interface names (`EventPublisher`), adapter naming (`sqssubscriber`), and service method names (`HandleVideoProcessingSucceeded`, `HandleVideoProcessingFailed`).

#### Scenario: Worker README uses EventPublisher
- **WHEN** `backend/worker/README.md` is read
- **THEN** the outbound SQS node is labeled `EventPublisher` and the sequence diagram shows a single `Publish(DomainEvent)` call

#### Scenario: API README uses current service method names
- **WHEN** `backend/api/README.md` is read
- **THEN** the processing application box lists `HandleVideoProcessingSucceeded` and `HandleVideoProcessingFailed`, not `ApplyResult`
