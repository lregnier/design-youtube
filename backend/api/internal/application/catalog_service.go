package application

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
)

type CatalogService interface {
	GetVideo(ctx context.Context, videoID string) (*video.Video, error)
	ListVideos(ctx context.Context) ([]*video.Video, error)
}

var _ CatalogService = (*catalogService)(nil)

type catalogService struct {
	repo  video.VideoRepository
	cache Cache
}

func NewCatalogService(repo video.VideoRepository, cache Cache) CatalogService {
	return &catalogService{repo: repo, cache: cache}
}

func (s *catalogService) GetVideo(ctx context.Context, videoID string) (*video.Video, error) {
	if videoID == "" {
		return nil, errors.New("videoID is required")
	}

	cacheKey := "video:" + videoID
	if data, err := s.cache.Get(ctx, cacheKey); err == nil {
		var vid video.Video
		if json.Unmarshal(data, &vid) == nil {
			return &vid, nil
		}
	}

	vid, err := s.repo.FindByID(ctx, video.VideoID(videoID))
	if err != nil {
		return nil, err
	}
	if vid == nil {
		return nil, nil
	}

	if data, err := json.Marshal(vid); err == nil {
		s.cache.Set(ctx, cacheKey, data)
	}

	return vid, nil
}

func (s *catalogService) ListVideos(ctx context.Context) ([]*video.Video, error) {
	return s.repo.List(ctx)
}
