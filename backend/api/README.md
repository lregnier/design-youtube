# API

REST API for the design-youtube platform. Handles video uploads via S3 presigned multipart URLs, tracks video metadata in DynamoDB, caches presigned URLs in Redis, and consumes processing results from the worker over SQS.

Built with Go using hexagonal (ports and adapters) architecture. HTTP routes are generated from the OpenAPI spec with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen).

## Architecture

```mermaid
graph LR
    subgraph Inbound["Inbound adapters"]
        HTTP["HTTP handler\n(chi router)"]
        SQSIn["SQS consumer\n(results queue)"]
    end

    subgraph App["Application"]
        Upload["upload\n· InitUpload\n· ConfirmChunk\n· CompleteUpload"]
        Catalog["catalog\n· GetVideo\n· ListVideos"]
        Processing["processing\n· ApplyResult"]
    end

    subgraph Outbound["Outbound adapters"]
        DynDB["DynamoDB\n(VideoRepository)"]
        S3["S3\n(ObjectStore)"]
        Redis["Redis\n(Cache)"]
        SQSOut["SQS\n(Queue, processing jobs)"]
    end

    HTTP --> Upload
    HTTP --> Catalog
    SQSIn --> Processing

    Upload --> DynDB
    Upload --> S3
    Upload --> SQSOut
    Catalog --> DynDB
    Catalog --> Redis
    Processing --> DynDB
```

## API Endpoints

### Catalog

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/videos` | List all ready videos |
| `GET` | `/videos/{videoId}` | Get video metadata and streaming URL |

### Upload

Requires `X-Upload-Secret` header on all upload endpoints.

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/videos/upload/init` | Initiate or resume a multipart upload |
| `POST` | `/videos/{videoId}/upload/confirm-chunk` | Confirm a chunk and receive the next presigned URL |
| `POST` | `/videos/{videoId}/upload/complete` | Finalise the multipart upload |

Full spec: [`api/openapi.yaml`](api/openapi.yaml)

## Upload Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant API as API
    participant S3 as Object Store
    participant Q as SQS (video-processing)

    C->>API: POST /videos/upload/init
    API->>S3: CreateMultipartUpload
    API->>API: Save video record (status: uploading)
    API-->>C: {videoId, uploadId, presignedUrl}

    loop for each chunk
        C->>S3: PUT presignedUrl (chunk bytes)
        S3-->>C: ETag
        C->>API: POST /videos/{id}/upload/confirm-chunk {partNumber, eTag}
        API->>API: Mark chunk uploaded
        API-->>C: {done, nextPresignedUrl}
    end

    C->>API: POST /videos/{id}/upload/complete
    API->>S3: CompleteMultipartUpload
    API->>API: Update status → processing
    API->>Q: SendMessage {videoId, s3Key}
    Note over Q: Worker polls and picks up the job
```

## Get / Stream Video Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant API as API
    participant Cache as Cache (Redis)
    participant DB as Database (DynamoDB)
    participant CDN as CDN

    C->>API: GET /videos/{videoId}
    API->>Cache: Get video metadata
    alt cache hit
        Cache-->>API: Video record
    else cache miss
        Cache-->>API: miss
        API->>DB: GetItem by videoId
        DB-->>API: Video record
        API->>Cache: Set video metadata
    end
    API-->>C: {videoId, title, status, manifestUrl, thumbnailUrl}

    Note over C,CDN: Client drives HLS playback directly against the CDN
    C->>CDN: GET manifestUrl (HLS master playlist)
    CDN-->>C: #EXTM3U with quality variants
    loop for each segment
        C->>CDN: GET segment URL
        CDN-->>C: .ts segment data
    end
```

## Video Status Lifecycle

```mermaid
stateDiagram-v2
    [*] --> uploading: InitUpload
    uploading --> processing: CompleteUpload
    processing --> ready: ApplyResult (success)
    processing --> failed: ApplyResult (failure)
    ready --> [*]
    failed --> [*]
```

## Configuration

| Variable | Description |
|----------|-------------|
| `UPLOAD_SECRET` | Shared secret required for upload endpoints |
| `AWS_REGION` | AWS region |
| `DYNAMODB_TABLE` | DynamoDB table name |
| `S3_BUCKET` | S3 bucket for video storage |
| `CLOUDFRONT_DOMAIN` | CloudFront domain for serving assets |
| `SQS_QUEUE_URL` | SQS URL for dispatching processing jobs |
| `RESULTS_QUEUE_URL` | SQS URL for consuming processing results |
| `REDIS_ADDR` | Redis address (`host:port`) |
| `AWS_ENDPOINT_URL` | Override AWS endpoint (LocalStack in dev) |
| `S3_USE_PATH_STYLE` | Use path-style S3 addressing (`true` for LocalStack, unset/`false` in production) |
| `S3_PUBLIC_ENDPOINT_URL` | Rewrite presigned URL host to this endpoint (browser-accessible LocalStack URL) |
| `CORS_ALLOWED_ORIGINS` | Comma-separated list of allowed CORS origins for upload endpoints |

## Development

Run the full stack:

```bash
# From repo root
docker compose up --build
```

Run tests:

```bash
go test ./...
```

Regenerate OpenAPI server code after editing `api/openapi.yaml`:

```bash
go generate ./api/...
```

Regenerate mocks after changing port interfaces:

```bash
mockery
```
