## ADDED Requirements

### Requirement: Home page displays videos in a responsive grid
The home page SHALL render videos in an auto-fill grid with a minimum card width of 300 px and a maximum page width of 1280 px.

#### Scenario: Videos fill the grid
- **WHEN** one or more videos are returned from the API
- **THEN** they SHALL be laid out in a CSS grid with `repeat(auto-fill, minmax(300px, 1fr))` columns

### Requirement: Home page shows an icon-based empty state
When no videos exist, the home page SHALL display a centred play-button icon with a message prompting the user to upload their first video.

#### Scenario: Empty state renders when no videos
- **WHEN** the API returns an empty list
- **THEN** the page SHALL display an icon and "No videos yet" heading with a sub-line "Upload your first video to get started"

### Requirement: Home page polls while videos are processing
The home page SHALL automatically refresh the video list every 5 seconds while at least one video has status `processing`. Polling SHALL stop as soon as no processing videos remain.

#### Scenario: Polling starts on processing video
- **WHEN** the initial load returns at least one video with status `processing`
- **THEN** the page SHALL begin polling the API every 5 seconds

#### Scenario: Polling stops when processing completes
- **WHEN** a poll response contains no videos with status `processing`
- **THEN** the interval SHALL be cleared and polling SHALL stop

#### Scenario: Polling cleans up on unmount
- **WHEN** the user navigates away from the home page while polling is active
- **THEN** the interval SHALL be cleared to prevent memory leaks
