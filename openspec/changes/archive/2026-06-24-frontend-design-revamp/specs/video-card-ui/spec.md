## ADDED Requirements

### Requirement: Video card is borderless with a rounded thumbnail
Each video card SHALL render without a visible card border. The thumbnail image SHALL have rounded corners (`border-radius: 12px`) and maintain a 16:9 aspect ratio.

#### Scenario: Ready video shows thumbnail
- **WHEN** a video has status `ready` and a thumbnail URL
- **THEN** the card SHALL display the thumbnail image at 16:9 aspect ratio with rounded corners

### Requirement: Video card title is clamped to two lines
The video title SHALL be clamped to a maximum of two lines with an ellipsis on overflow.

#### Scenario: Long title truncates
- **WHEN** a video title exceeds two lines of text
- **THEN** the title SHALL be cut with an ellipsis at the second line

### Requirement: Video card is status-aware
A video card SHALL render differently based on video status. Only `ready` videos SHALL be clickable and navigate to the video page.

#### Scenario: Processing video shows spinner
- **WHEN** a video has status `processing`
- **THEN** the thumbnail area SHALL show a dark placeholder with a spinning animation and "Processing…" label

#### Scenario: Failed video shows error state
- **WHEN** a video has status `failed`
- **THEN** the thumbnail area SHALL show a dark placeholder with a red "Processing failed" label and no spinner

#### Scenario: Ready video is clickable
- **WHEN** a video has status `ready`
- **THEN** the card SHALL be wrapped in a `<Link>` navigating to `/videos/:id`

#### Scenario: Non-ready video is not clickable
- **WHEN** a video has status `processing` or `failed`
- **THEN** the card SHALL NOT be a link and cursor SHALL be `default`
