package sqspublisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/lregnier/design-youtube/worker/internal/application"
	"github.com/lregnier/design-youtube/worker/internal/domain/processing"
)

var _ application.EventPublisher = (*publisher)(nil)

type publisher struct {
	client   *sqs.Client
	queueURL string
}

func NewPublisher(client *sqs.Client, queueURL string) application.EventPublisher {
	return &publisher{client: client, queueURL: queueURL}
}

type videoProcessedMessage struct {
	EventType    string `json:"eventType"`
	VideoID      string `json:"videoId"`
	ManifestURL  string `json:"manifestUrl"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type videoFailedMessage struct {
	EventType string `json:"eventType"`
	VideoID   string `json:"videoId"`
	Reason    string `json:"reason"`
}

func (p *publisher) Publish(ctx context.Context, event processing.DomainEvent) error {
	switch evt := event.(type) {
	case processing.VideoProcessingSucceededEvent:
		return p.emit(ctx, evt.VideoID, videoProcessedMessage{
			EventType:    "VideoProcessed",
			VideoID:      evt.VideoID,
			ManifestURL:  evt.ManifestURL,
			ThumbnailURL: evt.ThumbnailURL,
		})
	case processing.VideoProcessingFailedEvent:
		return p.emit(ctx, evt.VideoID, videoFailedMessage{
			EventType: "VideoFailed",
			VideoID:   evt.VideoID,
			Reason:    evt.Reason,
		})
	default:
		return fmt.Errorf("unknown event type: %T", event)
	}
}

func (p *publisher) emit(ctx context.Context, videoID string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	body := string(data)
	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:               &p.queueURL,
		MessageBody:            &body,
		MessageGroupId:         &videoID,
		MessageDeduplicationId: &videoID,
	})
	return err
}
