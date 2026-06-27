package http

import (
	"context"
	"time"

	"github.com/lregnier/design-youtube/api/internal/application"
	"github.com/lregnier/design-youtube/api/internal/domain/video"
	"github.com/lregnier/design-youtube/api/gen/api"
)

type Handler struct {
	upload  application.UploadService
	catalog application.CatalogService
}

func NewHandler(uploadSvc application.UploadService, catalogSvc application.CatalogService) *Handler {
	return &Handler{upload: uploadSvc, catalog: catalogSvc}
}

func (h *Handler) GetVideos(ctx context.Context, _ api.GetVideosRequestObject) (api.GetVideosResponseObject, error) {
	videos, err := h.catalog.ListVideos(ctx)
	if err != nil {
		return api.GetVideos500JSONResponse{Error: err.Error()}, nil
	}
	summaries := make([]api.VideoSummary, 0, len(videos))
	for _, v := range videos {
		summaries = append(summaries, toSummary(v))
	}
	return api.GetVideos200JSONResponse(summaries), nil
}

func (h *Handler) GetVideo(ctx context.Context, req api.GetVideoRequestObject) (api.GetVideoResponseObject, error) {
	v, err := h.catalog.GetVideo(ctx, req.VideoId)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return api.GetVideo404JSONResponse{Error: "video not found"}, nil
	}
	detail := api.VideoDetail{
		VideoId:     v.ID.String(),
		Title:       v.Title,
		Description: v.Description,
		Status:      api.VideoStatus(v.Status),
		UploadedAt:  v.UploadedAt,
	}
	if v.ManifestURL != "" {
		detail.ManifestUrl = &v.ManifestURL
	}
	if v.ThumbnailURL != "" {
		detail.ThumbnailUrl = &v.ThumbnailURL
	}
	return api.GetVideo200JSONResponse(detail), nil
}

func (h *Handler) InitUpload(ctx context.Context, req api.InitUploadRequestObject) (api.InitUploadResponseObject, error) {
	body := req.Body
	videoID := ""
	if body.VideoId != nil {
		videoID = *body.VideoId
	}
	result, err := h.upload.InitUpload(ctx, videoID, body.Title, body.Description, body.FileSize, body.TotalChunks)
	if err != nil {
		return api.InitUpload400JSONResponse{Error: err.Error()}, nil
	}
	chunks := make([]api.ChunkStatus, len(result.Chunks))
	for i, c := range result.Chunks {
		eTag := c.ETag
		chunks[i] = api.ChunkStatus{PartNumber: c.PartNumber, Uploaded: c.Uploaded, ETag: &eTag}
	}
	return api.InitUpload200JSONResponse{
		VideoId:        result.VideoID,
		UploadId:       result.UploadID,
		NextPartNumber: result.NextPartNumber,
		PresignedUrl:   result.PresignedURL,
		Chunks:         chunks,
	}, nil
}

func (h *Handler) ConfirmChunk(ctx context.Context, req api.ConfirmChunkRequestObject) (api.ConfirmChunkResponseObject, error) {
	result, err := h.upload.ConfirmChunk(ctx, req.VideoId, req.Body.PartNumber, req.Body.ETag)
	if err != nil {
		if err.Error() == "video not found" {
			return api.ConfirmChunk404JSONResponse{Error: err.Error()}, nil
		}
		return nil, err
	}
	resp := api.ConfirmChunkResponse{
		PartNumber: result.PartNumber,
		Done:       result.Done,
	}
	if !result.Done {
		resp.NextPartNumber = &result.NextPartNumber
		resp.PresignedUrl = &result.PresignedURL
	}
	return api.ConfirmChunk200JSONResponse(resp), nil
}

func (h *Handler) CompleteUpload(ctx context.Context, req api.CompleteUploadRequestObject) (api.CompleteUploadResponseObject, error) {
	err := h.upload.CompleteUpload(ctx, req.VideoId, req.Body.UploadId)
	if err != nil {
		if err.Error() == "video not found" {
			return api.CompleteUpload404JSONResponse{Error: err.Error()}, nil
		}
		return nil, err
	}
	return api.CompleteUpload200Response{}, nil
}

func toSummary(v *video.Video) api.VideoSummary {
	return api.VideoSummary{
		VideoId:      v.ID.String(),
		Title:        v.Title,
		ThumbnailUrl: v.ThumbnailURL,
		UploadedAt:   v.UploadedAt.UTC().Truncate(time.Second),
		Status:       api.VideoStatus(v.Status),
	}
}
