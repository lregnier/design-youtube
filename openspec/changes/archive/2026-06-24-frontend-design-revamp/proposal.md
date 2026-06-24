## Why

The original frontend had no coherent visual identity — plain browser defaults, no navigation chrome, and no brand presence. A design refresh makes the product feel polished and gives it a clear identity (YouFlick) while establishing consistent layout patterns for future pages.

## What Changes

- **App renamed to YouFlick** — teal-on-dark wordmark with a teal play-button SVG logo
- **Global CSS reset** — `#f9f9f9` page background, `#0f0f0f` primary text, `#606060` secondary text, `#0d9488` teal accent, Roboto/Arial font stack, box-sizing reset
- **New `Navbar` component** — sticky top header with YouFlick logo (links to `/`) and a pill-style Upload button
- **New `Layout` component** — React Router layout route wrapping all pages with `<Navbar>` + `<Outlet>`
- **`VideoCard` redesign** — borderless YouTube-style card; thumbnail with rounded corners; status-aware: spinner + "Processing…" for in-progress, red "Processing failed" for failed, clickable `<Link>` only when ready; 2-line title clamp
- **`HomePage` redesign** — removed inline header; max-width 1280 px grid; icon-based empty state; 5 s polling via `setInterval` / `useRef` that self-stops when no video is processing
- **`VideoPage` redesign** — SVG chevron back link; larger title; long-form date; description in `#f2f2f2` rounded block; spinner placeholder for non-ready videos
- **`UploadPage` redesign** — white card surface; labelled inputs with `#e5e5e5` borders; teal-coloured progress bar; styled submit button

## Capabilities

### New Capabilities

- `app-branding`: YouFlick brand identity — name, logo, colour palette, and global typography
- `app-layout`: Persistent navigation shell (Navbar + Layout route) shared across all pages
- `video-card-ui`: Status-aware video card component with thumbnail, title, date, and processing states
- `home-page-ui`: Home page grid with polling, empty state, and status-aware cards
- `video-page-ui`: Video detail page with player placeholder, metadata, and back navigation
- `upload-page-ui`: Upload form page with card surface, progress bar, and labelled inputs

### Modified Capabilities

<!-- No existing spec-level requirements are changing -->

## Impact

- `frontend/web/src/index.css` — full rewrite
- `frontend/web/src/main.tsx` — layout route added
- `frontend/web/src/components/Navbar.tsx` — new file
- `frontend/web/src/components/Layout.tsx` — new file
- `frontend/web/src/components/VideoCard.tsx` — redesigned
- `frontend/web/src/pages/HomePage.tsx` — redesigned
- `frontend/web/src/pages/VideoPage.tsx` — redesigned
- `frontend/web/src/pages/UploadPage.tsx` — redesigned
- No backend, API, or dependency changes
