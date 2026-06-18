## MODIFIED Requirements

### Requirement: Frontend player streams video using HLS.js with adaptive bitrate and manual quality selection
The frontend video page SHALL use HLS.js to load the master manifest URL and play the video. When `Hls.isSupported()` is true, the player SHALL expose a quality selector overlay that allows the viewer to choose a specific quality level (derived from `hls.levels`) or return to automatic adaptive selection. HLS.js SHALL continue to automatically select the appropriate bitrate when "Auto" is chosen. The player SHALL support play, pause, and seek.

#### Scenario: Player loads and starts playback
- **WHEN** a user navigates to a video page for a ready video
- **THEN** HLS.js fetches the master manifest, selects an initial quality level, and begins buffering and playing

#### Scenario: Player adapts quality on bandwidth change in Auto mode
- **WHEN** the client's available bandwidth drops significantly during playback and no manual quality is selected
- **THEN** HLS.js switches to a lower bitrate segment without interrupting playback

#### Scenario: Player displays thumbnail before playback starts
- **WHEN** a user navigates to a video page before pressing play
- **THEN** the thumbnail image is displayed as the video poster

#### Scenario: Quality selector shows available levels
- **WHEN** the master manifest has been parsed and `Hls.isSupported()` is true
- **THEN** the quality selector displays "Auto" plus one option per quality level (e.g. "1080p", "720p", "360p") derived from `hls.levels[i].height`

#### Scenario: Viewer selects a specific quality
- **WHEN** a viewer selects "720p" from the quality selector
- **THEN** `hls.currentLevel` is set to the index of the 720p level and playback continues at that quality

#### Scenario: Viewer returns to Auto
- **WHEN** a viewer selects "Auto" from the quality selector
- **THEN** `hls.currentLevel` is set to `-1` and HLS.js resumes adaptive bitrate selection

#### Scenario: Quality selector not shown on native HLS
- **WHEN** `Hls.isSupported()` is false (e.g. Safari using native HLS)
- **THEN** the quality selector overlay is not rendered
