package application

import (
	"context"

	"github.com/lregnier/design-youtube/worker/internal/domain/processing"
)

type VideoStorage interface {
	DownloadRaw(ctx context.Context, videoID, destPath string) error
	UploadSegments(ctx context.Context, videoID, segDir string) error
	UploadManifest(ctx context.Context, videoID string, content []byte) (cloudfrontURL string, err error)
	UploadThumbnail(ctx context.Context, videoID string, data []byte) (cloudfrontURL string, err error)
}

type Transcoder interface {
	Duration(ctx context.Context, inputPath string) (float64, error)
	TranscodeHLS(ctx context.Context, inputPath, outputDir, scale, bitrate string) error
	ExtractThumbnail(ctx context.Context, inputPath string, offset float64) ([]byte, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, event processing.DomainEvent) error
}
