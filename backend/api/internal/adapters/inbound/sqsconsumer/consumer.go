package sqsconsumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/lregnier/design-youtube/api/internal/application/processing"
)

type Consumer struct {
	sqsClient  *sqs.Client
	queueURL   string
	useCase    processing.ApplyProcessingResult
}

func NewConsumer(sqsClient *sqs.Client, queueURL string, uc processing.ApplyProcessingResult) *Consumer {
	return &Consumer{sqsClient: sqsClient, queueURL: queueURL, useCase: uc}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("results consumer started")
	for {
		select {
		case <-ctx.Done():
			log.Println("results consumer stopped")
			return
		default:
			c.poll(ctx)
		}
	}
}

func (c *Consumer) poll(ctx context.Context) {
	out, err := c.sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &c.queueURL,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
		VisibilityTimeout:   60,
	})
	if err != nil {
		log.Printf("results consumer receive error: %v", err)
		time.Sleep(5 * time.Second)
		return
	}

	for _, msg := range out.Messages {
		if err := c.handle(ctx, msg.Body); err != nil {
			log.Printf("results consumer handle error (message stays in queue): %v", err)
			continue
		}
		c.sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      &c.queueURL,
			ReceiptHandle: msg.ReceiptHandle,
		})
	}
}

type envelope struct {
	EventType string `json:"eventType"`
}

func (c *Consumer) handle(ctx context.Context, body *string) error {
	var env envelope
	if err := json.Unmarshal([]byte(*body), &env); err != nil {
		return err
	}

	switch env.EventType {
	case "VideoProcessed":
		var evt processing.VideoProcessedEvent
		if err := json.Unmarshal([]byte(*body), &evt); err != nil {
			return err
		}
		return c.useCase.OnProcessed(ctx, evt)

	case "VideoFailed":
		var evt processing.VideoFailedEvent
		if err := json.Unmarshal([]byte(*body), &evt); err != nil {
			return err
		}
		return c.useCase.OnFailed(ctx, evt)

	default:
		log.Printf("results consumer: unknown eventType %q, skipping", env.EventType)
		return nil
	}
}
