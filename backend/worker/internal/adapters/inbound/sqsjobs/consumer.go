package sqsjobs

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/lregnier/design-youtube/worker/internal/application"
	"github.com/lregnier/design-youtube/worker/internal/domain/processing"
)

type Consumer struct {
	sqsClient    *sqs.Client
	queueURL     string
	processVideo application.ProcessVideo
}

func NewConsumer(sqsClient *sqs.Client, queueURL string, pv application.ProcessVideo) *Consumer {
	return &Consumer{sqsClient: sqsClient, queueURL: queueURL, processVideo: pv}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("job consumer started, polling SQS")
	for {
		select {
		case <-ctx.Done():
			log.Println("job consumer stopped")
			return
		default:
			c.poll(ctx)
		}
	}
}

func (c *Consumer) poll(ctx context.Context) {
	out, err := c.sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &c.queueURL,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     20,
		VisibilityTimeout:   900,
	})
	if err != nil {
		log.Printf("job consumer receive error: %v", err)
		time.Sleep(5 * time.Second)
		return
	}

	for _, msg := range out.Messages {
		job, err := parseJob(msg.Body)
		if err != nil {
			log.Printf("job consumer parse error: %v — skipping", err)
			c.deleteMessage(ctx, msg.ReceiptHandle)
			continue
		}

		if err := c.processVideo.Execute(ctx, job); err != nil {
			log.Printf("job consumer process error for videoId=%s (message stays in queue): %v", job.VideoID, err)
			continue
		}

		c.deleteMessage(ctx, msg.ReceiptHandle)
	}
}

func (c *Consumer) deleteMessage(ctx context.Context, receiptHandle *string) {
	c.sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &c.queueURL,
		ReceiptHandle: receiptHandle,
	})
}

type rawJob struct {
	VideoID string `json:"videoId"`
	S3Key   string `json:"s3Key"`
}

func parseJob(body *string) (processing.ProcessingJob, error) {
	var j rawJob
	if err := json.Unmarshal([]byte(*body), &j); err != nil {
		return processing.ProcessingJob{}, err
	}
	return processing.ProcessingJob{VideoID: j.VideoID, S3Key: j.S3Key}, nil
}
