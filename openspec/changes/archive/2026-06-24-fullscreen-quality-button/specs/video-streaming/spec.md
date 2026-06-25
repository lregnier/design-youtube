## ADDED Requirements

### Requirement: Quality selector remains visible during fullscreen playback
When the video player enters fullscreen, the quality selector overlay SHALL remain visible and interactive. The player SHALL ensure the fullscreen element is the container div (not the bare `<video>` element), so that the overlay is included in the fullscreen layer.

#### Scenario: Quality button is visible in fullscreen
- **WHEN** the user enters fullscreen via the native video controls
- **THEN** the quality selector button SHALL be visible in the lower-right corner of the fullscreen view

#### Scenario: Quality selection works in fullscreen
- **WHEN** the user clicks the quality button while in fullscreen
- **THEN** the quality menu SHALL open and the user SHALL be able to select a quality level without exiting fullscreen

#### Scenario: Fullscreen swaps to container if video is fullscreen element
- **WHEN** `document.fullscreenElement` (or `document.webkitFullscreenElement`) is the `<video>` element
- **THEN** the player SHALL exit fullscreen and immediately re-enter fullscreen on the container div

#### Scenario: Quality selector works on webkit browsers in fullscreen
- **WHEN** the user enters fullscreen on Safari (webkit fullscreen API)
- **THEN** the quality selector SHALL be visible and functional
