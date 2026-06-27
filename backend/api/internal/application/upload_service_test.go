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
	"github.com/lregnier/design-youtube/api/internal/gen/mocks"
)

// helpers

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

// InitUpload tests

func TestUploadService_InitUpload_ValidRequest(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	publisher := mocks.NewMockEventPublisher(t)

	store.EXPECT().
		CreateMultipartUpload(mock.Anything, mock.AnythingOfType("string")).
		Return(&application.MultipartUpload{UploadID: "mpu-123", Key: "raw/vid/original"}, nil)
	store.EXPECT().
		PresignUploadPart(mock.Anything, mock.AnythingOfType("string"), "mpu-123", 1).
		Return(&application.PresignedURL{URL: "https://s3.example.com/part1", PartNumber: 1}, nil)
	repo.EXPECT().
		Save(mock.Anything, mock.AnythingOfType("*video.Video")).
		Return(nil)

	svc := application.NewUploadService(repo, store, publisher, "my-bucket")

	// Act
	result, err := svc.InitUpload(context.Background(), "", "My Video", "A test video", 10*1024*1024, 1)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result.VideoID)
	assert.Equal(t, "mpu-123", result.UploadID)
	assert.Equal(t, 1, result.NextPartNumber)
	assert.Equal(t, "https://s3.example.com/part1", result.PresignedURL)
	assert.Len(t, result.Chunks, 1)
}

func TestUploadService_InitUpload_FileTooLarge(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	publisher := mocks.NewMockEventPublisher(t)
	svc := application.NewUploadService(repo, store, publisher, "my-bucket")

	// Act
	result, err := svc.InitUpload(context.Background(), "", "Big Video", "", application.MaxFileSize+1, 11)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "100MB")
	assert.Empty(t, result.VideoID)
}

func TestUploadService_InitUpload_ResumeExistingUpload(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	publisher := mocks.NewMockEventPublisher(t)

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
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("existing-id")).Return(existingVideo, nil)
	store.EXPECT().
		PresignUploadPart(mock.Anything, "raw/existing-id/original", "existing-mpu", 2).
		Return(&application.PresignedURL{URL: "https://s3.example.com/part2", PartNumber: 2}, nil)

	svc := application.NewUploadService(repo, store, publisher, "my-bucket")

	// Act
	result, err := svc.InitUpload(context.Background(), "existing-id", "My Video", "", 20*1024*1024, 2)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "existing-id", result.VideoID)
	assert.Equal(t, 2, result.NextPartNumber)
	assert.Equal(t, "https://s3.example.com/part2", result.PresignedURL)
}

func TestUploadService_InitUpload_RepoSaveError(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	publisher := mocks.NewMockEventPublisher(t)

	store.EXPECT().
		CreateMultipartUpload(mock.Anything, mock.AnythingOfType("string")).
		Return(&application.MultipartUpload{UploadID: "mpu-123", Key: "raw/vid/original"}, nil)
	repo.EXPECT().
		Save(mock.Anything, mock.AnythingOfType("*video.Video")).
		Return(errors.New("dynamodb unavailable"))

	svc := application.NewUploadService(repo, store, publisher, "my-bucket")

	// Act
	_, err := svc.InitUpload(context.Background(), "", "My Video", "", 5*1024*1024, 1)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dynamodb unavailable")
}

// ConfirmChunk tests

func TestUploadService_ConfirmChunk_MorePartsRemaining(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	publisher := mocks.NewMockEventPublisher(t)

	vid := makeVideo([]video.Chunk{
		{PartNumber: 1, Uploaded: false},
		{PartNumber: 2, Uploaded: false},
		{PartNumber: 3, Uploaded: false},
	})
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*video.Video")).Return(nil)
	store.EXPECT().
		PresignUploadPart(mock.Anything, "raw/vid-1/original", "mpu-1", 2).
		Return(&application.PresignedURL{URL: "https://s3.example.com/part2", PartNumber: 2}, nil)

	svc := application.NewUploadService(repo, store, publisher, "my-bucket")

	// Act
	result, err := svc.ConfirmChunk(context.Background(), "vid-1", 1, "etag1")

	// Assert
	assert.NoError(t, err)
	assert.False(t, result.Done)
	assert.Equal(t, 2, result.NextPartNumber)
	assert.Equal(t, "https://s3.example.com/part2", result.PresignedURL)
}

