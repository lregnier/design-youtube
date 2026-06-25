package sqsjobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/lregnier/design-youtube/worker/internal/application"
	"github.com/lregnier/design-youtube/worker/internal/domain/processing"
	"github.com/lregnier/design-youtube/worker/internal/ports"
)

const maxReceiveCount = 3

type Consumer struct {
	sqsClient    *sqs.Client
	queueURL     string
	processVideo application.ProcessVideo
	publisher    ports.ResultPublisher
}

func NewConsumer(sqsClient *sqs.Client, queueURL string, pv application.ProcessVideo, publisher ports.ResultPublisher) *Consumer {
	return &Consumer{sqsClient: sqsClient, queueURL: queueURL, processVideo: pv, publisher: publisher}
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
		MessageSystemAttributeNames: []sqstypes.MessageSystemAttributeName{
			sqstypes.MessageSystemAttributeNameApproximateReceiveCount,
		},
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

		heartbeatCtx, stopHeartbeat := context.WithCancel(ctx)
		go c.startHeartbeat(heartbeatCtx, msg.ReceiptHandle)

		err = c.processVideo.Execute(ctx, job)
		stopHeartbeat()

		if err != nil {
			receiveCount := parseReceiveCount(msg.Attributes)
			if receiveCount >= maxReceiveCount {
				log.Printf("job consumer: videoId=%s reached max retries (%d), publishing failure", job.VideoID, maxReceiveCount)
				if pubErr := c.publisher.PublishFailed(ctx, job.VideoID, fmt.Sprintf("processing failed after %d attempts: %v", maxReceiveCount, err)); pubErr != nil {
					log.Printf("job consumer: failed to publish failure event for videoId=%s: %v", job.VideoID, pubErr)
				}
				c.deleteMessage(ctx, msg.ReceiptHandle)
			} else {
				log.Printf("job consumer: videoId=%s failed (attempt %d/%d), will retry: %v", job.VideoID, receiveCount, maxReceiveCount, err)
			}
			continue
		}

		c.deleteMessage(ctx, msg.ReceiptHandle)
	}
}

func (c *Consumer) startHeartbeat(ctx context.Context, receiptHandle *string) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, err := c.sqsClient.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
				QueueUrl:          &c.queueURL,
				ReceiptHandle:     receiptHandle,
				VisibilityTimeout: 900,
			})
			if err != nil {
				log.Printf("job consumer heartbeat error: %v", err)
			}
		}
	}
}

func (c *Consumer) deleteMessage(ctx context.Context, receiptHandle *string) {
	c.sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &c.queueURL,
		ReceiptHandle: receiptHandle,
	})
}

func parseReceiveCount(attrs map[string]string) int {
	v, ok := attrs["ApproximateReceiveCount"]
	if !ok {
		return 0
	}
	n, _ := strconv.Atoi(v)
	return n
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
