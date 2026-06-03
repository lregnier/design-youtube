package processing

import (
	"context"
	"errors"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
)

type VideoProcessedEvent struct {
	EventType    string `json:"eventType"`
	VideoID      string `json:"videoId"`
	ManifestURL  string `json:"manifestUrl"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type VideoFailedEvent struct {
	EventType string `json:"eventType"`
	VideoID   string `json:"videoId"`
	Reason    string `json:"reason"`
}

type ApplyProcessingResult struct {
	repo video.VideoRepository
}

func NewApplyProcessingResult(repo video.VideoRepository) ApplyProcessingResult {
	return ApplyProcessingResult{repo: repo}
}

func (uc ApplyProcessingResult) OnProcessed(ctx context.Context, evt VideoProcessedEvent) error {
	if evt.VideoID == "" {
		return errors.New("videoId is required")
	}
	vid, err := uc.repo.FindByID(ctx, video.VideoID(evt.VideoID))
	if err != nil {
		return err
	}
	if vid == nil {
		return errors.New("video not found")
	}
	// Idempotent: already ready is a no-op
	vid.MarkReady(evt.ManifestURL, evt.ThumbnailURL)
	return uc.repo.Save(ctx, vid)
}

func (uc ApplyProcessingResult) OnFailed(ctx context.Context, evt VideoFailedEvent) error {
	if evt.VideoID == "" {
		return errors.New("videoId is required")
	}
	vid, err := uc.repo.FindByID(ctx, video.VideoID(evt.VideoID))
	if err != nil {
		return err
	}
	if vid == nil {
		return errors.New("video not found")
	}
	vid.MarkFailed()
	return uc.repo.Save(ctx, vid)
}
