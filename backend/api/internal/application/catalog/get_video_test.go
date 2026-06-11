package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
	"github.com/lregnier/design-youtube/api/internal/gen/mocks"
)

func TestGetVideo_Execute_CacheHit(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	cache := mocks.NewMockCache(t)

	vid := &video.Video{ID: "vid-1", Title: "Cached", Status: video.StatusReady, UploadedAt: time.Now()}
	data, _ := json.Marshal(vid)
	cache.EXPECT().Get(mock.Anything, "video:vid-1").Return(data, nil)
	// repo must NOT be called on cache hit

	uc := NewGetVideo(repo, cache)

	// Act
	result, err := uc.Execute(context.Background(), GetVideoCommand{VideoID: "vid-1"})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, video.VideoID("vid-1"), result.ID)
}

func TestGetVideo_Execute_CacheMissPopulatesCache(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	cache := mocks.NewMockCache(t)

	vid := &video.Video{ID: "vid-1", Title: "From DB", Status: video.StatusReady, UploadedAt: time.Now()}
	cache.EXPECT().Get(mock.Anything, "video:vid-1").Return(nil, errors.New("cache miss"))
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	cache.EXPECT().Set(mock.Anything, "video:vid-1", mock.AnythingOfType("[]uint8")).Return(nil)

	uc := NewGetVideo(repo, cache)

	// Act
	result, err := uc.Execute(context.Background(), GetVideoCommand{VideoID: "vid-1"})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "From DB", result.Title)
}

func TestGetVideo_Execute_VideoNotFound(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	cache := mocks.NewMockCache(t)

	cache.EXPECT().Get(mock.Anything, "video:missing").Return(nil, errors.New("cache miss"))
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("missing")).Return(nil, nil)

	uc := NewGetVideo(repo, cache)

	// Act
	result, err := uc.Execute(context.Background(), GetVideoCommand{VideoID: "missing"})

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestGetVideo_Execute_EmptyVideoID(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	cache := mocks.NewMockCache(t)
	uc := NewGetVideo(repo, cache)

	// Act
	_, err := uc.Execute(context.Background(), GetVideoCommand{VideoID: ""})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}
