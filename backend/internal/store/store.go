package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/redis/go-redis/v9"

	"github.com/lregnier/design-youtube/backend/internal/config"
)

const (
	StatusUploading  = "uploading"
	StatusProcessing = "processing"
	StatusReady      = "ready"
	StatusFailed     = "failed"

	cacheTTL = 60 * time.Second
)

type VideoRecord struct {
	VideoID      string       `dynamodbav:"videoId" json:"videoId"`
	Title        string       `dynamodbav:"title" json:"title"`
	Description  string       `dynamodbav:"description" json:"description"`
	Status       string       `dynamodbav:"status" json:"status"`
	UploadedAt   string       `dynamodbav:"uploadedAt" json:"uploadedAt"`
	UploadID     string       `dynamodbav:"uploadId" json:"uploadId"`
	TotalChunks  int          `dynamodbav:"totalChunks" json:"totalChunks"`
	Chunks       []ChunkState `dynamodbav:"chunks" json:"chunks"`
	ManifestURL  string       `dynamodbav:"manifestUrl,omitempty" json:"manifestUrl,omitempty"`
	ThumbnailURL string       `dynamodbav:"thumbnailUrl,omitempty" json:"thumbnailUrl,omitempty"`
}

type ChunkState struct {
	PartNumber int    `dynamodbav:"partNumber" json:"partNumber"`
	Uploaded   bool   `dynamodbav:"uploaded" json:"uploaded"`
	ETag       string `dynamodbav:"eTag,omitempty" json:"eTag,omitempty"`
}

type Store struct {
	ddb       *dynamodb.Client
	rdb       *redis.Client
	tableName string
}

func New(cfg *config.Config) (*Store, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})

	return &Store{
		ddb:       dynamodb.NewFromConfig(awsCfg),
		rdb:       rdb,
		tableName: cfg.DynamoDBTable,
	}, nil
}

func (s *Store) PutVideo(ctx context.Context, v *VideoRecord) error {
	item, err := attributevalue.MarshalMap(v)
	if err != nil {
		return err
	}
	_, err = s.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &s.tableName,
		Item:      item,
	})
	return err
}

func (s *Store) GetVideo(ctx context.Context, videoID string) (*VideoRecord, error) {
	cacheKey := "video:" + videoID

	cached, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var v VideoRecord
		if jsonErr := json.Unmarshal([]byte(cached), &v); jsonErr == nil {
			return &v, nil
		}
	}

	out, err := s.ddb.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &s.tableName,
		Key: map[string]types.AttributeValue{
			"videoId": &types.AttributeValueMemberS{Value: videoID},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}

	var v VideoRecord
	if err := attributevalue.UnmarshalMap(out.Item, &v); err != nil {
		return nil, err
	}

	if data, err := json.Marshal(&v); err == nil {
		s.rdb.Set(ctx, cacheKey, data, cacheTTL)
	}

	return &v, nil
}

func (s *Store) UpdateVideoStatus(ctx context.Context, videoID, status string) error {
	_, err := s.ddb.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &s.tableName,
		Key: map[string]types.AttributeValue{
			"videoId": &types.AttributeValueMemberS{Value: videoID},
		},
		UpdateExpression: aws.String("SET #st = :st"),
		ExpressionAttributeNames:  map[string]string{"#st": "status"},
		ExpressionAttributeValues: map[string]types.AttributeValue{":st": &types.AttributeValueMemberS{Value: status}},
	})
	return err
}

func (s *Store) MarkChunkUploaded(ctx context.Context, videoID string, partNumber int, eTag string) error {
	v, err := s.GetVideo(ctx, videoID)
	if err != nil {
		return err
	}
	if v == nil {
		return errors.New("video not found")
	}
	for i := range v.Chunks {
		if v.Chunks[i].PartNumber == partNumber {
			v.Chunks[i].Uploaded = true
			v.Chunks[i].ETag = eTag
		}
	}
	return s.PutVideo(ctx, v)
}

func (s *Store) ListReadyVideos(ctx context.Context) ([]VideoRecord, error) {
	out, err := s.ddb.Query(ctx, &dynamodb.QueryInput{
		TableName:              &s.tableName,
		IndexName:              aws.String("status-uploadedAt-index"),
		KeyConditionExpression: aws.String("#st = :st"),
		ExpressionAttributeNames:  map[string]string{"#st": "status"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":st": &types.AttributeValueMemberS{Value: StatusReady},
		},
		ScanIndexForward: aws.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	var videos []VideoRecord
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &videos); err != nil {
		return nil, err
	}
	return videos, nil
}
