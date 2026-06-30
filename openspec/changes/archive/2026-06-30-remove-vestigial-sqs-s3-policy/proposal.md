## Why

`infra/aws/sqs.tf` contains an `aws_sqs_queue_policy` that grants `s3.amazonaws.com` permission to `SendMessage` on `video-processing.fifo`. S3 event notifications cannot publish directly to FIFO queues — a Lambda intermediary is required. We've decided to keep the current API-dispatches approach, so this policy is dead code that adds confusion.

## What Changes

- Remove `aws_sqs_queue_policy.video_processing` resource from `infra/aws/sqs.tf`
- Remove `data.aws_iam_policy_document.sqs_s3_publish` resource from `infra/aws/sqs.tf`

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None.

## Impact

- `infra/aws/sqs.tf` only — no application code changes
- Terraform will destroy the queue policy on next `apply` (the queue itself is unaffected)
