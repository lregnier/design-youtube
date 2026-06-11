package upload

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
	"github.com/lregnier/design-youtube/api/internal/gen/mocks"
	"github.com/lregnier/design-youtube/api/internal/ports"
)

func TestInitUpload_Execute_ValidRequest(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)

	store.EXPECT().
		CreateMultipartUpload(mock.Anything, mock.AnythingOfType("string")).
		Return(&ports.MultipartUpload{UploadID: "mpu-123", Key: "raw/vid/original"}, nil)
	store.EXPECT().
		PresignUploadPart(mock.Anything, mock.AnythingOfType("string"), "mpu-123", 1).
		Return(&ports.PresignedURL{URL: "https://s3.example.com/part1", PartNumber: 1}, nil)
	repo.EXPECT().
		Save(mock.Anything, mock.AnythingOfType("*video.Video")).
		Return(nil)

	uc := NewInitUpload(repo, store, "my-bucket")

	// Act
	result, err := uc.Execute(context.Background(), InitUploadCommand{
		Title:       "My Video",
		Description: "A test video",
		FileSize:    10 * 1024 * 1024, // 10MB
		TotalChunks: 1,
	})

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result.VideoID)
	assert.Equal(t, "mpu-123", result.UploadID)
	assert.Equal(t, 1, result.NextPartNumber)
	assert.Equal(t, "https://s3.example.com/part1", result.PresignedURL)
	assert.Len(t, result.Chunks, 1)
}

func TestInitUpload_Execute_FileTooLarge(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	uc := NewInitUpload(repo, store, "my-bucket")

	// Act
	result, err := uc.Execute(context.Background(), InitUploadCommand{
		Title:       "Big Video",
		FileSize:    MaxFileSize + 1,
		TotalChunks: 11,
	})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "100MB")
	assert.Empty(t, result.VideoID)
	// repo and store must not be called — verified by mockery expecter
}

func TestInitUpload_Execute_ResumeExistingUpload(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)

	existingVideo := &video.Video{
		ID:          "existing-id",
		Status:      video.StatusUploading,
		UploadID:    "existing-mpu",
		TotalChunks: 2,
		Chunks: []video.Chunk{
			{PartNumber: 1, Uploaded: true, ETag: "etag1"},
			{PartNumber: 2, Uploaded: false},
		},
	}
	repo.EXPECT().
		FindByID(mock.Anything, video.VideoID("existing-id")).
		Return(existingVideo, nil)
	store.EXPECT().
		PresignUploadPart(mock.Anything, "raw/existing-id/original", "existing-mpu", 2).
		Return(&ports.PresignedURL{URL: "https://s3.example.com/part2", PartNumber: 2}, nil)

	uc := NewInitUpload(repo, store, "my-bucket")

	// Act
	result, err := uc.Execute(context.Background(), InitUploadCommand{
		VideoID:     "existing-id",
		Title:       "My Video",
		FileSize:    20 * 1024 * 1024,
		TotalChunks: 2,
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "existing-id", result.VideoID)
	assert.Equal(t, 2, result.NextPartNumber)
	assert.Equal(t, "https://s3.example.com/part2", result.PresignedURL)
}

func TestInitUpload_Execute_RepoSaveError(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)

	store.EXPECT().
		CreateMultipartUpload(mock.Anything, mock.AnythingOfType("string")).
		Return(&ports.MultipartUpload{UploadID: "mpu-123", Key: "raw/vid/original"}, nil)
	repo.EXPECT().
		Save(mock.Anything, mock.AnythingOfType("*video.Video")).
		Return(errors.New("dynamodb unavailable"))

	uc := NewInitUpload(repo, store, "my-bucket")

	// Act
	_, err := uc.Execute(context.Background(), InitUploadCommand{
		Title:       "My Video",
		FileSize:    5 * 1024 * 1024,
		TotalChunks: 1,
	})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dynamodb unavailable")
}

// helper to build a video with n chunks, first n-1 uploaded
func videoWithChunks(n int, allUploaded bool) *video.Video {
	v := &video.Video{
		ID:          "vid-1",
		Status:      video.StatusUploading,
		UploadID:    "mpu-1",
		TotalChunks: n,
		UploadedAt:  time.Now(),
	}
	for i := 1; i <= n; i++ {
		v.Chunks = append(v.Chunks, video.Chunk{
			PartNumber: i,
			Uploaded:   allUploaded || i < n,
			ETag:       "etag",
		})
	}
	return v
}
