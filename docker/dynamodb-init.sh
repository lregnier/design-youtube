#!/bin/sh
set -e

aws dynamodb describe-table \
  --endpoint-url http://dynamodb-local:8000 \
  --table-name videos 2>/dev/null && exit 0

aws dynamodb create-table \
  --endpoint-url http://dynamodb-local:8000 \
  --table-name videos \
  --attribute-definitions \
    AttributeName=videoId,AttributeType=S \
    AttributeName=status,AttributeType=S \
    AttributeName=uploadedAt,AttributeType=S \
  --key-schema AttributeName=videoId,KeyType=HASH \
  --global-secondary-indexes \
    '[{"IndexName":"status-uploadedAt-index","KeySchema":[{"AttributeName":"status","KeyType":"HASH"},{"AttributeName":"uploadedAt","KeyType":"RANGE"}],"Projection":{"ProjectionType":"ALL"}}]' \
  --billing-mode PAY_PER_REQUEST
