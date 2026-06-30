## 1. Update elasticmq.conf

- [x] 1.1 Add `video-processing-results-dlq.fifo` queue to `docker/elasticmq.conf`
- [x] 1.2 Add `deadLettersQueue` block to `video-processing-results.fifo` with `name = "video-processing-results-dlq.fifo"` and `maxReceiveCount = 3`

## 2. Verify

- [ ] 2.1 Restart ElasticMQ (`docker compose restart elasticmq`) and confirm `video-processing-results-dlq.fifo` appears in `aws sqs list-queues`
