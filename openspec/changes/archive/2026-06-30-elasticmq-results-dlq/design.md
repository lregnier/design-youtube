## Context

ElasticMQ queues are defined statically in `docker/elasticmq.conf`, not derived from Terraform. Every time a queue is added to AWS infrastructure it must be mirrored here manually. The `sqs-results-dlq` change added `video-processing-results-dlq.fifo` to Terraform but did not update the local config.

## Goals / Non-Goals

**Goals:**
- Local ElasticMQ mirrors AWS SQS queue topology exactly
- `video-processing-results.fifo` enforces `maxReceiveCount=3` locally, same as in AWS

**Non-Goals:**
- Changing any application code or Terraform
- Adding monitoring or alerting for the local DLQ

## Decisions

**Mirror AWS config exactly** — same `maxReceiveCount=3`, same naming convention (`video-processing-results-dlq.fifo`). No reason to diverge locally.
