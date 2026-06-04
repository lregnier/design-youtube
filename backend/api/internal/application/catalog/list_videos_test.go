package catalog

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
	"github.com/lregnier/design-youtube/api/internal/mocks"
)

func TestListVideos_Execute_ReturnsList(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	videos := []*video.Video{
		{ID: "vid-1", Title: "First", Status: video.StatusReady, UploadedAt: time.Now()},
		{ID: "vid-2", Title: "Second", Status: video.StatusReady, UploadedAt: time.Now().Add(-time.Hour)},
	}
	repo.EXPECT().ListReady(mock.Anything).Return(videos, nil)

	uc := NewListVideos(repo)

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, video.VideoID("vid-1"), result[0].ID)
}

func TestListVideos_Execute_EmptyList(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	repo.EXPECT().ListReady(mock.Anything).Return([]*video.Video{}, nil)

	uc := NewListVideos(repo)

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}
