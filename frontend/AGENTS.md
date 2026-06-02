# Frontend

React 19 + TypeScript SPA built with Vite. Communicates with the Go backend API.

## Structure

```
src/
  api/
    openapi.yaml -> (source lives at ../backend/api/openapi.yaml)
    types.gen.ts  Generated TypeScript types — do not edit by hand
    client.ts     Typed fetch wrappers using generated types
  components/
    VideoCard.tsx
    VideoPlayer.tsx
  pages/
    HomePage.tsx    /
    VideoPage.tsx   /videos/:videoId
    UploadPage.tsx  /upload
```

## Commands

```bash
npm install --legacy-peer-deps

# Regenerate types after changing backend/api/openapi.yaml
npm run generate:api

# Dev server
npm run dev

# Production build
npm run build
```

## Key conventions

- `src/api/types.gen.ts` is auto-generated from `../backend/api/openapi.yaml`. Never edit it manually — run `npm run generate:api` instead.
- All API calls go through `src/api/client.ts`. Add new API functions there using the generated types.
- The API base URL is set via `VITE_API_URL` env var (defaults to `http://localhost:8080`).
- Upload endpoints need the `X-Upload-Secret` header — the user provides this in the upload form.
- Max file size enforced client-side before any network call: 100MB.
- Chunk size for multipart upload: 10MB.

## Routing

| Path                | Component    |
|---------------------|--------------|
| `/`                 | HomePage     |
| `/videos/:videoId`  | VideoPage    |
| `/upload`           | UploadPage   |

## HLS playback

`VideoPlayer.tsx` uses HLS.js when the browser supports it, falling back to native HLS (Safari). The manifest URL comes from the backend's `VideoDetail.manifestUrl` field and is a CloudFront URL.
