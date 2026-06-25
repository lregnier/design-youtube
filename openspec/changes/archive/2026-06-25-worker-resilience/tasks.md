## 1. Fix PublishFailed error handling

- [x] 1.1 In `process_video.go`, check the error returned by each `publisher.PublishFailed` call and return it (instead of swallowing it and returning `nil`)
- [x] 1.2 Update unit tests in `process_video_test.go` to assert that a `PublishFailed` error causes `Execute` to return an error

## 2. Visibility timeout heartbeat

- [x] 2.1 In `consumer.go`, add a `startHeartbeat(ctx, receiptHandle)` helper that calls `ChangeMessageVisibility` every 5 minutes with a 900 s extension
- [x] 2.2 In `poll`, spawn the heartbeat goroutine before calling `processVideo.Execute` and cancel it (via `context.WithCancel`) when `Execute` returns
- [x] 2.3 Grant `sqs:ChangeMessageVisibility` on `video-processing.fifo` to the worker ECS task role in `infra/aws/iam.tf`

## 3. Dead Letter Queue

- [x] 3.1 Add `aws_sqs_queue` resource `video_processing_dlq` (`video-processing-dlq.fifo`, FIFO, content-based dedup) in `infra/aws/sqs.tf`
- [x] 3.2 Add `redrive_policy` to `video_processing` queue referencing the DLQ ARN with `maxReceiveCount = 3`
- [x] 3.3 Grant the worker ECS task role `sqs:GetQueueAttributes` on the DLQ (required by AWS for redrive policies to function)
