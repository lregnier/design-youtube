## ADDED Requirements

### Requirement: Secret field has a visibility toggle
The Upload page secret field SHALL include a clickable eye icon button that toggles the field between hidden and visible states.

#### Scenario: Default state is hidden
- **WHEN** the Upload page loads
- **THEN** the secret field SHALL display characters as dots (type="password")

#### Scenario: User clicks eye icon to reveal
- **WHEN** the user clicks the eye icon while the field is hidden
- **THEN** the field SHALL switch to visible (type="text") and the icon SHALL change to eye-with-slash

#### Scenario: User clicks eye-slash icon to hide
- **WHEN** the user clicks the eye-slash icon while the field is visible
- **THEN** the field SHALL switch to hidden (type="password") and the icon SHALL change to eye

#### Scenario: Toggle is disabled during upload
- **WHEN** an upload is in progress
- **THEN** the toggle button SHALL be inert (input is disabled, no interaction possible)
