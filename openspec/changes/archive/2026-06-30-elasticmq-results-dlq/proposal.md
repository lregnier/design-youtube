## Why

`docker/elasticmq.conf` is missing `video-processing-results-dlq.fifo`, so the local dev environment doesn't mirror the AWS infrastructure where this DLQ was just added. The queue is silently absent — ElasticMQ won't enforce the redrive policy without it.

## What Changes

- Add `video-processing-results-dlq.fifo` queue to `docker/elasticmq.conf`
- Add `deadLettersQueue` redrive config to `video-processing-results.fifo` pointing to the new DLQ with `maxReceiveCount=3`

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None.

## Impact

- `docker/elasticmq.conf` only — no application or Terraform changes
