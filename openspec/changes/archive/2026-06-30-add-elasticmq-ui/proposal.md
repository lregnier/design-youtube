## Why

Inspecting SQS queues locally currently requires the AWS CLI, which adds friction during development. The local stack already has UI companions for MinIO (port 9001) and DynamoDB; an ElasticMQ UI completes the pattern and makes queue depth and message contents visible at a glance.

## What Changes

- Add `elasticmq-ui` service to `docker-compose.yml` using the ElasticMQ UI image
- Expose the UI on host port `9325` (container default is `3000`, remapped to avoid conflict with the `frontend` service on `3000:80`)
- Make the service depend on `elasticmq` being healthy

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None.

## Impact

- `docker-compose.yml` only — no application or infrastructure code changes
- ElasticMQ UI accessible at `http://localhost:9325` during local development
