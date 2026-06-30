## Context

`infra/aws/sqs.tf` was written with an `aws_sqs_queue_policy` granting S3 `SendMessage` rights on `video-processing.fifo`. AWS S3 event notifications only support standard (non-FIFO) SQS queues; sending to a FIFO queue requires a Lambda to receive the S3 notification and re-publish with a `MessageGroupId`. Since we ship without that Lambda and the API currently dispatches the job message directly after upload, the policy is unreachable dead code.

## Goals / Non-Goals

**Goals:**
- Remove the unused `aws_sqs_queue_policy` and its `aws_iam_policy_document` data source to eliminate misleading infrastructure

**Non-Goals:**
- Changing the API dispatch flow or any application code
- Adding S3→Lambda→SQS routing now or as a follow-up

## Decisions

**Delete both resources, no replacement** — the policy cannot be triggered (S3 → FIFO is blocked by AWS). Keeping it implies intent that doesn't exist and would confuse future readers. No migration needed: destroying a queue policy leaves the queue intact with its default (deny-all from external principals) policy.

## Risks / Trade-offs

- [Risk] Terraform plan shows a destroy, which may alarm reviewers → Mitigation: commit message and PR description explain why
