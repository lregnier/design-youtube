## Context

The Upload page has a secret field rendered as `<input type="password">`. There is no way to verify what was typed before submitting, which can cause upload failures due to typos. This is a contained, frontend-only change with no backend or API impact.

## Goals / Non-Goals

**Goals:**
- Let users toggle secret field visibility on the Upload page
- Provide clear visual feedback via an eye / eye-off icon

**Non-Goals:**
- Hold-to-reveal behaviour (not requested)
- Applying the pattern to other password fields in the app

## Decisions

**Inline toggle button over a separate checkbox**
An icon button placed inside the input is the standard UX pattern (used by most modern auth forms). A checkbox label ("Show secret") would work but is visually noisier and less conventional.

**SVG icons inlined in JSX over an icon library**
The project currently uses no icon library. Adding one for two icons adds unnecessary dependency weight. The two Feather-style SVGs (eye, eye-with-slash) are small and self-contained.

**`useState` for `showSecret` flag**
Simple boolean local state — no need for context or refs. The flag lives in `UploadPage` and drives both the `type` attribute and which icon is rendered.

## Risks / Trade-offs

[Autofill interference] Some browsers may not autofill a field that switches between `password` and `text` types → Low risk; the secret is not a password managers' target field.

[No migration needed] Pure additive UI change; no data model or API involvement.
