package config

import (
	"fmt"
	"os"
)

type Config struct {
	AWSRegion        string
	S3Bucket         string
	CloudFrontDomain string
	SQSQueueURL      string
	ResultsQueueURL  string
}

func Load() (*Config, error) {
	c := &Config{
		AWSRegion:        os.Getenv("AWS_REGION"),
		S3Bucket:         os.Getenv("S3_BUCKET"),
		CloudFrontDomain: os.Getenv("CLOUDFRONT_DOMAIN"),
		SQSQueueURL:      os.Getenv("SQS_QUEUE_URL"),
		ResultsQueueURL:  os.Getenv("RESULTS_QUEUE_URL"),
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

	return c, nil
}
