package sqsqueue

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/lregnier/design-youtube/api/internal/ports"
)

var _ ports.Queue = (*Queue)(nil)

type Queue struct {
	client   *sqs.Client
	queueURL string
}

func NewQueue(client *sqs.Client, queueURL string) *Queue {
	return &Queue{client: client, queueURL: queueURL}
}

func (q *Queue) SendMessage(ctx context.Context, body, messageGroupID string) error {
	_, err := q.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:               &q.queueURL,
		MessageBody:            &body,
		MessageGroupId:         &messageGroupID,
		MessageDeduplicationId: &messageGroupID,
	})
	return err
}
