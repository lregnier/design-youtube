package main

import (
	"context"
	"log"
	nethttp "net/http"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"

	httpadapter "github.com/lregnier/design-youtube/api/internal/adapters/inbound/http"
	"github.com/lregnier/design-youtube/api/internal/adapters/inbound/sqsconsumer"
	"github.com/lregnier/design-youtube/api/internal/adapters/outbound/dynamo"
	"github.com/lregnier/design-youtube/api/internal/adapters/outbound/rediscache"
	"github.com/lregnier/design-youtube/api/internal/adapters/outbound/s3store"
	"github.com/lregnier/design-youtube/api/internal/api"
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
	repo := dynamo.NewRepository(dynamodb.NewFromConfig(awsCfg), cfg.DynamoDBTable)
	store := s3store.NewStore(awss3.NewFromConfig(awsCfg), cfg.S3Bucket)
	cache := rediscache.NewCache(redis.NewClient(&redis.Options{Addr: cfg.RedisAddr}))
	sqsClient := sqs.NewFromConfig(awsCfg)

	// Use cases
	initUC := upload.NewInitUpload(repo, store, cfg.S3Bucket)
	confirmUC := upload.NewConfirmChunk(repo, store)
	completeUC := upload.NewCompleteUpload(repo, store)
	getVideoUC := catalog.NewGetVideo(repo, cache)
	listVideosUC := catalog.NewListVideos(repo)
	applyResultUC := processing.NewApplyProcessingResult(repo)

	// Inbound adapters
	h := httpadapter.NewHandler(initUC, confirmUC, completeUC, getVideoUC, listVideosUC)
	secretMw := httpadapter.UploadSecretMiddleware(cfg.UploadSecret)
	strictHandler := api.NewStrictHandlerWithOptions(h, []api.StrictMiddlewareFunc{secretMw}, api.StrictHTTPServerOptions{})

	consumer := sqsconsumer.NewConsumer(sqsClient, cfg.ResultsQueueURL, applyResultUC)
	go consumer.Start(context.Background())

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Get("/health", func(w nethttp.ResponseWriter, r *nethttp.Request) { w.WriteHeader(nethttp.StatusOK) })
	api.HandlerFromMux(strictHandler, r)

	log.Println("listening on :8080")
	if err := nethttp.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
