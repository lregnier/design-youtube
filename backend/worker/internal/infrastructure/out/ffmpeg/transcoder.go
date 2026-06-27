package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lregnier/design-youtube/worker/internal/application"
)

var _ application.Transcoder = (*transcoder)(nil)

type transcoder struct{}

func NewTranscoder() application.Transcoder { return &transcoder{} }

func (t *transcoder) Duration(_ context.Context, inputPath string) (float64, error) {
	out, err := exec.Command("ffprobe", "-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputPath).Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe: %w", err)
	}
	d, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		return 0, fmt.Errorf("parse duration: %w", err)
	}
	return d, nil
}

func (t *transcoder) TranscodeHLS(_ context.Context, inputPath, outputDir, scale, bitrate string) error {
	playlist := filepath.Join(outputDir, "media.m3u8")
	segment := filepath.Join(outputDir, "seg%03d.ts")
	cmd := exec.Command("ffmpeg", "-i", inputPath,
		"-vf", fmt.Sprintf("scale=%s", scale),
		"-c:v", "libx264", "-b:v", bitrate,
		"-c:a", "aac", "-b:a", "128k",
		"-hls_time", "6", "-hls_playlist_type", "vod",
		"-hls_segment_filename", segment,
		"-y", playlist,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg transcode: %w\n%s", err, out)
	}
	return nil
}

func (t *transcoder) ExtractThumbnail(_ context.Context, inputPath string, offset float64) ([]byte, error) {
	tmp, err := os.CreateTemp("", "thumb-*.jpg")
	if err != nil {
		return nil, fmt.Errorf("temp file: %w", err)
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	out, err := exec.Command("ffmpeg",
		"-ss", fmt.Sprintf("%.2f", offset),
		"-i", inputPath,
		"-frames:v", "1", "-q:v", "2",
		"-y", tmp.Name(),
	).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg thumbnail: %w\n%s", err, out)
	}
	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		return nil, fmt.Errorf("read thumbnail: %w", err)
	}
	return data, nil
}
