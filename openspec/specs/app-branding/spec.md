## ADDED Requirements

### Requirement: App is named YouFlick with a teal identity
The application SHALL be presented as "YouFlick" across all surfaces. The brand colour SHALL be teal (`#0d9488`). The page background SHALL be `#f9f9f9`, primary text `#0f0f0f`, and secondary text `#606060`.

#### Scenario: Brand name renders with split colouring
- **WHEN** the YouFlick wordmark is displayed
- **THEN** "You" SHALL render in teal (`#0d9488`) and "Flick" in dark (`#0f0f0f`)

#### Scenario: Logo icon is a teal play button
- **WHEN** the YouFlick logo is displayed
- **THEN** a teal rounded-rectangle SVG with a white play triangle SHALL appear to the left of the wordmark

### Requirement: Global typography and colour reset
The application SHALL apply a CSS reset that sets `box-sizing: border-box`, zero default margin/padding, Roboto/Arial font stack, and the brand background and text colours on `body`.

#### Scenario: Page background matches brand
- **WHEN** any page is rendered
- **THEN** the page background SHALL be `#f9f9f9` and body text SHALL default to `#0f0f0f`
