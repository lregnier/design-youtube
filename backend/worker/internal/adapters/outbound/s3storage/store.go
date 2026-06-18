package s3storage

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/lregnier/design-youtube/worker/internal/ports"
)

var _ ports.VideoStorage = (*Store)(nil)

type Store struct {
	client     *awss3.Client
	bucket     string
	urlBuilder PublicURLBuilder
}

func NewStore(client *awss3.Client, bucket string, urlBuilder PublicURLBuilder) *Store {
	return &Store{client: client, bucket: bucket, urlBuilder: urlBuilder}
}

func (s *Store) assetURL(key string) string {
	return s.urlBuilder.AssetURL(s.bucket, key)
}

func (s *Store) DownloadRaw(ctx context.Context, videoID, destPath string) error {
	key := fmt.Sprintf("raw/%s/original", videoID)
	out, err := s.client.GetObject(ctx, &awss3.GetObjectInput{Bucket: &s.bucket, Key: &key})
	if err != nil {
		return fmt.Errorf("s3 get %s: %w", key, err)
	}
	defer out.Body.Close()
	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := make([]byte, 1024*1024)
	for {
		n, rerr := out.Body.Read(buf)
		if n > 0 {
			f.Write(buf[:n])
		}
		if rerr != nil {
			break
		}
	}
	return nil
}

func (s *Store) UploadSegments(ctx context.Context, videoID, segDir string) error {
	return filepath.Walk(segDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(segDir, path)
		key := fmt.Sprintf("segments/%s/%s", videoID, rel)
		ct := "video/MP2T"
		if strings.HasSuffix(path, ".m3u8") {
			ct = "application/x-mpegURL"
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return s.putObject(ctx, key, data, ct)
	})
}

func (s *Store) UploadManifest(ctx context.Context, videoID string, content []byte) (string, error) {
	key := fmt.Sprintf("manifests/%s/master.m3u8", videoID)
	if err := s.putObject(ctx, key, content, "application/x-mpegURL"); err != nil {
		return "", err
	}
	return s.assetURL(key), nil
}

func (s *Store) UploadThumbnail(ctx context.Context, videoID string, data []byte) (string, error) {
	key := fmt.Sprintf("thumbnails/%s/thumb.jpg", videoID)
	if err := s.putObject(ctx, key, data, "image/jpeg"); err != nil {
		return "", err
	}
	return s.assetURL(key), nil
}

func (s *Store) putObject(ctx context.Context, key string, data []byte, contentType string) error {
	_, err := s.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: &contentType,
	})
	return err
}
