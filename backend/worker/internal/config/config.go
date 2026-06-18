package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AWSRegion          string
	S3Bucket           string
	CloudFrontDomain   string
	SQSQueueURL        string
	ResultsQueueURL    string
	S3UsePathStyle     bool
	LocalStackEnabled  bool
	LocalStackEndpoint string
}

func Load() (*Config, error) {
	localStack, _ := strconv.ParseBool(os.Getenv("LOCALSTACK_ENABLED"))
	s3UsePathStyle, _ := strconv.ParseBool(os.Getenv("S3_USE_PATH_STYLE"))
	c := &Config{
		AWSRegion:          os.Getenv("AWS_REGION"),
		S3Bucket:           os.Getenv("S3_BUCKET"),
		CloudFrontDomain:   os.Getenv("CLOUDFRONT_DOMAIN"),
		SQSQueueURL:        os.Getenv("SQS_QUEUE_URL"),
		ResultsQueueURL:    os.Getenv("RESULTS_QUEUE_URL"),
		S3UsePathStyle:     s3UsePathStyle,
		LocalStackEnabled:  localStack,
		LocalStackEndpoint: os.Getenv("LOCALSTACK_ENDPOINT"),
	}

	required := map[string]string{
		"AWS_REGION":        c.AWSRegion,
		"S3_BUCKET":         c.S3Bucket,
		"CLOUDFRONT_DOMAIN": c.CloudFrontDomain,
		"SQS_QUEUE_URL":     c.SQSQueueURL,
		"RESULTS_QUEUE_URL": c.ResultsQueueURL,
	}

	for name, val := range required {
		if val == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", name)
		}
	}

	if c.LocalStackEnabled && c.LocalStackEndpoint == "" {
		return nil, fmt.Errorf("LOCALSTACK_ENDPOINT must be set when LOCALSTACK_ENABLED is true")
	}

	return c, nil
}
