package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	UploadSecret       string
	AWSRegion          string
	DynamoDBTable      string
	S3Bucket           string
	CloudFrontDomain   string
	SQSQueueURL        string
	ResultsQueueURL    string
	RedisAddr           string
	S3PublicEndpointURL string
	S3UsePathStyle      bool
	CORSAllowedOrigins  string
	HTTPAddr            string
}

func Load() (*Config, error) {
	s3UsePathStyle, _ := strconv.ParseBool(os.Getenv("S3_USE_PATH_STYLE"))
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}
	c := &Config{
		UploadSecret:       os.Getenv("UPLOAD_SECRET"),
		AWSRegion:          os.Getenv("AWS_REGION"),
		DynamoDBTable:      os.Getenv("DYNAMODB_TABLE"),
		S3Bucket:           os.Getenv("S3_BUCKET"),
		CloudFrontDomain:   os.Getenv("CLOUDFRONT_DOMAIN"),
		SQSQueueURL:        os.Getenv("SQS_QUEUE_URL"),
		ResultsQueueURL:    os.Getenv("RESULTS_QUEUE_URL"),
		RedisAddr:           os.Getenv("REDIS_ADDR"),
		S3PublicEndpointURL: os.Getenv("S3_PUBLIC_ENDPOINT_URL"),
		S3UsePathStyle:      s3UsePathStyle,
		CORSAllowedOrigins:  os.Getenv("CORS_ALLOWED_ORIGINS"),
		HTTPAddr:            httpAddr,
	}

	required := map[string]string{
		"UPLOAD_SECRET":        c.UploadSecret,
		"AWS_REGION":           c.AWSRegion,
		"DYNAMODB_TABLE":       c.DynamoDBTable,
		"S3_BUCKET":            c.S3Bucket,
		"CLOUDFRONT_DOMAIN":    c.CloudFrontDomain,
		"SQS_QUEUE_URL":        c.SQSQueueURL,
		"RESULTS_QUEUE_URL":    c.ResultsQueueURL,
		"REDIS_ADDR":           c.RedisAddr,
		"CORS_ALLOWED_ORIGINS": c.CORSAllowedOrigins,
	}

	for name, val := range required {
		if val == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", name)
		}
	}

	return c, nil
}
