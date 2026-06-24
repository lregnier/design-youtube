## MODIFIED Requirements

### Requirement: Backend lists all visible videos
The backend SHALL expose `GET /videos`. The response SHALL return an array of video summaries for all videos with status `processing`, `ready`, or `failed`, ordered by upload timestamp descending. Each summary SHALL include: videoId, title, thumbnailUrl, uploadedAt, and status. Videos with status `uploading` SHALL be excluded (the upload flow is still in progress and the video is not yet visible to the viewer).

#### Scenario: Returns list including processing and ready videos
- **WHEN** a client requests `GET /videos` and videos exist with status `processing` and `ready`
- **THEN** the server returns 200 with summaries for both, ordered newest first

#### Scenario: Returns failed videos
- **WHEN** a client requests `GET /videos` and a video has status `failed`
- **THEN** the server includes the failed video in the response

#### Scenario: Excludes uploading videos
- **WHEN** a client requests `GET /videos` and a video has status `uploading`
- **THEN** the server does not include that video in the response

#### Scenario: Returns empty list when no visible videos exist
- **WHEN** a client requests `GET /videos` and all videos are in `uploading` state or none exist
- **THEN** the server returns 200 with an empty array

### Requirement: Frontend homepage displays all visible videos with status-aware cards
The frontend root route (`/`) SHALL render a grid of video cards for all videos returned by `GET /videos`. Each card SHALL display the video title and upload date. Cards for `ready` videos SHALL display the thumbnail image and be clickable, navigating to `/videos/{videoId}`. Cards for `processing` videos SHALL show a "Processing…" indicator in place of the thumbnail and SHALL NOT be clickable. Cards for `failed` videos SHALL show a "Failed" indicator and SHALL NOT be clickable.

#### Scenario: Ready video card is clickable with thumbnail
- **WHEN** a video with status `ready` is shown on the homepage
- **THEN** the card displays the thumbnail and navigating to it opens `/videos/{videoId}`

#### Scenario: Processing video card is not clickable
- **WHEN** a video with status `processing` is shown on the homepage
- **THEN** the card shows a "Processing…" indicator and clicking it does nothing

#### Scenario: Failed video card is not clickable
- **WHEN** a video with status `failed` is shown on the homepage
- **THEN** the card shows a "Failed" indicator and clicking it does nothing

#### Scenario: Homepage renders empty state
- **WHEN** `GET /videos` returns an empty array
- **THEN** an empty state message is displayed

### Requirement: Frontend homepage polls while videos are processing
The homepage SHALL automatically refresh the video list every 5 seconds while at least one video has status `processing`. Polling SHALL stop when no `processing` videos remain (all are `ready` or `failed`) and SHALL be cleaned up on component unmount.

#### Scenario: Polling starts when processing video is present
- **WHEN** the homepage loads and at least one video has status `processing`
- **THEN** the video list refreshes every 5 seconds without user interaction

#### Scenario: Polling stops when processing completes
- **WHEN** all videos transition out of `processing` status
- **THEN** the homepage stops polling and no further requests are made

#### Scenario: Polling cleans up on unmount
- **WHEN** the user navigates away from the homepage while polling is active
- **THEN** the interval is cleared and no further requests are made
