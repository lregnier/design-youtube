## Why

The Upload page secret field shows entered characters as dots with no way to verify what was typed, causing failed uploads due to typos in the secret. A visibility toggle lets users confirm their input before submitting.

## What Changes

- Add a clickable eye icon button inside the secret input field
- Toggle input type between `password` (hidden) and `text` (visible) on each click
- Icon switches between eye (show) and eye-with-slash (hide) to reflect current state
- No hold-to-reveal behaviour — state persists until toggled again

## Capabilities

### New Capabilities

- `upload-secret-visibility`: Toggle visibility of the upload secret field on the Upload page

### Modified Capabilities

<!-- No existing spec-level requirements are changing -->

## Impact

- `frontend/web/src/pages/UploadPage.tsx`: wrap secret `<input>` in a relative `<div>`, add absolute-positioned toggle `<button>` with SVG icons, add `showSecret` boolean state
