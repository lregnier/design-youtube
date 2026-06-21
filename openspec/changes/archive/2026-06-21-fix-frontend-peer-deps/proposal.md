## Why

The frontend was scaffolded with TypeScript `~6.0.2`, but `openapi-typescript` (code-gen tool) requires `typescript@^5.x` and `@typescript-eslint` has an upper bound of `<6.1.0`. The Dockerfile worked around this silently with `--legacy-peer-deps`; local `npm install` fails outright without it.

## What Changes

- Downgrade `typescript` from `~6.0.2` to `~5.8.0` in `frontend/web/package.json` and `package-lock.json`.
- Remove `--legacy-peer-deps` from `RUN npm ci` in `frontend/web/Dockerfile`.

## Capabilities

### New Capabilities

_None._

### Modified Capabilities

_None._ This is a dependency hygiene fix with no spec-level behavior changes.

## Impact

- `frontend/web/package.json` — `typescript` version constraint
- `frontend/web/package-lock.json` — regenerated lockfile
- `frontend/web/Dockerfile` — `npm ci` flag removed
