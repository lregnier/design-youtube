package upload

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/lregnier/design-youtube/backend/internal/domain/video"
	"github.com/lregnier/design-youtube/backend/internal/ports"
)

const MaxFileSize = 104857600 // 100MB

type InitUploadCommand struct {
	VideoID     string // empty for new uploads, set to resume
	Title       string
	Description string
	FileSize    int64
	TotalChunks int
}

type InitUploadResult struct {
	VideoID        string
	UploadID       string
	NextPartNumber int
	PresignedURL   string
	Chunks         []ChunkState
}

type ChunkState struct {
	PartNumber int
	Uploaded   bool
	ETag       string
}

type InitUpload struct {
	repo  video.VideoRepository
	store ports.ObjectStore
	bucket string
}

func NewInitUpload(repo video.VideoRepository, store ports.ObjectStore, bucket string) InitUpload {
	return InitUpload{repo: repo, store: store, bucket: bucket}
}

func (uc InitUpload) Execute(ctx context.Context, cmd InitUploadCommand) (InitUploadResult, error) {
	if cmd.FileSize > MaxFileSize {
		return InitUploadResult{}, fmt.Errorf("file size %d exceeds 100MB limit", cmd.FileSize)
	}

	// Resume existing upload
	if cmd.VideoID != "" {
		vid, err := uc.repo.FindByID(ctx, video.VideoID(cmd.VideoID))
		if err != nil {
			return InitUploadResult{}, err
		}
		if vid != nil && vid.Status == video.StatusUploading {
			return uc.buildResumeResult(ctx, vid)
		}
	}

	// New upload
	id := video.VideoID(uuid.New().String())
	key := fmt.Sprintf("raw/%s/original", id)

	mpu, err := uc.store.CreateMultipartUpload(ctx, key)
	if err != nil {
		return InitUploadResult{}, fmt.Errorf("create multipart upload: %w", err)
	}

	vid := video.New(id, cmd.Title, cmd.Description, cmd.TotalChunks, mpu.UploadID)
	if err := uc.repo.Save(ctx, vid); err != nil {
		return InitUploadResult{}, err
	}

	presigned, err := uc.store.PresignUploadPart(ctx, key, mpu.UploadID, 1)
	if err != nil {
		return InitUploadResult{}, err
	}

	return InitUploadResult{
		VideoID:        id.String(),
		UploadID:       mpu.UploadID,
		NextPartNumber: 1,
		PresignedURL:   presigned.URL,
		Chunks:         toChunkStates(vid.Chunks),
	}, nil
}

func (uc InitUpload) buildResumeResult(ctx context.Context, vid *video.Video) (InitUploadResult, error) {
	partNumber, ok := vid.NextMissingChunk()
	if !ok {
		return InitUploadResult{}, errors.New("all chunks already uploaded, call complete")
	}

	key := fmt.Sprintf("raw/%s/original", vid.ID)
	presigned, err := uc.store.PresignUploadPart(ctx, key, vid.UploadID, partNumber)
	if err != nil {
		return InitUploadResult{}, err
	}

	return InitUploadResult{
		VideoID:        vid.ID.String(),
		UploadID:       vid.UploadID,
		NextPartNumber: partNumber,
		PresignedURL:   presigned.URL,
		Chunks:         toChunkStates(vid.Chunks),
	}, nil
}

func toChunkStates(chunks []video.Chunk) []ChunkState {
	out := make([]ChunkState, len(chunks))
	for i, c := range chunks {
		out[i] = ChunkState{PartNumber: c.PartNumber, Uploaded: c.Uploaded, ETag: c.ETag}
	}
	return out
}
