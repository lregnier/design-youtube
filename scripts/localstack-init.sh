#!/bin/bash
set -e

echo ">>> creating S3 bucket"
awslocal s3 mb s3://design-youtube-video-prod

echo ">>> configuring S3 bucket CORS"
awslocal s3api put-bucket-cors \
  --bucket design-youtube-video-prod \
  --cors-configuration '{
    "CORSRules": [{
      "AllowedHeaders": ["*"],
      "AllowedMethods": ["PUT", "GET", "HEAD"],
      "AllowedOrigins": ["*"],
      "ExposeHeaders": ["ETag"],
      "MaxAgeSeconds": 3600
    }]
  }'

echo ">>> creating DynamoDB table"
awslocal dynamodb create-table \
  --table-name videos \
  --attribute-definitions \
    AttributeName=videoId,AttributeType=S \
    AttributeName=status,AttributeType=S \
    AttributeName=uploadedAt,AttributeType=S \
  --key-schema AttributeName=videoId,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --global-secondary-indexes '[{
    "IndexName": "status-uploadedAt-index",
    "KeySchema": [
      {"AttributeName": "status", "KeyType": "HASH"},
      {"AttributeName": "uploadedAt", "KeyType": "RANGE"}
    ],
    "Projection": {"ProjectionType": "ALL"}
  }]'

echo ">>> creating SQS queues"
awslocal sqs create-queue \
  --queue-name video-processing.fifo \
  --attributes FifoQueue=true,ContentBasedDeduplication=true,VisibilityTimeout=900

awslocal sqs create-queue \
  --queue-name video-processing-results.fifo \
  --attributes FifoQueue=true,ContentBasedDeduplication=true,VisibilityTimeout=900

echo ">>> configuring S3 event notification → SQS"
QUEUE_ARN=$(awslocal sqs get-queue-attributes \
  --queue-url http://localhost:4566/000000000000/video-processing.fifo \
  --attribute-names QueueArn \
  --query 'Attributes.QueueArn' \
  --output text)

BUCKET_ARN="arn:aws:s3:::design-youtube-video-prod"

awslocal sqs set-queue-attributes \
  --queue-url http://localhost:4566/000000000000/video-processing.fifo \
  --attributes "{
    \"Policy\": \"{\\\"Version\\\":\\\"2012-10-17\\\",\\\"Statement\\\":[{\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":{\\\"Service\\\":\\\"s3.amazonaws.com\\\"},\\\"Action\\\":\\\"sqs:SendMessage\\\",\\\"Resource\\\":\\\"$QUEUE_ARN\\\",\\\"Condition\\\":{\\\"ArnEquals\\\":{\\\"aws:SourceArn\\\":\\\"$BUCKET_ARN\\\"}}}]}\"
  }"

awslocal s3api put-bucket-notification-configuration \
  --bucket design-youtube-video-prod \
  --notification-configuration "{
    \"QueueConfigurations\": [{
      \"QueueArn\": \"$QUEUE_ARN\",
      \"Events\": [\"s3:ObjectCreated:CompleteMultipartUpload\"],
      \"Filter\": {
        \"Key\": {
          \"FilterRules\": [{\"Name\": \"prefix\", \"Value\": \"raw/\"}]
        }
      }
    }]
  }"

echo ">>> localstack init done"
