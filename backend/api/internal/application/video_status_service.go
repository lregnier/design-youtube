package application

import (
	"context"
	"errors"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
)

type VideoStatusService interface {
	MarkReady(ctx context.Context, evt video.VideoProcessingSucceededEvent) error
	MarkFailed(ctx context.Context, evt video.VideoProcessingFailedEvent) error
}

var _ VideoStatusService = (*videoStatusService)(nil)

type videoStatusService struct {
	repo video.VideoRepository
}

func NewVideoStatusService(repo video.VideoRepository) VideoStatusService {
	return &videoStatusService{repo: repo}
}

func (s *videoStatusService) MarkReady(ctx context.Context, evt video.VideoProcessingSucceededEvent) error {
	if evt.VideoID == "" {
		return errors.New("videoId is required")
	}
	vid, err := s.repo.FindByID(ctx, video.VideoID(evt.VideoID))
	if err != nil {
		return err
	}
	if vid == nil {
		return errors.New("video not found")
	}
	vid.MarkReady(evt.ManifestURL, evt.ThumbnailURL)
	return s.repo.Save(ctx, vid)
}

func (s *videoStatusService) MarkFailed(ctx context.Context, evt video.VideoProcessingFailedEvent) error {
	if evt.VideoID == "" {
		return errors.New("videoId is required")
	}
	vid, err := s.repo.FindByID(ctx, video.VideoID(evt.VideoID))
	if err != nil {
		return err
	}
	if vid == nil {
		return errors.New("video not found")
	}
	vid.MarkFailed()
	return s.repo.Save(ctx, vid)
}
