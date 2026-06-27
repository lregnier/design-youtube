package application

import (
	"context"
	"errors"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
)

type ProcessingService interface {
	HandleVideoProcessingSucceeded(ctx context.Context, evt video.VideoProcessingSucceededEvent) error
	HandleVideoProcessingFailed(ctx context.Context, evt video.VideoProcessingFailedEvent) error
}

var _ ProcessingService = (*processingService)(nil)

type processingService struct {
	repo video.VideoRepository
}

func NewProcessingService(repo video.VideoRepository) ProcessingService {
	return &processingService{repo: repo}
}

func (s *processingService) HandleVideoProcessingSucceeded(ctx context.Context, evt video.VideoProcessingSucceededEvent) error {
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

func (s *processingService) HandleVideoProcessingFailed(ctx context.Context, evt video.VideoProcessingFailedEvent) error {
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
