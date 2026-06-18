package s3storage

import "fmt"

type PublicURLBuilder interface {
	AssetURL(bucket, key string) string
}

type CloudFrontURLBuilder struct{ domain string }

func NewCloudFrontURLBuilder(domain string) *CloudFrontURLBuilder {
	return &CloudFrontURLBuilder{domain: domain}
}

func (b *CloudFrontURLBuilder) AssetURL(_, key string) string {
	return fmt.Sprintf("https://%s/%s", b.domain, key)
}

type LocalStackURLBuilder struct{ endpoint string }

func NewLocalStackURLBuilder(endpoint string) *LocalStackURLBuilder {
	return &LocalStackURLBuilder{endpoint: endpoint}
}

func (b *LocalStackURLBuilder) AssetURL(bucket, key string) string {
	return fmt.Sprintf("%s/%s/%s", b.endpoint, bucket, key)
}
