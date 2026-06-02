package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/lregnier/design-youtube/backend/internal/config"
)

type sqsEvent struct {
	VideoID string `json:"videoId"`
	S3Key   string `json:"s3Key"`
}

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

	sqsClient := sqs.NewFromConfig(awsCfg)
	s3Client := s3.NewFromConfig(awsCfg)
	ddbClient := dynamodb.NewFromConfig(awsCfg)

	log.Println("worker started, polling SQS")
	poll(context.Background(), cfg, sqsClient, s3Client, ddbClient)
}

func poll(ctx context.Context, cfg *config.Config, sqsClient *sqs.Client, s3Client *s3.Client, ddbClient *dynamodb.Client) {
	for {
		out, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &cfg.SQSQueueURL,
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     20,
			VisibilityTimeout:   900,
		})
		if err != nil {
			log.Printf("sqs receive: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, msg := range out.Messages {
			if err := process(ctx, cfg, s3Client, ddbClient, msg.Body); err != nil {
				log.Printf("process error (message stays in queue): %v", err)
				continue
			}
			sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      &cfg.SQSQueueURL,
				ReceiptHandle: msg.ReceiptHandle,
			})
		}
	}
}

func process(ctx context.Context, cfg *config.Config, s3Client *s3.Client, ddbClient *dynamodb.Client, body *string) error {
	var event sqsEvent
	// S3 event notifications wrap the record in a Records array
	var s3Notification struct {
		Records []struct {
			S3 struct {
				Object struct {
					Key string `json:"key"`
				} `json:"object"`
			} `json:"s3"`
		} `json:"Records"`
	}
	if err := json.Unmarshal([]byte(*body), &s3Notification); err == nil && len(s3Notification.Records) > 0 {
		s3Key := s3Notification.Records[0].S3.Object.Key
		// key format: raw/{videoId}/original
		parts := strings.Split(s3Key, "/")
		if len(parts) >= 2 {
			event.VideoID = parts[1]
			event.S3Key = s3Key
		}
	} else {
		if err := json.Unmarshal([]byte(*body), &event); err != nil {
			return fmt.Errorf("parse message: %w", err)
		}
	}

	if event.VideoID == "" {
		return fmt.Errorf("could not determine videoId from message")
	}

	log.Printf("processing videoId=%s", event.VideoID)

	tmpDir, err := os.MkdirTemp("", "video-"+event.VideoID+"-*")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Download raw video
	rawKey := fmt.Sprintf("raw/%s/original", event.VideoID)
	rawPath := filepath.Join(tmpDir, "original")
	if err := downloadS3(ctx, s3Client, cfg.S3Bucket, rawKey, rawPath); err != nil {
		markFailed(ctx, ddbClient, cfg.DynamoDBTable, event.VideoID, err)
		return nil // delete message, don't retry download failures
	}

	// Get video duration for thumbnail
	duration, err := videoDuration(rawPath)
	if err != nil {
		markFailed(ctx, ddbClient, cfg.DynamoDBTable, event.VideoID, err)
		return nil
	}

	// Transcode to HLS at three quality levels
	qualities := []struct {
		name    string
		scale   string
		bitrate string
	}{
		{"1080p", "1920:1080", "4000k"},
		{"720p", "1280:720", "2500k"},
		{"360p", "640:360", "800k"},
	}

	segDir := filepath.Join(tmpDir, "segments")
	if err := os.MkdirAll(segDir, 0755); err != nil {
		return err
	}

	for _, q := range qualities {
		outDir := filepath.Join(segDir, q.name)
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return err
		}
		if err := transcode(rawPath, outDir, q.scale, q.bitrate); err != nil {
			markFailed(ctx, ddbClient, cfg.DynamoDBTable, event.VideoID, err)
			return nil
		}
		// Upload segments and media playlist
		if err := uploadDir(ctx, s3Client, cfg.S3Bucket, outDir, fmt.Sprintf("segments/%s/%s/", event.VideoID, q.name)); err != nil {
			return err
		}
	}

	// Generate master manifest
	masterPath := filepath.Join(tmpDir, "master.m3u8")
	cloudfrontBase := fmt.Sprintf("https://%s/segments/%s", cfg.CloudFrontDomain, event.VideoID)
	masterContent := buildMasterManifest(cloudfrontBase, qualities)
	if err := os.WriteFile(masterPath, []byte(masterContent), 0644); err != nil {
		return err
	}
	manifestKey := fmt.Sprintf("manifests/%s/master.m3u8", event.VideoID)
	if err := uploadFile(ctx, s3Client, cfg.S3Bucket, manifestKey, masterPath, "application/x-mpegURL"); err != nil {
		return err
	}

	// Extract thumbnail at midpoint
	thumbPath := filepath.Join(tmpDir, "thumb.jpg")
	midpoint := duration / 2
	if err := extractThumbnail(rawPath, thumbPath, midpoint); err != nil {
		log.Printf("thumbnail extraction failed (non-fatal): %v", err)
	} else {
		thumbKey := fmt.Sprintf("thumbnails/%s/thumb.jpg", event.VideoID)
		uploadFile(ctx, s3Client, cfg.S3Bucket, thumbKey, thumbPath, "image/jpeg")
	}

	// Update DynamoDB record to ready
	manifestURL := fmt.Sprintf("https://%s/%s", cfg.CloudFrontDomain, manifestKey)
	thumbURL := fmt.Sprintf("https://%s/thumbnails/%s/thumb.jpg", cfg.CloudFrontDomain, event.VideoID)
	updateReady(ctx, ddbClient, cfg.DynamoDBTable, event.VideoID, manifestURL, thumbURL)

	log.Printf("completed videoId=%s", event.VideoID)
	return nil
}

