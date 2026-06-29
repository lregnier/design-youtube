## 1. Update worker/README.md

- [x] 1.1 Architecture diagram: rename `SQSOut["SQS\n(ResultPublisher)"]` → `SQSOut["SQS\n(EventPublisher)"]`
- [x] 1.2 Architecture diagram: rename inbound node `SQSIn["SQS consumer\n(video-processing.fifo)"]` → `SQSIn["SQS subscriber\n(video-processing.fifo)"]`
- [x] 1.3 Sequence diagram: replace `W->>R: PublishProcessed {videoId, manifestUrl, thumbnailUrl}` and failure note with a single `W->>R: Publish(VideoProcessingSucceededEvent)` call and a note that `Publish(VideoProcessingFailedEvent)` is sent on failure

## 2. Update api/README.md

- [x] 2.1 Architecture diagram: rename outbound SQS node `SQSOut["SQS\n(Queue, processing jobs)"]` → `SQSOut["SQS\n(EventPublisher)"]`
- [x] 2.2 Architecture diagram: rename inbound SQS node `SQSIn["SQS consumer\n(results queue)"]` → `SQSIn["SQS subscriber\n(results queue)"]`
- [x] 2.3 Application box: rename `Processing["processing\n· ApplyResult"]` → `Processing["processing\n· HandleVideoProcessingSucceeded\n· HandleVideoProcessingFailed"]`

## 3. Verify

- [x] 3.1 Review both rendered diagrams for accuracy
