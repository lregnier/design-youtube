package dynamo

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/lregnier/design-youtube/backend/internal/domain/video"
)

var _ video.VideoRepository = (*Repository)(nil)

type record struct {
	VideoID      string       `dynamodbav:"videoId"`
	Title        string       `dynamodbav:"title"`
	Description  string       `dynamodbav:"description"`
	Status       string       `dynamodbav:"status"`
	UploadedAt   string       `dynamodbav:"uploadedAt"`
	UploadID     string       `dynamodbav:"uploadId"`
	TotalChunks  int          `dynamodbav:"totalChunks"`
	Chunks       []chunkRecord `dynamodbav:"chunks"`
	ManifestURL  string       `dynamodbav:"manifestUrl,omitempty"`
	ThumbnailURL string       `dynamodbav:"thumbnailUrl,omitempty"`
}

type chunkRecord struct {
	PartNumber int    `dynamodbav:"partNumber"`
	Uploaded   bool   `dynamodbav:"uploaded"`
	ETag       string `dynamodbav:"eTag,omitempty"`
}

type Repository struct {
	client    *dynamodb.Client
	tableName string
}

func NewRepository(client *dynamodb.Client, tableName string) *Repository {
	return &Repository{client: client, tableName: tableName}
}

func (r *Repository) Save(ctx context.Context, v *video.Video) error {
	rec := toRecord(v)
	item, err := attributevalue.MarshalMap(rec)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	})
	return err
}

func (r *Repository) FindByID(ctx context.Context, id video.VideoID) (*video.Video, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"videoId": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}
	var rec record
	if err := attributevalue.UnmarshalMap(out.Item, &rec); err != nil {
		return nil, err
	}
	return toDomain(&rec), nil
}

func (r *Repository) ListReady(ctx context.Context) ([]*video.Video, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              &r.tableName,
		IndexName:              aws.String("status-uploadedAt-index"),
		KeyConditionExpression: aws.String("#st = :st"),
		ExpressionAttributeNames:  map[string]string{"#st": "status"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":st": &types.AttributeValueMemberS{Value: string(video.StatusReady)},
		},
		ScanIndexForward: aws.Bool(false),
	})
	if err != nil {
		return nil, err
	}
	var records []record
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &records); err != nil {
		return nil, err
	}
	videos := make([]*video.Video, len(records))
	for i, rec := range records {
		videos[i] = toDomain(&rec)
	}
	return videos, nil
}

func toRecord(v *video.Video) record {
	chunks := make([]chunkRecord, len(v.Chunks))
	for i, c := range v.Chunks {
		chunks[i] = chunkRecord{PartNumber: c.PartNumber, Uploaded: c.Uploaded, ETag: c.ETag}
	}
	return record{
		VideoID:      v.ID.String(),
		Title:        v.Title,
		Description:  v.Description,
		Status:       string(v.Status),
		UploadedAt:   v.UploadedAt.UTC().Format(time.RFC3339),
		UploadID:     v.UploadID,
		TotalChunks:  v.TotalChunks,
		Chunks:       chunks,
		ManifestURL:  v.ManifestURL,
		ThumbnailURL: v.ThumbnailURL,
	}
}

func toDomain(rec *record) *video.Video {
	chunks := make([]video.Chunk, len(rec.Chunks))
	for i, c := range rec.Chunks {
		chunks[i] = video.Chunk{PartNumber: c.PartNumber, Uploaded: c.Uploaded, ETag: c.ETag}
	}
	t, _ := time.Parse(time.RFC3339, rec.UploadedAt)
	return &video.Video{
		ID:           video.VideoID(rec.VideoID),
		Title:        rec.Title,
		Description:  rec.Description,
		Status:       video.VideoStatus(rec.Status),
		UploadedAt:   t,
		UploadID:     rec.UploadID,
		TotalChunks:  rec.TotalChunks,
		Chunks:       chunks,
		ManifestURL:  rec.ManifestURL,
		ThumbnailURL: rec.ThumbnailURL,
	}
}
