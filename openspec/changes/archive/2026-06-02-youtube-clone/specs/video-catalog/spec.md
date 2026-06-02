## ADDED Requirements

### Requirement: Backend lists all ready videos
The backend SHALL expose `GET /videos`. The response SHALL return an array of video summaries for all videos with status `ready`, ordered by upload timestamp descending. Each summary SHALL include: videoId, title, thumbnail URL, and upload timestamp. Videos with status `uploading`, `processing`, or `failed` SHALL be excluded.

#### Scenario: Returns list of ready videos
- **WHEN** a client requests `GET /videos` and ready videos exist
- **THEN** the server returns 200 with an array of video summaries ordered newest first

#### Scenario: Returns empty list when no videos are ready
- **WHEN** a client requests `GET /videos` and no videos have status `ready`
- **THEN** the server returns 200 with an empty array

### Requirement: Frontend homepage displays the video catalog
The frontend root route (`/`) SHALL render a grid of video cards for all ready videos. Each card SHALL display the thumbnail image, video title, and upload date. Clicking a card SHALL navigate to the video page for that videoId.

#### Scenario: Homepage renders video grid
- **WHEN** a user visits the homepage and ready videos exist
- **THEN** a grid of video cards is displayed, each showing thumbnail, title, and date

#### Scenario: Homepage renders empty state
- **WHEN** a user visits the homepage and no ready videos exist
- **THEN** an empty state message is displayed (e.g., "No videos yet")

#### Scenario: Clicking a card navigates to the video page
- **WHEN** a user clicks a video card on the homepage
- **THEN** the browser navigates to `/videos/{videoId}`

### Requirement: Frontend upload page is accessible to secret holders
The frontend SHALL expose an `/upload` route with a form for selecting a video file and entering a title and description. The form SHALL accept an upload secret input field. On submit, the frontend SHALL execute the multipart upload flow against the backend. Upload progress SHALL be displayed per chunk.

#### Scenario: Successful upload flow
- **WHEN** a user fills the upload form with a valid secret and a file ≤ 100MB and submits
- **THEN** the frontend completes the multipart upload and redirects to the homepage on completion

#### Scenario: File too large is rejected client-side
- **WHEN** a user selects a file larger than 100MB
- **THEN** the frontend displays an error before making any network request

#### Scenario: Wrong secret shows error
- **WHEN** a user submits the upload form with an incorrect secret
- **THEN** the backend returns 401 and the frontend displays an "invalid secret" error message
