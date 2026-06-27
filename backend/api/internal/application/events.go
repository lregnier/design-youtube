package application

type VideoProcessedEvent struct {
	EventType    string `json:"eventType"`
	VideoID      string `json:"videoId"`
	ManifestURL  string `json:"manifestUrl"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type VideoFailedEvent struct {
	EventType string `json:"eventType"`
	VideoID   string `json:"videoId"`
	Reason    string `json:"reason"`
}
