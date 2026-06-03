package main

import (
	"context"
	"log"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/lregnier/design-youtube/worker/internal/adapters/inbound/sqsjobs"
	"github.com/lregnier/design-youtube/worker/internal/adapters/outbound/ffmpeg"
	"github.com/lregnier/design-youtube/worker/internal/adapters/outbound/s3storage"
	"github.com/lregnier/design-youtube/worker/internal/adapters/outbound/sqspublisher"
	"github.com/lregnier/design-youtube/worker/internal/application"
	"github.com/lregnier/design-youtube/worker/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		log.Fatalf("aws config: %v", err)
	}

	// Outbound adapters
	store := s3storage.NewStore(awss3.NewFromConfig(awsCfg), cfg.S3Bucket, cfg.CloudFrontDomain)
	transcoder := ffmpeg.NewTranscoder()
	publisher := sqspublisher.NewPublisher(sqs.NewFromConfig(awsCfg), cfg.ResultsQueueURL)

	// Use case
	processVideo := application.NewProcessVideo(store, transcoder, publisher)

	// Inbound adapter
	consumer := sqsjobs.NewConsumer(sqs.NewFromConfig(awsCfg), cfg.SQSQueueURL, processVideo)
	consumer.Start(context.Background())
}
