## Context

Both `backend/api` and `backend/worker` currently read `S3_PUBLIC_ENDPOINT_URL` from the environment. A non-empty value implicitly signals LocalStack mode and provides the endpoint URL in one variable. This works but conflates a mode flag with a value, requiring readers to understand the implicit convention.

The strategy types (`PresignedURLTransformer` in api, `PublicURLBuilder` in worker) introduced in the previous change already accept an explicit implementation at construction time — the only remaining implicit coupling is in `Config` and the `main.go` wiring.

## Goals / Non-Goals

**Goals:**
- Replace `S3_PUBLIC_ENDPOINT_URL` with `LOCALSTACK_ENABLED` (bool) and `LOCALSTACK_ENDPOINT` (string) in both modules
- Make `Config.LocalStack bool` the single decision point for mode branching in `main.go`
- Update `docker-compose.yml` to the new names

**Non-Goals:**
- Changing the strategy types or their logic
- Cross-module config sharing
- Validating that `LOCALSTACK_ENDPOINT` is set when `LOCALSTACK_ENABLED=true` (keep it simple — if endpoint is empty, strategy gets an empty string, same as before)

## Decisions

**Two separate env vars, not one**

`LOCALSTACK_ENABLED=true` + `LOCALSTACK_ENDPOINT=http://localhost:4566` separates the boolean gate from the value. Each has a single clear purpose. An operator setting up a new environment knows exactly what each variable controls.

Alternative considered: keep `S3_PUBLIC_ENDPOINT_URL` but add `LOCALSTACK_ENABLED` as a derived override — rejected because it keeps the confusingly-named variable around.

**`Config.LocalStack bool`, not `Config.LocalStackEnabled bool`**

Shorter name; "LocalStack" is already a proper noun that implies the flag is about that tool. `Enabled` suffix is redundant.

**Config fields in both modules, independently**

The two modules have separate `go.mod` files and separate config packages — no shared config struct. Each gets its own `LocalStack bool` and `LocalStackEndpoint string`. This is consistent with existing patterns in the codebase.

**`CLOUDFRONT_DOMAIN` stays required in api config**

Even when LocalStack mode is active, `CLOUDFRONT_DOMAIN` is required by the api config validation today. In LocalStack mode the api doesn't use it (the `NoOpTransformer` doesn't touch the domain), but removing it from `required` is out of scope here.

## Risks / Trade-offs

- **Breaking change for existing local setups** → anyone with `S3_PUBLIC_ENDPOINT_URL` in their `.env` or shell must update to `LOCALSTACK_ENABLED=true` + `LOCALSTACK_ENDPOINT=<url>`. The `docker-compose.yml` update is the canonical reference.
- **No validation that endpoint is set when enabled** → if `LOCALSTACK_ENABLED=true` but `LOCALSTACK_ENDPOINT` is unset, the strategy receives an empty string and will produce malformed URLs. This matches current behaviour (empty `S3_PUBLIC_ENDPOINT_URL` = no rewrite). Acceptable for a dev-only path.
