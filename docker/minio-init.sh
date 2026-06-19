#!/bin/sh
set -e

mc alias set local http://minio:9000 minioadmin minioadmin
mc mb --ignore-existing local/design-youtube-video-prod
mc anonymous set download local/design-youtube-video-prod
