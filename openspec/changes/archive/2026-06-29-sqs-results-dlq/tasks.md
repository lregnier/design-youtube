## 1. Terraform

- [x] 1.1 Add `aws_sqs_queue.video_processing_results_dlq` resource to `infra/aws/sqs.tf` (`video-processing-results-dlq.fifo`, FIFO, content-based deduplication)
- [x] 1.2 Add `redrive_policy` to `aws_sqs_queue.video_processing_results` pointing to `video_processing_results_dlq.arn` with `maxReceiveCount=3`

## 2. Verify

- [x] 2.1 Run `terraform validate` from `infra/aws/` and confirm clean
- [x] 2.2 Run `terraform plan` from `infra/aws/` and confirm only the two expected resources appear (`aws_sqs_queue.video_processing_results_dlq` added, `aws_sqs_queue.video_processing_results` modified)
