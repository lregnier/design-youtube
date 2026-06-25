package video

type DomainEvent interface {
	domainEvent()
}

type VideoUploadedEvent struct {
	VideoID string
	S3Key   string
}

func (VideoUploadedEvent) domainEvent() {}
