## Context

The Vite scaffold set `typescript@~6.0.2`. TypeScript 6.0 was brand new at project creation and the broader ecosystem (notably `openapi-typescript@7.x` and `@typescript-eslint@8.x`) had not yet released versions supporting it. The Dockerfile masked the conflict with `--legacy-peer-deps`; local development hit the error on a clean `npm install`.

## Goals / Non-Goals

**Goals:**
- Clean `npm install` with no flags, locally and in Docker.
- All peer dependency constraints satisfied without suppression.

**Non-Goals:**
- Upgrading any other dependency.
- Adopting TypeScript 6-specific language features (none are used).

## Decisions

### Downgrade TypeScript to `~5.8.0` rather than keeping 6.x

**Rationale**: The project uses no TypeScript 6-specific features. TypeScript 5.8 is the latest stable 5.x release with full ecosystem support. Keeping 6.x would require either `--legacy-peer-deps` everywhere or waiting indefinitely for `openapi-typescript` to ship a compatible release.

**Alternatives considered:**
- **Keep TypeScript 6.x + `--legacy-peer-deps` everywhere**: Suppresses the warning without fixing it. A future `npm install` by a new contributor would fail or silently install an incompatible tree.
- **Replace `openapi-typescript`**: Overkill — the conflict is a timing issue, not a tool problem.
- **Pin TypeScript to `5.x` with `^5.0.0`**: Wider range than needed; `~5.8.0` is more precise and avoids accidentally pulling in a future 5.x with breaking changes.

## Risks / Trade-offs

- **[Risk] TypeScript 5.8 vs 6.x language differences** → Mitigation: the project compiles cleanly on 5.8 (verified). No 6.x-only syntax is in use.
- **[Risk] Future upgrades** → When `openapi-typescript` and `@typescript-eslint` release TypeScript 6.x support, upgrading is a one-line `package.json` change.
