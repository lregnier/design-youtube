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
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/lregnier/design-youtube/worker/internal/config"
	"github.com/lregnier/design-youtube/worker/internal/event"
	"github.com/lregnier/design-youtube/worker/internal/queue"
)

type processingJob struct {
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
	s3Client := awss3.NewFromConfig(awsCfg)
	publisher := queue.NewPublisher(sqsClient, cfg.ResultsQueueURL)

	log.Println("worker started, polling SQS")
	poll(context.Background(), cfg, sqsClient, s3Client, publisher)
}

func poll(ctx context.Context, cfg *config.Config, sqsClient *sqs.Client, s3Client *awss3.Client, pub *queue.Publisher) {
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
			if err := process(ctx, cfg, s3Client, pub, msg.Body); err != nil {
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

func process(ctx context.Context, cfg *config.Config, s3Client *awss3.Client, pub *queue.Publisher, body *string) error {
	var job processingJob
	var s3Notification struct {
		Records []struct {
			S3 struct {
				Object struct{ Key string `json:"key"` } `json:"object"`
			} `json:"s3"`
		} `json:"Records"`
	}
	if err := json.Unmarshal([]byte(*body), &s3Notification); err == nil && len(s3Notification.Records) > 0 {
		s3Key := s3Notification.Records[0].S3.Object.Key
		parts := strings.Split(s3Key, "/")
		if len(parts) >= 2 {
			job.VideoID = parts[1]
			job.S3Key = s3Key
		}
	} else {
		if err := json.Unmarshal([]byte(*body), &job); err != nil {
			return fmt.Errorf("parse message: %w", err)
		}
	}

	if job.VideoID == "" {
		return fmt.Errorf("could not determine videoId from message")
	}

	log.Printf("processing videoId=%s", job.VideoID)

	tmpDir, err := os.MkdirTemp("", "video-"+job.VideoID+"-*")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	rawKey := fmt.Sprintf("raw/%s/original", job.VideoID)
	rawPath := filepath.Join(tmpDir, "original")

	if err := downloadS3(ctx, s3Client, cfg.S3Bucket, rawKey, rawPath); err != nil {
		emitFailed(ctx, pub, job.VideoID, fmt.Sprintf("download failed: %v", err))
		return nil
	}

	duration, err := videoDuration(rawPath)
	if err != nil {
		emitFailed(ctx, pub, job.VideoID, fmt.Sprintf("ffprobe failed: %v", err))
		return nil
	}

	qualities := []struct{ name, scale, bitrate string }{
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
			emitFailed(ctx, pub, job.VideoID, fmt.Sprintf("transcode %s failed: %v", q.name, err))
			return nil
		}
		if err := uploadDir(ctx, s3Client, cfg.S3Bucket, outDir, fmt.Sprintf("segments/%s/%s/", job.VideoID, q.name)); err != nil {
			return err
		}
	}

	manifestKey := fmt.Sprintf("manifests/%s/master.m3u8", job.VideoID)
	cloudfrontBase := fmt.Sprintf("https://%s/segments/%s", cfg.CloudFrontDomain, job.VideoID)
	masterContent := buildMasterManifest(cloudfrontBase, qualities)
	if err := uploadBytes(ctx, s3Client, cfg.S3Bucket, manifestKey, []byte(masterContent), "application/x-mpegURL"); err != nil {
		return err
	}

	thumbKey := fmt.Sprintf("thumbnails/%s/thumb.jpg", job.VideoID)
	thumbPath := filepath.Join(tmpDir, "thumb.jpg")
	if err := extractThumbnail(rawPath, thumbPath, duration/2); err != nil {
		log.Printf("thumbnail extraction failed (non-fatal): %v", err)
	} else if data, err := os.ReadFile(thumbPath); err == nil {
		uploadBytes(ctx, s3Client, cfg.S3Bucket, thumbKey, data, "image/jpeg")
	}

	manifestURL := fmt.Sprintf("https://%s/%s", cfg.CloudFrontDomain, manifestKey)
	thumbURL := fmt.Sprintf("https://%s/%s", cfg.CloudFrontDomain, thumbKey)

	if err := pub.Emit(ctx, job.VideoID, event.VideoProcessed{
		EventType:    event.TypeVideoProcessed,
		VideoID:      job.VideoID,
		ManifestURL:  manifestURL,
		ThumbnailURL: thumbURL,
	}); err != nil {
		return fmt.Errorf("emit VideoProcessed: %w", err)
	}

	log.Printf("completed videoId=%s", job.VideoID)
	return nil
}

func emitFailed(ctx context.Context, pub *queue.Publisher, videoID, reason string) {
	log.Printf("emitting VideoFailed for videoId=%s: %s", videoID, reason)
	pub.Emit(ctx, videoID, event.VideoFailed{
		EventType: event.TypeVideoFailed,
		VideoID:   videoID,
		Reason:    reason,
	})
}

func downloadS3(ctx context.Context, s3c *awss3.Client, bucket, key, dest string) error {
	out, err := s3c.GetObject(ctx, &awss3.GetObjectInput{Bucket: &bucket, Key: &key})
	if err != nil {
		return err
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

func uploadBytes(ctx context.Context, s3c *awss3.Client, bucket, key string, data []byte, ct string) error {
	_, err := s3c.PutObject(ctx, &awss3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        strings.NewReader(string(data)),
		ContentType: &ct,
	})
	return err
}

func uploadDir(ctx context.Context, s3c *awss3.Client, bucket, dir, prefix string) error {
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
		return uploadBytes(ctx, s3c, bucket, key, data, ct)
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
		"-hls_time", "6", "-hls_playlist_type", "vod",
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
		"-i", input, "-frames:v", "1", "-q:v", "2", "-y", dest,
	).Run()
}

func buildMasterManifest(cloudfrontBase string, qualities []struct{ name, scale, bitrate string }) string {
	var sb strings.Builder
	sb.WriteString("#EXTM3U\n")
	bw := map[string]int{"1080p": 4500000, "720p": 2800000, "360p": 1000000}
	res := map[string]string{"1080p": "1920x1080", "720p": "1280x720", "360p": "640x360"}
	for _, q := range qualities {
		sb.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%s\n", bw[q.name], res[q.name]))
		sb.WriteString(fmt.Sprintf("%s/%s/media.m3u8\n", cloudfrontBase, q.name))
	}
	return sb.String()
}
