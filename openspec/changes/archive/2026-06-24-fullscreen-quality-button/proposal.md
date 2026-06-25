## Why

When the user enters fullscreen via the native video controls, the browser makes only the `<video>` element fullscreen — the quality selector overlay (rendered in the parent container div) disappears entirely. Users lose quality control in the mode they're most likely to care about it.

## What Changes

- Listen for `fullscreenchange` / `webkitfullscreenchange` events on the document
- When the `<video>` element itself is detected as `document.fullscreenElement`, immediately swap to fullscreen on the container `<div>` so the quality overlay is included
- Add a `ref` to the container `<div>` to support `requestFullscreen()` calls on it
- Track `isFullscreen` state to allow the quality button position to adapt if needed

## Capabilities

### New Capabilities

<!-- No new capabilities — this is a fix to existing player behaviour -->

### Modified Capabilities

- `video-streaming`: Add requirement that the quality selector remains visible and functional during fullscreen playback

## Impact

- `frontend/web/src/components/VideoPlayer.tsx` — add container ref, fullscreen event listener, and fullscreen-swap logic; no API or backend changes
