## 1. Add ElasticMQ UI service

- [x] 1.1 Add `elasticmq-ui` service to `docker-compose.yml` using image `softwaremill/elasticmq-ui`, port mapping `9325:3000`, env var `ELASTICMQ_SERVER_URL=http://elasticmq:9324`, and `depends_on: elasticmq: condition: service_healthy`

## 2. Update documentation

- [x] 2.1 Add ElasticMQ UI entry (port 9325) to the local services table in the root `README.md`
