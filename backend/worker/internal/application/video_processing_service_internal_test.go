package application

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildMasterManifest_AllQualities(t *testing.T) {
	// Arrange
	videoID := "vid-1"

	// Act
	result := buildMasterManifest(videoID, qualities)

	// Assert
	assert.True(t, strings.HasPrefix(result, "#EXTM3U\n"))
	assert.Equal(t, 3, strings.Count(result, "#EXT-X-STREAM-INF:"))
	assert.Contains(t, result, "BANDWIDTH=4500000,RESOLUTION=1920x1080")
	assert.Contains(t, result, "BANDWIDTH=2800000,RESOLUTION=1280x720")
	assert.Contains(t, result, "BANDWIDTH=1000000,RESOLUTION=640x360")
	assert.Contains(t, result, "../../segments/vid-1/1080p/media.m3u8")
	assert.Contains(t, result, "../../segments/vid-1/720p/media.m3u8")
	assert.Contains(t, result, "../../segments/vid-1/360p/media.m3u8")
}
