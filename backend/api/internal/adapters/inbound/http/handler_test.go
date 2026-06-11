package http

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lregnier/design-youtube/api/internal/api"
	"github.com/lregnier/design-youtube/api/internal/application/catalog"
	"github.com/lregnier/design-youtube/api/internal/application/upload"
	"github.com/lregnier/design-youtube/api/internal/domain/video"
	"github.com/lregnier/design-youtube/api/internal/mocks"
	"github.com/lregnier/design-youtube/api/internal/ports"
)

func newTestHandler(
	repo *mocks.MockVideoRepository,
	store *mocks.MockObjectStore,
	cache *mocks.MockCache,
	queue *mocks.MockQueue,
) *Handler {
	return NewHandler(
		upload.NewInitUpload(repo, store, "my-bucket"),
		upload.NewConfirmChunk(repo, store),
		upload.NewCompleteUpload(repo, store, queue),
		catalog.NewGetVideo(repo, cache),
		catalog.NewListVideos(repo),
	)
}

func TestHandler_GetVideos_Success(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	cache := mocks.NewMockCache(t)
	queue := mocks.NewMockQueue(t)

	uploadedAt := time.Now().UTC().Truncate(time.Second)
	repo.EXPECT().ListReady(mock.Anything).Return([]*video.Video{
		{ID: "vid-1", Title: "First", Status: video.StatusReady, ThumbnailURL: "https://cdn.example.com/vid-1/thumb.jpg", UploadedAt: uploadedAt},
	}, nil)

	h := newTestHandler(repo, store, cache, queue)

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
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	cache := mocks.NewMockCache(t)
	queue := mocks.NewMockQueue(t)

	repo.EXPECT().ListReady(mock.Anything).Return(nil, errors.New("dynamodb unavailable"))

	h := newTestHandler(repo, store, cache, queue)

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
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	cache := mocks.NewMockCache(t)
	queue := mocks.NewMockQueue(t)

	vid := &video.Video{
		ID:           "vid-1",
		Title:        "My Video",
		Description:  "A test video",
		Status:       video.StatusReady,
		UploadedAt:   time.Now(),
		ManifestURL:  "https://cdn.example.com/vid-1/master.m3u8",
		ThumbnailURL: "https://cdn.example.com/vid-1/thumb.jpg",
	}
	data, _ := json.Marshal(vid)
	cache.EXPECT().Get(mock.Anything, "video:vid-1").Return(data, nil)

	h := newTestHandler(repo, store, cache, queue)

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
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	cache := mocks.NewMockCache(t)
	queue := mocks.NewMockQueue(t)

	cache.EXPECT().Get(mock.Anything, "video:missing").Return(nil, errors.New("cache miss"))
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("missing")).Return(nil, nil)

	h := newTestHandler(repo, store, cache, queue)

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
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	cache := mocks.NewMockCache(t)
	queue := mocks.NewMockQueue(t)

	store.EXPECT().
		CreateMultipartUpload(mock.Anything, mock.AnythingOfType("string")).
		Return(&ports.MultipartUpload{UploadID: "mpu-123", Key: "raw/vid/original"}, nil)
	store.EXPECT().
		PresignUploadPart(mock.Anything, mock.AnythingOfType("string"), "mpu-123", 1).
		Return(&ports.PresignedURL{URL: "https://s3.example.com/part1", PartNumber: 1}, nil)
	repo.EXPECT().
		Save(mock.Anything, mock.AnythingOfType("*video.Video")).
		Return(nil)

	h := newTestHandler(repo, store, cache, queue)

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
	assert.NotEmpty(t, result.VideoId)
	assert.Equal(t, "mpu-123", result.UploadId)
	assert.Equal(t, 1, result.NextPartNumber)
	assert.Equal(t, "https://s3.example.com/part1", result.PresignedUrl)
	assert.Len(t, result.Chunks, 1)
}

func TestHandler_InitUpload_FileTooLarge(t *testing.T) {
	// Arrange
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	cache := mocks.NewMockCache(t)
	queue := mocks.NewMockQueue(t)

	h := newTestHandler(repo, store, cache, queue)

	body := api.UploadInitRequest{
		Title:       "Big Video",
		FileSize:    upload.MaxFileSize + 1,
		TotalChunks: 11,
	}

	// Act
	resp, err := h.InitUpload(context.Background(), api.InitUploadRequestObject{Body: &body})

	// Assert
	assert.NoError(t, err)
	errResp, ok := resp.(api.InitUpload400JSONResponse)
	assert.True(t, ok)
	assert.Contains(t, errResp.Error, "100MB")
	// repo and store must not be called — verified by mockery expecter
}

func TestHandler_ConfirmChunk_Success(t *testing.T) {
	// Arrange — last remaining chunk confirmed → done
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	cache := mocks.NewMockCache(t)
	queue := mocks.NewMockQueue(t)

	vid := &video.Video{
		ID:          "vid-1",
		Status:      video.StatusUploading,
		UploadID:    "mpu-1",
		TotalChunks: 1,
		UploadedAt:  time.Now(),
		Chunks:      []video.Chunk{{PartNumber: 1, Uploaded: false}},
	}
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*video.Video")).Return(nil)

	h := newTestHandler(repo, store, cache, queue)

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
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	cache := mocks.NewMockCache(t)
	queue := mocks.NewMockQueue(t)

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("missing")).Return(nil, nil)

	h := newTestHandler(repo, store, cache, queue)

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
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	cache := mocks.NewMockCache(t)
	queue := mocks.NewMockQueue(t)

	vid := &video.Video{
		ID:          "vid-1",
		Status:      video.StatusUploading,
		UploadID:    "mpu-1",
		TotalChunks: 1,
		UploadedAt:  time.Now(),
		Chunks:      []video.Chunk{{PartNumber: 1, Uploaded: true, ETag: "etag1"}},
	}
	repo.EXPECT().FindByID(mock.Anything, video.VideoID("vid-1")).Return(vid, nil)
	store.EXPECT().
		CompleteMultipartUpload(mock.Anything, "raw/vid-1/original", "mpu-1", mock.AnythingOfType("[]ports.CompletedPart")).
		Return(nil)
	repo.EXPECT().Save(mock.Anything, mock.MatchedBy(func(v *video.Video) bool {
		return v.Status == video.StatusProcessing
	})).Return(nil)
	queue.EXPECT().
		SendMessage(mock.Anything, `{"videoId":"vid-1","s3Key":"raw/vid-1/original"}`, "vid-1").
		Return(nil)

	h := newTestHandler(repo, store, cache, queue)

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
	repo := mocks.NewMockVideoRepository(t)
	store := mocks.NewMockObjectStore(t)
	cache := mocks.NewMockCache(t)
	queue := mocks.NewMockQueue(t)

	repo.EXPECT().FindByID(mock.Anything, video.VideoID("missing")).Return(nil, nil)

	h := newTestHandler(repo, store, cache, queue)

	body := api.CompleteUploadRequest{UploadId: "mpu-1"}

	// Act
	resp, err := h.CompleteUpload(context.Background(), api.CompleteUploadRequestObject{VideoId: "missing", Body: &body})

	// Assert
	assert.NoError(t, err)
	errResp, ok := resp.(api.CompleteUpload404JSONResponse)
	assert.True(t, ok)
	assert.Equal(t, "video not found", errResp.Error)
}
