package s3storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudFrontURLBuilder_AssetURL(t *testing.T) {
	// Arrange
	b := NewCloudFrontURLBuilder("cdn.example.com")

	// Act
	url := b.AssetURL("my-bucket", "manifests/vid-1/master.m3u8")

	// Assert
	assert.Equal(t, "https://cdn.example.com/manifests/vid-1/master.m3u8", url)
}

func TestEndpointURLBuilder_AssetURL(t *testing.T) {
	// Arrange
	b := NewEndpointURLBuilder("http://localhost:4566")

	// Act
	url := b.AssetURL("my-bucket", "thumbnails/vid-1/thumb.jpg")

	// Assert
	assert.Equal(t, "http://localhost:4566/my-bucket/thumbnails/vid-1/thumb.jpg", url)
}
