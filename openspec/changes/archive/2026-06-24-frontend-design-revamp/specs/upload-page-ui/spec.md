## ADDED Requirements

### Requirement: Upload form is presented on a white card surface
The upload form SHALL be rendered inside a white rounded card (`background: #fff; border-radius: 16px`) centred on the `#f9f9f9` page background with a max width of 600 px.

#### Scenario: Form card is visually distinct from page background
- **WHEN** the upload page is rendered
- **THEN** the form SHALL appear inside a white card container clearly separated from the grey page background

### Requirement: Form inputs use labelled, bordered styling
Each form field SHALL have a visible label above it and an input with a 1 px `#e5e5e5` border and rounded corners.

#### Scenario: Inputs are clearly labelled
- **WHEN** the upload page is rendered
- **THEN** each field (file, title, description, secret) SHALL have a label rendered above the input

### Requirement: Upload progress bar is visible during upload
While an upload is in progress, the form SHALL display a progress bar and a "Uploading part N of M…" status line. The submit button SHALL be disabled and show "Uploading…".

#### Scenario: Progress bar updates per chunk
- **WHEN** a chunk upload completes
- **THEN** the progress bar width SHALL update to reflect `current / total * 100%`

#### Scenario: Submit button is disabled during upload
- **WHEN** an upload is in progress
- **THEN** the submit button SHALL be disabled with a greyed-out appearance and "Uploading…" label
