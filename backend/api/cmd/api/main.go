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

	httpadapter "github.com/lregnier/design-youtube/api/internal/adapters/inbound/http"
	"github.com/lregnier/design-youtube/api/internal/adapters/inbound/sqsconsumer"
	"github.com/lregnier/design-youtube/api/internal/adapters/outbound/dynamo"
	"github.com/lregnier/design-youtube/api/internal/adapters/outbound/rediscache"
	"github.com/lregnier/design-youtube/api/internal/adapters/outbound/s3store"
	"github.com/lregnier/design-youtube/api/internal/adapters/outbound/sqsqueue"
	"github.com/lregnier/design-youtube/api/internal/application/catalog"
	"github.com/lregnier/design-youtube/api/internal/application/processing"
	"github.com/lregnier/design-youtube/api/internal/application/upload"
	"github.com/lregnier/design-youtube/api/internal/config"
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
	dynamoOpts := []func(*dynamodb.Options){}
	if cfg.DynamoDBEndpoint != "" {
		dynamoOpts = append(dynamoOpts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(cfg.DynamoDBEndpoint)
		})
	}
	repo := dynamo.NewRepository(dynamodb.NewFromConfig(awsCfg, dynamoOpts...), cfg.DynamoDBTable)

	s3Opts := []func(*awss3.Options){}
	if cfg.S3Endpoint != "" {
		s3Opts = append(s3Opts, func(o *awss3.Options) {
			o.BaseEndpoint = aws.String(cfg.S3Endpoint)
			o.UsePathStyle = true
		})
	}
	var transformer s3store.PresignedURLTransformer = s3store.NoOpTransformer{}
	if cfg.S3PublicURL != "" {
		transformer = s3store.NewEndpointTransformer(cfg.S3PublicURL)
	}
	store := s3store.NewStore(awss3.NewFromConfig(awsCfg, s3Opts...), cfg.S3Bucket, transformer)

	cache := rediscache.NewCache(redis.NewClient(&redis.Options{Addr: cfg.RedisAddr}))

	sqsOpts := []func(*sqs.Options){}
	if cfg.SQSEndpoint != "" {
		sqsOpts = append(sqsOpts, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(cfg.SQSEndpoint)
		})
	}
	sqsClient := sqs.NewFromConfig(awsCfg, sqsOpts...)
	processingQueue := sqsqueue.NewQueue(sqsClient, cfg.SQSQueueURL)

	// Use cases
	initUC := upload.NewInitUpload(repo, store, cfg.S3Bucket)
	confirmUC := upload.NewConfirmChunk(repo, store)
	completeUC := upload.NewCompleteUpload(repo, store, processingQueue)
	getVideoUC := catalog.NewGetVideo(repo, cache)
	listVideosUC := catalog.NewListVideos(repo)
	applyResultUC := processing.NewApplyProcessingResult(repo)

	// Inbound adapters
	h := httpadapter.NewHandler(initUC, confirmUC, completeUC, getVideoUC, listVideosUC)
	srv := httpadapter.NewServer(h, cfg.UploadSecret, strings.Split(cfg.CORSAllowedOrigins, ","), cfg.HTTPAddr)

	consumer := sqsconsumer.NewConsumer(sqs.NewFromConfig(awsCfg, sqsOpts...), cfg.ResultsQueueURL, applyResultUC)
	go consumer.Start(context.Background())

	if err := srv.Start(); err != nil {
		log.Fatalf("server: %v", err)
	}
}
