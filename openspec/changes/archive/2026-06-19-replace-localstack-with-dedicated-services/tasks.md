## 0. Pre-flight

- [x] 0.1 Archive `persist-localstack-data` change (it is superseded by this change)

## 1. Rename adapter types

- [x] 1.1 In `backend/api/internal/adapters/outbound/s3store/url_transformer.go`: rename `LocalStackTransformer` → `EndpointTransformer` and `NewLocalStackTransformer` → `NewEndpointTransformer`
- [x] 1.2 In `backend/worker/internal/adapters/outbound/s3storage/url_builder.go`: rename `LocalStackURLBuilder` → `EndpointURLBuilder` and `NewLocalStackURLBuilder` → `NewEndpointURLBuilder`

## 2. Update API config

- [x] 2.1 In `backend/api/internal/config/config.go`: replace `LocalStackEnabled bool`, `LocalStackEndpoint string`, `S3UsePathStyle bool` fields with `S3Endpoint string`, `S3PublicURL string`, `DynamoDBEndpoint string`, `SQSEndpoint string`
- [x] 2.2 Remove `strconv` import and the `LOCALSTACK_ENABLED` / `S3_USE_PATH_STYLE` parsing
- [x] 2.3 Read `S3_ENDPOINT_URL`, `S3_PUBLIC_URL`, `DYNAMODB_ENDPOINT_URL`, `SQS_ENDPOINT_URL` from env
- [x] 2.4 Remove the `if c.LocalStackEnabled && c.LocalStackEndpoint == "" { ... }` validation block

## 3. Update worker config

- [x] 3.1 In `backend/worker/internal/config/config.go`: replace `LocalStackEnabled bool`, `LocalStackEndpoint string`, `S3UsePathStyle bool` fields with `S3Endpoint string`, `S3PublicURL string`, `SQSEndpoint string`
- [x] 3.2 Remove `strconv` import and the `LOCALSTACK_ENABLED` / `S3_USE_PATH_STYLE` parsing
- [x] 3.3 Read `S3_ENDPOINT_URL`, `S3_PUBLIC_URL`, `SQS_ENDPOINT_URL` from env
- [x] 3.4 Remove the `if c.LocalStackEnabled && c.LocalStackEndpoint == "" { ... }` validation block

## 4. Update API main.go

- [x] 4.1 In `backend/api/cmd/api/main.go`: replace S3 path-style block with `BaseEndpoint` option when `cfg.S3Endpoint != ""`
- [x] 4.2 Add DynamoDB endpoint option: inject `o.BaseEndpoint = aws.String(cfg.DynamoDBEndpoint)` when `cfg.DynamoDBEndpoint != ""`; pass opts to `dynamodb.NewFromConfig`
- [x] 4.3 Add SQS endpoint option: inject `o.BaseEndpoint = aws.String(cfg.SQSEndpoint)` when `cfg.SQSEndpoint != ""`; pass opts to both `sqs.NewFromConfig` calls
- [x] 4.4 Change transformer selection: `if cfg.S3PublicURL != ""` → `s3store.NewEndpointTransformer(cfg.S3PublicURL)`; else `s3store.NoOpTransformer{}`
- [x] 4.5 Add `"github.com/aws/aws-sdk-go-v2/aws"` import; remove unused imports if any

## 5. Update worker main.go

- [x] 5.1 In `backend/worker/cmd/worker/main.go`: replace S3 path-style block with `BaseEndpoint` option when `cfg.S3Endpoint != ""`; enable `UsePathStyle = true` when endpoint is set
- [x] 5.2 Add SQS endpoint option: inject `o.BaseEndpoint = aws.String(cfg.SQSEndpoint)` when `cfg.SQSEndpoint != ""`; pass opts to both `sqs.NewFromConfig` calls
- [x] 5.3 Change URL builder selection: `if cfg.S3PublicURL != ""` → `s3storage.NewEndpointURLBuilder(cfg.S3PublicURL)`; else `s3storage.NewCloudFrontURLBuilder(cfg.CloudFrontDomain)`
- [x] 5.4 Add `"github.com/aws/aws-sdk-go-v2/aws"` import; remove unused imports if any

