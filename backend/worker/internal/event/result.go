package event

const (
	TypeVideoProcessed = "VideoProcessed"
	TypeVideoFailed    = "VideoFailed"
)

type VideoProcessed struct {
	EventType    string `json:"eventType"`
	VideoID      string `json:"videoId"`
	ManifestURL  string `json:"manifestUrl"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type VideoFailed struct {
	EventType string `json:"eventType"`
	VideoID   string `json:"videoId"`
	Reason    string `json:"reason"`
}
