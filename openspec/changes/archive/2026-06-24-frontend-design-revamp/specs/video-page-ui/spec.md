## ADDED Requirements

### Requirement: Video page has a chevron back navigation link
The video detail page SHALL display a back link with a left-chevron SVG icon that navigates to `/`.

#### Scenario: Back link is visible and functional
- **WHEN** the video page is rendered
- **THEN** a "Back" link with a chevron icon SHALL be visible at the top of the page and SHALL navigate to `/` when clicked

### Requirement: Video page shows a spinner placeholder for non-ready videos
When a video is not yet ready, the video area SHALL show a dark 16:9 placeholder. If the video is processing, a spinner animation SHALL be shown inside the placeholder.

#### Scenario: Processing video shows spinner in placeholder
- **WHEN** a video has status `processing`
- **THEN** the video area SHALL render a dark `#1a1a1a` placeholder at 16:9 with a spinning animation and "Processing…" label

#### Scenario: Failed video shows error label in placeholder
- **WHEN** a video has status `failed`
- **THEN** the video area SHALL render the dark placeholder with a "Processing failed" label and no spinner

### Requirement: Video page displays metadata below the player
Below the player (or placeholder), the page SHALL display the video title, upload date in long format, and (if present) the description in a distinct rounded block.

#### Scenario: Metadata renders correctly
- **WHEN** a ready video is displayed
- **THEN** the title SHALL appear in a large font, the date in long format (e.g., "June 9, 2026"), and the description in a `#f2f2f2` rounded block
