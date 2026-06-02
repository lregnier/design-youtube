resource "aws_dynamodb_table" "videos" {
  name         = "videos"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "videoId"

  attribute {
    name = "videoId"
    type = "S"
  }

  attribute {
    name = "status"
    type = "S"
  }

  attribute {
    name = "uploadedAt"
    type = "S"
  }

  global_secondary_index {
    name            = "status-uploadedAt-index"
    hash_key        = "status"
    range_key       = "uploadedAt"
    projection_type = "ALL"
  }
}
