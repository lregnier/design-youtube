package sqssubscriber

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
)

const maxReceiveCount = 3

type Subscriber struct {
	sqsClient  *sqs.Client
	queueURL   string
	svc        application.VideoProcessingService
	publisher  application.EventPublisher
}

func NewSubscriber(sqsClient *sqs.Client, queueURL string, svc application.VideoProcessingService, publisher application.EventPublisher) *Subscriber {
	return &Subscriber{sqsClient: sqsClient, queueURL: queueURL, svc: svc, publisher: publisher}
}

func (s *Subscriber) Start(ctx context.Context) {
	log.Println("job subscriber started, polling SQS")
	for {
		select {
		case <-ctx.Done():
			log.Println("job subscriber stopped")
			return
		default:
			s.poll(ctx)
		}
	}
}

func (s *Subscriber) poll(ctx context.Context) {
	out, err := s.sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &s.queueURL,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     20,
		VisibilityTimeout:   900,
		MessageSystemAttributeNames: []sqstypes.MessageSystemAttributeName{
			sqstypes.MessageSystemAttributeNameApproximateReceiveCount,
		},
	})
	if err != nil {
		log.Printf("job subscriber receive error: %v", err)
		time.Sleep(5 * time.Second)
		return
	}

	for _, msg := range out.Messages {
		job, err := parseJob(msg.Body)
		if err != nil {
			log.Printf("job subscriber parse error: %v — skipping", err)
			s.deleteMessage(ctx, msg.ReceiptHandle)
			continue
		}

		heartbeatCtx, stopHeartbeat := context.WithCancel(ctx)
		go s.startHeartbeat(heartbeatCtx, msg.ReceiptHandle)

		err = s.svc.Process(ctx, job)
		stopHeartbeat()

		if err != nil {
			receiveCount := parseReceiveCount(msg.Attributes)
			if receiveCount >= maxReceiveCount {
				log.Printf("job subscriber: videoId=%s reached max retries (%d), publishing failure", job.VideoID, maxReceiveCount)
				if pubErr := s.publisher.Publish(ctx, processing.VideoProcessingFailedEvent{
					VideoID: job.VideoID,
					Reason:  fmt.Sprintf("processing failed after %d attempts: %v", maxReceiveCount, err),
				}); pubErr != nil {
					log.Printf("job subscriber: failed to publish failure event for videoId=%s: %v", job.VideoID, pubErr)
				}
				s.deleteMessage(ctx, msg.ReceiptHandle)
			} else {
				log.Printf("job subscriber: videoId=%s failed (attempt %d/%d), will retry: %v", job.VideoID, receiveCount, maxReceiveCount, err)
			}
			continue
		}

		s.deleteMessage(ctx, msg.ReceiptHandle)
	}
}

func (s *Subscriber) startHeartbeat(ctx context.Context, receiptHandle *string) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, err := s.sqsClient.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
				QueueUrl:          &s.queueURL,
				ReceiptHandle:     receiptHandle,
				VisibilityTimeout: 900,
			})
			if err != nil {
				log.Printf("job subscriber heartbeat error: %v", err)
			}
		}
	}
}

func (s *Subscriber) deleteMessage(ctx context.Context, receiptHandle *string) {
	s.sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &s.queueURL,
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