func downloadS3(ctx context.Context, s3c *s3.Client, bucket, key, dest string) error {
	out, err := s3c.GetObject(ctx, &s3.GetObjectInput{Bucket: &bucket, Key: &key})
	if err != nil {
		return fmt.Errorf("s3 get %s: %w", key, err)
	}
	defer out.Body.Close()
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := make([]byte, 1024*1024)
	for {
		n, rerr := out.Body.Read(buf)
		if n > 0 {
			f.Write(buf[:n])
		}
		if rerr != nil {
			break
		}
	}
	return nil
}

func videoDuration(path string) (float64, error) {
	out, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1", path).Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe: %w", err)
	}
	d, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		return 0, fmt.Errorf("parse duration: %w", err)
	}
	return d, nil
}

func transcode(input, outDir, scale, bitrate string) error {
	playlist := filepath.Join(outDir, "media.m3u8")
	segment := filepath.Join(outDir, "seg%03d.ts")
	cmd := exec.Command("ffmpeg", "-i", input,
		"-vf", fmt.Sprintf("scale=%s", scale),
		"-c:v", "libx264", "-b:v", bitrate,
		"-c:a", "aac", "-b:a", "128k",
		"-hls_time", "6",
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", segment,
		"-y", playlist,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func extractThumbnail(input, dest string, offset float64) error {
	cmd := exec.Command("ffmpeg",
		"-ss", fmt.Sprintf("%.2f", offset),
		"-i", input,
		"-frames:v", "1",
		"-q:v", "2",
		"-y", dest,
	)
	return cmd.Run()
}

func buildMasterManifest(cloudfrontBase string, qualities []struct{ name, scale, bitrate string }) string {
	var sb strings.Builder
	sb.WriteString("#EXTM3U\n")
	bandwidths := map[string]int{"1080p": 4500000, "720p": 2800000, "360p": 1000000}
	resolutions := map[string]string{"1080p": "1920x1080", "720p": "1280x720", "360p": "640x360"}
	for _, q := range qualities {
		sb.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%s\n", bandwidths[q.name], resolutions[q.name]))
		sb.WriteString(fmt.Sprintf("%s/%s/media.m3u8\n", cloudfrontBase, q.name))
	}
	return sb.String()
}

func uploadDir(ctx context.Context, s3c *s3.Client, bucket, dir, prefix string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		key := prefix + rel
		ct := "video/MP2T"
		if strings.HasSuffix(path, ".m3u8") {
			ct = "application/x-mpegURL"
		}
		return uploadFile(ctx, s3c, bucket, key, path, ct)
	})
}

func uploadFile(ctx context.Context, s3c *s3.Client, bucket, key, path, contentType string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = s3c.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        f,
		ContentType: &contentType,
	})
	return err
}

func markFailed(ctx context.Context, ddb *dynamodb.Client, table, videoID string, reason error) {
	log.Printf("marking videoId=%s failed: %v", videoID, reason)
	ddb.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &table,
		Key:       map[string]types.AttributeValue{"videoId": &types.AttributeValueMemberS{Value: videoID}},
		UpdateExpression: aws.String("SET #st = :st"),
		ExpressionAttributeNames:  map[string]string{"#st": "status"},
		ExpressionAttributeValues: map[string]types.AttributeValue{":st": &types.AttributeValueMemberS{Value: "failed"}},
	})
}

func updateReady(ctx context.Context, ddb *dynamodb.Client, table, videoID, manifestURL, thumbURL string) {
	ddb.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &table,
		Key:       map[string]types.AttributeValue{"videoId": &types.AttributeValueMemberS{Value: videoID}},
		UpdateExpression: aws.String("SET #st = :st, manifestUrl = :m, thumbnailUrl = :t"),
		ExpressionAttributeNames: map[string]string{"#st": "status"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":st": &types.AttributeValueMemberS{Value: "ready"},
			":m":  &types.AttributeValueMemberS{Value: manifestURL},
			":t":  &types.AttributeValueMemberS{Value: thumbURL},
		},
	})
}
