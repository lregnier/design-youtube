## Why

The custom HLS.js player uses a React overlay for the quality selector which cannot be embedded in the native browser control bar and conflicts with the native fullscreen button. Plyr provides a polished, accessible control bar with native-feeling quality selection and correct fullscreen support out of the box — eliminating the overlay hack entirely.

## What Changes

- Add `plyr` and `@types/plyr` as dependencies
- Remove the custom overlay (quality dropdown + custom fullscreen button) from `VideoPlayer.tsx`
- Initialise Plyr on the `<video>` element, passing the existing `hls.js` instance for HLS playback
- Wire `hls.levels` into Plyr's quality menu so the viewer can switch quality from within the Plyr control bar
- Remove `controlsList`, `containerRef`, fullscreen state, and all overlay JSX from the component
- Remove the `fullscreen-quality-button` OpenSpec change (superseded by this change)

## Capabilities

### New Capabilities

<!-- No new user-facing capabilities — this is a player implementation swap -->

### Modified Capabilities

- `video-streaming`: The player control bar implementation changes from native `<video controls>` + React overlay to Plyr; quality selection and fullscreen behaviour requirements are otherwise unchanged

## Impact

- `frontend/web/package.json` — add `plyr`, `@types/plyr`
- `frontend/web/src/components/VideoPlayer.tsx` — full rewrite
- No backend, API, or routing changes
