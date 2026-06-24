## ADDED Requirements

### Requirement: All pages share a persistent navigation shell
The application SHALL render a sticky `Navbar` at the top of every page via a React Router layout route. Page content SHALL render below the navbar via `<Outlet>`.

#### Scenario: Navbar is visible on every route
- **WHEN** the user navigates to `/`, `/videos/:id`, or `/upload`
- **THEN** the Navbar SHALL be visible at the top of the viewport

#### Scenario: Navbar sticks during scroll
- **WHEN** the user scrolls down a long page
- **THEN** the Navbar SHALL remain fixed at the top (`position: sticky; top: 0`)

### Requirement: Navbar contains logo and upload action
The Navbar SHALL display the YouFlick logo (linking to `/`) on the left and a pill-style Upload button (linking to `/upload`) on the right.

#### Scenario: Logo navigates home
- **WHEN** the user clicks the YouFlick logo
- **THEN** the user SHALL be navigated to `/`

#### Scenario: Upload button navigates to upload page
- **WHEN** the user clicks the Upload button in the Navbar
- **THEN** the user SHALL be navigated to `/upload`
