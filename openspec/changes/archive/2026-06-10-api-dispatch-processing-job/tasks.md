## 1. Extend the Queue port and adapter

- [x] 1.1 Add `SendMessage(ctx context.Context, body, messageGroupID string) error` to `ports.Queue` in `backend/api/internal/ports/outbound.go`
- [x] 1.2 Implement `SendMessage` in `backend/api/internal/adapters/outbound/sqsqueue/queue.go` — call `sqs.SendMessage` with `MessageGroupId` and `MessageDeduplicationId` both set to `messageGroupID`

## 2. Update the CompleteUpload use case

- [x] 2.1 Add a `queue ports.Queue` field and include it in `NewCompleteUpload` constructor in `backend/api/internal/application/upload/complete.go`
- [x] 2.2 After `CompleteMultipartUpload` and `repo.Save` succeed, marshal `{"videoId":"...","s3Key":"..."}` and call `queue.SendMessage` with the videoId as the message group ID
- [x] 2.3 Return an error if `SendMessage` fails

## 3. Wire up in main

- [x] 3.1 In `backend/api/cmd/api/main.go`, construct `sqsqueue.NewQueue(sqsClient, cfg.SQSQueueURL)` and pass it as the third argument to `upload.NewCompleteUpload`

## 4. Remove S3 event notification from init script

- [x] 4.1 In `scripts/localstack-init.sh`, remove the `sqs set-queue-attributes` block (queue policy for S3) and the `s3api put-bucket-notification-configuration` block — keep both `sqs create-queue` calls

## 5. Update documentation diagrams

- [x] 5.1 In root `README.md`, remove the `ObjectStore -- "upload-complete event" --> Queue1` edge from the architecture diagram (S3 no longer notifies the queue directly)
- [x] 5.2 In `backend/api/README.md`, add the processing queue as an outbound adapter in the architecture diagram (`Upload --> SQS producer`)
- [x] 5.3 In `backend/api/README.md`, update the Upload Flow sequence diagram to replace the "Object store fires event → queue → worker picks up" note with an explicit `API->>Queue: SendMessage {videoId, s3Key}` step

## 6. Remove dead S3 notification parsing from worker

- [x] 6.1 In `backend/worker/internal/adapters/inbound/sqsjobs/consumer.go`, remove the `s3Notification` type and its branch in `parseJob` — jobs always arrive as raw `{videoId, s3Key}` JSON now

## 7. Remove unused DeleteMessage from the Queue port

- [x] 7.1 Remove `DeleteMessage` from `ports.Queue` (`backend/api/internal/ports/outbound.go`) and its implementation in `backend/api/internal/adapters/outbound/sqsqueue/queue.go` — never called; the results-queue consumer acks via the raw SQS client directly. Regenerate mocks with `mockery`
