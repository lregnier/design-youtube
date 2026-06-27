package http

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lregnier/design-youtube/api/internal/application"
	"github.com/lregnier/design-youtube/api/internal/domain/video"
	"github.com/lregnier/design-youtube/api/gen/api"
	"github.com/lregnier/design-youtube/api/gen/mocks"
)

func newTestHandler(
	uploadSvc *mocks.MockUploadService,
	catalogSvc *mocks.MockCatalogService,
) *Handler {
	return NewHandler(uploadSvc, catalogSvc)
}

func TestHandler_GetVideos_Success(t *testing.T) {
	// Arrange
	uploadSvc := mocks.NewMockUploadService(t)
	catalogSvc := mocks.NewMockCatalogService(t)

	uploadedAt := time.Now().UTC().Truncate(time.Second)
	catalogSvc.EXPECT().ListVideos(mock.Anything).Return([]*video.Video{
		{ID: "vid-1", Title: "First", Status: video.StatusReady, ThumbnailURL: "https://cdn.example.com/vid-1/thumb.jpg", UploadedAt: uploadedAt},
	}, nil)

	h := newTestHandler(uploadSvc, catalogSvc)

	// Act
	resp, err := h.GetVideos(context.Background(), api.GetVideosRequestObject{})

	// Assert
	assert.NoError(t, err)
	summaries, ok := resp.(api.GetVideos200JSONResponse)
	assert.True(t, ok)
	assert.Len(t, summaries, 1)
	assert.Equal(t, "vid-1", summaries[0].VideoId)
	assert.Equal(t, "https://cdn.example.com/vid-1/thumb.jpg", summaries[0].ThumbnailUrl)
	assert.Equal(t, api.VideoStatus(video.StatusReady), summaries[0].Status)
}

func TestHandler_GetVideos_Error(t *testing.T) {
	// Arrange
	uploadSvc := mocks.NewMockUploadService(t)
	catalogSvc := mocks.NewMockCatalogService(t)

	catalogSvc.EXPECT().ListVideos(mock.Anything).Return(nil, errors.New("dynamodb unavailable"))

	h := newTestHandler(uploadSvc, catalogSvc)

	// Act
	resp, err := h.GetVideos(context.Background(), api.GetVideosRequestObject{})

	// Assert
	assert.NoError(t, err)
	errResp, ok := resp.(api.GetVideos500JSONResponse)
	assert.True(t, ok)
	assert.Equal(t, "dynamodb unavailable", errResp.Error)
}

func TestHandler_GetVideo_Found(t *testing.T) {
	// Arrange
	uploadSvc := mocks.NewMockUploadService(t)
	catalogSvc := mocks.NewMockCatalogService(t)

	vid := &video.Video{
		ID:           "vid-1",
		Title:        "My Video",
		Description:  "A test video",
		Status:       video.StatusReady,
		UploadedAt:   time.Now(),
		ManifestURL:  "https://cdn.example.com/vid-1/master.m3u8",
		ThumbnailURL: "https://cdn.example.com/vid-1/thumb.jpg",
	}
	catalogSvc.EXPECT().GetVideo(mock.Anything, "vid-1").Return(vid, nil)

	h := newTestHandler(uploadSvc, catalogSvc)

	// Act
	resp, err := h.GetVideo(context.Background(), api.GetVideoRequestObject{VideoId: "vid-1"})

	// Assert
	assert.NoError(t, err)
	detail, ok := resp.(api.GetVideo200JSONResponse)
	assert.True(t, ok)
	assert.Equal(t, "vid-1", detail.VideoId)
	assert.Equal(t, "My Video", detail.Title)
	assert.Equal(t, "https://cdn.example.com/vid-1/master.m3u8", *detail.ManifestUrl)
	assert.Equal(t, "https://cdn.example.com/vid-1/thumb.jpg", *detail.ThumbnailUrl)
}

func TestHandler_GetVideo_NotFound(t *testing.T) {
	// Arrange
	uploadSvc := mocks.NewMockUploadService(t)
	catalogSvc := mocks.NewMockCatalogService(t)

	catalogSvc.EXPECT().GetVideo(mock.Anything, "missing").Return(nil, nil)

	h := newTestHandler(uploadSvc, catalogSvc)

	// Act
	resp, err := h.GetVideo(context.Background(), api.GetVideoRequestObject{VideoId: "missing"})

	// Assert
	assert.NoError(t, err)
	errResp, ok := resp.(api.GetVideo404JSONResponse)
	assert.True(t, ok)
	assert.Equal(t, "video not found", errResp.Error)
}

