package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"

	"github.com/lregnier/design-youtube/backend/internal/api"
	"github.com/lregnier/design-youtube/backend/internal/config"
	"github.com/lregnier/design-youtube/backend/internal/store"
)

const maxFileSize = 104857600 // 100MB
const chunkMaxSize = 10485760 // 10MB

type Handler struct {
	cfg *config.Config
	db  *store.Store
	s3  *s3.Client
}

func New(cfg *config.Config, db *store.Store) *Handler {
	awsCfg, _ := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.AWSRegion),
	)
	return &Handler{cfg: cfg, db: db, s3: s3.NewFromConfig(awsCfg)}
}

func (h *Handler) GetVideos(ctx context.Context, _ api.GetVideosRequestObject) (api.GetVideosResponseObject, error) {
	records, err := h.db.ListReadyVideos(ctx)
	if err != nil {
		return api.GetVideos500JSONResponse{Error: err.Error()}, nil
	}

	summaries := make([]api.VideoSummary, 0, len(records))
	for _, r := range records {
		ts, _ := time.Parse(time.RFC3339, r.UploadedAt)
		summaries = append(summaries, api.VideoSummary{
			VideoId:      r.VideoID,
			Title:        r.Title,
			ThumbnailUrl: r.ThumbnailURL,
			UploadedAt:   ts,
			Status:       api.VideoStatus(r.Status),
		})
	}
	return api.GetVideos200JSONResponse(summaries), nil
}

func (h *Handler) GetVideo(ctx context.Context, req api.GetVideoRequestObject) (api.GetVideoResponseObject, error) {
	rec, err := h.db.GetVideo(ctx, req.VideoId)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return api.GetVideo404JSONResponse{Error: "video not found"}, nil
	}

	ts, _ := time.Parse(time.RFC3339, rec.UploadedAt)
	detail := api.VideoDetail{
		VideoId:     rec.VideoID,
		Title:       rec.Title,
		Description: rec.Description,
		Status:      api.VideoStatus(rec.Status),
		UploadedAt:  ts,
	}
	if rec.ManifestURL != "" {
		detail.ManifestUrl = &rec.ManifestURL
	}
	if rec.ThumbnailURL != "" {
		detail.ThumbnailUrl = &rec.ThumbnailURL
	}
	return api.GetVideo200JSONResponse(detail), nil
}

func (h *Handler) InitUpload(ctx context.Context, req api.InitUploadRequestObject) (api.InitUploadResponseObject, error) {
	body := req.Body
	if body.FileSize > maxFileSize {
		msg := fmt.Sprintf("file size exceeds 100MB limit (%d bytes)", maxFileSize)
		return api.InitUpload400JSONResponse{Error: msg}, nil
	}

	// Resume existing upload
	if body.VideoId != nil && *body.VideoId != "" {
		rec, err := h.db.GetVideo(ctx, *body.VideoId)
		if err != nil {
			return nil, err
		}
		if rec != nil && rec.Status == store.StatusUploading {
			return h.buildResumeResponse(ctx, rec)
		}
	}

	// New upload
	videoID := uuid.New().String()
	s3Key := fmt.Sprintf("raw/%s/original", videoID)

	mpu, err := h.s3.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket: &h.cfg.S3Bucket,
		Key:    &s3Key,
	})
	if err != nil {
		return nil, fmt.Errorf("create multipart upload: %w", err)
	}

	chunks := make([]store.ChunkState, body.TotalChunks)
	for i := range chunks {
		chunks[i] = store.ChunkState{PartNumber: i + 1}
	}

	rec := &store.VideoRecord{
		VideoID:     videoID,
		Title:       body.Title,
		Description: body.Description,
		Status:      store.StatusUploading,
		UploadedAt:  time.Now().UTC().Format(time.RFC3339),
		UploadID:    *mpu.UploadId,
		TotalChunks: body.TotalChunks,
		Chunks:      chunks,
	}
	if err := h.db.PutVideo(ctx, rec); err != nil {
		return nil, err
	}

	presigned, err := h.presignPart(ctx, s3Key, *mpu.UploadId, 1)
	if err != nil {
		return nil, err
	}

	apiChunks := make([]api.ChunkStatus, len(chunks))
	for i, c := range chunks {
		apiChunks[i] = api.ChunkStatus{PartNumber: c.PartNumber, Uploaded: c.Uploaded}
	}

	return api.InitUpload200JSONResponse{
		VideoId:       videoID,
		UploadId:      *mpu.UploadId,
		NextPartNumber: 1,
		PresignedUrl:  presigned,
		Chunks:        apiChunks,
	}, nil
}

