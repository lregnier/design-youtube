package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Publisher struct {
	client   *sqs.Client
	queueURL string
}

func NewPublisher(client *sqs.Client, queueURL string) *Publisher {
	return &Publisher{client: client, queueURL: queueURL}
}

func (p *Publisher) Emit(ctx context.Context, videoID string, payload any) error {
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
