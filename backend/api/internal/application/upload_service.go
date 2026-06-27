package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
)

const MaxFileSize = 104857600 // 100MB

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

type ConfirmChunkResult struct {
	PartNumber     int
	Done           bool
	NextPartNumber int
	PresignedURL   string
}

type UploadService interface {
	InitUpload(ctx context.Context, videoID, title, description string, fileSize int64, totalChunks int) (InitUploadResult, error)
	ConfirmChunk(ctx context.Context, videoID string, partNumber int, eTag string) (ConfirmChunkResult, error)
	CompleteUpload(ctx context.Context, videoID, uploadID string) error
}

var _ UploadService = (*uploadService)(nil)

type uploadService struct {
	repo      video.VideoRepository
	store     ObjectStore
	publisher EventPublisher
	bucket    string
}

func NewUploadService(repo video.VideoRepository, store ObjectStore, publisher EventPublisher, bucket string) UploadService {
	return &uploadService{repo: repo, store: store, publisher: publisher, bucket: bucket}
}

func (s *uploadService) InitUpload(ctx context.Context, videoID, title, description string, fileSize int64, totalChunks int) (InitUploadResult, error) {
	if fileSize > MaxFileSize {
		return InitUploadResult{}, fmt.Errorf("file size %d exceeds 100MB limit", fileSize)
	}

	if videoID != "" {
		vid, err := s.repo.FindByID(ctx, video.VideoID(videoID))
		if err != nil {
			return InitUploadResult{}, err
		}
		if vid != nil && vid.Status == video.StatusUploading {
			return s.buildResumeResult(ctx, vid)
		}
	}

	id := video.VideoID(uuid.New().String())
	key := fmt.Sprintf("raw/%s/original", id)

	mpu, err := s.store.CreateMultipartUpload(ctx, key)
	if err != nil {
		return InitUploadResult{}, fmt.Errorf("create multipart upload: %w", err)
	}

	vid := video.New(id, title, description, totalChunks, mpu.UploadID)
	if err := s.repo.Save(ctx, vid); err != nil {
		return InitUploadResult{}, err
	}

	presigned, err := s.store.PresignUploadPart(ctx, key, mpu.UploadID, 1)
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

func (s *uploadService) ConfirmChunk(ctx context.Context, videoID string, partNumber int, eTag string) (ConfirmChunkResult, error) {
	vid, err := s.repo.FindByID(ctx, video.VideoID(videoID))
	if err != nil {
		return ConfirmChunkResult{}, err
	}
	if vid == nil {
		return ConfirmChunkResult{}, errors.New("video not found")
	}

	if err := vid.MarkChunkUploaded(partNumber, eTag); err != nil {
		return ConfirmChunkResult{}, err
	}
	if err := s.repo.Save(ctx, vid); err != nil {
		return ConfirmChunkResult{}, err
	}

	next, hasNext := vid.NextMissingChunk()
	if !hasNext {
		return ConfirmChunkResult{PartNumber: partNumber, Done: true}, nil
	}

	key := fmt.Sprintf("raw/%s/original", vid.ID)
	presigned, err := s.store.PresignUploadPart(ctx, key, vid.UploadID, next)
	if err != nil {
		return ConfirmChunkResult{}, err
	}

	return ConfirmChunkResult{
		PartNumber:     partNumber,
		Done:           false,
		NextPartNumber: next,
		PresignedURL:   presigned.URL,
	}, nil
}

func (s *uploadService) CompleteUpload(ctx context.Context, videoID, uploadID string) error {
	vid, err := s.repo.FindByID(ctx, video.VideoID(videoID))
	if err != nil {
		return err
	}
	if vid == nil {
		return errors.New("video not found")
	}

	parts := make([]CompletedPart, 0, len(vid.Chunks))
	for _, c := range vid.Chunks {
		if c.Uploaded {
			parts = append(parts, CompletedPart{PartNumber: c.PartNumber, ETag: c.ETag})
		}
	}

	key := fmt.Sprintf("raw/%s/original", vid.ID)
	if err := s.store.CompleteMultipartUpload(ctx, key, uploadID, parts); err != nil {
		return fmt.Errorf("complete multipart upload: %w", err)
	}

	vid.MarkProcessing()
	if err := s.repo.Save(ctx, vid); err != nil {
		return err
	}

	if err := s.publisher.Publish(ctx, video.VideoUploadedEvent{VideoID: string(vid.ID), S3Key: key}); err != nil {
		return fmt.Errorf("publish video uploaded: %w", err)
	}

	return nil
}

func (s *uploadService) buildResumeResult(ctx context.Context, vid *video.Video) (InitUploadResult, error) {
	partNumber, ok := vid.NextMissingChunk()
	if !ok {
		return InitUploadResult{}, errors.New("all chunks already uploaded, call complete")
	}

	key := fmt.Sprintf("raw/%s/original", vid.ID)
	presigned, err := s.store.PresignUploadPart(ctx, key, vid.UploadID, partNumber)
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
