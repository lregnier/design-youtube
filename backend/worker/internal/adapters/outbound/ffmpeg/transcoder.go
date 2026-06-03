package ffmpeg

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lregnier/design-youtube/worker/internal/ports"
)

var _ ports.Transcoder = (*Transcoder)(nil)

type Transcoder struct{}

func NewTranscoder() *Transcoder { return &Transcoder{} }

func (t *Transcoder) Duration(_ context.Context, inputPath string) (float64, error) {
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

func (t *Transcoder) TranscodeHLS(_ context.Context, inputPath, outputDir, scale, bitrate string) error {
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

func (t *Transcoder) ExtractThumbnail(_ context.Context, inputPath, outputPath string, offset float64) error {
	out, err := exec.Command("ffmpeg",
		"-ss", fmt.Sprintf("%.2f", offset),
		"-i", inputPath,
		"-frames:v", "1", "-q:v", "2",
		"-y", outputPath,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg thumbnail: %w\n%s", err, out)
	}
	return nil
}
