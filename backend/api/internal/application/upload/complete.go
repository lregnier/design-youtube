package upload

import (
	"context"
	"errors"
	"fmt"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
	"github.com/lregnier/design-youtube/api/internal/ports"
)

type CompleteUploadCommand struct {
	VideoID  string
	UploadID string
}

type CompleteUpload struct {
	repo      video.VideoRepository
	store     ports.ObjectStore
	publisher ports.EventPublisher
}

func NewCompleteUpload(repo video.VideoRepository, store ports.ObjectStore, publisher ports.EventPublisher) CompleteUpload {
	return CompleteUpload{repo: repo, store: store, publisher: publisher}
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
	if err := uc.repo.Save(ctx, vid); err != nil {
		return err
	}

	if err := uc.publisher.Publish(ctx, video.VideoUploadedEvent{VideoID: string(vid.ID), S3Key: key}); err != nil {
		return fmt.Errorf("publish video uploaded: %w", err)
	}

	return nil
}
