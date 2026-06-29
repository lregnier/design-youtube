package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lregnier/design-youtube/api/internal/application"
	"github.com/lregnier/design-youtube/api/internal/domain/video"
	"github.com/lregnier/design-youtube/api/gen/mocks"
)

func processingVideo() *video.Video {
	return &video.Video{
		ID:         "vid-1",
		Status:     video.StatusProcessing,
		UploadedAt: time.Now(),
	}
}

func TestVideoStatusService_MarkReady_Success(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	vid := processingVideo()

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	repo.EXPECT().Save(mock.Anything, mock.MatchedBy(func(v *video.Video) bool {
		return v.Status == video.StatusReady &&
			v.ManifestURL == "https://cdn.example.com/manifest.m3u8" &&
			v.ThumbnailURL == "https://cdn.example.com/thumb.jpg"
	})).Return(nil)

	svc := application.NewVideoStatusService(repo)

	// Act
	err := svc.MarkReady(context.Background(), video.VideoProcessingSucceededEvent{
		VideoID:      "vid-1",
		ManifestURL:  "https://cdn.example.com/manifest.m3u8",
		ThumbnailURL: "https://cdn.example.com/thumb.jpg",
	})

	// Assert
	assert.NoError(t, err)
}

func TestVideoStatusService_MarkReady_IdempotentOnReadyVideo(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	vid := &video.Video{
		ID:           "vid-1",
		Status:       video.StatusReady,
		ManifestURL:  "https://cdn.example.com/old.m3u8",
		ThumbnailURL: "https://cdn.example.com/old.jpg",
		UploadedAt:   time.Now(),
	}

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*video.Video")).Return(nil)

	svc := application.NewVideoStatusService(repo)

	// Act
	err := svc.MarkReady(context.Background(), video.VideoProcessingSucceededEvent{
		VideoID:      "vid-1",
		ManifestURL:  "https://cdn.example.com/new.m3u8",
		ThumbnailURL: "https://cdn.example.com/new.jpg",
	})

	// Assert
	assert.NoError(t, err)
}

func TestVideoStatusService_MarkReady_VideoNotFound(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("missing")).Return(nil, nil)

	svc := application.NewVideoStatusService(repo)

	// Act
	err := svc.MarkReady(context.Background(), video.VideoProcessingSucceededEvent{VideoID: "missing"})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestVideoStatusService_MarkFailed_Success(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	vid := processingVideo()

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	repo.EXPECT().Save(mock.Anything, mock.MatchedBy(func(v *video.Video) bool {
		return v.Status == video.StatusFailed
	})).Return(nil)

	svc := application.NewVideoStatusService(repo)

	// Act
	err := svc.MarkFailed(context.Background(), video.VideoProcessingFailedEvent{
		VideoID: "vid-1",
		Reason:  "ffmpeg decode error",
	})

	// Assert
	assert.NoError(t, err)
}

func TestVideoStatusService_MarkFailed_VideoNotFound(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("missing")).Return(nil, nil)

	svc := application.NewVideoStatusService(repo)

	// Act
	err := svc.MarkFailed(context.Background(), video.VideoProcessingFailedEvent{VideoID: "missing", Reason: "error"})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestVideoStatusService_MarkReady_RepoSaveError(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	vid := processingVideo()

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("db error"))

	svc := application.NewVideoStatusService(repo)

	// Act
	err := svc.MarkReady(context.Background(), video.VideoProcessingSucceededEvent{
		VideoID:      "vid-1",
		ManifestURL:  "https://cdn.example.com/manifest.m3u8",
		ThumbnailURL: "https://cdn.example.com/thumb.jpg",
	})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}
