## 1. Refactor hls instance into a ref

- [x] 1.1 In `VideoPlayer.tsx`, move the `hls` instance into a `useRef<Hls | null>` so it persists across renders and is accessible outside the `useEffect`

## 2. Capture quality levels from HLS.js

- [x] 2.1 Add `useState<Hls['levels']>` (or `Level[]`) to hold the parsed levels list, initialised to `[]`
- [x] 2.2 Add `useState<number>` for the selected level index, initialised to `-1` (Auto)
- [x] 2.3 In the `useEffect`, listen for `Hls.Events.MANIFEST_PARSED` and set the levels state from `hls.levels`

## 3. Build the quality selector UI

- [x] 3.1 Add `useState<boolean>` to track whether the picker dropdown is open
- [x] 3.2 Wrap the `<video>` element in a `position: relative` container `<div>`
- [x] 3.3 Render a trigger button (absolute-positioned, bottom-right) that shows the active label ("Auto" when `selectedLevel === -1`, otherwise `"${hls.levels[selectedLevel].height}p"`); only render when `levels.length > 0`
- [x] 3.4 Render the dropdown list above the trigger when `isOpen` is true: one item for "Auto" (`currentLevel = -1`) and one per level (`"${level.height}p"`, index as `currentLevel`)
- [x] 3.5 On item click: set `hls.currentLevel` to the chosen index, update `selectedLevel` state, close the dropdown
- [x] 3.6 Highlight the active option in the dropdown list

## 4. Verify

- [x] 4.1 TypeScript compiles with no errors (`tsc --noEmit`)
- [x] 4.2 Manual test: open a ready video, confirm quality selector appears with Auto/1080p/720p/360p options
- [x] 4.3 Manual test: select "360p", confirm `hls.currentLevel` changes (visible in browser devtools HLS.js debug or network tab)
- [x] 4.4 Manual test: select "Auto", confirm ABR resumes
