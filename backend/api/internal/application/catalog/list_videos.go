package catalog

import (
	"context"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
)

type ListVideos struct {
	repo video.VideoRepository
}

func NewListVideos(repo video.VideoRepository) ListVideos {
	return ListVideos{repo: repo}
}

func (uc ListVideos) Execute(ctx context.Context) ([]*video.Video, error) {
	return uc.repo.ListReady(ctx)
}
