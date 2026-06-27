package video

type DomainEvent interface {
	domainEvent()
}

type VideoUploadedEvent struct {
	VideoID string
	S3Key   string
}

func (VideoUploadedEvent) domainEvent() {}

type VideoProcessingSucceededEvent struct {
	VideoID      string
	ManifestURL  string
	ThumbnailURL string
}

type VideoProcessingFailedEvent struct {
	VideoID string
	Reason  string
}
