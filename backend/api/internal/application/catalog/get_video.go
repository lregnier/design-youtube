package catalog

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
	"github.com/lregnier/design-youtube/api/internal/ports"
)

type GetVideoCommand struct {
	VideoID string
}

type GetVideo struct {
	repo  video.VideoRepository
	cache ports.Cache
}

func NewGetVideo(repo video.VideoRepository, cache ports.Cache) GetVideo {
	return GetVideo{repo: repo, cache: cache}
}

func (uc GetVideo) Execute(ctx context.Context, cmd GetVideoCommand) (*video.Video, error) {
	if cmd.VideoID == "" {
		return nil, errors.New("videoID is required")
	}

	cacheKey := "video:" + cmd.VideoID
	if data, err := uc.cache.Get(ctx, cacheKey); err == nil {
		var vid video.Video
		if json.Unmarshal(data, &vid) == nil {
			return &vid, nil
		}
	}

	vid, err := uc.repo.FindByID(ctx, video.VideoID(cmd.VideoID))
	if err != nil {
		return nil, err
	}
	if vid == nil {
		return nil, nil
	}

	if data, err := json.Marshal(vid); err == nil {
		uc.cache.Set(ctx, cacheKey, data)
	}

	return vid, nil
}
