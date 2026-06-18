package s3store

import "net/url"

type PresignedURLTransformer interface {
	Transform(presignedURL string) string
}

type NoOpTransformer struct{}

func (NoOpTransformer) Transform(u string) string { return u }

type LocalStackTransformer struct{ publicEndpoint string }

func NewLocalStackTransformer(endpoint string) *LocalStackTransformer {
	return &LocalStackTransformer{publicEndpoint: endpoint}
}

func (t *LocalStackTransformer) Transform(presignedURL string) string {
	pub, err := url.Parse(t.publicEndpoint)
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