func (h *Handler) buildResumeResponse(ctx context.Context, rec *store.VideoRecord) (api.InitUploadResponseObject, error) {
	s3Key := fmt.Sprintf("raw/%s/original", rec.VideoID)

	apiChunks := make([]api.ChunkStatus, len(rec.Chunks))
	firstMissing := -1
	for i, c := range rec.Chunks {
		apiChunks[i] = api.ChunkStatus{PartNumber: c.PartNumber, Uploaded: c.Uploaded, ETag: &c.ETag}
		if !c.Uploaded && firstMissing < 0 {
			firstMissing = c.PartNumber
		}
	}

	presigned, err := h.presignPart(ctx, s3Key, rec.UploadID, firstMissing)
	if err != nil {
		return nil, err
	}

	return api.InitUpload200JSONResponse{
		VideoId:       rec.VideoID,
		UploadId:      rec.UploadID,
		NextPartNumber: firstMissing,
		PresignedUrl:  presigned,
		Chunks:        apiChunks,
	}, nil
}

func (h *Handler) ConfirmChunk(ctx context.Context, req api.ConfirmChunkRequestObject) (api.ConfirmChunkResponseObject, error) {
	rec, err := h.db.GetVideo(ctx, req.VideoId)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return api.ConfirmChunk404JSONResponse{Error: "video not found"}, nil
	}

	if err := h.db.MarkChunkUploaded(ctx, req.VideoId, req.Body.PartNumber, req.Body.ETag); err != nil {
		return nil, err
	}

	// Re-fetch to get updated state
	rec, err = h.db.GetVideo(ctx, req.VideoId)
	if err != nil {
		return nil, err
	}

	nextPart := -1
	for _, c := range rec.Chunks {
		if !c.Uploaded {
			nextPart = c.PartNumber
			break
		}
	}

	done := nextPart < 0
	resp := api.ConfirmChunkResponse{
		PartNumber: req.Body.PartNumber,
		Done:       done,
	}

	if !done {
		s3Key := fmt.Sprintf("raw/%s/original", rec.VideoID)
		presigned, err := h.presignPart(ctx, s3Key, rec.UploadID, nextPart)
		if err != nil {
			return nil, err
		}
		resp.NextPartNumber = &nextPart
		resp.PresignedUrl = &presigned
	}

	return api.ConfirmChunk200JSONResponse(resp), nil
}

func (h *Handler) CompleteUpload(ctx context.Context, req api.CompleteUploadRequestObject) (api.CompleteUploadResponseObject, error) {
	rec, err := h.db.GetVideo(ctx, req.VideoId)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return api.CompleteUpload404JSONResponse{Error: "video not found"}, nil
	}

	parts := make([]s3types.CompletedPart, 0, len(rec.Chunks))
	for _, c := range rec.Chunks {
		if c.Uploaded {
			pn := int32(c.PartNumber)
			parts = append(parts, s3types.CompletedPart{
				PartNumber: &pn,
				ETag:       aws.String(c.ETag),
			})
		}
	}

	s3Key := fmt.Sprintf("raw/%s/original", rec.VideoID)
	_, err = h.s3.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   &h.cfg.S3Bucket,
		Key:      &s3Key,
		UploadId: &req.Body.UploadId,
		MultipartUpload: &s3types.CompletedMultipartUpload{
			Parts: parts,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("complete multipart upload: %w", err)
	}

	if err := h.db.UpdateVideoStatus(ctx, req.VideoId, store.StatusProcessing); err != nil {
		return nil, err
	}

	return api.CompleteUpload200Response{}, nil
}

func (h *Handler) presignPart(ctx context.Context, s3Key, uploadID string, partNumber int) (string, error) {
	presigner := s3.NewPresignClient(h.s3)
	pn := int32(partNumber)
	out, err := presigner.PresignUploadPart(ctx, &s3.UploadPartInput{
		Bucket:     &h.cfg.S3Bucket,
		Key:        &s3Key,
		UploadId:   &uploadID,
		PartNumber: &pn,
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", fmt.Errorf("presign part %d: %w", partNumber, err)
	}
	return out.URL, nil
}
