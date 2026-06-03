package upload

import (
	"context"
	"errors"
	"fmt"

	"github.com/lregnier/design-youtube/backend/internal/domain/video"
	"github.com/lregnier/design-youtube/backend/internal/ports"
)

type ConfirmChunkCommand struct {
	VideoID    string
	PartNumber int
	ETag       string
}

type ConfirmChunkResult struct {
	PartNumber     int
	Done           bool
	NextPartNumber int
	PresignedURL   string
}

type ConfirmChunk struct {
	repo  video.VideoRepository
	store ports.ObjectStore
}

func NewConfirmChunk(repo video.VideoRepository, store ports.ObjectStore) ConfirmChunk {
	return ConfirmChunk{repo: repo, store: store}
}

func (uc ConfirmChunk) Execute(ctx context.Context, cmd ConfirmChunkCommand) (ConfirmChunkResult, error) {
	vid, err := uc.repo.FindByID(ctx, video.VideoID(cmd.VideoID))
	if err != nil {
		return ConfirmChunkResult{}, err
	}
	if vid == nil {
		return ConfirmChunkResult{}, errors.New("video not found")
	}

	if err := vid.MarkChunkUploaded(cmd.PartNumber, cmd.ETag); err != nil {
		return ConfirmChunkResult{}, err
	}
	if err := uc.repo.Save(ctx, vid); err != nil {
		return ConfirmChunkResult{}, err
	}

	next, hasNext := vid.NextMissingChunk()
	if !hasNext {
		return ConfirmChunkResult{PartNumber: cmd.PartNumber, Done: true}, nil
	}

	key := fmt.Sprintf("raw/%s/original", vid.ID)
	presigned, err := uc.store.PresignUploadPart(ctx, key, vid.UploadID, next)
	if err != nil {
		return ConfirmChunkResult{}, err
	}

	return ConfirmChunkResult{
		PartNumber:     cmd.PartNumber,
		Done:           false,
		NextPartNumber: next,
		PresignedURL:   presigned.URL,
	}, nil
}
