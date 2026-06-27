package application_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lregnier/design-youtube/api/internal/application"
	"github.com/lregnier/design-youtube/api/internal/domain/video"
	"github.com/lregnier/design-youtube/api/gen/mocks"
)

func TestCatalogService_GetVideo_CacheHit(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	cache := mocks.NewMockCache(t)

	vid := &video.Video{ID: "vid-1", Title: "Cached", Status: video.StatusReady, UploadedAt: time.Now()}
	data, _ := json.Marshal(vid)
	cache.EXPECT().Get(mock.Anything, "video:vid-1").Return(data, nil)

	svc := application.NewCatalogService(repo, cache)

	// Act
	result, err := svc.GetVideo(context.Background(), "vid-1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, video.VideoID("vid-1"), result.ID)
}

func TestCatalogService_GetVideo_CacheMissPopulatesCache(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	cache := mocks.NewMockCache(t)

	vid := &video.Video{ID: "vid-1", Title: "From DB", Status: video.StatusReady, UploadedAt: time.Now()}
	cache.EXPECT().Get(mock.Anything, "video:vid-1").Return(nil, errors.New("cache miss"))
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	cache.EXPECT().Set(mock.Anything, "video:vid-1", mock.AnythingOfType("[]uint8")).Return(nil)

	svc := application.NewCatalogService(repo, cache)

	// Act
	result, err := svc.GetVideo(context.Background(), "vid-1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "From DB", result.Title)
}

func TestCatalogService_GetVideo_VideoNotFound(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	cache := mocks.NewMockCache(t)

	cache.EXPECT().Get(mock.Anything, "video:missing").Return(nil, errors.New("cache miss"))
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("missing")).Return(nil, nil)

	svc := application.NewCatalogService(repo, cache)

	// Act
	result, err := svc.GetVideo(context.Background(), "missing")

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestCatalogService_GetVideo_EmptyVideoID(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	cache := mocks.NewMockCache(t)
	svc := application.NewCatalogService(repo, cache)

	// Act
	_, err := svc.GetVideo(context.Background(), "")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}

func TestCatalogService_ListVideos_ReturnsList(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	cache := mocks.NewMockCache(t)

	videos := []*video.Video{
		{ID: "vid-1", Title: "First", Status: video.StatusReady, UploadedAt: time.Now()},
		{ID: "vid-2", Title: "Second", Status: video.StatusReady, UploadedAt: time.Now().Add(-time.Hour)},
	}
	repo.EXPECT().List(mock.Anything).Return(videos, nil)

	svc := application.NewCatalogService(repo, cache)

	// Act
	result, err := svc.ListVideos(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, video.VideoID("vid-1"), result[0].ID)
}

func TestCatalogService_ListVideos_EmptyList(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	cache := mocks.NewMockCache(t)
	repo.EXPECT().List(mock.Anything).Return([]*video.Video{}, nil)

	svc := application.NewCatalogService(repo, cache)

	// Act
	result, err := svc.ListVideos(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}
