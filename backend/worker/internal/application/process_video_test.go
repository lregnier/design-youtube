package application

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lregnier/design-youtube/worker/internal/domain/processing"
	"github.com/lregnier/design-youtube/worker/internal/mocks"
)

var testJob = processing.ProcessingJob{VideoID: "vid-1", S3Key: "raw/vid-1/original"}

func TestProcessVideo_Execute_SuccessfulPipeline(t *testing.T) {
	// Arrange
	storage := mocks.NewMockVideoStorage(t)
	transcoder := mocks.NewMockTranscoder(t)
	publisher := mocks.NewMockResultPublisher(t)

	storage.EXPECT().DownloadRaw(mock.Anything, "vid-1", mock.AnythingOfType("string")).Return(nil)
	transcoder.EXPECT().Duration(mock.Anything, mock.AnythingOfType("string")).Return(60.0, nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "1920:1080", "4000k").Return(nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "1280:720", "2500k").Return(nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "640:360", "800k").Return(nil)
	storage.EXPECT().UploadSegments(mock.Anything, "vid-1", mock.AnythingOfType("string")).Return(nil)
	storage.EXPECT().UploadManifest(mock.Anything, "vid-1", mock.AnythingOfType("[]uint8")).Return("https://cdn.example.com/manifest.m3u8", nil)
	transcoder.EXPECT().ExtractThumbnail(mock.Anything, mock.Anything, 30.0).Return([]byte("jpeg-data"), nil)
	storage.EXPECT().UploadThumbnail(mock.Anything, "vid-1", []byte("jpeg-data")).Return("https://cdn.example.com/thumb.jpg", nil)
	publisher.EXPECT().PublishProcessed(mock.Anything, "vid-1", "https://cdn.example.com/manifest.m3u8", "https://cdn.example.com/thumb.jpg").Return(nil)

	uc := NewProcessVideo(storage, transcoder, publisher)

	// Act
	err := uc.Execute(context.Background(), testJob)

	// Assert
	assert.NoError(t, err)
}

func TestProcessVideo_Execute_DownloadFailure(t *testing.T) {
	// Arrange
	storage := mocks.NewMockVideoStorage(t)
	transcoder := mocks.NewMockTranscoder(t)
	publisher := mocks.NewMockResultPublisher(t)

	storage.EXPECT().DownloadRaw(mock.Anything, "vid-1", mock.Anything).Return(errors.New("s3 not found"))
	publisher.EXPECT().PublishFailed(mock.Anything, "vid-1", mock.MatchedBy(func(r string) bool {
		return r != ""
	})).Return(nil)

	uc := NewProcessVideo(storage, transcoder, publisher)

	// Act
	err := uc.Execute(context.Background(), testJob)

	// Assert — error is swallowed (failure published as event), no error returned
	assert.NoError(t, err)
}

func TestProcessVideo_Execute_TranscodeFailure(t *testing.T) {
	// Arrange
	storage := mocks.NewMockVideoStorage(t)
	transcoder := mocks.NewMockTranscoder(t)
	publisher := mocks.NewMockResultPublisher(t)

	storage.EXPECT().DownloadRaw(mock.Anything, "vid-1", mock.Anything).Return(nil)
	transcoder.EXPECT().Duration(mock.Anything, mock.Anything).Return(30.0, nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "1920:1080", "4000k").Return(errors.New("ffmpeg error"))
	publisher.EXPECT().PublishFailed(mock.Anything, "vid-1", mock.MatchedBy(func(r string) bool {
		return r != ""
	})).Return(nil)

	uc := NewProcessVideo(storage, transcoder, publisher)

	// Act
	err := uc.Execute(context.Background(), testJob)

	// Assert
	assert.NoError(t, err)
}

func TestProcessVideo_Execute_ThumbnailFailureNonFatal(t *testing.T) {
	// Arrange
	storage := mocks.NewMockVideoStorage(t)
	transcoder := mocks.NewMockTranscoder(t)
	publisher := mocks.NewMockResultPublisher(t)

	storage.EXPECT().DownloadRaw(mock.Anything, "vid-1", mock.Anything).Return(nil)
	transcoder.EXPECT().Duration(mock.Anything, mock.Anything).Return(60.0, nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "1920:1080", "4000k").Return(nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "1280:720", "2500k").Return(nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "640:360", "800k").Return(nil)
	storage.EXPECT().UploadSegments(mock.Anything, "vid-1", mock.Anything).Return(nil)
	storage.EXPECT().UploadManifest(mock.Anything, "vid-1", mock.Anything).Return("https://cdn.example.com/manifest.m3u8", nil)
	transcoder.EXPECT().ExtractThumbnail(mock.Anything, mock.Anything, 30.0).Return(nil, errors.New("ffmpeg thumbnail error"))
	// Thumbnail upload NOT called — but PublishProcessed still called with empty thumbnail
	publisher.EXPECT().PublishProcessed(mock.Anything, "vid-1", "https://cdn.example.com/manifest.m3u8", "").Return(nil)

	uc := NewProcessVideo(storage, transcoder, publisher)

	// Act
	err := uc.Execute(context.Background(), testJob)

	// Assert — thumbnail failure is non-fatal
	assert.NoError(t, err)
}

func TestProcessVideo_Execute_PublishProcessedError(t *testing.T) {
	// Arrange
	storage := mocks.NewMockVideoStorage(t)
	transcoder := mocks.NewMockTranscoder(t)
	publisher := mocks.NewMockResultPublisher(t)

	storage.EXPECT().DownloadRaw(mock.Anything, "vid-1", mock.Anything).Return(nil)
	transcoder.EXPECT().Duration(mock.Anything, mock.Anything).Return(10.0, nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "1920:1080", "4000k").Return(nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "1280:720", "2500k").Return(nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "640:360", "800k").Return(nil)
	storage.EXPECT().UploadSegments(mock.Anything, "vid-1", mock.Anything).Return(nil)
	storage.EXPECT().UploadManifest(mock.Anything, "vid-1", mock.Anything).Return("https://cdn.example.com/manifest.m3u8", nil)
	transcoder.EXPECT().ExtractThumbnail(mock.Anything, mock.Anything, mock.AnythingOfType("float64")).Return(nil, errors.New("no thumb"))
	publisher.EXPECT().PublishProcessed(mock.Anything, "vid-1", mock.Anything, mock.Anything).Return(errors.New("sqs unavailable"))

	uc := NewProcessVideo(storage, transcoder, publisher)

	// Act
	err := uc.Execute(context.Background(), testJob)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "publish processed")
}
