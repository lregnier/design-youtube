## 1. backend/api — introduce PresignedURLTransformer

- [x] 1.1 Create `backend/api/internal/adapters/outbound/s3store/url_transformer.go` with `PresignedURLTransformer` interface, `NoOpTransformer` struct, and `LocalStackTransformer` struct (move `rewriteHost` logic into `LocalStackTransformer.Rewrite`)
- [x] 1.2 Update `Store` struct: replace `s3PublicEndpointURL string` field with `transformer PresignedURLTransformer`
- [x] 1.3 Update `NewStore` signature: replace `s3PublicEndpointURL string` param with `transformer PresignedURLTransformer`
- [x] 1.4 Update `PresignUploadPart`: replace the `if s.s3PublicEndpointURL != ""` block with `presignedURL = s.transformer.Transform(out.URL)`
- [x] 1.5 Delete the `rewriteHost` free function from `store.go`; remove `"net/url"` import if no longer used

## 2. backend/api — wire in cmd/api/main.go

- [x] 2.1 In `cmd/api/main.go`, construct the transformer based on config (`NoOpTransformer{}` if `S3PublicEndpointURL == ""`, else `NewLocalStackTransformer(cfg.S3PublicEndpointURL)`) and pass it to `s3store.NewStore`

## 3. backend/worker — introduce PublicURLBuilder

- [x] 3.1 Create `backend/worker/internal/adapters/outbound/s3storage/url_builder.go` with `PublicURLBuilder` interface, `CloudFrontURLBuilder` struct, and `LocalStackURLBuilder` struct
- [x] 3.2 Update `Store` struct: replace `cloudfrontDomain string` and `s3PublicEndpointURL string` fields with `urlBuilder PublicURLBuilder`
- [x] 3.3 Update `NewStore` signature: replace `cloudfrontDomain, s3PublicEndpointURL string` params with `urlBuilder PublicURLBuilder`
- [x] 3.4 Update `assetURL`: replace the `if/else` with `return s.urlBuilder.AssetURL(s.bucket, key)`

## 4. backend/worker — wire in cmd/worker/main.go

- [x] 4.1 In `cmd/worker/main.go`, construct the URL builder based on config (`LocalStackURLBuilder` if `S3PublicEndpointURL != ""`, else `CloudFrontURLBuilder`) and pass it to `s3storage.NewStore`

## 5. Update tests

- [x] 5.1 Update `backend/worker` unit tests (if any reference `NewStore` with the old signature) to pass a stub `PublicURLBuilder`
- [x] 5.2 Update `backend/api` unit tests (if any reference `NewStore` with the old signature) to pass a stub `PresignedURLTransformer`

## 6. Verify

- [x] 6.1 `go build ./...` and `go vet ./...` pass in both `backend/api` and `backend/worker`
- [x] 6.2 `go test ./...` passes in both modules
- [x] 6.3 Grep for remaining references to `s3PublicEndpointURL` or `cloudfrontDomain` fields in store files (should be zero)
