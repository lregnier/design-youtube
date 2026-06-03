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

func (q *Queue) DeleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := q.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &q.queueURL,
		ReceiptHandle: &receiptHandle,
	})
	return err
}
