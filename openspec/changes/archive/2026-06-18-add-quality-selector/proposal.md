## Why

The player currently runs HLS.js in full ABR mode with no user control — viewers cannot choose a quality level manually. The backend already produces a three-variant HLS master manifest (1080p / 720p / 360p), so no backend work is needed; the gap is purely in the player UI.

## What Changes

- Add a quality picker overlay to `VideoPlayer.tsx` that appears on the video element
- After `MANIFEST_PARSED` fires, read `hls.levels` to populate the options list: "Auto" (ABR) plus one entry per level using `level.height` (e.g. "1080p", "720p", "360p")
- Selecting a level sets `hls.currentLevel` to the level index; selecting "Auto" sets it to `-1`
- The currently active level is highlighted in the picker; the selected label is shown on the trigger button

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `video-streaming`: the player SHALL expose a manual quality selector in addition to automatic ABR — viewers can override the adaptive selection at any time

## Impact

- Affected file: `frontend/web/src/components/VideoPlayer.tsx` only
- No backend changes, no API changes, no new npm dependencies (HLS.js already installed)
