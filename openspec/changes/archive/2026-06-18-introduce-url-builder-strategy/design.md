## Context

Two S3 adapters currently carry deployment-environment knowledge inside their methods:

- **`backend/api/internal/adapters/outbound/s3store/store.go`**: `PresignUploadPart` conditionally rewrites the presigned URL's host via `rewriteHost()` when `s3PublicEndpointURL != ""`. The `Store` struct holds `s3PublicEndpointURL string`.
- **`backend/worker/internal/adapters/outbound/s3storage/store.go`**: `assetURL(key)` branches between a LocalStack path-style URL (`endpoint/bucket/key`) and a CloudFront URL (`https://domain/key`). The `Store` struct holds both `cloudfrontDomain` and `s3PublicEndpointURL`.

Both env-selection decisions currently happen at construction time (the `!= ""` check gates which path is taken at call time), but the branching sits inside the adapter method rather than at the composition root.

## Goals / Non-Goals

**Goals:**
- Extract the URL-building decision into dedicated types injected at construction time
- Leave `Store.PresignUploadPart` and `Store.assetURL` with no if/else — a single unconditional call to the injected strategy
- Keep implementations package-local to each adapter (not a shared library, since they serve different contracts)
- Make each URL strategy unit-testable independently of AWS SDK

**Non-Goals:**
- Sharing a common URL-builder abstraction across `backend/api` and `backend/worker` (different Go modules, different contracts — `api` rewrites presigned URLs, `worker` builds public asset URLs from keys)
- Changing port interfaces, application layer, or any domain code
- Introducing a third URL strategy (e.g., staging CDN) — this is prep, not pre-implementation

## Decisions

**`backend/api` — `PresignedURLTransformer` interface**

```go
// url_transformer.go (package s3store)
type PresignedURLTransformer interface {
    Transform(presignedURL string) string
}

type NoOpTransformer struct{}
func (NoOpTransformer) Transform(u string) string { return u }

type LocalStackTransformer struct{ publicEndpoint string }
func NewLocalStackTransformer(endpoint string) *LocalStackTransformer { ... }
func (r *LocalStackTransformer) Transform(presignedURL string) string {
    // existing rewriteHost logic moved here
}
```

`Store` changes:
```go
type Store struct {
    client   *awss3.Client
    bucket   string
    transformer PresignedURLTransformer   // replaces s3PublicEndpointURL
}

func NewStore(client *awss3.Client, bucket string, transformer PresignedURLTransformer) *Store { ... }

// PresignUploadPart — no if/else:
presignedURL = s.transformer.Transform(out.URL)
```

`cmd/api/main.go` selection:
```go
var transformer s3store.PresignedURLTransformer = s3store.NoOpTransformer{}
if cfg.S3PublicEndpointURL != "" {
    transformer = s3store.NewLocalStackTransformer(cfg.S3PublicEndpointURL)
}
store := s3store.NewStore(awss3.NewFromConfig(awsCfg, s3Opts...), cfg.S3Bucket, transformer)
```

---

**`backend/worker` — `PublicURLBuilder` interface**

```go
// url_builder.go (package s3storage)
type PublicURLBuilder interface {
    AssetURL(bucket, key string) string
}

type CloudFrontURLBuilder struct{ domain string }
func NewCloudFrontURLBuilder(domain string) *CloudFrontURLBuilder { ... }
func (b *CloudFrontURLBuilder) AssetURL(_, key string) string {
    return fmt.Sprintf("https://%s/%s", b.domain, key)
}

type LocalStackURLBuilder struct{ endpoint string }
func NewLocalStackURLBuilder(endpoint string) *LocalStackURLBuilder { ... }
func (b *LocalStackURLBuilder) AssetURL(bucket, key string) string {
    return fmt.Sprintf("%s/%s/%s", b.endpoint, bucket, key)
}
```

`Store` changes:
```go
type Store struct {
    client     *awss3.Client
    bucket     string
    urlBuilder PublicURLBuilder   // replaces cloudfrontDomain + s3PublicEndpointURL
}

func NewStore(client *awss3.Client, bucket string, urlBuilder PublicURLBuilder) *Store { ... }

// assetURL — no if/else:
func (s *Store) assetURL(key string) string {
    return s.urlBuilder.AssetURL(s.bucket, key)
}
```

`cmd/worker/main.go` selection:
```go
var urlBuilder s3storage.PublicURLBuilder
if cfg.S3PublicEndpointURL != "" {
    urlBuilder = s3storage.NewLocalStackURLBuilder(cfg.S3PublicEndpointURL)
} else {
    urlBuilder = s3storage.NewCloudFrontURLBuilder(cfg.CloudFrontDomain)
}
store := s3storage.NewStore(awss3.NewFromConfig(awsCfg, s3Opts...), cfg.S3Bucket, urlBuilder)
```

---

**Why not a closure (`func(key string) string`) instead of an interface?**

An interface is marginally more verbose but gives the strategy a name (`PresignedURLTransformer` / `PublicURLBuilder`) that documents the intent and is straightforward to mock in tests. A closure works equally well but is harder to identify in stack traces and `fmt.Sprintf("%T", ...)` debugging.

**Why keep interfaces package-local to each adapter?**

The two contracts are distinct (`Rewrite(url) string` vs `AssetURL(bucket, key) string`). Sharing them would require a common internal package across separate Go modules — unnecessary coupling for types that serve different purposes.

## Risks / Trade-offs

- **More files, same number of paths** → the if/else moves from the store method to `main.go`; total branching doesn't decrease, but it moves to the right layer (composition root).
- **Tests for URL strategies are trivial** → `NoOpTransformer`, `LocalStackTransformer`, `CloudFrontURLBuilder`, `LocalStackURLBuilder` are pure functions over strings; test coverage is cheap and high-value.
