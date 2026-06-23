package application

import (
	"context"
	"errors"
	"strings"
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

func TestProcessVideo_Execute_DurationFailure(t *testing.T) {
	// Arrange
	storage := mocks.NewMockVideoStorage(t)
	transcoder := mocks.NewMockTranscoder(t)
	publisher := mocks.NewMockResultPublisher(t)

	storage.EXPECT().DownloadRaw(mock.Anything, "vid-1", mock.Anything).Return(nil)
	transcoder.EXPECT().Duration(mock.Anything, mock.Anything).Return(0.0, errors.New("ffprobe error"))
	publisher.EXPECT().PublishFailed(mock.Anything, "vid-1", mock.MatchedBy(func(r string) bool {
		return r != ""
	})).Return(nil)

	uc := NewProcessVideo(storage, transcoder, publisher)

	// Act
	err := uc.Execute(context.Background(), testJob)

	// Assert
	assert.NoError(t, err)
}

func TestProcessVideo_Execute_TranscodeFailure720p(t *testing.T) {
	// Arrange
	storage := mocks.NewMockVideoStorage(t)
	transcoder := mocks.NewMockTranscoder(t)
	publisher := mocks.NewMockResultPublisher(t)

	storage.EXPECT().DownloadRaw(mock.Anything, "vid-1", mock.Anything).Return(nil)
	transcoder.EXPECT().Duration(mock.Anything, mock.Anything).Return(30.0, nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "1920:1080", "4000k").Return(nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "1280:720", "2500k").Return(errors.New("ffmpeg error"))
	publisher.EXPECT().PublishFailed(mock.Anything, "vid-1", mock.MatchedBy(func(r string) bool {
		return r != ""
	})).Return(nil)

	uc := NewProcessVideo(storage, transcoder, publisher)

	// Act
	err := uc.Execute(context.Background(), testJob)

	// Assert
	assert.NoError(t, err)
}

func TestProcessVideo_Execute_TranscodeFailure360p(t *testing.T) {
	// Arrange
	storage := mocks.NewMockVideoStorage(t)
	transcoder := mocks.NewMockTranscoder(t)
	publisher := mocks.NewMockResultPublisher(t)

	storage.EXPECT().DownloadRaw(mock.Anything, "vid-1", mock.Anything).Return(nil)
	transcoder.EXPECT().Duration(mock.Anything, mock.Anything).Return(30.0, nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "1920:1080", "4000k").Return(nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "1280:720", "2500k").Return(nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "640:360", "800k").Return(errors.New("ffmpeg error"))
	publisher.EXPECT().PublishFailed(mock.Anything, "vid-1", mock.MatchedBy(func(r string) bool {
		return r != ""
	})).Return(nil)

	uc := NewProcessVideo(storage, transcoder, publisher)

	// Act
	err := uc.Execute(context.Background(), testJob)

	// Assert
	assert.NoError(t, err)
}

func TestProcessVideo_Execute_UploadSegmentsError(t *testing.T) {
	// Arrange
	storage := mocks.NewMockVideoStorage(t)
	transcoder := mocks.NewMockTranscoder(t)
	publisher := mocks.NewMockResultPublisher(t)

	storage.EXPECT().DownloadRaw(mock.Anything, "vid-1", mock.Anything).Return(nil)
	transcoder.EXPECT().Duration(mock.Anything, mock.Anything).Return(30.0, nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "1920:1080", "4000k").Return(nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "1280:720", "2500k").Return(nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "640:360", "800k").Return(nil)
	storage.EXPECT().UploadSegments(mock.Anything, "vid-1", mock.Anything).Return(errors.New("s3 write error"))

	uc := NewProcessVideo(storage, transcoder, publisher)

	// Act
	err := uc.Execute(context.Background(), testJob)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "upload segments")
}

func TestProcessVideo_Execute_UploadManifestError(t *testing.T) {
	// Arrange
	storage := mocks.NewMockVideoStorage(t)
	transcoder := mocks.NewMockTranscoder(t)
	publisher := mocks.NewMockResultPublisher(t)

	storage.EXPECT().DownloadRaw(mock.Anything, "vid-1", mock.Anything).Return(nil)
	transcoder.EXPECT().Duration(mock.Anything, mock.Anything).Return(30.0, nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "1920:1080", "4000k").Return(nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "1280:720", "2500k").Return(nil)
	transcoder.EXPECT().TranscodeHLS(mock.Anything, mock.Anything, mock.Anything, "640:360", "800k").Return(nil)
	storage.EXPECT().UploadSegments(mock.Anything, "vid-1", mock.Anything).Return(nil)
	storage.EXPECT().UploadManifest(mock.Anything, "vid-1", mock.Anything).Return("", errors.New("s3 write error"))

	uc := NewProcessVideo(storage, transcoder, publisher)

	// Act
	err := uc.Execute(context.Background(), testJob)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "upload manifest")
}

func TestProcessVideo_Execute_UploadThumbnailFailureNonFatal(t *testing.T) {
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
	transcoder.EXPECT().ExtractThumbnail(mock.Anything, mock.Anything, 30.0).Return([]byte("jpeg-data"), nil)
	storage.EXPECT().UploadThumbnail(mock.Anything, "vid-1", []byte("jpeg-data")).Return("", errors.New("s3 write error"))
	// PublishProcessed still called with empty thumbnail URL — upload failure is non-fatal
	publisher.EXPECT().PublishProcessed(mock.Anything, "vid-1", "https://cdn.example.com/manifest.m3u8", "").Return(nil)

	uc := NewProcessVideo(storage, transcoder, publisher)

	// Act
	err := uc.Execute(context.Background(), testJob)

	// Assert
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

func TestBuildMasterManifest_AllQualities(t *testing.T) {
	// Arrange
	videoID := "vid-1"

	// Act
	result := buildMasterManifest(videoID, qualities)

	// Assert
	assert.True(t, strings.HasPrefix(result, "#EXTM3U\n"))
	assert.Equal(t, 3, strings.Count(result, "#EXT-X-STREAM-INF:"))
	assert.Contains(t, result, "BANDWIDTH=4500000,RESOLUTION=1920x1080")
	assert.Contains(t, result, "BANDWIDTH=2800000,RESOLUTION=1280x720")
	assert.Contains(t, result, "BANDWIDTH=1000000,RESOLUTION=640x360")
	assert.Contains(t, result, "../../segments/vid-1/1080p/media.m3u8")
	assert.Contains(t, result, "../../segments/vid-1/720p/media.m3u8")
	assert.Contains(t, result, "../../segments/vid-1/360p/media.m3u8")
}
