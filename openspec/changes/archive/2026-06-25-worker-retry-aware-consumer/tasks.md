## 1. Strip PublishFailed from the application layer

- [x] 1.1 In `process_video.go`, replace each `return uc.publisher.PublishFailed(...)` call with a plain `return fmt.Errorf(...)` — download, ffprobe, and transcode failure sites
- [x] 1.2 In `process_video_test.go`, remove `publisher.EXPECT().PublishFailed(...)` expectations from domain-failure tests and change assertions from `assert.NoError` to `assert.Error` (tests now expect an error returned, no publisher interaction)

## 2. Make the consumer retry-aware

- [x] 2.1 In `consumer.go`, add `publisher ports.ResultPublisher` field to `Consumer` and update `NewConsumer` to accept and store it
- [x] 2.2 In `poll`, add `MessageSystemAttributeNames: []string{"ApproximateReceiveCount"}` (or `AttributeNames`) to `ReceiveMessageInput` so the attribute is populated on each message
- [x] 2.3 In `poll`, after `processVideo.Execute` returns an error, parse `msg.Attributes["ApproximateReceiveCount"]` and compare against `maxReceiveCount` (const 3): if at max call `c.publisher.PublishFailed` then delete the message; otherwise log and continue (no delete)

## 3. Wire publisher into consumer

- [x] 3.1 In `cmd/worker/main.go`, pass `publisher` as an argument to `sqsjobs.NewConsumer`

## 4. Local environment

- [x] 4.1 In `docker/elasticmq.conf`, add `video-processing-dlq.fifo` queue and wire it to `video-processing.fifo` via `deadLettersQueue` with `maxReceiveCount = 3`
