resource "aws_sqs_queue" "video_processing_results" {
  name                        = "video-processing-results.fifo"
  fifo_queue                  = true
  content_based_deduplication = true
  visibility_timeout_seconds  = 900

  tags = {
    Name = "video-processing-results"
  }
}

resource "aws_sqs_queue" "video_processing" {
  name                        = "video-processing.fifo"
  fifo_queue                  = true
  content_based_deduplication = true
  visibility_timeout_seconds  = 900

  tags = {
    Name = "video-processing"
  }
}

resource "aws_sqs_queue_policy" "video_processing" {
  queue_url = aws_sqs_queue.video_processing.id
  policy    = data.aws_iam_policy_document.sqs_s3_publish.json
}

data "aws_iam_policy_document" "sqs_s3_publish" {
  statement {
    principals {
      type        = "Service"
      identifiers = ["s3.amazonaws.com"]
    }
    actions   = ["sqs:SendMessage"]
    resources = [aws_sqs_queue.video_processing.arn]
    condition {
      test     = "ArnEquals"
      variable = "aws:SourceArn"
      values   = [aws_s3_bucket.video.arn]
    }
  }
}
