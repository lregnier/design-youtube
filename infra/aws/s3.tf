resource "aws_s3_bucket" "video" {
  bucket = "${var.project}-video-${var.environment}"
}

resource "aws_s3_bucket_versioning" "video" {
  bucket = aws_s3_bucket.video.id
  versioning_configuration {
    status = "Disabled"
  }
}

resource "aws_s3_bucket_public_access_block" "video" {
  bucket                  = aws_s3_bucket.video.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_cors_configuration" "video" {
  bucket = aws_s3_bucket.video.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "GET", "HEAD"]
    allowed_origins = ["*"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3600
  }
}

resource "aws_s3_bucket_notification" "video_processing" {
  bucket = aws_s3_bucket.video.id

  queue {
    queue_arn     = aws_sqs_queue.video_processing.arn
    events        = ["s3:ObjectCreated:CompleteMultipartUpload"]
    filter_prefix = "raw/"
  }

  depends_on = [aws_sqs_queue_policy.video_processing]
}
