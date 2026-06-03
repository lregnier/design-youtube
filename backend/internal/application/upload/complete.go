package upload

import (
	"context"
	"errors"
	"fmt"

	"github.com/lregnier/design-youtube/backend/internal/domain/video"
	"github.com/lregnier/design-youtube/backend/internal/ports"
)

type CompleteUploadCommand struct {
	VideoID  string
	UploadID string
}

type CompleteUpload struct {
	repo  video.VideoRepository
	store ports.ObjectStore
}

func NewCompleteUpload(repo video.VideoRepository, store ports.ObjectStore) CompleteUpload {
	return CompleteUpload{repo: repo, store: store}
}

func (uc CompleteUpload) Execute(ctx context.Context, cmd CompleteUploadCommand) error {
	vid, err := uc.repo.FindByID(ctx, video.VideoID(cmd.VideoID))
	if err != nil {
		return err
	}
	if vid == nil {
		return errors.New("video not found")
	}

	parts := make([]ports.CompletedPart, 0, len(vid.Chunks))
	for _, c := range vid.Chunks {
		if c.Uploaded {
			parts = append(parts, ports.CompletedPart{PartNumber: c.PartNumber, ETag: c.ETag})
		}
	}

	key := fmt.Sprintf("raw/%s/original", vid.ID)
	if err := uc.store.CompleteMultipartUpload(ctx, key, cmd.UploadID, parts); err != nil {
		return fmt.Errorf("complete multipart upload: %w", err)
	}

	vid.MarkProcessing()
	return uc.repo.Save(ctx, vid)
}
