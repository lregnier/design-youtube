## Context

`VideoPlayer.tsx` renders a `<video controls>` element inside a container `<div>`. The quality selector is a React overlay (`position: absolute`) on the container. When the native fullscreen button (inside the browser's built-in controls) is clicked, the browser calls `videoElement.requestFullscreen()`, making only the `<video>` node the fullscreen element. The container div — and therefore the quality overlay — is not part of the fullscreen layer and becomes invisible.

## Goals / Non-Goals

**Goals:**
- Quality selector is visible and interactive during fullscreen playback
- Fullscreen transition is initiated by the native controls button (no custom fullscreen button required)

**Non-Goals:**
- Replacing native browser video controls with custom ones
- Changing quality selector appearance or behaviour outside fullscreen
- Supporting Picture-in-Picture

## Decisions

**Remove native fullscreen button + add custom fullscreen button in the overlay**
Add `controlsList="nofullscreen"` to the `<video>` element to hide the browser's native fullscreen button. Add a custom fullscreen toggle button alongside the quality selector in the overlay. The button calls `container.requestFullscreen()` directly — a synchronous call from a real user gesture — so the browser always allows it. Both our quality button and fullscreen button are children of the container and are therefore visible in fullscreen.

Alternatives considered and rejected:
- *`requestFullscreen` JS property override*: Native controls call the fullscreen API at the C++ browser level, bypassing JavaScript property overrides entirely. The override never fires.
- *`fullscreenchange` listener + async swap (`exitFullscreen().then(enter)`)*: The async gap between exit and re-enter breaks the user gesture chain. The browser rejects the second `requestFullscreen()` call — fullscreen flashes on and immediately exits.
- *React portal into fullscreen element*: Can't portal into a `<video>` element — it doesn't accept children.
- *`position: fixed` overlay*: Fixed elements are clipped to the fullscreen layer boundary; they won't appear over a fullscreen `<video>`.

**`document` event listener over video element listener**
`fullscreenchange` is fired on `document` (and bubbles from the element). Listening on `document` covers both standard and webkit-prefixed variants in one place and avoids needing an extra ref dependency.

**Webkit prefix**
Safari fires `webkitfullscreenchange` and exposes `document.webkitFullscreenElement`. Both must be handled for cross-browser correctness.

## Risks / Trade-offs

[Brief visual flash on swap] `exitFullscreen` + `requestFullscreen` are two sequential async operations. Most browsers handle this in one frame, but a sub-frame black flash is possible on slower hardware → Acceptable; this matches how other players (e.g. Vimeo) handle the same problem.

[Webkit fullscreen API differences] `document.webkitFullscreenElement` and `webkitRequestFullscreen` are required for Safari support → Handle both standard and webkit variants in the event listener and the swap logic.

[User may use keyboard shortcut or F key] Some browsers allow entering fullscreen without clicking the controls button. The `fullscreenchange` listener fires regardless of how fullscreen was triggered, so the swap will work in all cases.
