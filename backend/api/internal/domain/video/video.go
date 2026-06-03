package video

import (
	"errors"
	"time"
)

type VideoID string

func (id VideoID) String() string { return string(id) }

type VideoStatus string

const (
	StatusUploading  VideoStatus = "uploading"
	StatusProcessing VideoStatus = "processing"
	StatusReady      VideoStatus = "ready"
	StatusFailed     VideoStatus = "failed"
)

type Chunk struct {
	PartNumber int
	Uploaded   bool
	ETag       string
}

type Video struct {
	ID           VideoID
	Title        string
	Description  string
	Status       VideoStatus
	UploadedAt   time.Time
	UploadID     string
	TotalChunks  int
	Chunks       []Chunk
	ManifestURL  string
	ThumbnailURL string
}

func New(id VideoID, title, description string, totalChunks int, uploadID string) *Video {
	chunks := make([]Chunk, totalChunks)
	for i := range chunks {
		chunks[i] = Chunk{PartNumber: i + 1}
	}
	return &Video{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      StatusUploading,
		UploadedAt:  time.Now().UTC(),
		UploadID:    uploadID,
		TotalChunks: totalChunks,
		Chunks:      chunks,
	}
}

func (v *Video) IsReady() bool { return v.Status == StatusReady }

func (v *Video) NextMissingChunk() (int, bool) {
	for _, c := range v.Chunks {
		if !c.Uploaded {
			return c.PartNumber, true
		}
	}
	return 0, false
}

func (v *Video) MarkChunkUploaded(partNumber int, eTag string) error {
	for i := range v.Chunks {
		if v.Chunks[i].PartNumber == partNumber {
			v.Chunks[i].Uploaded = true
			v.Chunks[i].ETag = eTag
			return nil
		}
	}
	return errors.New("chunk not found")
}

func (v *Video) AllChunksUploaded() bool {
	for _, c := range v.Chunks {
		if !c.Uploaded {
			return false
		}
	}
	return true
}

func (v *Video) MarkProcessing() { v.Status = StatusProcessing }
func (v *Video) MarkReady(manifestURL, thumbnailURL string) {
	v.Status = StatusReady
	v.ManifestURL = manifestURL
	v.ThumbnailURL = thumbnailURL
}
func (v *Video) MarkFailed() { v.Status = StatusFailed }
