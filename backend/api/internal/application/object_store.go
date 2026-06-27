package application

import "context"

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