func TestUploadService_ConfirmChunk_LastChunk(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	publisher := mocks.NewMockEventPublisher(t)

	vid := makeVideo([]video.Chunk{
		{PartNumber: 1, Uploaded: true, ETag: "etag1"},
		{PartNumber: 2, Uploaded: false},
	})
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*video.Video")).Return(nil)

	svc := application.NewUploadService(repo, store, publisher, "my-bucket")

	// Act
	result, err := svc.ConfirmChunk(context.Background(), "vid-1", 2, "etag2")

	// Assert
	assert.NoError(t, err)
	assert.True(t, result.Done)
	assert.Empty(t, result.PresignedURL)
}

func TestUploadService_ConfirmChunk_VideoNotFound(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	publisher := mocks.NewMockEventPublisher(t)

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("missing")).Return(nil, nil)

	svc := application.NewUploadService(repo, store, publisher, "my-bucket")

	// Act
	_, err := svc.ConfirmChunk(context.Background(), "missing", 1, "etag1")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestUploadService_ConfirmChunk_RepoError(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	publisher := mocks.NewMockEventPublisher(t)

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(nil, errors.New("db error"))

	svc := application.NewUploadService(repo, store, publisher, "my-bucket")

	// Act
	_, err := svc.ConfirmChunk(context.Background(), "vid-1", 1, "e")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

// CompleteUpload tests

func TestUploadService_CompleteUpload_Success(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	publisher := mocks.NewMockEventPublisher(t)

	vid := videoWithChunks(2, true)
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	store.EXPECT().
		CompleteMultipartUpload(mock.Anything, "raw/vid-1/original", "mpu-1", mock.AnythingOfType("[]application.CompletedPart")).
		Return(nil)
	repo.EXPECT().Save(mock.Anything, mock.MatchedBy(func(v *video.Video) bool {
		return v.Status == video.StatusProcessing
	})).Return(nil)
	publisher.EXPECT().
		Publish(mock.Anything, video.VideoUploadedEvent{VideoID: "vid-1", S3Key: "raw/vid-1/original"}).
		Return(nil)

	svc := application.NewUploadService(repo, store, publisher, "my-bucket")

	// Act
	err := svc.CompleteUpload(context.Background(), "vid-1", "mpu-1")

	// Assert
	assert.NoError(t, err)
}

func TestUploadService_CompleteUpload_VideoNotFound(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	publisher := mocks.NewMockEventPublisher(t)

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("missing")).Return(nil, nil)
	svc := application.NewUploadService(repo, store, publisher, "my-bucket")

	// Act
	err := svc.CompleteUpload(context.Background(), "missing", "mpu")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestUploadService_CompleteUpload_S3Error(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	publisher := mocks.NewMockEventPublisher(t)

	vid := videoWithChunks(1, true)
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	store.EXPECT().
		CompleteMultipartUpload(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("s3 error"))

	svc := application.NewUploadService(repo, store, publisher, "my-bucket")

	// Act
	err := svc.CompleteUpload(context.Background(), "vid-1", "mpu-1")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "s3 error")
}

func TestUploadService_CompleteUpload_PublishError(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	publisher := mocks.NewMockEventPublisher(t)

	vid := videoWithChunks(1, true)
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	store.EXPECT().
		CompleteMultipartUpload(mock.Anything, "raw/vid-1/original", "mpu-1", mock.AnythingOfType("[]application.CompletedPart")).
		Return(nil)
	repo.EXPECT().Save(mock.Anything, mock.MatchedBy(func(v *video.Video) bool {
		return v.Status == video.StatusProcessing
	})).Return(nil)
	publisher.EXPECT().Publish(mock.Anything, mock.Anything).Return(errors.New("sqs error"))

	svc := application.NewUploadService(repo, store, publisher, "my-bucket")

	// Act
	err := svc.CompleteUpload(context.Background(), "vid-1", "mpu-1")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sqs error")
}
