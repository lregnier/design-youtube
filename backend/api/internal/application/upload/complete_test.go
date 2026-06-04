package upload

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
	"github.com/lregnier/design-youtube/api/internal/mocks"
	"github.com/lregnier/design-youtube/api/internal/ports"
)

func TestCompleteUpload_Execute_Success(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)

	vid := videoWithChunks(2, true)
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	store.EXPECT().
		CompleteMultipartUpload(mock.Anything, "raw/vid-1/original", "mpu-1", mock.AnythingOfType("[]ports.CompletedPart")).
		Return(nil)
	repo.EXPECT().Save(mock.Anything, mock.MatchedBy(func(v *video.Video) bool {
		return v.Status == video.StatusProcessing
	})).Return(nil)

	uc := NewCompleteUpload(repo, store)

	// Act
	err := uc.Execute(context.Background(), CompleteUploadCommand{
		VideoID:  "vid-1",
		UploadID: "mpu-1",
	})

	// Assert
	assert.NoError(t, err)
}

func TestCompleteUpload_Execute_VideoNotFound(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("missing")).Return(nil, nil)
	uc := NewCompleteUpload(repo, store)

	// Act
	err := uc.Execute(context.Background(), CompleteUploadCommand{VideoID: "missing", UploadID: "mpu"})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestCompleteUpload_Execute_S3Error(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)

	vid := videoWithChunks(1, true)
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	store.EXPECT().
		CompleteMultipartUpload(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("s3 error"))

	uc := NewCompleteUpload(repo, store)

	// Act
	err := uc.Execute(context.Background(), CompleteUploadCommand{VideoID: "vid-1", UploadID: "mpu-1"})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "s3 error")
}

// ensure CompleteMultipartUpload receives correctly typed slice
var _ ports.CompletedPart = ports.CompletedPart{}
