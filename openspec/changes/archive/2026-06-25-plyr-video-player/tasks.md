## 1. Dependencies

- [x] 1.1 Install `plyr` and `@types/plyr` in `frontend/web`

## 2. VideoPlayer rewrite

- [x] 2.1 Import `Plyr` and `plyr/dist/plyr.css` in `VideoPlayer.tsx`
- [x] 2.2 Replace `containerRef` / overlay JSX with a single `useRef<HTMLDivElement>` wrapper and a plain `<video>` element (no `controls` attribute)
- [x] 2.3 In `useEffect`, create the `Hls` instance, attach it to the video element, and initialise `Plyr` on the video element with `quality` options derived from `hls.levels` after `MANIFEST_PARSED`
- [x] 2.4 Wire Plyr's `qualitychange` event to set `hls.currentLevel` (map height → level index; `0` → `-1` for Auto)
- [x] 2.5 Destroy both `hls` and `plyr` instances in the `useEffect` cleanup
- [x] 2.6 Remove all overlay JSX, `isFullscreen` state, `toggleFullscreen`, `controlsList`, and `itemStyle` helper

## 3. Cleanup

- [x] 3.1 Delete or archive the `fullscreen-quality-button` OpenSpec change (superseded)
