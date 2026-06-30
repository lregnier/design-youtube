## Context

The local stack uses `softwaremill/elasticmq-native` on port 9324. The ElasticMQ project ships a companion UI image (`softwaremill/elasticmq-ui`) that connects to the ElasticMQ REST endpoint and displays queues, message counts, and message contents. Its default container port is 3000, which conflicts with the `frontend` service (`3000:80`).

## Goals / Non-Goals

**Goals:**
- Add the ElasticMQ UI as a dev-only container accessible at `http://localhost:9325`

**Non-Goals:**
- Changing the ElasticMQ configuration or queue setup
- Exposing the UI in any non-local environment

## Decisions

**Host port 9325, container port 3000** — 9325 is the natural neighbour to ElasticMQ's 9324 and is unused in the current stack. The container-internal port stays at 3000 (image default) so no custom config is needed.

**`depends_on: elasticmq: condition: service_healthy`** — mirrors the pattern used by `api` and `worker`; prevents the UI from starting before ElasticMQ is ready.

**No `ports` exposure for ElasticMQ internal UI (9325 is the only addition)** — the UI only needs to reach ElasticMQ over the internal Docker network (`http://elasticmq:9324`).

## Risks / Trade-offs

- [Risk] Image name or default env vars for the ElasticMQ UI may differ from assumptions → Mitigation: verify image name `softwaremill/elasticmq-ui` and required env var (likely `ELASTICMQ_SERVER_URL=http://elasticmq:9324`) against the softwaremill/elasticmq README before applying
