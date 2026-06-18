## Context

`VideoPlayer.tsx` creates an `Hls` instance, loads the manifest, and attaches it to a `<video>` element. The `hls` instance is currently scoped to the `useEffect` cleanup function and not exposed to the rest of the component. HLS.js fires `Hls.Events.MANIFEST_PARSED` after parsing the master manifest, at which point `hls.levels` contains an array of `Level` objects (each with `height`, `width`, `bitrate`). Setting `hls.currentLevel` to a level index locks playback to that rendition; `-1` re-enables ABR.

## Goals / Non-Goals

**Goals:**
- Expose quality options derived from `hls.levels` in a picker overlaid on the player
- Let the user switch quality at any time during playback without interruption
- Show which quality is currently selected; label the trigger button with the active selection

**Non-Goals:**
- Custom player chrome (no full player redesign)
- Persisting quality preference across sessions
- Showing current network bandwidth or buffer state
- Supporting native HLS (Safari without HLS.js) — quality selection only applies when `Hls.isSupported()`

## Decisions

**Store `hls` instance and levels in component state/refs**

The `hls` instance needs to live beyond the `useEffect` so the picker can call `hls.currentLevel = index`. Use a `useRef` for the instance (no re-render on assignment) and `useState` for the levels array and selected level index (trigger re-render when they change).

**Derive quality labels from `level.height`**

`hls.levels[i].height` gives the vertical resolution (e.g. 1080, 720, 360). Formatting as `"${height}p"` matches what users expect and matches the labels the worker embeds in the manifest. No need to parse the manifest URI.

**Overlay picker, not a separate controls bar**

A small absolute-positioned button in the bottom-right corner of the player container opens a dropdown list above it. This avoids interfering with the native `<video controls>` bar and keeps the implementation contained to one component. The native controls handle play/pause/seek/volume.

**`currentLevel` vs `nextLevel`**

`hls.currentLevel` switches immediately (mid-segment); `hls.nextLevel` waits for the current segment to finish. Immediate switching gives faster feedback; the brief quality dip during the switch is acceptable for a manual selection action.

## Risks / Trade-offs

- **Native HLS (Safari)**: `Hls.isSupported()` is false on Safari, which uses native HLS. Quality selection is silently unavailable there — the picker simply won't render. Acceptable given Safari's own adaptive streaming handles this transparently.
- **Overlay z-index conflicts**: The picker dropdown sits above the video element. If the browser's native controls overlap, the dropdown may be partially obscured. Mitigated by positioning the trigger above the native controls bar.
- **Level index stability**: `hls.levels` is populated once on `MANIFEST_PARSED` and does not change during playback for VOD content, so storing level index is safe.
