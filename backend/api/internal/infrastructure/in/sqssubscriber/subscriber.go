package sqssubscriber

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/lregnier/design-youtube/api/internal/application"
)

type Subscriber struct {
	sqsClient *sqs.Client
	queueURL  string
	svc       application.ProcessingService
}

func NewSubscriber(sqsClient *sqs.Client, queueURL string, svc application.ProcessingService) *Subscriber {
	return &Subscriber{sqsClient: sqsClient, queueURL: queueURL, svc: svc}
}

func (s *Subscriber) Start(ctx context.Context) {
	log.Println("results subscriber started")
	for {
		select {
		case <-ctx.Done():
			log.Println("results subscriber stopped")
			return
		default:
			s.poll(ctx)
		}
	}
}

func (s *Subscriber) poll(ctx context.Context) {
	out, err := s.sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &s.queueURL,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
		VisibilityTimeout:   60,
	})
	if err != nil {
		log.Printf("results subscriber receive error: %v", err)
		time.Sleep(5 * time.Second)
		return
	}

	for _, msg := range out.Messages {
		if err := s.handle(ctx, msg.Body); err != nil {
			log.Printf("results subscriber handle error (message stays in queue): %v", err)
			continue
		}
		s.sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      &s.queueURL,
			ReceiptHandle: msg.ReceiptHandle,
		})
	}
}

type envelope struct {
	EventType string `json:"eventType"`
}

func (s *Subscriber) handle(ctx context.Context, body *string) error {
	var env envelope
	if err := json.Unmarshal([]byte(*body), &env); err != nil {
		return err
	}

	switch env.EventType {
	case "VideoProcessed":
		var evt application.VideoProcessedEvent
		if err := json.Unmarshal([]byte(*body), &evt); err != nil {
			return err
		}
		return s.svc.OnProcessed(ctx, evt)

	case "VideoFailed":
		var evt application.VideoFailedEvent
		if err := json.Unmarshal([]byte(*body), &evt); err != nil {
			return err
		}
		return s.svc.OnFailed(ctx, evt)

	default:
		log.Printf("results subscriber: unknown eventType %q, skipping", env.EventType)
		return nil
	}
}
