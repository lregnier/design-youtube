package main

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/redis/go-redis/v9"

	"github.com/lregnier/design-youtube/api/internal/application"
	"github.com/lregnier/design-youtube/api/internal/domain/video"
	httpadapter "github.com/lregnier/design-youtube/api/internal/infrastructure/in/http"
	"github.com/lregnier/design-youtube/api/internal/infrastructure/in/sqssubscriber"
	"github.com/lregnier/design-youtube/api/internal/infrastructure/out/dynamo"
	"github.com/lregnier/design-youtube/api/internal/infrastructure/out/rediscache"
	"github.com/lregnier/design-youtube/api/internal/infrastructure/out/s3store"
	"github.com/lregnier/design-youtube/api/internal/infrastructure/out/sqspublisher"
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

	// Infrastructure — outbound
	repo := newRepository(cfg, awsCfg)
	store := newStore(cfg, awsCfg)
	cache := rediscache.NewCache(redis.NewClient(&redis.Options{Addr: cfg.RedisAddr}))
	publisher := newPublisher(cfg, awsCfg)

	// Application services
	uploadSvc := application.NewUploadService(repo, store, publisher, cfg.S3Bucket)
	catalogSvc := application.NewCatalogService(repo, cache)
	processingSvc := application.NewVideoStatusService(repo)

	// Infrastructure — inbound
	h := httpadapter.NewHandler(uploadSvc, catalogSvc)
	srv := httpadapter.NewServer(h, cfg.UploadSecret, strings.Split(cfg.CORSAllowedOrigins, ","), cfg.HTTPAddr)

	subscriber := sqssubscriber.NewSubscriber(newSQSClient(cfg, awsCfg), cfg.ResultsQueueURL, processingSvc)
	go subscriber.Start(context.Background())

	if err := srv.Start(); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func newRepository(cfg *Config, awsCfg aws.Config) video.VideoRepository {
	opts := []func(*dynamodb.Options){}
	if cfg.DynamoDBEndpoint != "" {
		opts = append(opts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(cfg.DynamoDBEndpoint)
		})
	}
	return dynamo.NewRepository(dynamodb.NewFromConfig(awsCfg, opts...), cfg.DynamoDBTable)
}

func newStore(cfg *Config, awsCfg aws.Config) application.ObjectStore {
	opts := []func(*awss3.Options){}
	if cfg.S3Endpoint != "" {
		opts = append(opts, func(o *awss3.Options) {
			o.BaseEndpoint = aws.String(cfg.S3Endpoint)
			o.UsePathStyle = true
		})
	}
	var transformer s3store.PresignedURLTransformer = s3store.NoOpTransformer{}
	if cfg.S3PublicURL != "" {
		transformer = s3store.NewEndpointTransformer(cfg.S3PublicURL)
	}
	return s3store.NewStore(awss3.NewFromConfig(awsCfg, opts...), cfg.S3Bucket, transformer)
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
	return sqspublisher.NewPublisher(newSQSClient(cfg, awsCfg), cfg.SQSQueueURL)
}
