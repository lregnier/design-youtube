package config

import (
	"fmt"
	"os"
)

type Config struct {
	UploadSecret    string
	AWSRegion       string
	DynamoDBTable   string
	S3Bucket        string
	CloudFrontDomain string
	SQSQueueURL     string
	RedisAddr       string
}

func Load() (*Config, error) {
	c := &Config{
		UploadSecret:    os.Getenv("UPLOAD_SECRET"),
		AWSRegion:       os.Getenv("AWS_REGION"),
		DynamoDBTable:   os.Getenv("DYNAMODB_TABLE"),
		S3Bucket:        os.Getenv("S3_BUCKET"),
		CloudFrontDomain: os.Getenv("CLOUDFRONT_DOMAIN"),
		SQSQueueURL:     os.Getenv("SQS_QUEUE_URL"),
		RedisAddr:       os.Getenv("REDIS_ADDR"),
	}

	required := map[string]string{
		"UPLOAD_SECRET":     c.UploadSecret,
		"AWS_REGION":        c.AWSRegion,
		"DYNAMODB_TABLE":    c.DynamoDBTable,
		"S3_BUCKET":         c.S3Bucket,
		"CLOUDFRONT_DOMAIN": c.CloudFrontDomain,
		"SQS_QUEUE_URL":     c.SQSQueueURL,
		"REDIS_ADDR":        c.RedisAddr,
	}

	for name, val := range required {
		if val == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", name)
		}
	}

	return c, nil
}
