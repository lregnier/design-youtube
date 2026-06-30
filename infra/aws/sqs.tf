resource "aws_sqs_queue" "video_processing_results_dlq" {
  name                        = "video-processing-results-dlq.fifo"
  fifo_queue                  = true
  content_based_deduplication = true

  tags = {
    Name = "video-processing-results-dlq"
  }
}

resource "aws_sqs_queue" "video_processing_results" {
  name                        = "video-processing-results.fifo"
  fifo_queue                  = true
  content_based_deduplication = true
  visibility_timeout_seconds  = 900

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.video_processing_results_dlq.arn
    maxReceiveCount     = 3
  })

  tags = {
    Name = "video-processing-results"
  }
}

resource "aws_sqs_queue" "video_processing_dlq" {
  name                        = "video-processing-dlq.fifo"
  fifo_queue                  = true
  content_based_deduplication = true

  tags = {
    Name = "video-processing-dlq"
  }
}

resource "aws_sqs_queue" "video_processing" {
  name                        = "video-processing.fifo"
  fifo_queue                  = true
  content_based_deduplication = true
  visibility_timeout_seconds  = 900

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.video_processing_dlq.arn
    maxReceiveCount     = 3
  })

  tags = {
    Name = "video-processing"
  }
}

