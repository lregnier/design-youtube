package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/lregnier/design-youtube/worker/internal/application"
	"github.com/lregnier/design-youtube/worker/internal/infrastructure/in/sqssubscriber"
	"github.com/lregnier/design-youtube/worker/internal/infrastructure/out/ffmpeg"
	"github.com/lregnier/design-youtube/worker/internal/infrastructure/out/s3storage"
	"github.com/lregnier/design-youtube/worker/internal/infrastructure/out/sqspublisher"
)

func main() {
	cfg, err := Load()
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
	store := newStore(cfg, awsCfg)
	transcoder := ffmpeg.NewTranscoder()
	publisher := newPublisher(cfg, awsCfg)

	// Use case
	svc := application.NewVideoProcessingService(store, transcoder, publisher)

	// Inbound adapter
	subscriber := sqssubscriber.NewSubscriber(newSQSClient(cfg, awsCfg), cfg.SQSQueueURL, svc, publisher)
	subscriber.Start(context.Background())
}

func newStore(cfg *Config, awsCfg aws.Config) application.VideoStorage {
	opts := []func(*awss3.Options){}
	if cfg.S3Endpoint != "" {
		opts = append(opts, func(o *awss3.Options) {
			o.BaseEndpoint = aws.String(cfg.S3Endpoint)
			o.UsePathStyle = true
		})
	}
	var urlBuilder s3storage.PublicURLBuilder
	if cfg.S3PublicURL != "" {
		urlBuilder = s3storage.NewEndpointURLBuilder(cfg.S3PublicURL)
	} else {
		urlBuilder = s3storage.NewCloudFrontURLBuilder(cfg.CloudFrontDomain)
	}
	return s3storage.NewStore(awss3.NewFromConfig(awsCfg, opts...), cfg.S3Bucket, urlBuilder)
}

func newSQSClient(cfg *Config, awsCfg aws.Config) *sqs.Client {
	opts := []func(*sqs.Options){}
	if cfg.SQSEndpoint != "" {
		opts = append(opts, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(cfg.SQSEndpoint)
		})
	}
	return sqs.NewFromConfig(awsCfg, opts...)
}

func newPublisher(cfg *Config, awsCfg aws.Config) application.EventPublisher {
	return sqspublisher.NewPublisher(newSQSClient(cfg, awsCfg), cfg.ResultsQueueURL)
}
