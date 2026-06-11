## MODIFIED Requirements

### Requirement: Worker generates an HLS master manifest
After all segments are uploaded, the worker SHALL generate an HLS master playlist (`.m3u8`) referencing the three quality-level media playlists. The manifest SHALL be uploaded to S3 under `manifests/{videoId}/master.m3u8` and served via CloudFront. Variant playlist references in the master manifest SHALL be relative paths that resolve, against the master manifest's own location, to the variant playlists' actual location at `segments/{videoId}/{quality}/media.m3u8`.

#### Scenario: Master manifest created
- **WHEN** all segment uploads complete successfully
- **THEN** a valid HLS master manifest is present in S3 and accessible via CloudFront

#### Scenario: Variant playlist references resolve correctly
- **WHEN** a player resolves a variant playlist reference from the master manifest relative to the manifest's own URL
- **THEN** the resolved URL points to the existing `segments/{videoId}/{quality}/media.m3u8` object
