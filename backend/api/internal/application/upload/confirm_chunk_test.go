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

func makeVideo(chunks []video.Chunk) *video.Video {
	return &video.Video{
		ID:          "vid-1",
		Status:      video.StatusUploading,
		UploadID:    "mpu-1",
		TotalChunks: len(chunks),
		UploadedAt:  time.Now(),
		Chunks:      chunks,
	}
}

func TestConfirmChunk_Execute_MorePartsRemaining(t *testing.T) {
	// Arrange — vid with 3 chunks, none uploaded; confirming chunk 1 → next missing is chunk 2
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)

	vid := makeVideo([]video.Chunk{
		{PartNumber: 1, Uploaded: false},
		{PartNumber: 2, Uploaded: false},
		{PartNumber: 3, Uploaded: false},
	})

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*video.Video")).Return(nil)
	store.EXPECT().
		PresignUploadPart(mock.Anything, "raw/vid-1/original", "mpu-1", 2).
		Return(&ports.PresignedURL{URL: "https://s3.example.com/part2", PartNumber: 2}, nil)

	uc := NewConfirmChunk(repo, store)

	// Act
	result, err := uc.Execute(context.Background(), ConfirmChunkCommand{
		VideoID:    "vid-1",
		PartNumber: 1,
		ETag:       "etag1",
	})

	// Assert
	assert.NoError(t, err)
	assert.False(t, result.Done)
	assert.Equal(t, 2, result.NextPartNumber)
	assert.Equal(t, "https://s3.example.com/part2", result.PresignedURL)
}

func TestConfirmChunk_Execute_LastChunk(t *testing.T) {
	// Arrange — vid with 2 chunks; chunk 1 already uploaded, confirming chunk 2 (last)
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)

	vid := makeVideo([]video.Chunk{
		{PartNumber: 1, Uploaded: true, ETag: "etag1"},
		{PartNumber: 2, Uploaded: false},
	})

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*video.Video")).Return(nil)
	// No PresignUploadPart call expected — all chunks done

	uc := NewConfirmChunk(repo, store)

	// Act
	result, err := uc.Execute(context.Background(), ConfirmChunkCommand{
		VideoID:    "vid-1",
		PartNumber: 2,
		ETag:       "etag2",
	})

	// Assert
	assert.NoError(t, err)
	assert.True(t, result.Done)
	assert.Empty(t, result.PresignedURL)
}

func TestConfirmChunk_Execute_VideoNotFound(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("missing")).Return(nil, nil)

	uc := NewConfirmChunk(repo, store)

	// Act
	_, err := uc.Execute(context.Background(), ConfirmChunkCommand{
		VideoID:    "missing",
		PartNumber: 1,
		ETag:       "etag1",
	})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestConfirmChunk_Execute_RepoError(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(nil, errors.New("db error"))

	uc := NewConfirmChunk(repo, store)

	// Act
	_, err := uc.Execute(context.Background(), ConfirmChunkCommand{VideoID: "vid-1", PartNumber: 1, ETag: "e"})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}
