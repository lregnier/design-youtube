resource "aws_ecs_cluster" "main" {
  name = var.project

  setting {
    name  = "containerInsights"
    value = "disabled"
  }
}

resource "aws_lb" "api" {
  name               = "${var.project}-api-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = aws_subnet.public[*].id
}

resource "aws_lb_target_group" "api" {
  name        = "${var.project}-api-tg"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip"

  health_check {
    path                = "/health"
    healthy_threshold   = 2
    unhealthy_threshold = 3
    interval            = 30
  }
}

resource "aws_lb_listener" "api_http" {
  load_balancer_arn = aws_lb.api.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }
}

resource "aws_ecs_task_definition" "api" {
  family                   = "${var.project}-api"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn            = aws_iam_role.api_task.arn

  container_definitions = jsonencode([{
    name      = "api"
    image     = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.aws_region}.amazonaws.com/${var.project}-api:latest"
    essential = true
    portMappings = [{
      containerPort = 8080
      protocol      = "tcp"
    }]
    environment = [
      { name = "AWS_REGION", value = var.aws_region },
      { name = "DYNAMODB_TABLE", value = aws_dynamodb_table.videos.name },
      { name = "S3_BUCKET", value = aws_s3_bucket.video.bucket },
      { name = "CLOUDFRONT_DOMAIN", value = aws_cloudfront_distribution.video.domain_name },
      { name = "SQS_QUEUE_URL", value = aws_sqs_queue.video_processing.url },
      { name = "REDIS_ADDR", value = "${aws_elasticache_cluster.redis.cache_nodes[0].address}:6379" },
      { name = "RESULTS_QUEUE_URL", value = aws_sqs_queue.video_processing_results.url },
    ]
    secrets = [{
      name      = "UPLOAD_SECRET"
      valueFrom = aws_ssm_parameter.upload_secret.arn
    }]
    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = "/ecs/${var.project}-api"
        "awslogs-region"        = var.aws_region
        "awslogs-stream-prefix" = "api"
      }
    }
  }])
}

resource "aws_ecs_task_definition" "worker" {
  family                   = "${var.project}-worker"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "1024"
  memory                   = "2048"
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn            = aws_iam_role.worker_task.arn

  container_definitions = jsonencode([{
    name      = "worker"
    image     = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.aws_region}.amazonaws.com/${var.project}-worker:latest"
    essential = true
    environment = [
      { name = "AWS_REGION", value = var.aws_region },
      { name = "S3_BUCKET", value = aws_s3_bucket.video.bucket },
      { name = "CLOUDFRONT_DOMAIN", value = aws_cloudfront_distribution.video.domain_name },
      { name = "SQS_QUEUE_URL", value = aws_sqs_queue.video_processing.url },
      { name = "RESULTS_QUEUE_URL", value = aws_sqs_queue.video_processing_results.url },
    ]
    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = "/ecs/${var.project}-worker"
        "awslogs-region"        = var.aws_region
        "awslogs-stream-prefix" = "worker"
      }
    }
  }])
}

resource "aws_ecs_service" "api" {
  name            = "${var.project}-api"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.api.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.private[*].id
    security_groups  = [aws_security_group.api.id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.api.arn
    container_name   = "api"
    container_port   = 8080
  }

  depends_on = [aws_lb_listener.api_http]
}

resource "aws_cloudwatch_log_group" "api" {
  name              = "/ecs/${var.project}-api"
  retention_in_days = 7
}

resource "aws_cloudwatch_log_group" "worker" {
  name              = "/ecs/${var.project}-worker"
  retention_in_days = 7
}

data "aws_caller_identity" "current" {}

resource "aws_ssm_parameter" "upload_secret" {
  name  = "/${var.project}/upload-secret"
  type  = "SecureString"
  value = "changeme"

  lifecycle {
    ignore_changes = [value]
  }
}
