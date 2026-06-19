package s3store

import "net/url"

type PresignedURLTransformer interface {
	Transform(presignedURL string) string
}

type NoOpTransformer struct{}

func (NoOpTransformer) Transform(u string) string { return u }

type EndpointTransformer struct{ publicEndpoint string }

func NewEndpointTransformer(endpoint string) *EndpointTransformer {
	return &EndpointTransformer{publicEndpoint: endpoint}
}

func (t *EndpointTransformer) Transform(presignedURL string) string {
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
