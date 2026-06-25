package application

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/lregnier/design-youtube/worker/internal/domain/processing"
	"github.com/lregnier/design-youtube/worker/internal/ports"
)

var qualities = []struct{ name, scale, bitrate string }{
	{"1080p", "1920:1080", "4000k"},
	{"720p", "1280:720", "2500k"},
	{"360p", "640:360", "800k"},
}

type ProcessVideo struct {
	storage   ports.VideoStorage
	transcoder ports.Transcoder
	publisher  ports.ResultPublisher
}

func NewProcessVideo(storage ports.VideoStorage, transcoder ports.Transcoder, publisher ports.ResultPublisher) ProcessVideo {
	return ProcessVideo{storage: storage, transcoder: transcoder, publisher: publisher}
}

func (uc ProcessVideo) Execute(ctx context.Context, job processing.ProcessingJob) error {
	log.Printf("processing videoId=%s", job.VideoID)

	tmpDir, err := os.MkdirTemp("", "video-"+job.VideoID+"-*")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	rawPath := tmpDir + "/original"
	if err := uc.storage.DownloadRaw(ctx, job.VideoID, rawPath); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	duration, err := uc.transcoder.Duration(ctx, rawPath)
	if err != nil {
		return fmt.Errorf("ffprobe failed: %w", err)
	}

	segDir := tmpDir + "/segments"
	if err := os.MkdirAll(segDir, 0755); err != nil {
		return err
	}

	for _, q := range qualities {
		outDir := segDir + "/" + q.name
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return err
		}
		if err := uc.transcoder.TranscodeHLS(ctx, rawPath, outDir, q.scale, q.bitrate); err != nil {
			return fmt.Errorf("transcode %s failed: %w", q.name, err)
		}
	}

	if err := uc.storage.UploadSegments(ctx, job.VideoID, segDir); err != nil {
		return fmt.Errorf("upload segments: %w", err)
	}

	manifestContent := buildMasterManifest(job.VideoID, qualities)
	manifestURL, err := uc.storage.UploadManifest(ctx, job.VideoID, []byte(manifestContent))
	if err != nil {
		return fmt.Errorf("upload manifest: %w", err)
	}

	var thumbnailURL string
	if thumbData, err := uc.transcoder.ExtractThumbnail(ctx, rawPath, duration/2); err != nil {
		log.Printf("thumbnail extraction failed (non-fatal): %v", err)
	} else {
		thumbnailURL, _ = uc.storage.UploadThumbnail(ctx, job.VideoID, thumbData)
	}

	if err := uc.publisher.PublishProcessed(ctx, job.VideoID, manifestURL, thumbnailURL); err != nil {
		return fmt.Errorf("publish processed: %w", err)
	}

	log.Printf("completed videoId=%s", job.VideoID)
	return nil
}

func buildMasterManifest(videoID string, qs []struct{ name, scale, bitrate string }) string {
	bw := map[string]int{"1080p": 4500000, "720p": 2800000, "360p": 1000000}
	res := map[string]string{"1080p": "1920x1080", "720p": "1280x720", "360p": "640x360"}
	s := "#EXTM3U\n"
	for _, q := range qs {
		s += fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%s\n../../segments/%s/%s/media.m3u8\n",
			bw[q.name], res[q.name], videoID, q.name)
	}
	return s
}
