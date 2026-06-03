package sqspublisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/lregnier/design-youtube/worker/internal/event"
	"github.com/lregnier/design-youtube/worker/internal/ports"
)

var _ ports.ResultPublisher = (*Publisher)(nil)

type Publisher struct {
	client   *sqs.Client
	queueURL string
}

func NewPublisher(client *sqs.Client, queueURL string) *Publisher {
	return &Publisher{client: client, queueURL: queueURL}
}

func (p *Publisher) PublishProcessed(ctx context.Context, videoID, manifestURL, thumbnailURL string) error {
	return p.emit(ctx, videoID, event.VideoProcessed{
		EventType:    event.TypeVideoProcessed,
		VideoID:      videoID,
		ManifestURL:  manifestURL,
		ThumbnailURL: thumbnailURL,
	})
}

func (p *Publisher) PublishFailed(ctx context.Context, videoID, reason string) error {
	return p.emit(ctx, videoID, event.VideoFailed{
		EventType: event.TypeVideoFailed,
		VideoID:   videoID,
		Reason:    reason,
	})
}

func (p *Publisher) emit(ctx context.Context, videoID string, payload any) error {
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
