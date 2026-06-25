## Context

`VideoPlayer.tsx` currently uses `hls.js` directly with a `<video controls>` element and a React overlay div for quality selection and fullscreen. The overlay cannot be embedded in the native control bar, leading to UI conflicts (duplicate/greyed-out fullscreen buttons) and unreliable fullscreen behaviour across browsers.

Plyr is a lightweight (~28 kB gzipped) player library that accepts an externally-managed `hls.js` instance and renders a fully custom, accessible control bar — including a quality menu — while correctly targeting the player container for fullscreen.

## Goals / Non-Goals

**Goals:**
- Single, coherent control bar with play/pause, scrubber, time, volume, quality, and fullscreen
- Quality selector visible and functional in fullscreen
- No duplicate or greyed-out browser control buttons
- Plyr's default theme used as-is (no custom CSS theming for now)

**Non-Goals:**
- Custom Plyr CSS theme matching YouFlick teal branding (can be a follow-up)
- Subtitles / caption support
- Picture-in-Picture

## Decisions

**Pass hls.js instance to Plyr rather than letting Plyr manage hls.js internally**
Plyr supports an `html5.hlsjs` option but the recommended pattern for full control is to create the `Hls` instance separately, attach it to the `<video>` element, then pass the instance to Plyr via `options.hlsjs`. This keeps hls.js version management in our hands and allows quality level wiring via `hls.currentLevel`.

**Quality wiring via `hls.on(LEVEL_SWITCHED)` + `plyr.on('qualitychange')`**
Plyr fires `qualitychange` when the user picks a quality. We map that back to `hls.currentLevel`. Conversely, when hls.js auto-switches levels, we don't need to reflect that back into Plyr's UI (Plyr shows the selected option, not the currently playing level).

**Import Plyr CSS from the npm package**
`import 'plyr/dist/plyr.css'` brings in Plyr's default stylesheet. Vite handles this automatically. No additional build config needed.

**Destroy Plyr instance on unmount**
`plyr.destroy()` also stops the hls.js instance if wired correctly. We destroy both explicitly in the `useEffect` cleanup to avoid memory leaks on navigation.

## Risks / Trade-offs

[Plyr default styling may clash with YouFlick design] Plyr's control bar uses its own colour scheme → Acceptable for now; Plyr CSS variables make theming straightforward in a follow-up.

[Plyr + hls.js version compatibility] Plyr targets hls.js v1.x which matches our current `^1.6.16` → Low risk; pin versions in package.json if needed.

[Bundle size increase] Plyr adds ~28 kB gzipped → Acceptable for the UX improvement gained.
