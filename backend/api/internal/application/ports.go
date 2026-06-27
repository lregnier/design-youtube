package application

import (
	"context"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event video.DomainEvent) error
}

type MultipartUpload struct {
	UploadID string
	Key      string
}

type PresignedURL struct {
	URL        string
	PartNumber int
}

type CompletedPart struct {
	PartNumber int
	ETag       string
}

type ObjectStore interface {
	CreateMultipartUpload(ctx context.Context, key string) (*MultipartUpload, error)
	PresignUploadPart(ctx context.Context, key, uploadID string, partNumber int) (*PresignedURL, error)
	CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []CompletedPart) error
}
