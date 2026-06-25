## 1. Container ref

- [x] 1.1 Add `containerRef` (`useRef<HTMLDivElement>`) to `VideoPlayer` and attach it to the outer `<div>`

## 2. Fullscreen swap logic

- [x] 2.1 Add `controlsList="nofullscreen"` to `<video>` to hide the native fullscreen button
- [x] 2.2 Add `isFullscreen` state tracked via `fullscreenchange` / `webkitfullscreenchange` listener on `document`
- [x] 2.3 Add `toggleFullscreen` function that calls `container.requestFullscreen()` or `exitFullscreen()` based on current state (with webkit fallback)
- [x] 2.4 Add custom fullscreen button in the overlay alongside the quality selector; icon toggles between enter/exit SVG