## 6. Create ElasticMQ config

- [x] 6.1 Create `docker/elasticmq.conf` with Hocon configuration declaring `video-processing.fifo` and `video-processing-results.fifo` queues, `accountId = "000000000000"`, and `region = "us-east-1"`

## 7. Update docker-compose.yml

- [x] 7.1 Remove the `localstack` service and the `localstack-data` volume entry
- [x] 7.2 Add `minio` service: image `minio/minio`, command `server /data --console-address :9001`, ports `9000:9000` and `9001:9001`, volume `minio-data:/data`, env `MINIO_ROOT_USER: minioadmin` / `MINIO_ROOT_PASSWORD: minioadmin`, healthcheck via `mc ready local`
- [x] 7.3 Add `minio-init` one-shot service: image `minio/mc`, depends on `minio` healthy, command sets alias, creates bucket with `--ignore-existing`, sets anonymous download policy and CORS allowing PUT from `*`
- [x] 7.4 Add `dynamodb-local` service: image `amazon/dynamodb-local`, port `8000:8000`, volume `dynamodb-data:/home/dynamodblocal/data`, command `"-jar DynamoDBLocal.jar -sharedDb -dbPath /home/dynamodblocal/data"`, healthcheck via `curl -s http://localhost:8000`
- [x] 7.5 Add `dynamodb-init` one-shot service: image `amazon/aws-cli`, depends on `dynamodb-local` healthy, env `AWS_ACCESS_KEY_ID: test` / `AWS_SECRET_ACCESS_KEY: test` / `AWS_DEFAULT_REGION: us-east-1`, command runs `aws dynamodb create-table` (with `--endpoint-url http://dynamodb-local:8000`) guarded with `|| true`
- [x] 7.6 Add `elasticmq` service: image `softwaremill/elasticmq-native`, port `9324:9324`, volume mount `./docker/elasticmq.conf:/opt/elasticmq.conf:ro`, healthcheck via `curl -s http://localhost:9324`
- [x] 7.7 Update `api` service: remove `AWS_ENDPOINT_URL`, `LOCALSTACK_ENABLED`, `LOCALSTACK_ENDPOINT`, `S3_USE_PATH_STYLE`; add `S3_ENDPOINT_URL: http://minio:9000`, `S3_PUBLIC_URL: http://localhost:9000`, `DYNAMODB_ENDPOINT_URL: http://dynamodb-local:8000`, `SQS_ENDPOINT_URL: http://elasticmq:9324`; update `SQS_QUEUE_URL` and `RESULTS_QUEUE_URL` to use `elasticmq:9324`; update `depends_on` to `minio-init`, `dynamodb-init`, `elasticmq`, `redis`
- [x] 7.8 Update `worker` service: same env var swap; add `S3_ENDPOINT_URL: http://minio:9000`, `S3_PUBLIC_URL: http://localhost:9000`; update queue URLs; update `depends_on` to `minio-init`, `dynamodb-init`, `elasticmq`
- [x] 7.9 Update top-level `volumes`: replace `localstack-data:` with `minio-data:` and `dynamodb-data:`

## 8. Remove localstack-init.sh

- [x] 8.1 Delete `docker/localstack-init.sh`

## 9. Verify

- [x] 9.1 Run `go vet ./...` in both `backend/api` and `backend/worker` — no errors
- [x] 9.2 Run `go test ./...` in both — all pass
- [x] 9.3 Run `docker compose down -v && docker compose up --build` — all services start healthy
- [x] 9.4 Upload a video through the UI — confirm S3 upload, transcoding, and playback work end-to-end
- [x] 9.5 Run `docker compose stop && docker compose start` — confirm uploaded video is still visible (persistence check)
- [x] 9.6 Update `README.md`, `backend/api/README.md`, and `backend/worker/README.md` — replace LocalStack references with MinIO, DynamoDB Local, and ElasticMQ; update env var tables
