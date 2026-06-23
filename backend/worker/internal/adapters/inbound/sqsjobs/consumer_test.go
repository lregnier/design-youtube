package sqsjobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseJob_ValidJSON(t *testing.T) {
	// Arrange
	body := `{"videoId":"vid-1","s3Key":"raw/vid-1/original"}`

	// Act
	job, err := parseJob(&body)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "vid-1", job.VideoID)
	assert.Equal(t, "raw/vid-1/original", job.S3Key)
}

func TestParseJob_InvalidJSON(t *testing.T) {
	// Arrange
	body := `not-json`

	// Act
	_, err := parseJob(&body)

	// Assert
	assert.Error(t, err)
}
