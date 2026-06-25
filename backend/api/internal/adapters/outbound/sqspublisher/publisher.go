package sqspublisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
	"github.com/lregnier/design-youtube/api/internal/ports"
)

var _ ports.EventPublisher = (*Publisher)(nil)

type Publisher struct {
	client   *sqs.Client
	queueURL string
}

func NewPublisher(client *sqs.Client, queueURL string) *Publisher {
	return &Publisher{client: client, queueURL: queueURL}
}

type processingJob struct {
	VideoID string `json:"videoId"`
	S3Key   string `json:"s3Key"`
}

func (p *Publisher) Publish(ctx context.Context, event video.DomainEvent) error {
	switch e := event.(type) {
	case video.VideoUploadedEvent:
		return p.publishVideoUploaded(ctx, e)
	default:
		return fmt.Errorf("unknown event type: %T", event)
	}
}

func (p *Publisher) publishVideoUploaded(ctx context.Context, e video.VideoUploadedEvent) error {
	body, err := json.Marshal(processingJob{VideoID: e.VideoID, S3Key: e.S3Key})
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:               &p.queueURL,
		MessageBody:            aws(string(body)),
		MessageGroupId:         aws(e.VideoID),
		MessageDeduplicationId: aws(e.VideoID),
	})
	return err
}

func aws(s string) *string { return &s }
