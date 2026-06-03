package s3store

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/lregnier/design-youtube/backend/internal/ports"
)

var _ ports.ObjectStore = (*Store)(nil)

type Store struct {
	client *awss3.Client
	bucket string
}

func NewStore(client *awss3.Client, bucket string) *Store {
	return &Store{client: client, bucket: bucket}
}

func (s *Store) CreateMultipartUpload(ctx context.Context, key string) (*ports.MultipartUpload, error) {
	out, err := s.client.CreateMultipartUpload(ctx, &awss3.CreateMultipartUploadInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("create multipart upload: %w", err)
	}
	return &ports.MultipartUpload{UploadID: *out.UploadId, Key: key}, nil
}

func (s *Store) PresignUploadPart(ctx context.Context, key, uploadID string, partNumber int) (*ports.PresignedURL, error) {
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
	return &ports.PresignedURL{URL: out.URL, PartNumber: partNumber}, nil
}

func (s *Store) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []ports.CompletedPart) error {
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

func (s *Store) PutObject(ctx context.Context, key string, data []byte, contentType string) error {
	_, err := s.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: &contentType,
	})
	return err
}

func (s *Store) GetObject(ctx context.Context, key string) ([]byte, error) {
	out, err := s.client.GetObject(ctx, &awss3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()
	return io.ReadAll(out.Body)
}
