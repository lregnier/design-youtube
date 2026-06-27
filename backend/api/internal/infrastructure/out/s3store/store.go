package s3store

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/lregnier/design-youtube/api/internal/application"
)

var _ application.ObjectStore = (*store)(nil)

type store struct {
	client      *awss3.Client
	bucket      string
	transformer PresignedURLTransformer
}

func NewStore(client *awss3.Client, bucket string, transformer PresignedURLTransformer) application.ObjectStore {
	return &store{client: client, bucket: bucket, transformer: transformer}
}

func (s *store) CreateMultipartUpload(ctx context.Context, key string) (*application.MultipartUpload, error) {
	out, err := s.client.CreateMultipartUpload(ctx, &awss3.CreateMultipartUploadInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("create multipart upload: %w", err)
	}
	return &application.MultipartUpload{UploadID: *out.UploadId, Key: key}, nil
}

func (s *store) PresignUploadPart(ctx context.Context, key, uploadID string, partNumber int) (*application.PresignedURL, error) {
	presigner := awss3.NewPresignClient(s.client)
	pn := int32(partNumber)
	out, err := presigner.PresignUploadPart(ctx, &awss3.UploadPartInput{
		Bucket:     &s.bucket,
		Key:        &key,
		UploadId:   &uploadID,
		PartNumber: &pn,
	}, awss3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return nil, fmt.Errorf("presign part %d: %w", partNumber, err)
	}
	presignedURL := s.transformer.Transform(out.URL)
	return &application.PresignedURL{URL: presignedURL, PartNumber: partNumber}, nil
}

func (s *store) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []application.CompletedPart) error {
	completed := make([]s3types.CompletedPart, len(parts))
	for i, p := range parts {
		pn := int32(p.PartNumber)
		completed[i] = s3types.CompletedPart{
			PartNumber: &pn,
			ETag:       aws.String(p.ETag),
		}
	}
	_, err := s.client.CompleteMultipartUpload(ctx, &awss3.CompleteMultipartUploadInput{
		Bucket:   &s.bucket,
		Key:      &key,
		UploadId: &uploadID,
		MultipartUpload: &s3types.CompletedMultipartUpload{
			Parts: completed,
		},
	})
	return err
}
