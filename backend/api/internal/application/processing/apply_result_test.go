package processing

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
	"github.com/lregnier/design-youtube/api/internal/gen/mocks"
)

func readyVideo() *video.Video {
	return &video.Video{
		ID:         "vid-1",
		Status:     video.StatusProcessing,
		UploadedAt: time.Now(),
	}
}

func TestApplyProcessingResult_OnProcessed_Success(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	vid := readyVideo()

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	repo.EXPECT().Save(mock.Anything, mock.MatchedBy(func(v *video.Video) bool {
		return v.Status == video.StatusReady &&
			v.ManifestURL == "https://cdn.example.com/manifest.m3u8" &&
			v.ThumbnailURL == "https://cdn.example.com/thumb.jpg"
	})).Return(nil)

	uc := NewApplyProcessingResult(repo)

	// Act
	err := uc.OnProcessed(context.Background(), VideoProcessedEvent{
		VideoID:      "vid-1",
		ManifestURL:  "https://cdn.example.com/manifest.m3u8",
		ThumbnailURL: "https://cdn.example.com/thumb.jpg",
	})

	// Assert
	assert.NoError(t, err)
}

func TestApplyProcessingResult_OnProcessed_IdempotentOnReadyVideo(t *testing.T) {
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

	uc := NewApplyProcessingResult(repo)

	// Act — applying again with new URLs on an already-ready video
	err := uc.OnProcessed(context.Background(), VideoProcessedEvent{
		VideoID:      "vid-1",
		ManifestURL:  "https://cdn.example.com/new.m3u8",
		ThumbnailURL: "https://cdn.example.com/new.jpg",
	})

	// Assert — no error, operation is idempotent
	assert.NoError(t, err)
}

func TestApplyProcessingResult_OnProcessed_VideoNotFound(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("missing")).Return(nil, nil)

	uc := NewApplyProcessingResult(repo)

	// Act
	err := uc.OnProcessed(context.Background(), VideoProcessedEvent{VideoID: "missing"})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestApplyProcessingResult_OnFailed_Success(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	vid := readyVideo()

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	repo.EXPECT().Save(mock.Anything, mock.MatchedBy(func(v *video.Video) bool {
		return v.Status == video.StatusFailed
	})).Return(nil)

	uc := NewApplyProcessingResult(repo)

	// Act
	err := uc.OnFailed(context.Background(), VideoFailedEvent{
		VideoID: "vid-1",
		Reason:  "ffmpeg decode error",
	})

	// Assert
	assert.NoError(t, err)
}

func TestApplyProcessingResult_OnFailed_VideoNotFound(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("missing")).Return(nil, nil)

	uc := NewApplyProcessingResult(repo)

	// Act
	err := uc.OnFailed(context.Background(), VideoFailedEvent{VideoID: "missing", Reason: "error"})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestApplyProcessingResult_OnProcessed_RepoSaveError(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	vid := readyVideo()

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("db error"))

	uc := NewApplyProcessingResult(repo)

	// Act
	err := uc.OnProcessed(context.Background(), VideoProcessedEvent{
		VideoID:      "vid-1",
		ManifestURL:  "https://cdn.example.com/manifest.m3u8",
		ThumbnailURL: "https://cdn.example.com/thumb.jpg",
	})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}
