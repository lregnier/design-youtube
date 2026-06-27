package main

import (
	"fmt"
	"os"
)

type Config struct {
	UploadSecret       string
	AWSRegion          string
	DynamoDBTable      string
	S3Bucket           string
	CloudFrontDomain   string
	SQSQueueURL        string
	ResultsQueueURL    string
	RedisAddr          string
	S3Endpoint         string
	S3PublicURL        string
	DynamoDBEndpoint   string
	SQSEndpoint        string
	CORSAllowedOrigins string
	HTTPAddr           string
}

func Load() (*Config, error) {
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
		RedisAddr:          os.Getenv("REDIS_ADDR"),
		S3Endpoint:         os.Getenv("S3_ENDPOINT_URL"),
		S3PublicURL:        os.Getenv("S3_PUBLIC_URL"),
		DynamoDBEndpoint:   os.Getenv("DYNAMODB_ENDPOINT_URL"),
		SQSEndpoint:        os.Getenv("SQS_ENDPOINT_URL"),
		CORSAllowedOrigins: os.Getenv("CORS_ALLOWED_ORIGINS"),
		HTTPAddr:           httpAddr,
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
