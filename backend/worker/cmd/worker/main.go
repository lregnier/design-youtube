package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	s3Opts := []func(*awss3.Options){}
	if cfg.S3Endpoint != "" {
		s3Opts = append(s3Opts, func(o *awss3.Options) {
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
	store := s3storage.NewStore(awss3.NewFromConfig(awsCfg, s3Opts...), cfg.S3Bucket, urlBuilder)

	transcoder := ffmpeg.NewTranscoder()

	sqsOpts := []func(*sqs.Options){}
	if cfg.SQSEndpoint != "" {
		sqsOpts = append(sqsOpts, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(cfg.SQSEndpoint)
		})
	}
	publisher := sqspublisher.NewPublisher(sqs.NewFromConfig(awsCfg, sqsOpts...), cfg.ResultsQueueURL)

	// Use case
	processVideo := application.NewProcessVideo(store, transcoder, publisher)

	// Inbound adapter
	consumer := sqsjobs.NewConsumer(sqs.NewFromConfig(awsCfg, sqsOpts...), cfg.SQSQueueURL, processVideo, publisher)
	consumer.Start(context.Background())
}
