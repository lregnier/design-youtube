output "cloudfront_domain" {
  description = "CloudFront distribution domain for serving video assets"
  value       = aws_cloudfront_distribution.video.domain_name
}

output "alb_dns_name" {
  description = "ALB DNS name for the backend API"
  value       = aws_lb.api.dns_name
}
