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
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/lregnier/design-youtube/backend/internal/adapters/outbound/dynamo"
	"github.com/lregnier/design-youtube/backend/internal/adapters/outbound/s3store"
	"github.com/lregnier/design-youtube/backend/internal/adapters/outbound/sqsqueue"
	"github.com/lregnier/design-youtube/backend/internal/config"
	"github.com/lregnier/design-youtube/backend/internal/domain/video"
	"github.com/lregnier/design-youtube/backend/internal/ports"
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

	repo := dynamo.NewRepository(dynamodb.NewFromConfig(awsCfg), cfg.DynamoDBTable)
	store := s3store.NewStore(awss3.NewFromConfig(awsCfg), cfg.S3Bucket)
	queue := sqsqueue.NewQueue(sqs.NewFromConfig(awsCfg), cfg.SQSQueueURL)

	log.Println("worker started, polling SQS")
	poll(context.Background(), cfg, repo, store, queue, sqs.NewFromConfig(awsCfg))
}

func poll(ctx context.Context, cfg *config.Config, repo *dynamo.Repository, store *s3store.Store, queue *sqsqueue.Queue, sqsClient *sqs.Client) {
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
			if err := process(ctx, cfg, repo, store, msg.Body); err != nil {
				log.Printf("process error (message stays in queue): %v", err)
				continue
			}
			queue.DeleteMessage(ctx, *msg.ReceiptHandle)
		}
	}
}

func process(ctx context.Context, cfg *config.Config, repo *dynamo.Repository, store *s3store.Store, body *string) error {
	var event sqsEvent
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

	rawKey := fmt.Sprintf("raw/%s/original", event.VideoID)
	rawPath := filepath.Join(tmpDir, "original")

	rawData, err := store.GetObject(ctx, rawKey)
	if err != nil {
		markFailed(ctx, repo, event.VideoID, err)
		return nil
	}
	if err := os.WriteFile(rawPath, rawData, 0644); err != nil {
		return err
	}

	duration, err := videoDuration(rawPath)
	if err != nil {
		markFailed(ctx, repo, event.VideoID, err)
		return nil
	}

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
			markFailed(ctx, repo, event.VideoID, err)
			return nil
		}
		if err := uploadDir(ctx, store, outDir, fmt.Sprintf("segments/%s/%s/", event.VideoID, q.name)); err != nil {
			return err
		}
	}

	masterPath := filepath.Join(tmpDir, "master.m3u8")
	cloudfrontBase := fmt.Sprintf("https://%s/segments/%s", cfg.CloudFrontDomain, event.VideoID)
	masterContent := buildMasterManifest(cloudfrontBase, qualities)
	if err := os.WriteFile(masterPath, []byte(masterContent), 0644); err != nil {
		return err
	}
	manifestKey := fmt.Sprintf("manifests/%s/master.m3u8", event.VideoID)
	if err := store.PutObject(ctx, manifestKey, []byte(masterContent), "application/x-mpegURL"); err != nil {
		return err
	}

	thumbPath := filepath.Join(tmpDir, "thumb.jpg")
	if err := extractThumbnail(rawPath, thumbPath, duration/2); err != nil {
		log.Printf("thumbnail extraction failed (non-fatal): %v", err)
	} else {
		thumbData, err := os.ReadFile(thumbPath)
		if err == nil {
			thumbKey := fmt.Sprintf("thumbnails/%s/thumb.jpg", event.VideoID)
			store.PutObject(ctx, thumbKey, thumbData, "image/jpeg")
		}
	}

	manifestURL := fmt.Sprintf("https://%s/%s", cfg.CloudFrontDomain, manifestKey)
	thumbURL := fmt.Sprintf("https://%s/thumbnails/%s/thumb.jpg", cfg.CloudFrontDomain, event.VideoID)

	vid, err := repo.FindByID(ctx, video.VideoID(event.VideoID))
	if err != nil || vid == nil {
		log.Printf("could not find video to mark ready: %v", err)
		return nil
	}
	vid.MarkReady(manifestURL, thumbURL)
	if err := repo.Save(ctx, vid); err != nil {
		log.Printf("failed to mark video ready: %v", err)
	}

	log.Printf("completed videoId=%s", event.VideoID)
	return nil
}

func markFailed(ctx context.Context, repo *dynamo.Repository, videoID string, reason error) {
	log.Printf("marking videoId=%s failed: %v", videoID, reason)
	vid, err := repo.FindByID(ctx, video.VideoID(videoID))
	if err != nil || vid == nil {
		return
	}
	vid.MarkFailed()
	repo.Save(ctx, vid)
}

func uploadDir(ctx context.Context, store ports.ObjectStore, dir, prefix string) error {
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
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return store.PutObject(ctx, key, data, ct)
	})
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
	return exec.Command("ffmpeg",
		"-ss", fmt.Sprintf("%.2f", offset),
		"-i", input,
		"-frames:v", "1",
		"-q:v", "2",
		"-y", dest,
	).Run()
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
