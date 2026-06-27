package processing

type DomainEvent interface {
	domainEvent()
}

type VideoProcessingSucceededEvent struct {
	VideoID      string
	ManifestURL  string
	ThumbnailURL string
}

func (VideoProcessingSucceededEvent) domainEvent() {}

type VideoProcessingFailedEvent struct {
	VideoID string
	Reason  string
}

func (VideoProcessingFailedEvent) domainEvent() {}
