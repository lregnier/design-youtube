## Context

The frontend was built as a minimal proof-of-concept with no design system — raw HTML elements, no shared layout, and no brand identity. The revamp introduces a consistent visual language across all pages without adding external UI libraries or a CSS-in-JS framework.

## Goals / Non-Goals

**Goals:**
- Establish a coherent colour palette and typography baseline via a single global CSS file
- Introduce a persistent navigation shell (Navbar + Layout route) reusable by all current and future pages
- Make each page feel polished: correct spacing, readable typography, meaningful empty and loading states
- Keep the app name "YouFlick" with a recognisable teal brand colour

**Non-Goals:**
- Dark mode
- Responsive / mobile breakpoints (desktop-first for now)
- Animation beyond the processing spinner
- Replacing inline styles with a CSS module or Tailwind system (deferred)

## Decisions

**Inline styles over CSS modules or a utility framework**
The project already uses inline styles throughout. Introducing a new styling strategy mid-project would require a full migration. Inline styles keep the diff contained and avoid a new dependency.

**Layout route (`<Outlet>`) over per-page `<Navbar>` imports**
A React Router layout route renders the Navbar once and avoids duplicating the import on every page. It also makes future layout additions (sidebars, footers) a single-file change.

**`setInterval` + `useRef` for home-page polling**
A ref holds the interval handle so the cleanup function and the polling callback share the same reference without stale closures. The interval is created only when a processing video is detected and torn down when none remain — avoiding unnecessary API traffic.

**Teal `#0d9488` as the single accent colour**
Teal is visually distinct from YouTube's red, making the YouFlick brand instantly differentiated. It works well against both `#fff` and `#f9f9f9` backgrounds at WCAG AA contrast for UI chrome (logos, buttons) even if not for small body text.

**No icon library — inline SVGs**
The revamp uses ~5 distinct icons (play, upload, chevron, eye, eye-slash). Shipping a full icon library for five glyphs is unnecessary overhead. Each SVG is small and typed inline in JSX.

## Risks / Trade-offs

[Inline styles at scale] As pages grow, inline styles become hard to maintain → Mitigated by keeping styles co-located with the component that owns them; a future CSS-modules migration is straightforward.

[Polling on the home page] Every client polls every 5 s while a video is processing → Acceptable at current scale; a WebSocket or server-sent events solution can replace it later without changing the component API.

[No responsive design] The 1280 px max-width grid wraps gracefully but was not tested on small screens → Acceptable for the current development phase.
