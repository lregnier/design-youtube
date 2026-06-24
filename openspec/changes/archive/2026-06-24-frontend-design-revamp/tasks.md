## 1. Global Styles & Branding

- [x] 1.1 Rewrite `index.css` with box-sizing reset, `#f9f9f9` background, `#0f0f0f` text, Roboto/Arial font stack
- [x] 1.2 Define teal `#0d9488` as the accent colour in global styles

## 2. Navigation Shell

- [x] 2.1 Create `Navbar.tsx` with YouFlick teal logo SVG + wordmark and pill Upload button
- [x] 2.2 Create `Layout.tsx` layout route rendering `<Navbar>` + `<Outlet>`
- [x] 2.3 Wrap all routes in `main.tsx` with the `<Layout>` route

## 3. VideoCard

- [x] 3.1 Redesign `VideoCard.tsx` as a borderless card with rounded 16:9 thumbnail
- [x] 3.2 Add status-aware rendering: spinner for `processing`, red error label for `failed`
- [x] 3.3 Wrap card in `<Link>` only when status is `ready`; use `<div>` otherwise
- [x] 3.4 Clamp title to 2 lines with `-webkit-line-clamp`

## 4. HomePage

- [x] 4.1 Remove inline page header; set max-width 1280 px container
- [x] 4.2 Add icon-based empty state (play icon + "No videos yet" messaging)
- [x] 4.3 Implement 5 s polling via `setInterval` / `useRef` that starts on processing videos and self-stops when none remain

## 5. VideoPage

- [x] 5.1 Add SVG chevron back link navigating to `/`
- [x] 5.2 Render dark 16:9 placeholder with spinner for `processing`, error label for `failed`
- [x] 5.3 Display title, long-form date, and description in a `#f2f2f2` rounded block

## 6. UploadPage

- [x] 6.1 Wrap form in a white rounded card surface
- [x] 6.2 Style all inputs with `#e5e5e5` borders and labelled layout
- [x] 6.3 Add progress bar and per-chunk status line; disable submit button during upload