func TestHandler_InitUpload_Success(t *testing.T) {
	// Arrange
	uploadSvc := mocks.NewMockUploadService(t)
	catalogSvc := mocks.NewMockCatalogService(t)

	uploadSvc.EXPECT().
		InitUpload(mock.Anything, "", "My Video", "A test video", int64(10*1024*1024), 1).
		Return(application.InitUploadResult{
			VideoID:        "vid-1",
			UploadID:       "mpu-123",
			NextPartNumber: 1,
			PresignedURL:   "https://s3.example.com/part1",
			Chunks:         []application.ChunkState{{PartNumber: 1, Uploaded: false}},
		}, nil)

	h := newTestHandler(uploadSvc, catalogSvc)

	body := api.UploadInitRequest{
		Title:       "My Video",
		Description: "A test video",
		FileSize:    10 * 1024 * 1024,
		TotalChunks: 1,
	}

	// Act
	resp, err := h.InitUpload(context.Background(), api.InitUploadRequestObject{Body: &body})

	// Assert
	assert.NoError(t, err)
	result, ok := resp.(api.InitUpload200JSONResponse)
	assert.True(t, ok)
	assert.Equal(t, "vid-1", result.VideoId)
	assert.Equal(t, "mpu-123", result.UploadId)
	assert.Equal(t, 1, result.NextPartNumber)
	assert.Equal(t, "https://s3.example.com/part1", result.PresignedUrl)
	assert.Len(t, result.Chunks, 1)
}

func TestHandler_InitUpload_Error(t *testing.T) {
	// Arrange
	uploadSvc := mocks.NewMockUploadService(t)
	catalogSvc := mocks.NewMockCatalogService(t)

	uploadSvc.EXPECT().
		InitUpload(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(application.InitUploadResult{}, errors.New("file size 200 exceeds 100MB limit"))

	h := newTestHandler(uploadSvc, catalogSvc)

	body := api.UploadInitRequest{Title: "Big Video", FileSize: 200, TotalChunks: 1}

	// Act
	resp, err := h.InitUpload(context.Background(), api.InitUploadRequestObject{Body: &body})

	// Assert
	assert.NoError(t, err)
	errResp, ok := resp.(api.InitUpload400JSONResponse)
	assert.True(t, ok)
	assert.Contains(t, errResp.Error, "100MB")
}

func TestHandler_ConfirmChunk_Success(t *testing.T) {
	// Arrange
	uploadSvc := mocks.NewMockUploadService(t)
	catalogSvc := mocks.NewMockCatalogService(t)

	uploadSvc.EXPECT().
		ConfirmChunk(mock.Anything, "vid-1", 1, "etag1").
		Return(application.ConfirmChunkResult{PartNumber: 1, Done: true}, nil)

	h := newTestHandler(uploadSvc, catalogSvc)

	body := api.ConfirmChunkRequest{PartNumber: 1, ETag: "etag1"}

	// Act
	resp, err := h.ConfirmChunk(context.Background(), api.ConfirmChunkRequestObject{VideoId: "vid-1", Body: &body})

	// Assert
	assert.NoError(t, err)
	result, ok := resp.(api.ConfirmChunk200JSONResponse)
	assert.True(t, ok)
	assert.True(t, result.Done)
	assert.Nil(t, result.NextPartNumber)
	assert.Nil(t, result.PresignedUrl)
}

func TestHandler_ConfirmChunk_VideoNotFound(t *testing.T) {
	// Arrange
	uploadSvc := mocks.NewMockUploadService(t)
	catalogSvc := mocks.NewMockCatalogService(t)

	uploadSvc.EXPECT().
		ConfirmChunk(mock.Anything, "missing", 1, "etag1").
		Return(application.ConfirmChunkResult{}, errors.New("video not found"))

	h := newTestHandler(uploadSvc, catalogSvc)

	body := api.ConfirmChunkRequest{PartNumber: 1, ETag: "etag1"}

	// Act
	resp, err := h.ConfirmChunk(context.Background(), api.ConfirmChunkRequestObject{VideoId: "missing", Body: &body})

	// Assert
	assert.NoError(t, err)
	errResp, ok := resp.(api.ConfirmChunk404JSONResponse)
	assert.True(t, ok)
	assert.Equal(t, "video not found", errResp.Error)
}

func TestHandler_CompleteUpload_Success(t *testing.T) {
	// Arrange
	uploadSvc := mocks.NewMockUploadService(t)
	catalogSvc := mocks.NewMockCatalogService(t)

	uploadSvc.EXPECT().CompleteUpload(mock.Anything, "vid-1", "mpu-1").Return(nil)

	h := newTestHandler(uploadSvc, catalogSvc)

	body := api.CompleteUploadRequest{UploadId: "mpu-1"}

	// Act
	resp, err := h.CompleteUpload(context.Background(), api.CompleteUploadRequestObject{VideoId: "vid-1", Body: &body})

	// Assert
	assert.NoError(t, err)
	_, ok := resp.(api.CompleteUpload200Response)
	assert.True(t, ok)
}

func TestHandler_CompleteUpload_VideoNotFound(t *testing.T) {
	// Arrange
	uploadSvc := mocks.NewMockUploadService(t)
	catalogSvc := mocks.NewMockCatalogService(t)

	uploadSvc.EXPECT().CompleteUpload(mock.Anything, "missing", "mpu-1").Return(errors.New("video not found"))

	h := newTestHandler(uploadSvc, catalogSvc)

	body := api.CompleteUploadRequest{UploadId: "mpu-1"}

	// Act
	resp, err := h.CompleteUpload(context.Background(), api.CompleteUploadRequestObject{VideoId: "missing", Body: &body})

	// Assert
	assert.NoError(t, err)
	errResp, ok := resp.(api.CompleteUpload404JSONResponse)
	assert.True(t, ok)
	assert.Equal(t, "video not found", errResp.Error)
}
