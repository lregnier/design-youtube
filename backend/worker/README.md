# Worker

Async video processing service for the design-youtube platform. Polls the `video-processing` SQS FIFO queue, downloads raw uploads from S3, transcodes them to HLS at three resolutions using FFmpeg, uploads the segments and master manifest back to S3, and publishes the result to `video-processing-results` for the API to consume.

Built with Go using hexagonal (ports and adapters) architecture.

## Architecture

```mermaid
graph LR
    subgraph Inbound["Inbound adapter"]
        SQSIn["SQS consumer\n(video-processing.fifo)"]
    end

    subgraph App["Application"]
        PV["ProcessVideo"]
    end

    subgraph Outbound["Outbound adapters"]
        S3["S3\n(VideoStorage)"]
        FFmpeg["FFmpeg\n(Transcoder)"]
        SQSOut["SQS\n(ResultPublisher)"]
    end

    SQSIn --> PV
    PV --> S3
    PV --> FFmpeg
    PV --> SQSOut
```

## Processing Pipeline

```mermaid
sequenceDiagram
    participant SQS as SQS (video-processing)
    participant W as Worker
    participant S3 as Object Store
    participant FF as FFmpeg
    participant R as SQS (results)

    SQS->>W: ProcessingJob {videoId, s3Key}
    W->>S3: DownloadRaw → /tmp/original
    W->>FF: ffprobe — get duration

    loop 1080p · 720p · 360p
        W->>FF: TranscodeHLS → /tmp/segments/{quality}/
    end

    W->>S3: UploadSegments
    W->>S3: UploadManifest (HLS master playlist)
    W->>FF: ExtractThumbnail (frame at duration/2)
    W->>S3: UploadThumbnail
    W->>R: PublishProcessed {videoId, manifestUrl, thumbnailUrl}

    Note over W,R: On any failure, PublishFailed is sent instead
```

## Output Qualities

| Quality | Resolution | Bitrate |
|---------|-----------|---------|
| 1080p | 1920×1080 | 4000 kbps |
| 720p | 1280×720 | 2500 kbps |
| 360p | 640×360 | 800 kbps |

## Configuration

| Variable | Description |
|----------|-------------|
| `AWS_REGION` | AWS region |
| `S3_BUCKET` | S3 bucket for video storage |
| `CLOUDFRONT_DOMAIN` | CloudFront domain used in published asset URLs |
| `SQS_QUEUE_URL` | SQS URL to poll for processing jobs |
| `RESULTS_QUEUE_URL` | SQS URL to publish results to |
| `S3_ENDPOINT_URL` | Override S3 endpoint (e.g. `http://minio:9000` for local dev); unset in production |
| `S3_PUBLIC_URL` | Publish `manifestUrl`/`thumbnailUrl` as `{S3_PUBLIC_URL}/{bucket}/{key}` (e.g. `http://localhost:9000`); unset in production to use `CLOUDFRONT_DOMAIN` |
| `SQS_ENDPOINT_URL` | Override SQS endpoint (e.g. `http://elasticmq:9324`); unset in production |

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

Regenerate mocks after changing port interfaces:

```bash
mockery
```
