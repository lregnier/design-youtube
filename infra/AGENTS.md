# Infra

Terraform configuration for AWS. All resources are in a single root module.

## Files

| File              | Contents                                      |
|-------------------|-----------------------------------------------|
| `main.tf`         | Provider, backend (S3 state), default tags    |
| `variables.tf`    | `aws_region`, `environment`, `project`        |
| `outputs.tf`      | `cloudfront_domain`, `alb_dns_name`           |
| `vpc.tf`          | VPC, subnets, IGW, NAT, security groups       |
| `s3.tf`           | Video bucket, CORS, S3 event notification     |
| `cloudfront.tf`   | CloudFront distribution (OAC), S3 bucket policy |
| `dynamodb.tf`     | `videos` table + GSI (`status-uploadedAt-index`) |
| `sqs.tf`          | FIFO queue `video-processing.fifo`, SQS policy |
| `elasticache.tf`  | Redis single-node cluster (cache.t3.micro)    |
| `ecs.tf`          | ECS cluster, task definitions, ALB, ECS service, SSM parameter |
| `iam.tf`          | Execution role, API task role, worker task role |

## Commands

```bash
# Init (provide backend config)
terraform init -backend-config=backend.tfbackend

# Plan
terraform plan

# Apply
terraform apply
```

## Backend config (backend.tfbackend — gitignored)

```hcl
bucket = "your-tf-state-bucket"
key    = "design-youtube/terraform.tfstate"
region = "us-east-1"
```

## Key conventions

- All resources are tagged via `provider.default_tags` — do not add duplicate `tags` blocks unless adding resource-specific tags.
- The `UPLOAD_SECRET` is stored in SSM Parameter Store at `/${var.project}/upload-secret` as a `SecureString`. Update the value via AWS console or CLI after first apply — Terraform ignores value changes (`lifecycle.ignore_changes`).
- CloudFront serves only `segments/`, `manifests/`, and `thumbnails/` prefixes. The `raw/` prefix is never exposed via CloudFront.
- Worker task definition has higher CPU/memory (1024/2048) than the API (256/512) — ffmpeg is CPU-intensive.
- ECR repositories for `design-youtube-api`, `design-youtube-worker`, and `design-youtube-frontend` must exist before the first ECS deploy (create manually or add to Terraform).

## Security groups

| SG        | Inbound from        | Purpose                  |
|-----------|---------------------|--------------------------|
| `alb`     | 0.0.0.0/0 (80/443)  | Public ALB               |
| `api`     | ALB SG (8080)       | ECS API tasks            |
| `worker`  | none                | ECS worker tasks (egress only) |
| `redis`   | API SG (6379)       | ElastiCache               |
