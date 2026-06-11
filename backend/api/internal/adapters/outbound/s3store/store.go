package s3store

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/lregnier/design-youtube/api/internal/ports"
)

var _ ports.ObjectStore = (*Store)(nil)

type Store struct {
	client              *awss3.Client
	bucket              string
	s3PublicEndpointURL string
}

func NewStore(client *awss3.Client, bucket, s3PublicEndpointURL string) *Store {
	return &Store{client: client, bucket: bucket, s3PublicEndpointURL: s3PublicEndpointURL}
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
	presignedURL := out.URL
	if s.s3PublicEndpointURL != "" {
		presignedURL = rewriteHost(presignedURL, s.s3PublicEndpointURL)
	}
	return &ports.PresignedURL{URL: presignedURL, PartNumber: partNumber}, nil
}

func rewriteHost(presignedURL, publicEndpoint string) string {
	pub, err := url.Parse(publicEndpoint)
	if err != nil {
		return presignedURL
	}
	parsed, err := url.Parse(presignedURL)
	if err != nil {
		return presignedURL
	}
	parsed.Scheme = pub.Scheme
	parsed.Host = pub.Host
	return parsed.String()
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
